package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

func (itubo *Itubot) marginShort(tradeSchedule TradeSchedule, baseAsset, quoteAsset binance.IsolatedUserAsset) error {
	freeBaseAsset, _ := strconv.ParseFloat(baseAsset.Free, 64)
	borrowedBaseAsset, _ := strconv.ParseFloat(baseAsset.Borrowed, 64)
	borrowedBaseAsset -= freeBaseAsset
	// margin short base asset
	if freeBaseAsset > 0 {
		freeBaseAsset = math.Ceil(freeBaseAsset*100) / 100 // round to two decimals

		buyPrice := freeBaseAsset * tradeSchedule.TargetPrice
		fmt.Printf("\tWill create a sell order for %f %s @ %f %s per %s Totalling %f %s\n", freeBaseAsset, baseAsset.Asset, tradeSchedule.TargetPrice,
			quoteAsset.Asset, baseAsset.Asset, buyPrice, quoteAsset.Asset)

		order, err := itubo.createMarginOrder(tradeSchedule.Symbol, tradeSchedule.TargetPrice, freeBaseAsset, binance.SideTypeSell)
		if err != nil {
			return fmt.Errorf("error creating sell order: %v", err)
		}

		fmt.Printf("\tOrder successful, ID: %d\n", order.OrderID)
	} else {

		borrowedBaseAsset = math.Ceil(borrowedBaseAsset*100) / 100 // round to two decimals

		buyPrice := borrowedBaseAsset * tradeSchedule.EntryPrice
		fmt.Printf("\tWill create a buy order for %f %s @ %f %s per %s Totalling %f %s\n", borrowedBaseAsset, baseAsset.Asset, tradeSchedule.EntryPrice,
			quoteAsset.Asset, baseAsset.Asset, buyPrice, quoteAsset.Asset)

		order, err := itubo.createMarginOrder(tradeSchedule.Symbol, tradeSchedule.EntryPrice, borrowedBaseAsset, binance.SideTypeBuy)
		if err != nil {
			return fmt.Errorf("error creating buy order: %v", err)
		}

		fmt.Printf("\tOrder successful, ID: %d\n", order.OrderID)

		// // will borrow base assets and sell for target usdt price
		// fmt.Printf("\tWill borrow base assets and sell for target asset(%s) if target price has never hit\n", quoteAsset.Asset)
		// max, err := itubo.client.NewGetMaxBorrowableService().Asset(baseAsset.Asset).IsolatedSymbol(tradeSchedule.Symbol).Do(context.Background())
		// if err != nil {
		// 	fmt.Println("error getting max borrow:", err)
		// 	continue
		// }

		// maxBorrow, _ := strconv.ParseFloat(max.Amount, 64)
		// if maxBorrow > 0 {
		// 	fmt.Printf("\tMax borrow: %f, Will borrow %.2f%% : %f %s\n", maxBorrow, tradeSchedule.BorrowPercent*100,
		// 		maxBorrow*tradeSchedule.BorrowPercent, baseAsset.Asset)
		// } else {
		// 	fmt.Printf("\tCannot borrow, Fund %s to borrow %s\n", quoteAsset.Asset, baseAsset.Asset)
		// }
	}

	return nil
}
