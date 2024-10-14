package worker

import (
	"context"
	"demo1/pdf"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"

	"github.com/riverqueue/river"
)

type InvoiceGeneratorWorker struct {
}

func (s InvoiceGeneratorWorker) NextRetry(job *river.Job[InvoiceGeneratorWorkerArgs]) time.Time {
	return time.Now().Add(10 * time.Second)
}

func (s InvoiceGeneratorWorker) Timeout(job *river.Job[InvoiceGeneratorWorkerArgs]) time.Duration {
	return 30 * time.Second
}

func (s InvoiceGeneratorWorker) Work(ctx context.Context, job *river.Job[InvoiceGeneratorWorkerArgs]) error {
	fmt.Printf("Generating invoice %s  ...\n", job.Args.Name)
	pdf.Complex1Report(job.Args.Name)

	if job.Args.SendMail {
		client := river.ClientFromContext[pgx.Tx](ctx)
		_, err := client.Insert(ctx, SendMailWorkerArgs{Username: job.Args.Username}, &river.InsertOpts{
			Queue: "mailing",
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type InvoiceGeneratorWorkerArgs struct {
	Name     string
	Username string
	SendMail bool
}

func (s InvoiceGeneratorWorkerArgs) Kind() string {
	return "InvoiceGeneratorWorkerArgs"
}
