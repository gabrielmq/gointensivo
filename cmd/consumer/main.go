package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gabrielmq/gointensivo/internal/order/infra/database"
	"github.com/gabrielmq/gointensivo/internal/order/usecase"
	"github.com/gabrielmq/gointensivo/pkg/rabbitmq"
	_ "github.com/mattn/go-sqlite3"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	db, err := sql.Open("sqlite3", "./orders.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	gateway := database.NewOrderRepository(db)
	uc := usecase.NewCalculateFinalPriceUseCase(gateway)

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	out := make(chan amqp.Delivery, 10)

	go rabbitmq.Consume(ch, out)

	workers := 10
	for i := 1; i <= workers; i++ {
		go worker(i, out, uc)
	}

	http.HandleFunc("/total", func(w http.ResponseWriter, r *http.Request) {
		uc := usecase.NewTotalUseCase(gateway)
		total, err := uc.Execute()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		json.NewEncoder(w).Encode(total)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func worker(workerID int, deliveryMessage <-chan amqp.Delivery, uc *usecase.CalculateFinalPriceUseCase) {
	for msg := range deliveryMessage {
		var input usecase.OrderInputDTO
		if err := json.Unmarshal(msg.Body, &input); err != nil {
			panic(err)
		}

		output, err := uc.Execute(input)
		if err != nil {
			panic(err)
		}

		msg.Ack(false)
		fmt.Printf("Worker %d has processed order %s\n", workerID, output.ID)
		time.Sleep(100 * time.Millisecond)
	}
}
