package report

import "time"

func aggregation(records []Record) Record {
	var s time.Duration
	for _, r := range records {
		s += r.Elapse
	}
	n := len(records)
	return Record{
		Start:  records[n/2].Start,
		Elapse: s / time.Duration(n),
	}
}

func sampling(records []Record, size int) []Record {
	n := len(records)
	if n <= size {
		return records
	}

	freq := n / size
	samples := make([]Record, size)
	for i := 0; i < size; i++ {
		samples[i] = aggregation(records[i*freq : (i+1)*freq])
	}
	return samples
}
