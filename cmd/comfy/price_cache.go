package main

import (
	"errors"
	"log"
	"strconv"
	"time"
)

type FetchPrices = func() []HourlyPrice

type PriceCache struct {
	prices        []HourlyPrice
	pricesByKey   map[string]*HourlyPrice
	fetchCallback FetchPrices
}

func initPriceCache(fetchCallback FetchPrices) *PriceCache {
	var priceCache = &PriceCache{prices: nil, fetchCallback: fetchCallback}

	priceCache.refreshPrices(priceCache.fetchCallback())

	return priceCache
}

func (priceCache *PriceCache) refreshPrices(prices []HourlyPrice) {
	log.Printf("Refreshing prices!")
	priceCache.prices = prices
	priceCache.pricesByKey = make(map[string]*HourlyPrice, len(prices))
	for index, value := range prices {
		var key string = priceCache.getLookupKey(value.StartsAt.UTC())
		priceCache.pricesByKey[key] = &prices[index]
	}
}

func (priceCache *PriceCache) getHourlyPrice(currentTime time.Time) (*HourlyPrice, error) {
	var key string = priceCache.getLookupKey(currentTime.UTC())

	if !priceCache.cacheContainsKey(key) {
		priceCache.refreshPrices(priceCache.fetchCallback())
	}
	var hourlyPrice *HourlyPrice = priceCache.pricesByKey[key]

	if hourlyPrice == nil {
		return nil, errors.New("could not find Hourly Price")
	}
	return hourlyPrice, nil
}

func (*PriceCache) getLookupKey(currentTime time.Time) string {
	var day = currentTime.Day()
	var hour = currentTime.Hour()
	return convertToDoubleDigitString(day) + convertToDoubleDigitString(hour)
}

func (priceCache *PriceCache) cacheContainsKey(key string) bool {
	return priceCache.pricesByKey[key] != nil
}

func convertToDoubleDigitString(value int) string {
	if value >= 0 && value < 10 {
		return "0" + strconv.Itoa(value)
	}
	return strconv.Itoa(value)
}
