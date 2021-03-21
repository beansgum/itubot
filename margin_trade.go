package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

const minBuyPrice = 11 // usdt

func (itubo *Itubot) marginShort(tradeSchedule TradeSchedule, baseAsset, quoteAsset binance.IsolatedUserAsset) error {
	freeBaseAsset, _ := strconv.ParseFloat(baseAsset.Free, 64)
	borrowedBaseAsset, _ := strconv.ParseFloat(baseAsset.Borrowed, 64)

	borrowedBaseAsset -= freeBaseAsset // deduct available assets from borrowed assets

	freeBaseAsset = math.Floor(freeBaseAsset*tradeSchedule.DecimalPlaces) / tradeSchedule.DecimalPlaces // round down to decimal places
	buyPriceInUSDT := freeBaseAsset * tradeSchedule.TargetPrice

	fmt.Println("\t Decimal Places:", tradeSchedule.DecimalPlaces)
	// margin short base asset
	if freeBaseAsset > 0 && buyPriceInUSDT >= minBuyPrice {
		// sell assets for usdt

		fmt.Printf("\tWill create a sell order for %f %s @ %f %s per %s Totalling %f %s\n", freeBaseAsset, baseAsset.Asset, tradeSchedule.TargetPrice,
			quoteAsset.Asset, baseAsset.Asset, buyPriceInUSDT, quoteAsset.Asset)

		order, err := itubo.createMarginOrder(tradeSchedule.Symbol, tradeSchedule.TargetPrice, freeBaseAsset, binance.SideTypeSell)
		if err != nil {
			return fmt.Errorf("error creating sell order: %v", err)
		}

		fmt.Printf("\tOrder successful, ID: %d\n", order.OrderID)
	} else {
		// sell usdt to buy base assets
		borrowedBaseAsset = math.Ceil(borrowedBaseAsset*tradeSchedule.DecimalPlaces) / tradeSchedule.DecimalPlaces // round up to decimal places

		buyPriceInUSDT := borrowedBaseAsset * tradeSchedule.EntryPrice
		if buyPriceInUSDT <= minBuyPrice {
			return fmt.Errorf("borrowed & free assets are less than min buy price")
		}

		fmt.Printf("\tWill create a buy order for %f %s @ %f %s per %s Totalling %f %s\n", borrowedBaseAsset, baseAsset.Asset, tradeSchedule.EntryPrice,
			quoteAsset.Asset, baseAsset.Asset, buyPriceInUSDT, quoteAsset.Asset)

		order, err := itubo.createMarginOrder(tradeSchedule.Symbol, tradeSchedule.EntryPrice, borrowedBaseAsset, binance.SideTypeBuy)
		if err != nil {
			return fmt.Errorf("error creating buy order: %v", err)
		}

		fmt.Printf("\tOrder successful, ID: %d\n", order.OrderID)

	}

	return nil
}
