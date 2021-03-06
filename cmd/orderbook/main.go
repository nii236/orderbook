package main

import (
	"context"
	"fmt"
	"math/rand"
	"orderbook"
	"orderbook/agent"
	"orderbook/ui"
	"time"

	"orderbook/stub"

	"github.com/oklog/run"
)

func main() {
	// conn, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	b := &stub.Stub{}
	ob, err := orderbook.New(b, "NIKE-AIRFORCEONES")
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	g := &run.Group{}
	for i := 0; i < 500; i++ {
		g.Add(func() error {
			aBuyer := agent.NewLimitAgent(time.Duration(rand.Intn(1000)*100+1)*time.Millisecond, ob)
			return aBuyer.Start(ctx, orderbook.BUY)
		}, func(error) {
			cancel()
		})

	}
	for i := 0; i < 500; i++ {
		g.Add(func() error {
			aSeller := agent.NewLimitAgent(time.Duration(rand.Intn(1000)*100+1)*time.Millisecond, ob)
			return aSeller.Start(ctx, orderbook.SELL)
		}, func(error) {
			cancel()
		})
	}
	g.Add(func() error {
		return ui.Run(ob)
	}, func(error) {
		cancel()
	})
	fmt.Println(g.Run())

}
