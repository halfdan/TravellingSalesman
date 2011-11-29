// GeneticAlgorithm.go
package ga

type GeneticAlgorithm interface {	
	Init([][]int)
	Mutate()
	Crossover()
	Replicate()
	Run()
}
