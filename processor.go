package orderbook

import (
	"fmt"
)

// Process an order and return the trades generated before adding the remaining amount to the market
func (ob *OrderBook) Process(order *Order) ([]Trade, error) {
	if order.Side == BUY {
		result, err := ob.processLimitBuy(order)
		if err != nil {
			return nil, err
		}
		if len(result) > 0 {
			fmt.Println(ob.Status(), len(result), "trades matched")
		}
		return result, nil
	}

	result, err := ob.processLimitSell(order)
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		fmt.Println(ob.Status(), len(result), "trades matched")
	}
	return result, nil
}

func (ob *OrderBook) fulfilBuy(sellOrders []*Order, newOrder *Order, trades []Trade) (*Order, []Trade, error) {
	for i, matchingOrder := range sellOrders {
		// Process full
		if matchingOrder.Amount >= newOrder.Amount {
			trades = append(trades, Trade{newOrder.ID, matchingOrder.ID, newOrder.Amount, matchingOrder.Price})
			matchingOrder.Amount -= newOrder.Amount
			if matchingOrder.Amount == 0 {
				ob.RemoveSell(i)
			}
			return newOrder, trades, nil
		}
		// Process partial and re-run
		if matchingOrder.Amount < newOrder.Amount {
			trades = append(trades, Trade{newOrder.ID, matchingOrder.ID, matchingOrder.Amount, matchingOrder.Price})
			newOrder.Amount -= matchingOrder.Amount
			ob.RemoveSell(i)
			return ob.fulfilBuy(ob.SellOrders, newOrder, trades)
		}
	}
	return newOrder, []Trade{}, nil
}

func (ob *OrderBook) fulfilSell(buyOrders []*Order, newOrder *Order, trades []Trade) (*Order, []Trade, error) {
	for i, matchingOrder := range buyOrders {
		// fill the entire order
		if matchingOrder.Amount >= newOrder.Amount {
			trades = append(trades, Trade{newOrder.ID, matchingOrder.ID, newOrder.Amount, matchingOrder.Price})
			matchingOrder.Amount -= newOrder.Amount
			if matchingOrder.Amount == 0 {
				ob.RemoveBuy(i)
			}
			return newOrder, trades, nil
		}
		// fill a partial order and continue
		if matchingOrder.Amount < newOrder.Amount {
			trades = append(trades, Trade{newOrder.ID, matchingOrder.ID, matchingOrder.Amount, matchingOrder.Price})
			newOrder.Amount -= matchingOrder.Amount
			ob.RemoveBuy(i)
			fmt.Println("recurse")
			return ob.fulfilSell(ob.BuyOrders, newOrder, trades)
		}
	}
	return newOrder, []Trade{}, nil
}

// Process a limit buy order
func (ob *OrderBook) processLimitBuy(order *Order) ([]Trade, error) {
	var err error
	n := len(ob.SellOrders)
	if n == 0 {
		err = ob.Append(order)
		if err != nil {
			return nil, err
		}
		return []Trade{}, nil
	}

	if ob.SellOrders[0].Price > order.Price {
		err = ob.Append(order)
		if err != nil {
			return nil, err
		}
		return []Trade{}, nil
	}
	trades := []Trade{}
	// check if we have at least one matching order
	// traverse all orders that match
	finalOrder := &Order{}
	finalOrder, trades, err = ob.fulfilBuy(ob.SellOrders, order, trades)
	if err != nil {
		return nil, err
	}
	// finally add the remaining order to the list
	err = ob.Append(finalOrder)
	if err != nil {
		return nil, err
	}
	return trades, nil
}

// Process a limit sell order
func (ob *OrderBook) processLimitSell(order *Order) ([]Trade, error) {
	var err error
	n := len(ob.BuyOrders)
	if n == 0 {
		ob.AddSell(order)
		return []Trade{}, nil
	}

	if ob.BuyOrders[0].Price < order.Price {
		ob.AddSell(order)
		return []Trade{}, nil
	}
	trades := []Trade{}
	// check if we have at least one matching order
	// traverse all orders that match
	finalOrder := &Order{}
	finalOrder, trades, err = ob.fulfilSell(ob.BuyOrders, order, trades)
	if err != nil {
		return nil, err
	}
	fmt.Println(finalOrder)
	// finally add the remaining order to the list
	err = ob.Append(finalOrder)
	if err != nil {
		return nil, err
	}
	return trades, nil
}
