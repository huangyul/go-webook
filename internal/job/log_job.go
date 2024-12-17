package job

import "fmt"

type LogJob struct{}

func NewLogJob() *LogJob {
	return &LogJob{}
}

func (l *LogJob) Name() string {
	return "LogJob"
}

func (l *LogJob) Run() error {
	fmt.Println("LogJob Run")
	return nil
}
