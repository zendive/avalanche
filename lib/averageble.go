package lib

import (
	"fmt"
	"math"
	"sync"
)

type Averageble struct {
	population bool
	mutex      sync.Mutex
	count      uint64
	variance   float64
	sumDelta2  float64
	Mean       float64
	Min        float64
	Max        float64
	StdDev     float64
}

type AveragebleOptions struct {
	Population bool
}

func NewAverageble(options AveragebleOptions) *Averageble {
	a := &Averageble{}
	a.population = options.Population
	return a
}

func (avg *Averageble) Add(value float64) {
	avg.mutex.Lock()
	avg.count++

	if avg.count == 1 {
		avg.Min = value
		avg.Max = value
	} else {
		avg.Min = math.Min(avg.Min, value)
		avg.Max = math.Max(avg.Max, value)
	}

	var delta float64 = value - avg.Mean
	avg.Mean += delta / float64(avg.count)
	avg.sumDelta2 += delta * (value - avg.Mean)

	if avg.population {
		avg.variance = avg.sumDelta2 / float64(avg.count)
	} else {
		avg.variance = avg.sumDelta2 / (float64(avg.count) - 1)
	}

	avg.StdDev = math.Sqrt(avg.variance)

	avg.mutex.Unlock()
}

func (avg *Averageble) String() string {
	return fmt.Sprintf("avg=%.3f,\tmin=%.3f,\tmax=%.3f,\tstddev=%.3f",
		avg.Mean, avg.Min, avg.Max, avg.StdDev)
}
