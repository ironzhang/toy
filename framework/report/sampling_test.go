package report

import (
	"reflect"
	"testing"
	"time"
)

func TestAggregation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		records []Record
		want    Record
	}{
		{
			records: []Record{
				{Start: now.Add(0 * time.Second), Elapse: 0 * time.Second},
			},
			want: Record{Start: now.Add(0 * time.Second), Elapse: 0 * time.Second},
		},
		{
			records: []Record{
				{Start: now.Add(0 * time.Second), Elapse: 0 * time.Second},
				{Start: now.Add(1 * time.Second), Elapse: 1 * time.Second},
			},
			want: Record{Start: now.Add(1 * time.Second), Elapse: 1 * time.Second / 2},
		},
		{
			records: []Record{
				{Start: now.Add(0 * time.Second), Elapse: 0 * time.Second},
				{Start: now.Add(1 * time.Second), Elapse: 1 * time.Second},
				{Start: now.Add(2 * time.Second), Elapse: 2 * time.Second},
			},
			want: Record{Start: now.Add(1 * time.Second), Elapse: 3 * time.Second / 3},
		},
		{
			records: []Record{
				{Start: now.Add(0 * time.Second), Elapse: 0 * time.Second},
				{Start: now.Add(1 * time.Second), Elapse: 1 * time.Second},
				{Start: now.Add(2 * time.Second), Elapse: 2 * time.Second},
				{Start: now.Add(4 * time.Second), Elapse: 4 * time.Second},
			},
			want: Record{Start: now.Add(2 * time.Second), Elapse: 7 * time.Second / 4},
		},
		{
			records: []Record{
				{Start: now.Add(1 * time.Second), Elapse: 1 * time.Second},
				{Start: now.Add(2 * time.Second), Elapse: 2 * time.Second},
				{Start: now.Add(3 * time.Second), Elapse: 3 * time.Second},
				{Start: now.Add(4 * time.Second), Elapse: 4 * time.Second},
				{Start: now.Add(5 * time.Second), Elapse: 5 * time.Second},
			},
			want: Record{Start: now.Add(3 * time.Second), Elapse: 15 * time.Second / 5},
		},
	}
	for i, tt := range tests {
		if got, want := aggregation(tt.records), tt.want; got != want {
			t.Errorf("case%d: %v != %v", i, got, want)
		}
	}
}

func TestSampling(t *testing.T) {
	now := time.Now()
	records := []Record{
		{Start: now.Add(1 * time.Second), Elapse: 1 * time.Second},
		{Start: now.Add(2 * time.Second), Elapse: 2 * time.Second},
		{Start: now.Add(3 * time.Second), Elapse: 3 * time.Second},
		{Start: now.Add(4 * time.Second), Elapse: 4 * time.Second},
		{Start: now.Add(5 * time.Second), Elapse: 5 * time.Second},
		{Start: now.Add(6 * time.Second), Elapse: 6 * time.Second},
		{Start: now.Add(7 * time.Second), Elapse: 7 * time.Second},
		{Start: now.Add(8 * time.Second), Elapse: 8 * time.Second},
		{Start: now.Add(9 * time.Second), Elapse: 9 * time.Second},
		{Start: now.Add(10 * time.Second), Elapse: 10 * time.Second},
	}
	tests := []struct {
		size int
		want []Record
	}{
		{
			size: 10,
			want: records,
		},
		{
			size: 5,
			want: []Record{
				{Start: now.Add(2 * time.Second), Elapse: (1 + 2) * time.Second / 2},
				{Start: now.Add(4 * time.Second), Elapse: (3 + 4) * time.Second / 2},
				{Start: now.Add(6 * time.Second), Elapse: (5 + 6) * time.Second / 2},
				{Start: now.Add(8 * time.Second), Elapse: (7 + 8) * time.Second / 2},
				{Start: now.Add(10 * time.Second), Elapse: (9 + 10) * time.Second / 2},
			},
		},
		{
			size: 3,
			want: []Record{
				{Start: now.Add(2 * time.Second), Elapse: (1 + 2 + 3) * time.Second / 3},
				{Start: now.Add(5 * time.Second), Elapse: (4 + 5 + 6) * time.Second / 3},
				{Start: now.Add(8 * time.Second), Elapse: (7 + 8 + 9) * time.Second / 3},
			},
		},
	}
	for i, tt := range tests {
		if got, want := sampling(records, tt.size), tt.want; !reflect.DeepEqual(got, want) {
			t.Errorf("case%d: %v != %v", i, got, want)
		}
	}
}
