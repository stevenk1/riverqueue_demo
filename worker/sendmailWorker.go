package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/riverqueue/river"
)

type SendMailWorker struct {
}

func (s SendMailWorker) NextRetry(job *river.Job[SendMailWorkerArgs]) time.Time {
	return time.Now().Add(10 * time.Second)
}

func (s SendMailWorker) Timeout(job *river.Job[SendMailWorkerArgs]) time.Duration {
	return 30 * time.Second
}

func (s SendMailWorker) Work(ctx context.Context, job *river.Job[SendMailWorkerArgs]) error {
	fmt.Printf("Sending mail to %s  ...\n", job.Args.Username)

	return nil
}

type SendMailWorkerArgs struct {
	Username string
}

func (s SendMailWorkerArgs) Kind() string {
	return "SendMailWorkerArgs"
}
