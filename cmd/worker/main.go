package main

import (
	"log"

	"github.com/sadia-54/qstack-backend/internal/queue"
	"github.com/sadia-54/qstack-backend/internal/workers"
)

func main() {

	if err := queue.Connect(); err != nil {
		log.Fatal(err)
	}
	defer queue.Close()

	queue.StartConsumer(workers.EmailWorker)
}