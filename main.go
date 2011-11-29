// TravellingSalesman project main.go
package main

import (
	"fmt"
	"flag"
	"os"
	"io"
	"bufio"
	"strings"
	"math/rand"
	"strconv"
	"time"
	"runtime"
	"./ga/travelling_salesman"
)

// Program parameters
var fileName *string = flag.String("file", "", "File to parse")
var pc *float64 = flag.Float64("pc", 0.5, "Crossover probability")
var pm *float64 = flag.Float64("pm", 0.01, "Mutation probability")
var ngenes *int = flag.Int("ngenes", 100, "Number of genes")
var maxGenerations *int = flag.Int("maxgener", 200, "Max. Generations")
var protectBest *bool = flag.Bool("protect_best", false, "Protect best (true/false)")
var runs *int = flag.Int("runs", 10, "Number of runs to average over")
var cores *int = flag.Int("cores", 1, "Number of cores for processing")


func main() {
	flag.Parse()
	
	// Explicitly tell go to use *cores core
	runtime.GOMAXPROCS(*cores)
	
	if *fileName != "" {
		fmt.Errorf("Filename not given!\n")
	}
	
	// Initialize random seed
	rand.Seed(time.Nanoseconds())
	
	//fmt.Println("Travelling Salesman!")
	
	f, err := os.Open(*fileName)	
	defer f.Close()
	
    if f == nil {
        fmt.Printf("can't open file; err=%s\n", err)
        os.Exit(1)
    }

	// Generate input matrix
	var matrix = make([][]int, 0)
	reader := bufio.NewReader(f)	
	var line string
	line, err = reader.ReadString('\n')	
	for err != io.EOF {
		values := strings.Fields(line)
		ints := getIntArray(values)
		matrix = append(matrix, ints)
		line, err = reader.ReadString('\n')
	}		
	
	var generationSum int = 0
	/*
	for i:=0; i < *runs; i++ {
		salesman := ga.NewGA(*pm, *pc, *ngenes)
		salesman.Init(matrix)
		salesman.Run(*maxGenerations)
		generationSum += salesman.CurrentGeneration()	
	}*/
	
	c := make(chan int)
	for i:=0; i < *runs; i++ {
		go func() {
			salesman := ga.NewGA(*pm, *pc, *ngenes)
			salesman.Init(matrix)
			salesman.Run(*maxGenerations)
			c <- salesman.CurrentGeneration()
		}()
	}
	
	for i:=0; i < *runs; i++ {
		generationSum += <-c
	}
	fmt.Printf("%2f\t%2f\t%d\n", *pc, *pm, generationSum/(*runs))
}

func getIntArray(values []string) []int {
	intArray := make([]int, len(values))
	for i:= range values {
		intArray[i], _ = strconv.Atoi(values[i])
	}
	return intArray
}