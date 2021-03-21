package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func readSchedules() (map[string]TradeSchedule, error) {

	schedulesFile, err := os.Open("schedules.json")
	if err != nil {
		return nil, fmt.Errorf("error opening trade schedules json: %v", err)
	}

	var tradeSchedules []TradeSchedule

	jsonParser := json.NewDecoder(schedulesFile)
	if err = jsonParser.Decode(&tradeSchedules); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	tradeSchedulesMap := make(map[string]TradeSchedule)
	for _, tradeSchedule := range tradeSchedules {
		tradeSchedulesMap[tradeSchedule.Symbol] = tradeSchedule
	}

	return tradeSchedulesMap, nil
}

// func findSchedule(schedules []TradeSchedule, symbol string) (*TradeSchedule, error) {
// 	for _, schedule := range schedules {
// 		if schedule.Symbol == symbol {
// 			return &schedule, nil
// 		}
// 	}

// 	return nil, fmt.Errorf("trade schedule not found for '%s'", symbol)
// }
