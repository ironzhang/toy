package report

import (
	"reflect"
	"testing"
	"time"
)

func TestStatsHistogram(t *testing.T) {
	tests := []struct {
		lats    []time.Duration
		buckets []time.Duration
		counts  []int
	}{
		// case0
		{
			lats: []time.Duration{1 * time.Second},
			buckets: []time.Duration{
				1 * time.Second, //1
				1 * time.Second, //2
				1 * time.Second, //3
				1 * time.Second, //4
				1 * time.Second, //5
				1 * time.Second, //6
				1 * time.Second, //7
				1 * time.Second, //8
				1 * time.Second, //9
				1 * time.Second, //10
				1 * time.Second, //11
			},
			counts: []int{
				1, //1
				0, //2
				0, //3
				0, //4
				0, //5
				0, //6
				0, //7
				0, //8
				0, //9
				0, //10
				0, //11
			},
		},

		// case1
		{
			lats: []time.Duration{1 * time.Second, 11 * time.Second},
			buckets: []time.Duration{
				1 * time.Second,  //1
				2 * time.Second,  //2
				3 * time.Second,  //3
				4 * time.Second,  //4
				5 * time.Second,  //5
				6 * time.Second,  //6
				7 * time.Second,  //7
				8 * time.Second,  //8
				9 * time.Second,  //9
				10 * time.Second, //10
				11 * time.Second, //11
			},
			counts: []int{
				1, //1
				0, //2
				0, //3
				0, //4
				0, //5
				0, //6
				0, //7
				0, //8
				0, //9
				0, //10
				1, //11
			},
		},

		// case2
		{
			lats: []time.Duration{
				// 1s
				1 * time.Second,
				// 2s
				2 * time.Second,
				2 * time.Second,
				// 3s
				3 * time.Second,
				3 * time.Second,
				3 * time.Second,
				// 4s
				3*time.Second + 1*time.Millisecond,
				3*time.Second + 2*time.Millisecond,
				// 5s
				4*time.Second + 1*time.Millisecond,
				5 * time.Second,
				// 6s
				6 * time.Second,
				6 * time.Second,
				// 9s
				8*time.Second + 1*time.Millisecond,
				// 11s
				10*time.Second + time.Millisecond,
				11 * time.Second,
			},
			buckets: []time.Duration{
				1 * time.Second,  //1
				2 * time.Second,  //2
				3 * time.Second,  //3
				4 * time.Second,  //4
				5 * time.Second,  //5
				6 * time.Second,  //6
				7 * time.Second,  //7
				8 * time.Second,  //8
				9 * time.Second,  //9
				10 * time.Second, //10
				11 * time.Second, //11
			},
			counts: []int{
				1, //1
				2, //2
				3, //3
				2, //4
				2, //5
				2, //6
				0, //7
				0, //8
				1, //9
				0, //10
				2, //11
			},
		},
	}
	for i, tt := range tests {
		s := Stats{
			Fastest: tt.lats[0],
			Slowest: tt.lats[len(tt.lats)-1],
			Lats:    tt.lats,
		}
		buckets, counts := s.Histogram()
		if got, want := buckets, tt.buckets; !reflect.DeepEqual(got, want) {
			t.Errorf("case%d: buckets: %v != %v", i, got, want)
		}
		if got, want := counts, tt.counts; !reflect.DeepEqual(got, want) {
			t.Errorf("case%d: buckets: %v != %v", i, got, want)
		}
	}
}

func makeTestDuration(n int) []time.Duration {
	ds := make([]time.Duration, 0, n)
	for i := 1; i <= n; i++ {
		ds = append(ds, time.Duration(i)*time.Second)
	}
	return ds
}

func TestStatsLatencies(t *testing.T) {
	tests := []struct {
		lats []time.Duration
		pcs  []float64
		data []time.Duration
	}{
		// case0
		{
			lats: nil,
			pcs:  nil,
			data: nil,
		},

		// case1
		{
			lats: []time.Duration{1 * time.Second},
			pcs:  []float64{99.9},
			data: []time.Duration{1 * time.Second},
		},

		// case2
		{
			lats: []time.Duration{1 * time.Second, 10 * time.Second},
			pcs:  []float64{25, 99.9},
			data: []time.Duration{1 * time.Second, 10 * time.Second},
		},

		// case3
		{
			lats: makeTestDuration(10),
			pcs:  []float64{10, 25, 50, 75, 99.9},
			data: []time.Duration{
				2 * time.Second,
				3 * time.Second,
				6 * time.Second,
				8 * time.Second,
				10 * time.Second,
			},
		},

		// case4
		{
			lats: makeTestDuration(100),
			pcs:  []float64{10, 25, 50, 75, 90, 95, 99.9},
			data: []time.Duration{
				11 * time.Second,
				26 * time.Second,
				51 * time.Second,
				76 * time.Second,
				91 * time.Second,
				96 * time.Second,
				100 * time.Second,
			},
		},

		// case5
		{
			lats: makeTestDuration(1000),
			pcs:  []float64{10, 25, 50, 75, 90, 95, 99, 99.9},
			data: []time.Duration{
				101 * time.Second,
				251 * time.Second,
				501 * time.Second,
				751 * time.Second,
				901 * time.Second,
				951 * time.Second,
				991 * time.Second,
				1000 * time.Second,
			},
		},

		// case6
		{
			lats: makeTestDuration(10000),
			pcs:  []float64{10, 25, 50, 75, 90, 95, 99, 99.9},
			data: []time.Duration{
				1001 * time.Second,
				2501 * time.Second,
				5001 * time.Second,
				7501 * time.Second,
				9001 * time.Second,
				9501 * time.Second,
				9901 * time.Second,
				9991 * time.Second,
			},
		},
	}
	for i, tt := range tests {
		s := Stats{
			Lats: tt.lats,
		}
		pcs, data := s.Latencies()
		if got, want := pcs, tt.pcs; !reflect.DeepEqual(got, want) {
			t.Errorf("case%d: pcs: %v != %v", i, got, want)
		}
		if got, want := data, tt.data; !reflect.DeepEqual(got, want) {
			t.Errorf("case%d: data: %v != %v", i, got, want)
		}
	}
}
