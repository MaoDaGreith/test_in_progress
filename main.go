package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"arkis_test/database"
	"arkis_test/processor"
	"arkis_test/queue"

	log "github.com/sirupsen/logrus"
)

func main() {
	ch := make(chan os.Signal)
	defer close(ch)
	signal.Notify(ch, syscall.SIGINT)

	ctxA := context.Background()
	ctxB := context.Background()

	inputQueueA, err := queue.New(os.Getenv("RABBITMQ_URL"), "input-A")
	if err != nil {
		log.WithError(err).Panic("Cannot create input queue A")
	}

	outputQueueA, err := queue.New(os.Getenv("RABBITMQ_URL"), "output-A")
	if err != nil {
		log.WithError(err).Panic("Cannot create output queue A")
	}

	inputQueueB, err := queue.New(os.Getenv("RABBITMQ_URL"), "input-B")
	if err != nil {
		log.WithError(err).Panic("Cannot create input queue B")
	}

	outputQueueB, err := queue.New(os.Getenv("RABBITMQ_URL"), "output-B")
	if err != nil {
		log.WithError(err).Panic("Cannot create output queue B")
	}

	log.Info("Application is ready to run")

	go func() {
		if err := processor.New(inputQueueA, outputQueueA, database.D{}).Run(ctxA); err != nil {
			log.WithError(err).Fatal("Error running processor A")
		}
	}()

	go func() {
		if err := processor.New(inputQueueB, outputQueueB, database.D{}).Run(ctxB); err != nil {
			log.WithError(err).Fatal("Error running processor B")
		}
	}()

	select {
	case <-ch:
		log.Info("Received interrupt signal, shutting down...")
		os.Exit(0) // Graceful exit
	}
}
