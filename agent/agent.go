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
				amount := uint64(rand.Intn(50))
				buyOrder := &orderbook.Order{
					Amount: amount,
					Price:  100,
					ID:     uuid.Must(uuid.NewV4()).String(),
					Side:   orderbook.BUY,
				}
				// fmt.Println("Agent buying", amount)
				lowestAsk, err := a.Orderbook.LowestAsk()
				if errors.Is(err, orderbook.ErrNoOrder) {
					a.Orderbook.Process(buyOrder)
					// fmt.Printf("successful trades: %+v\n", trades)
					continue
				} else if err != nil {
					fmt.Println(err)
					continue
				}

				buyOrder.Price = lowestAsk.Price - uint64(rand.Intn(10))

				a.Orderbook.Process(buyOrder)
				// fmt.Printf("successful trades: %+v\n", trades)
			}
			if oType == orderbook.SELL {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				amount := uint64(rand.Intn(50))
				sellOrder := &orderbook.Order{
					Amount: amount,
					Price:  100,
					ID:     uuid.Must(uuid.NewV4()).String(),
					Side:   orderbook.SELL,
				}
				// fmt.Println("Agent selling", amount)

				highestBid, err := a.Orderbook.HighestBid()
				if errors.Is(err, orderbook.ErrNoOrder) {
					a.Orderbook.Process(sellOrder)
					// fmt.Printf("successful trades: %+v\n", trades)
					continue
				} else if err != nil {
					fmt.Println(err)
					continue
				}

				sellOrder.Price = highestBid.Price + uint64(rand.Intn(10))

				a.Orderbook.Process(sellOrder)
				// fmt.Printf("successful trades: %+v\n", trades)

			}
		}
	}
}
