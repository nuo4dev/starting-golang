package main

import (
	"context"
	"time"
	"fmt"
)

func main() {
	ch := make(chan int, 6)
	defer close(ch)
	ctx, cancel := context.WithCancel(context.Background())
	// start producer
	go func (ch chan <- int) {
		producer := Producer{
			ctx: ctx,
			interval: time.Second,
			next: 0,
		}
		producer.produce(ch)
	}(ch)

	// start consumer
	go func (ch <- chan int) {
		consumer := Consumer{
			ctx: ctx, 
			interval: time.Second,
		}
		consumer.consume(ch)
	}(ch)

	time.Sleep(10 * time.Second)
	cancel()
	time.Sleep(2 * time.Second)
}

type Producer struct {
	ctx context.Context
	interval time.Duration
	next int
}

type Consumer struct {
	ctx context.Context
	interval time.Duration
}

func (p *Producer) produce(ch chan <- int) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <- p.ctx.Done():
			fmt.Println("[producer] process interrupted from main")
			return
		default:
			fmt.Println(" <- PRODUCE value: ", p.next)
			ch <- p.next
			p.next++
		}
	}
}

func (c *Consumer) consume(ch <- chan int) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <- c.ctx.Done():
			fmt.Println("[consumer] process interrupted from main")
			return
		default:
			fmt.Println(" -> CONSUME value: ", <-ch)
		}
	}
}