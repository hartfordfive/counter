package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	cnt "github.com/hartfordfive/counter"
)

func main() {

	fmt.Printf("-------------------------\nCreating counter\n------------------------\n")
	counter := cnt.NewCounter()

	fmt.Println("Incrementing counter 100000 times..")
	for i := 0; i < 100000; i++ {
		go func(c *cnt.Counter) {
			c.Incr(1)
		}(counter)
	}

	// Sleep a bit for all values to be counted
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("Count: %d\n", counter.Value())

	// Reset the counter
	fmt.Println("Resetting counter...")
	counter.Reset()

	fmt.Println("Incrementing counter 5000 times..")
	for i := 0; i < 5000; i++ {
		go func(c *cnt.Counter) {
			c.Incr(1)
		}(counter)
	}

	// Sleep a bit for all values to be counted
	time.Sleep(2000 * time.Millisecond)
	fmt.Printf("New count: %d\n\n", counter.Value())

	// Cancel the counter so that the accumulator chanel is closed and the count is reset
	counter.Cancel()

	// -------------------------------------------------------

	fmt.Printf("-------------------------\nCreating rate counter\n------------------------\n")

	rcounter := cnt.NewRateCounter(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Every second, get the current rate
	duration := time.Second * 1
	ticker := time.NewTicker(time.Second * 1)

	// Run a goroutine to print the current rate or exit when context notifies complete
	go func(rc *cnt.RateCounter, ticker *time.Ticker, duration time.Duration, ctx context.Context) {
		for {
			select {
			case <-ticker.C:
				fmt.Printf("Rate per %s: %d\n", duration.String(), rc.CurrRate())
			case <-ctx.Done():
				return

			}
		}
	}(rcounter, ticker, duration, ctx)

	// Create goroutines to increment the count by a random amount from 1 to 10
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stoping ticker and returning...")
			ticker.Stop()
			return
		default:
			go func(rc *cnt.RateCounter) {
				rc.Incr(int64(rand.Intn(2)))
			}(rcounter)
		}
	}

}
