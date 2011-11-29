package ga

import (
	"./distance_matrix"
)

type Gene struct {
	Data []int	
	Fitness float64
}

func (gene *Gene) fitness(distMatrix *util.DistanceMatrix) float64 {
	var fitness float64 = 0.0
	var length int = len(gene.Data)
	for i := 0; i < length-1; i++ {
		fieldX := gene.Data[i]-1
		fieldY := gene.Data[i+1]-1		
		fitness += distMatrix.GetDistance(fieldX, fieldY)
	}
	fitness += distMatrix.GetDistance(gene.Data[length-1]-1, gene.Data[0]-1)
	gene.Fitness = fitness
	return fitness
}

type GeneSlice []Gene

func (g GeneSlice) Len() int {
	return len(g)
}

func (g GeneSlice) Less(i,j int) bool {
	return g[i].Fitness < g[j].Fitness
}

func (g GeneSlice) Swap(i,j int) {
	g[i], g[j] = g[j], g[i]
}

func (g GeneSlice) CalculateFitness(distMatrix *util.DistanceMatrix) {	
    for i := 0; i<len(g); i++ {
		g[i].fitness(distMatrix)
    }
}

func (g GeneSlice) BestFitness() (fitness float64) {
	fitness = g[0].Fitness
	for i:=1; i < len(g); i++ {
		if g[i].Fitness < fitness {
			fitness = g[i].Fitness
		}
	}
	return
}