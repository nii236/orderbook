package orderbook_test

import (
	"orderbook"
	"sync"
	"testing"
)

func TestAddBuy(t *testing.T) {
	type fields struct {
		Conn       orderbook.KV
		Pair       orderbook.Pair
		BuyOrders  []*orderbook.Order
		SellOrders []*orderbook.Order
		RWMutex    *sync.RWMutex
	}
	type args struct {
		order *orderbook.Order
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		bookLength int
		id         string
		index      int
	}{
		{
			"nil input",
			fields{nil, "BTCUSD", []*orderbook.Order{}, nil, &sync.RWMutex{}},
			args{nil},
			1,
			"",
			0,
		},
		{
			"no orders",
			fields{nil, "BTCUSD", []*orderbook.Order{}, nil, &sync.RWMutex{}},
			args{&orderbook.Order{1, 10, "order-01", orderbook.BUY}},
			1,
			"order-01",
			0,
		},
		{
			"one orders same price",
			fields{nil, "BTCUSD", []*orderbook.Order{&orderbook.Order{1, 10, "order-01", orderbook.BUY}}, nil, &sync.RWMutex{}},
			args{&orderbook.Order{1, 10, "order-02", orderbook.BUY}},
			2,
			"order-02",
			0,
		},
		{
			"one orders higher price",
			fields{nil, "BTCUSD", []*orderbook.Order{&orderbook.Order{1, 10, "order-01", orderbook.BUY}}, nil, &sync.RWMutex{}},
			args{&orderbook.Order{1, 11, "order-02", orderbook.BUY}},
			2,
			"order-02",
			0,
		},
		{
			"one orders lower price",
			fields{nil, "BTCUSD", []*orderbook.Order{&orderbook.Order{1, 10, "order-01", orderbook.BUY}}, nil, &sync.RWMutex{}},
			args{&orderbook.Order{1, 9, "order-02", orderbook.BUY}},
			2,
			"order-02",
			1,
		},
		{
			"two orders same price",
			fields{nil, "BTCUSD", []*orderbook.Order{&orderbook.Order{1, 10, "order-01", orderbook.BUY}, &orderbook.Order{1, 10, "order-02", orderbook.BUY}}, nil, &sync.RWMutex{}},
			args{&orderbook.Order{1, 10, "order-03", orderbook.BUY}},
			3,
			"order-03",
			2,
		},
		{
			"two orders higher price",
			fields{nil, "BTCUSD", []*orderbook.Order{&orderbook.Order{1, 10, "order-01", orderbook.BUY}, &orderbook.Order{1, 10, "order-02", orderbook.BUY}}, nil, &sync.RWMutex{}},
			args{&orderbook.Order{1, 11, "order-03", orderbook.BUY}},
			3,
			"order-03",
			0,
		},
		{
			"two orders lower price",
			fields{nil, "BTCUSD", []*orderbook.Order{&orderbook.Order{1, 10, "order-01", orderbook.BUY}, &orderbook.Order{1, 10, "order-02", orderbook.BUY}}, nil, &sync.RWMutex{}},
			args{&orderbook.Order{1, 9, "order-03", orderbook.BUY}},
			3,
			"order-03",
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := &orderbook.OrderBook{
				Conn:       tt.fields.Conn,
				Pair:       tt.fields.Pair,
				BuyOrders:  tt.fields.BuyOrders,
				SellOrders: tt.fields.SellOrders,
				RWMutex:    tt.fields.RWMutex,
			}
			ob.AddBuy(tt.args.order)
			if len(ob.BuyOrders) != tt.bookLength {
				orders := []string{}
				for _, o := range ob.BuyOrders {
					orders = append(orders, o.ID)
				}
				t.Errorf("expected %d orders, got %d, %+v", tt.bookLength, len(ob.BuyOrders), orders)
			}

			for i, o := range ob.BuyOrders {
				if tt.id == "" {
					continue
				}
				if o.ID == tt.id && i != tt.index {
					t.Errorf("expected %d index, got %d", tt.index, i)
				}
			}
		})
	}
}

func TestAddSell(t *testing.T) {
	type fields struct {
		Conn       orderbook.KV
		Pair       orderbook.Pair
		BuyOrders  []*orderbook.Order
		SellOrders []*orderbook.Order
		RWMutex    *sync.RWMutex
	}
	type args struct {
		order *orderbook.Order
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		bookLength int
		id         string
		index      int
	}{
		{
			"nil input",
			fields{nil, "BTCUSD", nil, []*orderbook.Order{}, &sync.RWMutex{}},
			args{nil},
			1,
			"",
			0,
		},
		{
			"no orders",
			fields{nil, "BTCUSD", nil, []*orderbook.Order{}, &sync.RWMutex{}},
			args{&orderbook.Order{1, 10, "order-01", orderbook.SELL}},
			1,
			"order-01",
			0,
		},
		{
			"one orders same price",
			fields{nil, "BTCUSD", nil, []*orderbook.Order{&orderbook.Order{1, 10, "order-01", orderbook.SELL}}, &sync.RWMutex{}},
			args{&orderbook.Order{1, 10, "order-02", orderbook.SELL}},
			2,
			"order-02",
			0,
		},
		{
			"one orders higher price",
			fields{nil, "BTCUSD", nil, []*orderbook.Order{&orderbook.Order{1, 10, "order-01", orderbook.SELL}}, &sync.RWMutex{}},
			args{&orderbook.Order{1, 11, "order-02", orderbook.SELL}},
			2,
			"order-02",
			0,
		},
		{
			"one orders lower price",
			fields{nil, "BTCUSD", nil, []*orderbook.Order{&orderbook.Order{1, 10, "order-01", orderbook.SELL}}, &sync.RWMutex{}},
			args{&orderbook.Order{1, 9, "order-02", orderbook.SELL}},
			2,
			"order-02",
			1,
		},
		{
			"two orders same price",
			fields{nil, "BTCUSD", nil, []*orderbook.Order{&orderbook.Order{1, 10, "order-01", orderbook.SELL}, &orderbook.Order{1, 10, "order-02", orderbook.SELL}}, &sync.RWMutex{}},
			args{&orderbook.Order{1, 10, "order-03", orderbook.SELL}},
			3,
			"order-03",
			2,
		},
		{
			"two orders higher price",
			fields{nil, "BTCUSD", nil, []*orderbook.Order{&orderbook.Order{1, 10, "order-01", orderbook.SELL}, &orderbook.Order{1, 10, "order-02", orderbook.SELL}}, &sync.RWMutex{}},
			args{&orderbook.Order{1, 11, "order-03", orderbook.SELL}},
			3,
			"order-03",
			2,
		},
		{
			"two orders lower price",
			fields{nil, "BTCUSD", nil, []*orderbook.Order{&orderbook.Order{1, 10, "order-01", orderbook.SELL}, &orderbook.Order{1, 10, "order-02", orderbook.SELL}}, &sync.RWMutex{}},
			args{&orderbook.Order{1, 9, "order-03", orderbook.SELL}},
			3,
			"order-03",
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := &orderbook.OrderBook{
				Conn:       tt.fields.Conn,
				Pair:       tt.fields.Pair,
				BuyOrders:  tt.fields.BuyOrders,
				SellOrders: tt.fields.SellOrders,
				RWMutex:    tt.fields.RWMutex,
			}
			ob.AddSell(tt.args.order)
			if len(ob.SellOrders) != tt.bookLength {
				orders := []string{}
				for _, o := range ob.SellOrders {
					orders = append(orders, o.ID)
				}
				t.Errorf("expected %d orders, got %d, %+v", tt.bookLength, len(ob.SellOrders), orders)
			}

			for i, o := range ob.SellOrders {
				if tt.id == "" {
					continue
				}
				if o.ID == tt.id && i != tt.index {
					t.Errorf("expected %d index, got %d", tt.index, i)
				}
			}
		})
	}
}
