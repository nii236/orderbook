package orderbook

import "fmt"

type Pair string

func (p *Pair) BuyKey() string {
	return fmt.Sprintf("%s-%s", string(*p), BUY)
}
func (p *Pair) SellKey() string {
	return fmt.Sprintf("%s-%s", string(*p), SELL)
}
