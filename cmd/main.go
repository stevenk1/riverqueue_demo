package main

import (
	"context"
	"demo1/worker"
	"fmt"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"log"
	"strconv"
)

func main() {
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, "postgres://invoicing:invoicing@localhost:5432/test_db?")

	workers := river.NewWorkers()

	river.AddWorker(workers, &worker.InvoiceGeneratorWorker{})
	river.AddWorker(workers, &worker.SendMailWorker{})

	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 3},
			"mailing":          {MaxWorkers: 3},
		},
		Workers: workers,
	})
	if err != nil {
		fmt.Println("there was an error:", err)
	}

	// Run the client inline. All executed jobs will inherit from ctx:
	if err = riverClient.Start(ctx); err != nil {
		fmt.Println("there was an error:", err)
	}

	app := fiber.New()
	app.Post("/:amount", func(c *fiber.Ctx) error {
		amount, err := strconv.Atoi(c.Params("amount"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
		}
		for i := 0; i < amount; i++ {
			_, err = riverClient.Insert(ctx, &worker.InvoiceGeneratorWorkerArgs{
				Name:     fmt.Sprintf("invoice-%d", i),
				Username: petname.Generate(2, "-"),
				SendMail: true,
			}, nil)
			if err != nil {
				fmt.Println("there was an error:", err)
			}
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{})
	})
	log.Fatal(app.Listen(":3535"))
}
