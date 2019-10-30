package genetico

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"
)

type Chromosome struct {
	Chromosome   []uint8
	Aptitud      float64
	PointsValues []float64
}

type Generation struct {
	Population         []Chromosome
	aptFunc            aptFunc
	pointsFunc         pointsFunc
	maxNumberAlelo     int
	competePercentaje  float32
	sizeChromosome     int
	mutationPercentaje float32
	bitsMutate         int
	elitism            bool
}

type pointsFunc func(c []float64) []float64

type aptFunc func(p []float64) float64

func NewGenetic(
	population,
	sizeChromosome,
	maxNumberAlelo int,
	competePercentaje float32,
	mutationPercentaje float32,
	bitsMutate int,
	af aptFunc,
	pf pointsFunc,
) *Generation {
	ng := &Generation{
		Population:         make([]Chromosome, population),
		aptFunc:            af,
		pointsFunc:         pf,
		maxNumberAlelo:     maxNumberAlelo,
		competePercentaje:  competePercentaje,
		sizeChromosome:     sizeChromosome,
		mutationPercentaje: mutationPercentaje,
		bitsMutate:         bitsMutate,
		elitism:            true,
	}

	ng.generateRandom()
	ng.calculateAptitude()
	ng.sortPopulation()
	//ng.printPopulation()
	return ng
}

func (g *Generation) NextGeneration() *Generation {
	population := len(g.Population)
	ng := &Generation{
		Population:         make([]Chromosome, population),
		aptFunc:            g.aptFunc,
		pointsFunc:         g.pointsFunc,
		maxNumberAlelo:     g.maxNumberAlelo,
		competePercentaje:  g.competePercentaje,
		sizeChromosome:     g.sizeChromosome,
		mutationPercentaje: g.mutationPercentaje,
		bitsMutate:         g.bitsMutate,
		elitism:            g.elitism,
	}
	g.competeParents(ng)
	ng.mutate()
	//g.generateOperations(ng)
	ng.calculateAptitude()

	if g.elitism && g.ElitismRequired(ng) {
		g.sortPopulation()
		//g.printPopulation("OLD")
		ng.sortPopulation()
		//ng.printPopulation("NEW")
		fmt.Println("ELITISM")
		ng2 := Elitismo(g, ng)
		//ng2.printPopulation("NG2")
		return ng2
	}

	return ng
}

//
func (g *Generation) generateRandom() {
	//var wg sync.WaitGroup
	for i := range g.Population {
		//wg.Add(1)
		g.randomChromosome(i)
	}
	//wg.Wait()
}

func (g *Generation) randomChromosome(index int) {
	m := make(map[int]bool)
	var random int
	g.Population[index].Chromosome = make([]uint8, g.sizeChromosome)
	i := 0
	for i < g.sizeChromosome {
		rand.Seed(time.Now().UnixNano())
		random = rand.Intn(g.maxNumberAlelo) + 1
		if !m[random] {
			m[random] = true
			g.Population[index].Chromosome[i] = uint8(random)
			i++
		}
	}
}

func (g *Generation) competeSingle(newGeneration *Generation) {
	population := len(g.Population)
	var wg sync.WaitGroup

	for i := 0; i < population; i++ {
		wg.Add(1)
		go g.reproduceSingleChromosome(&newGeneration.Population[i], &wg)
	}

	wg.Wait()
}

func (g *Generation) reproduceSingleChromosome(
	newP *Chromosome,
	wg *sync.WaitGroup,
) {
	population := len(g.Population)
	p := float32(population) * g.competePercentaje
	bestApt := math.MaxFloat64
	var randomIndex, minIndex int
	for i := 0; i < int(p); i++ {
		rand.Seed(time.Now().UnixNano())
		randomIndex = rand.Intn(population)
		if g.Population[randomIndex].Aptitud < bestApt {
			bestApt = g.Population[randomIndex].Aptitud
			minIndex = randomIndex
		}
	}
	newChromosome := make([]uint8, g.sizeChromosome)
	copy(newChromosome, g.Population[minIndex].Chromosome)
	newP.Chromosome = newChromosome
	wg.Done()
}

func (g *Generation) competeParents(newGeneration *Generation) {
	population := len(g.Population)
	//var wg sync.WaitGroup

	for i := 0; i < population; i += 2 {
		//wg.Add(1)
		g.reproduceChildsChromosome(
			&newGeneration.Population[i],
			&newGeneration.Population[i+1],
			//&wg,
		)
	}

	//wg.Wait()
}
func (g *Generation) mutate() {
	population := len(g.Population)
	m := map[int]bool{}
	p := int(float32(population) * g.mutationPercentaje)
	var randomIndex int
	i := 0
	for i < p {
		rand.Seed(time.Now().UnixNano())
		randomIndex = rand.Intn(population)
		if !m[randomIndex] {
			Mutation(g.Population[randomIndex].Chromosome, g.bitsMutate)
			i++
		}

	}
}

func (g *Generation) reproduceChildsChromosome(
	newC1 *Chromosome,
	newC2 *Chromosome,
	//wg *sync.WaitGroup,
) {
	m := make(map[int]bool)
	population := len(g.Population)
	p := int(float32(population) * g.competePercentaje)
	bestApt1 := math.MaxFloat64
	bestApt2 := math.MaxFloat64
	var randomIndex, min1, min2 int
	i := 0
	for i < p {
		rand.Seed(time.Now().UnixNano())
		randomIndex = rand.Intn(population)
		if m[randomIndex] {
			continue
		}
		m[randomIndex] = true
		i++
		if g.Population[randomIndex].Aptitud < bestApt1 {
			bestApt1 = g.Population[randomIndex].Aptitud
			min1 = randomIndex
		}

	}

	i = 0
	m = make(map[int]bool)
	for i < p {
		rand.Seed(time.Now().UnixNano())
		randomIndex = rand.Intn(population)
		if m[randomIndex] {
			continue
		}
		if g.Population[randomIndex].Aptitud < bestApt2 {
			if min1 == randomIndex {
				continue
			}
			bestApt2 = g.Population[randomIndex].Aptitud
			min2 = randomIndex
		}
		m[randomIndex] = true
		i++
	}
	r1, r2 := Cruza(g.Population[min1].Chromosome, g.Population[min2].Chromosome)
	newC1.Chromosome = r1
	newC2.Chromosome = r2
	//wg.Done()
}

func (g *Generation) generateOperations() {
	var wg sync.WaitGroup
	countInversion := 0
	var r int
	for i := range g.Population {
		wg.Add(1)
		rand.Seed(time.Now().UnixNano())
		r = rand.Intn(2)
		if r == 0 && countInversion <= 50 {
			go Inversion(g.Population[i].Chromosome, &wg)
			countInversion++
		} else {
			go Intercambio(g.Population[i].Chromosome, &wg)
		}
	}
	wg.Wait()
}

func (g *Generation) calculatePoints() {
	for i, p := range g.Population {
		g.Population[i].PointsValues = g.pointsFunc(chromosomeRange(p.Chromosome))
	}
}

func chromosomeRange(c []uint8) []float64 {
	r := make([]float64, len(c))
	for i, v := range c {
		r[i] = float64(v / 3)
		if r[i] < 1 {
			r[i] += 1
		}
	}
	//for i := range c {
	//	r[i] = float64(c[i])
	//}
	return r
}

func (g *Generation) calculateAptitude() {
	g.calculatePoints()
	for i, p := range g.Population {
		g.Population[i].Aptitud = g.aptFunc(p.PointsValues)
	}
}

func (g *Generation) ElitismRequired(ng *Generation) bool {
	min := math.MaxFloat64
	m := map[float64]int{}
	popSize := len(g.Population)
	for i := 0; i < popSize; i++ {
		if min > g.Population[i].Aptitud {
			min = g.Population[i].Aptitud
		}
		if min > ng.Population[i].Aptitud {
			min = ng.Population[i].Aptitud
		}
		m[g.Population[i].Aptitud]++
		m[ng.Population[i].Aptitud]++
	}

	var max int
	for _, v := range m {
		if v > max {
			max = v
		}
	}
	g.elitism = (float64(max)/float64(popSize*2)) < 0.2 || m[min] < 3
	return g.elitism
}

func (g *Generation) printPopulation(head string) {
	f, err := os.OpenFile("dat.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("======%s===========\n", head)); err != nil {
		panic(err)
	}
	for i, v := range g.Population {
		_, err = f.WriteString(fmt.Sprintf("%d\t%v\t%v\n", i, v.Chromosome, v.Aptitud))
		if err != nil {
			panic(err)
		}
	}
	f.WriteString(fmt.Sprint("==============================\n"))
}

func (g *Generation) sortPopulation() {
	sort.SliceStable(g.Population, func(i, j int) bool {
		return g.Population[i].Aptitud < g.Population[j].Aptitud
	})
}
