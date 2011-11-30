package ga

import (
	"fmt"	
	"math"
	"sort"
	"math/rand"
	"../distance_matrix"
	"../gene"
)

type TravellingSalesman struct {
	distMatrix	util.DistanceMatrix
	genes		ga.GeneSlice
	pMutate		float64
	pCrossover	float64
	genCount	int
	geneLength  int
	nGenes		int
}

func NewGA(pMutate, pCrossover float64, nGenes int) *TravellingSalesman {
	ga := new(TravellingSalesman)
	ga.pMutate = pMutate
	ga.pCrossover = pCrossover
	ga.genCount = 0
	ga.nGenes = nGenes
	return ga
}

/*
 * matrix is the input matrix of integers
 *
 */
func (ts *TravellingSalesman) Init(matrix [][]int) {
	// Find locations and coordinates
	coords := map[int] [2]int{}
	for i := range matrix { 
		for j := range matrix[i] {			
			if matrix[i][j] != 0 {
				value := matrix[i][j]
				coords[value] = [2]int{i+1,j+1}
			}
		}
	}
	
	// Gene length equals the number of coordinates
	ts.geneLength = len(coords)
	
	// Calculate distances
	var size = ts.geneLength
	ts.distMatrix = util.NewDistanceMatrix(ts.geneLength)
	for i := 0; i < size*size; i++ {	
		row := i%size
		col := i/size
		area1 := coords[row+1] // X coordinate (Row)
		area2 := coords[col+1] // Y coordinate (Column)
		if area1[0] == area2[0] {
			// Same Column so distance eq diff(row2-row1)
			ts.distMatrix.SetDistance(
				row, 
				col, 
				math.Abs(float64(area2[1]-area1[1])))
		} else if area1[1] == area2[1] {
			// Same Row so distance eq diff(col2-col1)
			ts.distMatrix.SetDistance(
				row,
				col,
				math.Abs(float64(area2[0]-area1[0])))
		} else {
			a := math.Pow(float64(area2[1]-area1[1]), 2)
			b := math.Pow(float64(area2[0]-area1[0]), 2)
			c := math.Sqrt(a+b)
			ts.distMatrix.SetDistance(row, col, c)
		}
		//fmt.Printf("%d -> %d: %f\n", row+1, col+1, ts.distMatrix.GetDistance(row,col))
	}
	
	// Generate random genes
	ts.genes = make([]ga.Gene, ts.nGenes)
	var proto = make([]int, ts.geneLength)
	var i int = 0
	for key, _ := range coords {
		proto[i] = key
		i++
	}
	
	// Initialize all len(genes) and shuffle
	for i := 0; i < len(ts.genes); i++ {
		ts.genes[i].Data = make([]int, len(proto))		
		copy(ts.genes[i].Data, proto)
		shuffleArray(&(ts.genes[i].Data))
	}
}

func (ga *TravellingSalesman) String() string {
	return fmt.Sprintf("%f\t%f\t%d\n", ga.pCrossover, ga.pMutate, ga.genCount)
}

func (ga *TravellingSalesman) CurrentGeneration() int {
	return ga.genCount
}

func (ts *TravellingSalesman) Crossover() {
	// Crossover		
	var nCrossover = int(ts.pCrossover * float64(ts.nGenes))
	var newGenes = make([]ga.Gene, nCrossover)
	for i:=0; i < nCrossover; i++ {
		newGenes[i].Data = make([]int, ts.geneLength)
		n := rand.Intn(ts.nGenes)
		m := rand.Intn(ts.nGenes)
		
		currentCity := ts.genes[n].Data[0];
		newGenes[i].Data[0] = currentCity
		for k:=1; k < ts.geneLength; k++ {
			nextN := findNextCity(&(ts.genes[n].Data), currentCity)
			nextM := findNextCity(&(ts.genes[m].Data), currentCity)
			
			existN := isInArray(&(newGenes[i].Data), nextN)
			existM := isInArray(&(newGenes[i].Data), nextM)
			
			// n exists, m doesnt -> take m
			if existN && !existM {
				newGenes[i].Data[k] = nextM
				currentCity = nextM
			} else if !existN && existM {
				newGenes[i].Data[k] = nextN
				currentCity = nextN
			} else if existN && existM {
				nextRandom := findNextRandomCity(newGenes[i].Data[0:k], &(ts.genes[n].Data))
				newGenes[i].Data[k] = nextRandom
				currentCity = nextRandom
			} else {
				// If both didn't exist, take the shorter one				
				distN := ts.distMatrix.GetDistance(currentCity-1, nextN-1)
				distM := ts.distMatrix.GetDistance(currentCity-1, nextM-1)
				
				// Take the shorter route
				if distN < distM {
					newGenes[i].Data[k] = nextN
					currentCity = nextN
				} else {
					newGenes[i].Data[k] = nextM
					currentCity = nextM
				}
			}
		}
	}
		
	copy(ts.genes, newGenes)
}

func (ts *TravellingSalesman) Mutate() {
	// Mutate
	var nMutations = int(ts.pMutate * float64(ts.nGenes) * float64(ts.geneLength))
	//fmt.Printf("Number of Mutations: %d\n", nMutations)
	for i:=0; i < nMutations; i++ {
		n := rand.Intn(ts.nGenes)
		x := rand.Intn(ts.geneLength)
		y := rand.Intn(ts.geneLength)
		sort.IntSlice(ts.genes[n].Data).Swap(x,y)
	}
}

func (ts *TravellingSalesman) Replicate() {
	// Replicate
	// Sort by fitness		
	sort.Sort(ts.genes)
	
	for i:=0; i < (ts.nGenes/10)-1; i++ {
		lower := (i+1)*10
		upper := (i+2)*10
		copy(ts.genes[lower:upper], ts.genes[0:10])
	}
	
	/*
	for k,v := range ts.genes {
		fmt.Printf("Gene %3d, Fitness: %f\n", k, v.Fitness)
	}
	*/
}

func (ts *TravellingSalesman) Run(maxGenerations int) {
	var bestFitness float64 = float64(ts.geneLength+100)
	for int(bestFitness) > ts.geneLength && ts.genCount < maxGenerations {
		
		ts.Crossover()
		
		// Calculate current fitness for all Genes		
		ga.GeneSlice(ts.genes).CalculateFitness(&ts.distMatrix)
		bestFitness = ga.GeneSlice(ts.genes).BestFitness()
		
		ts.Mutate()
		
		// Calculate current fitness for all Genes
		ga.GeneSlice(ts.genes).CalculateFitness(&ts.distMatrix)
		bestFitness = ga.GeneSlice(ts.genes).BestFitness()
		
		ts.Replicate()

		ts.genCount++
	}
}

// Helper functions
func findNextCity(arr *[]int, currentCity int) int {
	var position int = 0
	for ; position < len(*arr); position++ {
		if (*arr)[position] == currentCity {
			break
		}		
	}	
	length := len(*arr)
	if position < length && (*arr)[position] == currentCity {		
		return (*arr)[(position+1)%length]
	}
	return -1
}

func findNextRandomCity(exist []int, comp *[]int) (city int) {
	for i:=0; i < len(*comp); i++ {
		cityExists := isInArray(&exist, (*comp)[i])
		if !cityExists {
			city = (*comp)[i]
			return
		}
	}
	
	return
}

func shuffleArray(arr *[]int) {
	// Shuffle the array len(*arr) times
	for i:=0; i < len(*arr); i++ {
		j := rand.Intn(len(*arr))
		(*arr)[i], (*arr)[j] = (*arr)[j], (*arr)[i]
	}
}

func isInArray(arr *[]int, value int) (exists bool) {
	exists = false
	for i:=0; i<len(*arr); i++ {
		if (*arr)[i] == value {
			exists = true
			return
		}
	}
	return
}
