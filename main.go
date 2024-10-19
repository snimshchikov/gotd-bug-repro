package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	// Create a new logger
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"debug.log"}
	log, err := cfg.Build()
	if err != nil {
		return err
	}
	defer func() { _ = log.Sync() }()

	// Start bot
	err = telegram.BotFromEnvironment(
		ctx,
		telegram.Options{Logger: log},
		nil,
		stressLoadBot,
	)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func stressLoadBot(ctx context.Context, client *telegram.Client) error {
	countChan := make(chan int)
	count := 0
	for i := 1; i <= 15; i++ {
		go loopGetDifference(ctx, client, i, countChan)
	}
	for {
		select {
		case <-countChan:
			count += 1
			if count%100 == 0 {
				fmt.Println("Workers made", count, "requests")
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func loopGetDifference(ctx context.Context, client *telegram.Client, workerID int, count chan int) {
	differenceRequest := &tg.UpdatesGetDifferenceRequest{Pts: 1, Qts: 1, Date: 1}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, err := client.API().UpdatesGetDifference(ctx, differenceRequest)
			if err != nil {
				fmt.Println("Error in worker", workerID, "while getting difference: ", err)
				continue
			}
			count <- 1
		}
	}
}
