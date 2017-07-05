package jobs

type BenchmarkJob struct {
	Verbose    int
	Ask        bool
	Record     bool
	RecordFile string
	RobotNum   int
	RobotPath  string
}

func (p *BenchmarkJob) Execute() error {
	return nil
}
