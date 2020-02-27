package orderbook

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"sync"
)

// OrderBook type
type OrderBook struct {
	Conn       KV
	Pair       Pair
	BuyOrders  []*Order
	SellOrders []*Order
	*sync.RWMutex
}

func New(conn KV, pair Pair) (*OrderBook, error) {
	ob := &OrderBook{
		conn,
		pair,
		[]*Order{},
		[]*Order{},
		&sync.RWMutex{},
	}
	return ob, nil

}

type OrderType string

const SELL OrderType = "SELL"
const BUY OrderType = "BUY"

func (ob *OrderBook) Status() string {
	bidStr := "---"
	bid, err := ob.HighestBid()
	if err != nil {
		bidStr = "---"
	} else {
		bidStr = fmt.Sprintf("%03d", bid.Price)
	}
	askStr := "---"
	ask, err := ob.LowestAsk()
	if err != nil {
		askStr = "---"
	} else {
		askStr = fmt.Sprintf("%03d", ask.Price)
	}
	return fmt.Sprintf("BID: %s USD ASK: %s USD TOTALBIDS:%03d TOTALASKS:%03d", bidStr, askStr, len(ob.BuyOrders), len(ob.SellOrders))
}
func (ob *OrderBook) Save() error {
	var err error
	ob.Lock()
	defer ob.Unlock()
	buyBuf := &bytes.Buffer{}
	buyEnc := gob.NewEncoder(buyBuf)
	err = buyEnc.Encode(ob.BuyOrders)
	if err != nil {
		return err
	}
	sellBuf := &bytes.Buffer{}
	sellEnc := gob.NewEncoder(sellBuf)
	err = sellEnc.Encode(ob.BuyOrders)
	if err != nil {
		return err
	}

	err = ob.Conn.Set([]byte(ob.Pair.BuyKey()), buyBuf.Bytes())
	if err != nil {
		return err
	}
	err = ob.Conn.Set([]byte(ob.Pair.SellKey()), sellBuf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (ob *OrderBook) Load() error {
	ob.Lock()
	defer ob.Unlock()
	var err error

	buyOrderBytes, err := ob.Conn.Get([]byte(ob.Pair.BuyKey()))
	buyDec := gob.NewDecoder(bytes.NewReader(buyOrderBytes))
	err = buyDec.Decode(&ob.BuyOrders)
	if err != nil {
		return err
	}

	sellOrderBytes, err := ob.Conn.Get([]byte(ob.Pair.SellKey()))
	sellDec := gob.NewDecoder(bytes.NewReader(sellOrderBytes))
	err = sellDec.Decode(&ob.SellOrders)
	if err != nil {
		return err
	}

	buyBuf := &bytes.Buffer{}
	err = buyDec.Decode(buyBuf)
	if err != nil {
		return err
	}
	sellBuf := &bytes.Buffer{}
	err = sellDec.Decode(sellBuf)
	if err != nil {
		return err
	}

	return nil
}

var ErrNoOrder = errors.New("no orders available")

func (ob *OrderBook) HighestBid() (*Order, error) {
	if len(ob.BuyOrders) == 0 {
		return nil, ErrNoOrder
	}
	return ob.BuyOrders[0], nil
}

func (ob *OrderBook) LowestAsk() (*Order, error) {
	if len(ob.SellOrders) == 0 {
		return nil, ErrNoOrder
	}
	return ob.SellOrders[0], nil
}

func (ob *OrderBook) Append(order *Order) error {
	if order.Side == BUY {
		return ob.AddBuy(order)
	}
	if order.Side == SELL {
		return ob.AddSell(order)
	}
	return errors.New("unsupported order type")
}
func (ob *OrderBook) AddBuy(order *Order) error {
	ob.Lock()
	defer ob.Unlock()
	n := len(ob.BuyOrders)
	appendToBook := false
	var i int
	for i := n - 1; i >= 0; i-- {
		buyOrder := ob.BuyOrders[i]
		if buyOrder.Price < order.Price {
			break
		}
	}
	if n == 0 || i == n-1 {
		appendToBook = true
	}
	if appendToBook {
		ob.BuyOrders = append(ob.BuyOrders, order)
		return nil
	}

	copy(ob.BuyOrders[i+1:], ob.BuyOrders[i:])
	ob.BuyOrders[i] = order
	return nil
}

// Add a sell order to the order ob
func (ob *OrderBook) AddSell(order *Order) error {
	n := len(ob.SellOrders)
	var i int
	for i := n - 1; i >= 0; i-- {
		sellOrder := ob.SellOrders[i]
		if sellOrder.Price > order.Price {
			break
		}
	}
	if n == 0 || i == n-1 {
		ob.SellOrders = append(ob.SellOrders, order)
	} else {
		copy(ob.SellOrders[i+1:], ob.SellOrders[i:])
		ob.SellOrders[i] = order
	}
	return nil
}

func (ob *OrderBook) RemoveBuy(index int) {
	ob.BuyOrders = append(ob.BuyOrders[:index], ob.BuyOrders[index+1:]...)
}

func (ob *OrderBook) RemoveSell(index int) {
	ob.SellOrders = append(ob.SellOrders[:index], ob.SellOrders[index+1:]...)
}
