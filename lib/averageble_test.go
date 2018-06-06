package lib

import "testing"

func TestAverageble(t *testing.T) {
	var avg = NewAverageble(AveragebleOptions{Population: true})
	var feed = [...]float64{-5, 1, 8, 7, 2}
	const MIN = -5.0
	const MAX = 8.0
	const MEAN = 2.6
	const STDDEV = 4.673328578219168

	for _, n := range feed {
		avg.Add(n)
	}

	if MIN != avg.Min {
		t.Errorf("Min was incorrect, got: %f, want: %f", avg.Min, MIN)
	}
	if MAX != avg.Max {
		t.Errorf("Max was incorrent, got: %f, want: %f", avg.Max, MAX)
	}
	if MEAN != avg.Mean {
		t.Errorf("Mean was incorrect, got: %f, want: %f", avg.Mean, MEAN)
	}
	if STDDEV != avg.StdDev {
		t.Errorf("Standart deviation was incorrect, got: %f, want: %f.", avg.StdDev, STDDEV)
	}
}
