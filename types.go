package main

type TradeSchedule struct {
	ID            int
	Symbol        string
	EntryPrice    float64
	TargetPrice   float64
	BorrowPercent float64
}
