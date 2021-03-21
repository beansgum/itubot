package main

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2"
)

func main() {

	itubo, err := newItubot()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = itubo.run()
	if err != nil {
		fmt.Println(err)
		return
	}

	select {}
}

type Itubot struct {
	client         *binance.Client
	tradeSchedules map[string]TradeSchedule
	closeTicker    chan bool
}

func newItubot() (*Itubot, error) {

	client := binance.NewClient(API_KEY, SECRET)

	fmt.Println("=> Initialized Binance client")

	tradeSchedules, err := readSchedules()
	if err != nil {
		return nil, fmt.Errorf("error reading trade schedules: %v", err)
	}

	return &Itubot{client: client, tradeSchedules: tradeSchedules, closeTicker: make(chan bool)}, nil
}

func (itubo *Itubot) run() error {
	fmt.Println("=> Itubo is running!")

	err := itubo.printTickerPrices()
	if err != nil {
		return err
	}

	err = itubo.printOpenMarginOrders()
	if err != nil {
		return err
	}

	fmt.Printf("=> Isolated margin accounts\n\n")

	for _, tradeSchedule := range itubo.tradeSchedules {
		isolatedAccount, err := itubo.client.NewGetIsolatedMarginAccountService().
			Symbols(tradeSchedule.Symbol).Do(context.Background())
		if err != nil {
			return err
		}

		for _, asset := range isolatedAccount.Assets {
			fmt.Printf("\t%s Price: $%s, Liquidate Price: $%s\n", asset.Symbol, asset.IndexPrice, asset.LiquidatePrice)

			fmt.Printf("\t%s Equity: %s, Borrowed: %s, Equity value: %s, Free: %s\n", asset.BaseAsset.Asset, asset.BaseAsset.NetAsset, asset.BaseAsset.Borrowed, asset.BaseAsset.TotalAsset, asset.BaseAsset.Free)
			fmt.Printf("\t%s Equity: %s, Borrowed: %s, Equity value: %s, Free: %s\n", asset.QuoteAsset.Asset, asset.QuoteAsset.NetAsset, asset.QuoteAsset.Borrowed, asset.QuoteAsset.TotalAsset, asset.QuoteAsset.Free)

			fmt.Println()
		}
	}

	// err = itubo.cancelMarginOrder(374773484, "ATOMUSDT")
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("Placing margin buy order")
	// order, err := itubo.createMarginOrder("ATOMUSDT", 19.2, 2, binance.SideTypeBuy)
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("ID: %d, Client ID: %s, Borrow amount: %s, Borrow asset: %s\n", order.OrderID, order.ClientOrderID, order.MarginBuyBorrowAmount, order.MarginBuyBorrowAsset)

	itubo.startTrading()

	return nil
}

func (itubo *Itubot) startTrading() {

	fmt.Println("=> Trader is running")

	go func() {

		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-itubo.closeTicker:
				ticker.Stop()
				fmt.Println("Trading channel closed")
				return
			case <-ticker.C:
				fmt.Println("=> Tick!")

				schedules, err := readSchedules()
				if err != nil {
					fmt.Println("unable to read trade schedules:", err)
					continue
				}

				for _, tradeSchedule := range schedules {

					fmt.Printf("\t%s\n\n", tradeSchedule.Symbol)

					price, err := itubo.fetchTickerPrice(tradeSchedule.Symbol)
					if err != nil {
						fmt.Println("error fetch ticker price:", err)
						continue
					}

					isolatedMarginAccount, err := itubo.client.NewGetIsolatedMarginAccountService().
						Symbols(tradeSchedule.Symbol).Do(context.Background())
					if err != nil {
						fmt.Println("error fetch isolated:", err)
						continue
					}

					// `isolatedMarginAccount` constains only one asset which is the symbol passed.
					baseAsset := isolatedMarginAccount.Assets[0].BaseAsset
					quoteAsset := isolatedMarginAccount.Assets[0].QuoteAsset

					fmt.Printf("\tEntry: %f, Target: %f, Current: %s, %s balance: %s debt: %s, %s balance: %s, Debt: %s\n", tradeSchedule.EntryPrice, tradeSchedule.TargetPrice, price.Price,
						baseAsset.Asset, baseAsset.Free, baseAsset.Borrowed, quoteAsset.Asset, quoteAsset.Free, quoteAsset.Borrowed)

					openOrders, err := itubo.fetchOpenMarginOrders(tradeSchedule.Symbol)
					if err != nil {
						fmt.Println("error fetching open orders:", err)
						continue
					}

					if len(openOrders) > 0 {
						fmt.Printf("\t=> Open %s orders\n", tradeSchedule.Symbol)

						for _, order := range openOrders {
							fmt.Printf("\t%s %s %s @ $%s status: %s\n", order.Side, order.OrigQuantity, order.Symbol, order.Price, order.Status)
						}
					} else {
						err = itubo.marginShort(tradeSchedule, baseAsset, quoteAsset)
						if err != nil {
							fmt.Println("\terror margin trading:", err)
							continue
						}
					}

					fmt.Printf("\n\n")
				}
			}
		}
	}()

}

func (itubo *Itubot) stopTrading() {
	itubo.closeTicker <- true
}

func (itubo *Itubot) printTickerPrices() error {
	prices, err := itubo.fetchTickersPrices()
	if err != nil {
		return err
	}

	fmt.Println("=> Current prices")
	for _, p := range prices {
		if _, ok := itubo.tradeSchedules[p.Symbol]; ok {
			fmt.Printf("\t%s: \t$%s\n", p.Symbol, p.Price)
		}
	}
	fmt.Println()

	return nil
}

func (itubo *Itubot) printOpenMarginOrders() error {
	fmt.Println("=> Open margin orders")

	for _, tradeSchedule := range itubo.tradeSchedules {
		openOrders, err := itubo.fetchOpenMarginOrders(tradeSchedule.Symbol)
		if err != nil {
			return err
		}

		for _, order := range openOrders {
			fmt.Printf("\t%s %s %s @ $%s status: %s\n", order.Side, order.OrigQuantity, order.Symbol, order.Price, order.Status)
		}
	}

	fmt.Println()

	return nil
}
