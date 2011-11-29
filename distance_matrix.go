package util

type DistanceMatrix struct {
    data []float64
	size int
}

func (d DistanceMatrix) SetDistance(x,y int, distance float64) {
	d.data[x*d.size+y] = distance
}

func (d DistanceMatrix) GetDistance(x,y int) float64 {
	return d.data[x*d.size+y]
}

func NewDistanceMatrix(size int) (DistanceMatrix) {
	return DistanceMatrix{make([]float64, size*size), size}
}