package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dustinkirkland/golang-petname"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

func main() {
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, "postgres://invoicing:invoicing@localhost:5432/test_db?")

	workers := river.NewWorkers()

	river.AddWorker(workers, &PdfGeneratorWorker{})
	river.AddWorker(workers, &SendMailWorker{})

	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 3},
		},
		Workers: workers,
	})
	if err != nil {
		// handle error
	}

	// Run the client inline. All executed jobs will inherit from ctx:
	if err = riverClient.Start(ctx); err != nil {
		// handle error
	}
	for i := 0; i < 10; i++ {
		_, err = riverClient.Insert(ctx, &PdfGeneratorWorkerArgs{
			Name:     fmt.Sprintf("invoice-%d", i),
			Template: "invoice",
			Username: petname.Generate(2, "-"),
			SendMail: true,
		}, nil)
		if err != nil {
			// handle error
		}
	}
	infiniteLoop()
}

func infiniteLoop() {
	for {
		// Introduce a delay to avoid excessive CPU usage
		time.Sleep(time.Second)
	}
}
