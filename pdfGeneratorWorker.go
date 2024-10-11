package main

import (
	"context"
	"demo1/pdf"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"

	"github.com/riverqueue/river"
)

type PdfGeneratorWorker struct {
}

func (s PdfGeneratorWorker) NextRetry(job *river.Job[PdfGeneratorWorkerArgs]) time.Time {
	return time.Now().Add(10 * time.Second)
}

func (s PdfGeneratorWorker) Timeout(job *river.Job[PdfGeneratorWorkerArgs]) time.Duration {
	return 30 * time.Second
}

func (s PdfGeneratorWorker) Work(ctx context.Context, job *river.Job[PdfGeneratorWorkerArgs]) error {
	if job.Args.Template == "invoice" {
		fmt.Printf("Generating invoice %s  ...\n", job.Args.Name)
		pdf.Complex1Report(job.Args.Name)
	}

	if job.Args.SendMail {
		client := river.ClientFromContext[pgx.Tx](ctx)
		_, err := client.Insert(ctx, SendMailWorkerArgs{Username: job.Args.Username}, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

type PdfGeneratorWorkerArgs struct {
	Name     string
	Username string
	Template string
	SendMail bool
}

func (s PdfGeneratorWorkerArgs) Kind() string {
	return "PdfGeneratorWorkerArgs"
}
