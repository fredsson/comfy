package main

import (
	"errors"
	"strconv"
	"time"
)

type PricesCache struct {
	prices        map[string]*HourlyPrice
	fetchCallback FetchPrices
}

type FetchPrices = func() []HourlyPrice

func initPricesCache(fetchCallback FetchPrices) *PricesCache {
	var pricesCache = &PricesCache{prices: nil, fetchCallback: fetchCallback}

	pricesCache.refreshPrices(pricesCache.fetchCallback())

	return pricesCache
}

func (pricesCache *PricesCache) refreshPrices(pricesToday []HourlyPrice) {
	pricesCache.prices = make(map[string]*HourlyPrice, 48)
	for _, value := range pricesToday {
		var key string = pricesCache.getLookupKey(value.StartsAt)
		pricesCache.prices[key] = &value
	}
}

func (pricesCache *PricesCache) getHourlyPrice(currentTime time.Time) (*HourlyPrice, error) {
	var key string = pricesCache.getLookupKey(currentTime)

	if !pricesCache.cacheContainsKey(key) {
		pricesCache.refreshPrices(pricesCache.fetchCallback())
	}
	var hourlyPrice *HourlyPrice = pricesCache.prices[key]

	if hourlyPrice == nil {
		return nil, errors.New("could not find Hourly Price")
	}
	return hourlyPrice, nil
}

func (*PricesCache) getLookupKey(currentTime time.Time) string {
	var day = currentTime.Day()
	var hour = currentTime.Hour()
	return strconv.Itoa(day) + strconv.Itoa(hour)
}

func (pricesCache *PricesCache) cacheContainsKey(key string) bool {
	return pricesCache.prices[key] != nil
}
