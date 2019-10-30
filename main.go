package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"genetico-polinomio/genetico"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math"
	"time"
)

type application struct {
	evolution         *canvas.Image
	equations         *canvas.Image
	window            fyne.Window
	generation        *genetico.Generation
	labelBest         *widget.Label
	plotEvolution     *plot.Plot
	xyEvolution       plotter.XYs
	xyCorrect         plotter.XYs
	bestChromosome    int
	bestSolution      float64
	bestAllChromosome []uint8
	bestAllSolution   float64
}

func main() {
	a := app.New()
	g := application{
		evolution: &canvas.Image{
			FillMode: canvas.ImageFillOriginal,
		},
		equations: &canvas.Image{
			FillMode: canvas.ImageFillOriginal,
		},
		window:          a.NewWindow("Problema Polinomio"),
		generation:      &genetico.Generation{},
		labelBest:       widget.NewLabel(fmt.Sprintf("%.3f", 0.0)),
		bestAllSolution: math.MaxFloat64,
	}
	g.window.SetContent(
		widget.NewVBox(
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2),
				g.evolution,
				g.equations,
			),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(3),
				widget.NewButton("Iniciar Applicacion", g.startApp),
				widget.NewLabel("Best Solution"),
				g.labelBest,
			),
		),
	)
	g.window.ShowAndRun()

}

func (a *application) startApp() {
	var err error
	a.plotEvolution, err = plot.New()
	if err != nil {
		panic(err)
	}

	correctSolution := pointsFunc([]float64{10, 25, 3, 70, 10, 6})

	a.plotEvolution.Title.Text = "Evolution"
	a.plotEvolution.Y.Label.Text = "Aptitude function"
	a.plotEvolution.X.Label.Text = "Generation"

	a.evolution.File = "image/evolution.png"
	a.equations.File = "image/equations.png"

	a.xyEvolution = plotter.XYs{}
	a.xyCorrect = createPlotXY(correctSolution)

	a.generation = genetico.NewGenetic(
		1000,
		6,
		255,
		0.1,
		0.8,
		3,
		getAptFunc(correctSolution),
		pointsFunc)
	//fmt.Println(a.generation.Population)
	time.Sleep(1 * time.Second)
	a.createGraph(0)
	for i := 0; i < 100; i++ {
		a.generation = a.generation.NextGeneration()
		//fmt.Println(a.generation.Population)
		//time.Sleep(1 * time.Second)
		a.createGraph(i + 1)
	}
}

func (a *application) createGraph(i int) {
	plotEquations, err := plot.New()
	if err != nil {
		panic(err)
	}

	plotEquations.Title.Text = "Equations"
	plotEquations.Y.Label.Text = "Y"
	plotEquations.X.Label.Text = "X"

	a.getBestChromosome()
	points := plotter.XY{}
	points.X = float64(i)
	points.Y = a.bestSolution

	a.xyEvolution = append(a.xyEvolution, points)
	err = plotutil.AddLinePoints(a.plotEvolution, a.xyEvolution)
	if err != nil {
		panic(err)
	}

	if err = a.plotEvolution.Save(7*vg.Inch, 7*vg.Inch, "image/evolution.png"); err != nil {
		panic(err)
	}

	err = plotutil.AddLinePoints(
		plotEquations,
		"Correct",
		a.xyCorrect,
		"Genetic",
		createPlotXY(a.generation.Population[a.bestChromosome].PointsValues),
	)
	if err != nil {
		panic(err)
	}

	if err := plotEquations.Save(7*vg.Inch, 7*vg.Inch, "image/equations.png"); err != nil {
		panic(err)
	}

	canvas.Refresh(a.evolution)
	canvas.Refresh(a.equations)
	a.labelBest.SetText(fmt.Sprintf("%.3f", a.bestAllSolution))
}

func (a *application) getBestChromosome() {
	minSolution := math.MaxFloat64
	index := 0
	for i, p := range a.generation.Population {
		if minSolution > p.Aptitud {
			index = i
			minSolution = p.Aptitud
		}
	}
	a.bestSolution = minSolution
	a.bestChromosome = index
	if a.bestAllSolution > minSolution {
		a.bestAllChromosome = a.generation.Population[index].Chromosome
		a.bestAllSolution = minSolution
	}

	fmt.Println("BEST chromosome >>>>>>>>>>>>", a.generation.Population[index].Chromosome)
	fmt.Println("BEST Aptitude $$$$$$$$$$$$$$", a.generation.Population[index].Aptitud)
	fmt.Println("minSolution ################", fmt.Sprintf("%.3f", minSolution))
}

func pointsFunc(ch []float64) []float64 {
	a := ch[0]
	b := ch[1]
	c := ch[2]
	d := ch[3]
	e := ch[4]
	f := ch[5]

	y := make([]float64, 1000)
	var x float64
	for i := 0; i < len(y); i++ {
		x = float64(i) / 10
		y[i] = a*(b*math.Sin(x/c)+
			d*math.Cos(x/e)) + f*x - d
	}
	return y
}

func createPlotXY(p []float64) plotter.XYs {
	pts := make(plotter.XYs, len(p))
	for i, v := range p {
		pts[i].X = float64(i) / 10
		pts[i].Y = v
	}
	return pts
}

func getAptFunc(correct []float64) func(p []float64) float64 {
	return func(p []float64) float64 {
		r := 0.0
		for i := 0; i < len(p); i++ {
			r += math.Abs(correct[i] - p[i])
		}
		return r
	}
}
