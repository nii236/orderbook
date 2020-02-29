package agent

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"orderbook"
	"time"

	"github.com/gofrs/uuid"
)

type LimitAgent struct {
	Orderbook *orderbook.OrderBook
	TickRate  time.Duration
}

func NewLimitAgent(rate time.Duration, ob *orderbook.OrderBook) *LimitAgent {
	a := &LimitAgent{
		ob,
		rate,
	}
	return a
}

func (a *LimitAgent) Start(ctx context.Context, oType orderbook.OrderType) error {
	t := time.NewTicker(a.TickRate)

	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		case <-t.C:
			if oType == orderbook.BUY {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				amount := uint64(rand.Intn(10)) + 1
				buyOrder := &orderbook.Order{
					Amount: amount,
					Price:  100,
					ID:     uuid.Must(uuid.NewV4()).String(),
					Side:   orderbook.BUY,
				}
				// fmt.Println("Agent buying", amount)
				lowestAsk, err := a.Orderbook.LowestAsk()
				if errors.Is(err, orderbook.ErrNoOrder) {
					_, err = a.Orderbook.Process(buyOrder)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// fmt.Printf("successful trades: %+v\n", trades)
					continue
				} else if err != nil {
					fmt.Println(err)
					continue
				}
				buyOrder.Price = lowestAsk.Price + 5 - uint64(rand.Intn(10))
				// fmt.Println("lowestAsk.Price", lowestAsk.Price)

				_, err = a.Orderbook.Process(buyOrder)
				if err != nil {
					fmt.Println(err)
					continue
				}
				// fmt.Printf("successful trades: %+v\n", trades)
			}
			if oType == orderbook.SELL {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				amount := uint64(rand.Intn(10)) + 1
				sellOrder := &orderbook.Order{
					Amount: amount,
					Price:  100,
					ID:     uuid.Must(uuid.NewV4()).String(),
					Side:   orderbook.SELL,
				}
				// fmt.Println("Agent selling", amount)

				highestBid, err := a.Orderbook.HighestBid()
				if errors.Is(err, orderbook.ErrNoOrder) {
					_, err = a.Orderbook.Process(sellOrder)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// fmt.Printf("successful trades: %+v\n", trades)
					continue
				} else if err != nil {
					fmt.Println(err)
					continue
				}

				sellOrder.Price = highestBid.Price - 5 + uint64(rand.Intn(10))
				// fmt.Println("highestBid.Price", highestBid.Price)
				_, err = a.Orderbook.Process(sellOrder)
				if err != nil {
					fmt.Println(err)
					continue
				}
				// fmt.Printf("successful trades: %+v\n", trades)

			}
		}
	}
}
