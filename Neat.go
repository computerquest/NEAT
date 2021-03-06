package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

//MAX 1000 innovation
/*
not going to speciate until after a couple of rounds
*/

type Neat struct {
	nodeMutate           float64   //odds for node mutation
	network              []Network //stores networks in species
	connectionInnovation [][]int   //stores innovation number and connection to and from ex: 1, fromNode:2, toNode: 5 [2,5]
	speciesThreshold     float64   //could adjust based upon average difference between networks
	species              []Species
	speciesId            int
}

func GetNeatInstance(numNetworks int, input int, output int, mutate float64, lr float64) Neat {
	n := Neat{speciesThreshold: .01,
		nodeMutate: mutate, network: make([]Network, numNetworks), connectionInnovation: make([][]int, input*output, 1000), species: make([]Species, 0, 5)}

	for i := output; i < input+output; i++ {
		for a := 0; a < output; a++ {
			n.connectionInnovation[(i-output)*output+a] = []int{i, a}
		}
	}

	for i := 0; i < len(n.network); i++ {
		n.network[i] = GetNetworkInstance(input, output, i, 0, lr, true)
	}

	return n
}

//this initializes all the species and mutates everything. this is a second method because the first does not return a reference.
func (n *Neat) initialize() {
	n.createSpecies(n.network[0 : len(n.network)%5+(len(n.network)/5)])
	for i := len(n.network)%5 + (len(n.network) / 5); i+(len(n.network)/5) <= len(n.network); i += (len(n.network) / 5) {
		n.createSpecies(n.network[i : i+(len(n.network)/5)])
	}

	for i := 0; i < len(n.species); i++ {
		for a := 0; a < len(n.species[i].network); a++ {
			n.species[i].mutateNetwork(n.species[i].network[a])
		}
	}
	n.speciateAll()
	n.checkSpecies()
}

//this is the actual neat training method. You provide the input, number of strikes, and the target accuracy
func (n *Neat) start(input [][][]float64, cutoff int, target float64) Network {
	strikes := cutoff
	var bestNet Network
	bestFit := 0.0
	var wg sync.WaitGroup

	//initial training
	for i := 0; i < len(n.species); i++ {
		wg.Add(1)
		go n.species[i].trainNetworks(input, &wg)
	}

	wg.Wait()

	for z := 0; strikes > 0 && bestFit < target; z++ {
		//mates
		for i := 0; i < len(n.species); i++ {
			wg.Add(1)
			go n.species[i].mateSpecies(&wg)
		}

		wg.Wait()

		//trains
		for i := 0; i < len(n.species); i++ {
			wg.Add(1)
			go n.species[i].trainNetworks(input, &wg)
		}
		wg.Wait()

		if z%5 == 0 {
			n.speciateAll()
			n.checkSpecies()
		}

		//determines the best
		bestIndex := -1
		for i := 0; i < len(n.network); i++ {
			if bestFit < n.network[i].fitness {
				bestFit = n.network[i].fitness
				bestIndex = i
			}
		}

		//compares the best
		if bestIndex != -1 {
			bestNet = clone(&n.network[bestIndex])
			strikes = cutoff
		} else {
			strikes--
			n.mutatePopulation()
			if z%5 != 0 {
				n.speciateAll()
				n.checkSpecies()
			}
		}

		fmt.Println("epoch:", z, bestFit)
	}

	return bestNet
}
func (n *Neat) printNeat() {
	fmt.Println()
	for i := 0; i < len(n.species); i++ {
		n.species[i].sortInnovation()
		fmt.Println("species id: ", n.species[i].id, " innovations: ", n.species[i].commonInnovation, " net: ", len(n.species[i].network))
		for a := 0; a < len(n.species[i].network); a++ {
			printNetwork(n.species[i].network[a])
			/*this prints the actual connections
			fmt.Println("expected connection: ")
			for b := 0; b < len(n.species[i].network[a].innovation); b++ {
				fmt.Print(n.getInnovation(n.species[i].network[a].innovation[b]))
			}*/
		}
	}
}

//mutates part of the population
func (n *Neat) mutatePopulation() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numNet := int((r.Int63n(int64(len(n.network)))-3)/5) + 3
	for i := 0; i < numNet; i++ {
		species := int(r.Int63n(int64(len(n.species))))

		n.species[species].mutateNetwork(n.species[species].getNetworkAt(int(r.Int63n(int64(len(n.species[species].network))))))
	}
}

/////////////////////////////////////////////////////////////SPECIATION
func (n *Neat) speciateAll() {
	for i := 0; i < len(n.network); i++ {
		n.speciate(&n.network[i])
	}
}
func (n *Neat) checkSpecies() {
	for i := 0; i < len(n.species); i++ {
		values := make([]float64, len(n.species))
		for a := 0; a < len(n.species); a++ {
			if a == i {
				values[a] = 100.0
				continue
			}

			values[a] = compareGenome(n.species[i].avgNode(), n.species[i].commonInnovation, n.species[a].avgNode(), n.species[a].commonInnovation)
		}

		lValue := 1000.0
		for a := 0; a < len(values); a++ {
			if values[a] < lValue {
				lValue = values[a]
			}
		}

		if lValue < n.speciesThreshold || len(n.species[i].network) < 2 { //switched direction if sign because %dif < difthreshold for it to be the same
			n.removeSpecies(n.species[i].id) //could say continue if similar so that the smaller does the hard workbut might screw with eliminating empties
			i--
		}
	}
}
func (n *Neat) speciate(network *Network) {
	values := make([]float64, len(n.species))

	for i := 0; i < len(n.species); i++ {
		values[i] = compareGenome(len(network.nodeList), network.innovation, n.species[i].avgNode(), n.species[i].commonInnovation)
	}

	//this should be faster than sorting the whole thing (it also retains position information)
	bestSpec := -1
	lValue := 1000.0
	for i := 0; i < len(values); i++ {
		if values[i] < lValue {
			bestSpec = n.species[i].id
			lValue = values[i]
		}
	}

	//s := n.getSpecies(network.species)
	if lValue > n.speciesThreshold { //&& s != nil && len(s.network) > 2 { //i flipped this sign i think it works better %different > differentThreshold
		//finds the position
		networkIndex := 0
		for i := 0; i < len(n.network); i++ {
			if n.network[i].networkId == network.networkId {
				networkIndex = i
			}
		}

		lastSpec := network.species
		newSpec := n.createSpecies(n.network[networkIndex : networkIndex+1])

		//remove from the old species
		s := n.getSpecies(lastSpec)
		if s != nil {
			//removes current and checks to see if the rest need to be speciated
			s.removeNetwork(network.networkId)
			for i := 0; i < len(s.network); i++ {
				if s.network[i].networkId != network.networkId && s.network[i].species == s.id {
					if compareGenome(len(s.network[i].nodeList), s.network[i].innovation, s.avgNode(), s.commonInnovation) > compareGenome(len(s.network[i].nodeList), s.network[i].innovation, newSpec.avgNode(), newSpec.commonInnovation) {
						newSpec.addNetwork(s.network[i])
						s.removeNetwork(s.network[i].networkId)
						i--
					}
				}
			}
		}

		//checks to see if new species meets size requirement
		if len(newSpec.network) < 2 {
			//reassign creator to next best in order to prevent a loop
			newSpec.removeNetwork(network.networkId)
			n.getSpecies(bestSpec).addNetwork(network) //could be problem because index changes when make new species (maybe because should be added to the end)

			n.removeSpecies(newSpec.id)
		}
	} else if network.species != bestSpec {
		lastSpec := n.getSpecies(network.species)
		n.getSpecies(bestSpec).addNetwork(network)

		if lastSpec != nil {
			lastSpec.removeNetwork(network.networkId)
		}
	}
}
func compareGenome(node int, innovation []int, nodeA int, innovationA []int) float64 {
	var larger []int
	var smaller []int

	if len(innovation) > len(innovationA) {
		larger = innovation
		smaller = innovationA
	} else {
		larger = innovationA
		smaller = innovation
	}

	missing := 0
	for b := 0; b < len(larger); b++ {
		if sort.SearchInts(smaller, larger[b]) == len(smaller) {
			missing++
		}
	}

	return float64(missing+int(math.Abs(float64(node-nodeA)))) / (float64(len(smaller)) + float64((node+nodeA)/2))
}

//////////////////////////////////////////////////////////INNOVATON
func (n *Neat) getInnovation(num int) []int {
	return n.connectionInnovation[num]
}
func (n *Neat) findInnovationNum(search []int) int {
	for i := 0; i < len(n.connectionInnovation); i++ {
		if n.connectionInnovation[i][0] == search[0] && n.connectionInnovation[i][1] == search[1] {
			return i
		}
	}

	return -1
}

//////////////////////////////////////////////////////SPECIES
func (n *Neat) getSpecies(id int) *Species {
	for i := 0; i < len(n.species); i++ {
		if n.species[i].id == id {
			return &n.species[i]
		}
	}

	return nil
}
func (n *Neat) createSpecies(possible []Network) *Species {
	for i := 0; i < len(possible); i++ {
		possible[i].species = n.speciesId
	}

	s := GetSpeciesInstance(n.speciesId, possible, &n.connectionInnovation, n.nodeMutate)
	if cap(n.species) <= len(n.species) {
		n.species = append(n.species, s)
	} else {
		n.species = n.species[0 : len(n.species)+1]
		n.species[len(n.species)-1] = s
	}

	n.speciesId++

	return &n.species[len(n.species)-1]
}
func (n *Neat) removeSpecies(id int) {
	for i := 0; i < len(n.species); i++ {
		if n.species[i].id == id {
			currentSpecies := n.species[i].network

			n.species = append(n.species[:i], n.species[i+1:]...)
			for a := 0; a < len(currentSpecies); a++ {
				if currentSpecies[a].species == id {
					n.speciate(currentSpecies[a])
				}
			}
		}
	}
}
