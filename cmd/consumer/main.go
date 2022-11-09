package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

	out := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, out)

	for msg := range out {
		var input usecase.OrderInputDTO
		if err := json.Unmarshal(msg.Body, &input); err != nil {
			panic(err)
		}

		output, err := uc.Execute(input)
		if err != nil {
			panic(err)
		}

		msg.Ack(false)
		fmt.Println("Consumed message", output)
		time.Sleep(500 * time.Millisecond)
	}
}
