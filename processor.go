package orderbook

import "fmt"

// Process an order and return the trades generated before adding the remaining amount to the market
func (ob *OrderBook) Process(order *Order) []Trade {
	if order.Side == BUY {
		result := ob.processLimitBuy(order)
		fmt.Println(ob.Status(), len(result), "trades matched")
		return result
	}
	result := ob.processLimitSell(order)
	fmt.Println(ob.Status(), len(result), "trades matched")
	return result
}

// Process a limit buy order
func (ob *OrderBook) processLimitBuy(order *Order) []Trade {
	trades := make([]Trade, 0, 1)
	n := len(ob.SellOrders)
	// check if we have at least one matching order
	if n != 0 {
		if ob.SellOrders[n-1].Price <= order.Price {
			// traverse all orders that match
			for i := n - 1; i >= 0; i-- {
				sellOrder := ob.SellOrders[i]
				if sellOrder.Price > order.Price {
					break
				}
				// fill the entire order
				if sellOrder.Amount >= order.Amount {
					trades = append(trades, Trade{order.ID, sellOrder.ID, order.Amount, sellOrder.Price})
					sellOrder.Amount -= order.Amount
					if sellOrder.Amount == 0 {
						ob.RemoveSell(i)
					}
					return trades
				}
				// fill a partial order and continue
				if sellOrder.Amount < order.Amount {
					trades = append(trades, Trade{order.ID, sellOrder.ID, sellOrder.Amount, sellOrder.Price})
					order.Amount -= sellOrder.Amount
					ob.RemoveSell(i)
					continue
				}
			}
		}
	}
	// finally add the remaining order to the list
	ob.AddBuy(order)
	return trades
}

// Process a limit sell order
func (ob *OrderBook) processLimitSell(order *Order) []Trade {
	trades := make([]Trade, 0, 1)
	n := len(ob.BuyOrders)
	// check if we have at least one matching order
	if n != 0 {
		if ob.BuyOrders[n-1].Price >= order.Price {
			// traverse all orders that match
			for i := n - 1; i >= 0; i-- {
				buyOrder := ob.BuyOrders[i]
				if buyOrder.Price < order.Price {
					break
				}
				// fill the entire order
				if buyOrder.Amount >= order.Amount {
					trades = append(trades, Trade{order.ID, buyOrder.ID, order.Amount, buyOrder.Price})
					buyOrder.Amount -= order.Amount
					if buyOrder.Amount == 0 {
						ob.RemoveBuy(i)
					}
					return trades
				}
				// fill a partial order and continue
				if buyOrder.Amount < order.Amount {
					trades = append(trades, Trade{order.ID, buyOrder.ID, buyOrder.Amount, buyOrder.Price})
					order.Amount -= buyOrder.Amount
					ob.RemoveBuy(i)
					continue
				}
			}
		}
	}
	// finally add the remaining order to the list
	ob.Append(order)
	return trades
}
