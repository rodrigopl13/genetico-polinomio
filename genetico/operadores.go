package genetico

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func Inversion(a []uint8, wg *sync.WaitGroup) {
	rand.Seed(time.Now().UnixNano())
	start := rand.Intn(len(a))
	size := rand.Intn(len(a) - 2)
	i := 0
	var tmp uint8
	var j, k int
	for i <= (size / 2) {
		j = start + i
		if j > len(a)-1 {
			j = j - len(a)
		}
		k = start + size - i
		if k > len(a)-1 {
			k = k - len(a)
		}
		tmp = a[k]
		a[k] = a[j]
		a[j] = tmp
		i++
	}
	wg.Done()
}

func Intercambio(a []uint8, wg *sync.WaitGroup) {
	rand.Seed(time.Now().UnixNano())
	size := rand.Intn((len(a) / 2)) + 1

	rand.Seed(time.Now().UnixNano())
	pos1 := rand.Intn((len(a)) - (size * 2) + 1)

	rand.Seed(time.Now().UnixNano())
	pos2 := rand.Intn(len(a)-pos1-(size*2)+1) + pos1 + size

	var tmp uint8
	for i := 0; i < size; i++ {
		tmp = a[pos1+i]
		a[pos1+i] = a[pos2+i]
		a[pos2+i] = tmp
	}
	wg.Done()
}

func Cruza(c1, c2 []uint8) ([]uint8, []uint8) {
	rand.Seed(time.Now().UnixNano())
	slicePoint := rand.Intn((len(c1)*8)-1) + 1
	var s1, s2 string
	for _, v := range c1 {
		s1 += fmt.Sprintf("%08b", v)
	}
	for _, v := range c2 {
		s2 += fmt.Sprintf("%08b", v)
	}
	t1 := fmt.Sprintf("%s%s", s1[:slicePoint], s2[slicePoint:])
	t2 := fmt.Sprintf("%s%s", s2[:slicePoint], s1[slicePoint:])
	r1 := make([]uint8, len(c1))
	r2 := make([]uint8, len(c2))
	for i := 0; i < len(t1)/8; i++ {
		if v, err := strconv.ParseUint(t1[i*8:(i+1)*8], 2, 8); err == nil {
			r1[i] = uint8(v)
		} else {
			fmt.Println("r1", err)
		}
		if v, err := strconv.ParseUint(t2[i*8:(i+1)*8], 2, 8); err == nil {
			r2[i] = uint8(v)
		} else {
			fmt.Println("r2", err)
		}
	}
	removeZero(r1)
	removeZero(r2)
	return r1, r2
}

func Mutation(c []uint8, numberBits int) {
	m := map[int]bool{}
	totalBits := 8 * len(c)
	i := 0
	for i < numberBits {
		rand.Seed(time.Now().UnixNano())
		bitRandom := rand.Intn(totalBits)
		if m[bitRandom] {
			continue
		}
		m[bitRandom] = true
		index := bitRandom / 8
		b := []byte(fmt.Sprintf("%08b", c[index]))
		if b[bitRandom%8] == 48 {
			b[bitRandom%8] = 49
		} else {
			b[bitRandom%8] = 48
		}
		if v, err := strconv.ParseUint(string(b), 2, 8); err == nil {
			c[index] = uint8(v)
		}
		i++
	}
	removeZero(c)
}

func Elitismo(g1, g2 *Generation) *Generation {
	population := len(g1.Population)
	var randomIndex int
	mIndex := make(map[int]bool)
	ng := &Generation{
		Population:         make([]Chromosome, population),
		aptFunc:            g1.aptFunc,
		pointsFunc:         g1.pointsFunc,
		maxNumberAlelo:     g1.maxNumberAlelo,
		competePercentaje:  g1.competePercentaje,
		sizeChromosome:     g1.sizeChromosome,
		mutationPercentaje: g1.mutationPercentaje,
		bitsMutate:         g1.bitsMutate,
		elitism:            g1.elitism,
	}
	var c Chromosome
	indexG1, indexG2 := 0, 0
	i := 0
	rand.Seed(time.Now().UnixNano())
	randomIndex = rand.Intn(population)
	for i < population {
		if g1.Population[indexG1].Aptitud <= g2.Population[indexG2].Aptitud {
			c = g1.Population[indexG1]
			indexG1++
		} else {
			c = g2.Population[indexG2]
			indexG2++
		}
		for mIndex[randomIndex] {
			rand.Seed(time.Now().UnixNano())
			randomIndex = rand.Intn(population)
		}
		ng.Population[randomIndex] = c
		mIndex[randomIndex] = true
		i++
	}
	return ng
}

func removeZero(c []uint8) {
	for i, v := range c {
		if v == 0 {
			c[i] = 1
		}
	}
}
