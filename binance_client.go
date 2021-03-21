package main

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
)

func (itubo *Itubot) fetchTickerPrice(symbol string) (*binance.SymbolPrice, error) {
	prices, err := itubo.fetchTickersPrices()
	if err != nil {
		return nil, err
	}

	for _, price := range prices {
		if price.Symbol == symbol {
			return price, nil
		}
	}

	return nil, fmt.Errorf("symbol not found in ticker symbol price slice")
}

func (itubo *Itubot) fetchTickersPrices() ([]*binance.SymbolPrice, error) {
	return itubo.client.NewListPricesService().Do(context.Background())
}

func (itubo *Itubot) fetchOpenMarginOrders(symbol string) ([]*binance.Order, error) {
	return itubo.client.NewListMarginOpenOrdersService().IsIsolated(true).Symbol(symbol).
		Do(context.Background())
}

func (itubo *Itubot) createMarginOrder(symbol string, price, quantity float64, side binance.SideType) (*binance.CreateOrderResponse, error) {
	fmt.Printf("\tPrice %f, Quantity: %f, symbol %s\n", price, quantity, symbol)
	order, err := itubo.client.NewCreateMarginOrderService().Symbol(symbol).Side(side).IsIsolated(true).
		Type(binance.OrderTypeLimit).TimeInForce(binance.TimeInForceTypeGTC).Quantity(fmt.Sprint(quantity)).
		Price(fmt.Sprint(price)).Do(context.Background())

	return order, err
}

func (itubo *Itubot) cancelMarginOrder(orderID int64, symbol string) error {
	_, err := itubo.client.NewCancelMarginOrderService().Symbol(symbol).IsIsolated(true).
		OrderID(orderID).Do(context.Background())

	return err
}
