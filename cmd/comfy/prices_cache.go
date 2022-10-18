package main

import (
	"errors"
	"time"
)

type PricesCache struct {
	prices        map[time.Time]HourlyPrice
	fetchCallback FetchPrices
}

type FetchPrices = func() []HourlyPrice

func initPricesCache(fetchCallback FetchPrices) *PricesCache {
	var pricesCache = &PricesCache{prices: nil, fetchCallback: fetchCallback}

	pricesCache.refreshPrices(pricesCache.fetchCallback())

	return pricesCache
}

func (pricesCache *PricesCache) refreshPrices(pricesToday []HourlyPrice) {
	pricesCache.prices = make(map[time.Time]HourlyPrice, 24)
	for _, value := range pricesToday {
		pricesCache.prices[value.StartsAt] = value
	}
}

func (pricesCache *PricesCache) getHourlyPrice(currentTime time.Time) (*HourlyPrice, error) {
	if currentTime.Hour() == 0 {
		pricesCache.refreshPrices(pricesCache.fetchCallback())
	}

	for key, value := range pricesCache.prices {
		if key.Day() == currentTime.Day() && key.UTC().Hour() == currentTime.UTC().Hour() {
			return &value, nil
		}
	}

	return nil, errors.New("could not find Hourly Price")
}
