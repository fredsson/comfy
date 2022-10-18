package main

import (
	"errors"
	"log"
	"strconv"
	"time"
)

type PriceCache struct {
	prices        map[string]*HourlyPrice
	fetchCallback FetchPrices
}

type FetchPrices = func() []HourlyPrice

func initPriceCache(fetchCallback FetchPrices) *PriceCache {
	var priceCache = &PriceCache{prices: nil, fetchCallback: fetchCallback}

	priceCache.refreshPrices(priceCache.fetchCallback())

	return priceCache
}

func (priceCache *PriceCache) refreshPrices(pricesToday []HourlyPrice) {
	log.Printf("Refreshing prices!")
	priceCache.prices = make(map[string]*HourlyPrice, 48)
	for _, value := range pricesToday {
		var key string = priceCache.getLookupKey(value.StartsAt)
		priceCache.prices[key] = &value
	}
}

func (priceCache *PriceCache) getHourlyPrice(currentTime time.Time) (*HourlyPrice, error) {
	var key string = priceCache.getLookupKey(currentTime)

	if !priceCache.cacheContainsKey(key) {
		priceCache.refreshPrices(priceCache.fetchCallback())
	}
	var hourlyPrice *HourlyPrice = priceCache.prices[key]

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
	return priceCache.prices[key] != nil
}

func convertToDoubleDigitString(value int) string {
	if value >= 0 && value < 10 {
		return "0" + strconv.Itoa(value)
	}
	return strconv.Itoa(value)
}
