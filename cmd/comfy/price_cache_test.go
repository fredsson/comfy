package main

import (
	"testing"
	"time"
)

func buildTime(datestr string, t *testing.T) time.Time {
	var layout = "2006-01-02T15:04:05-07:00"
	var time, err = time.Parse(layout, datestr)
	if err != nil {
		t.Error("Could not parse date: ", datestr)
	}

	return time
}

func TestInitCacheShouldFetchPrices(t *testing.T) {
	var called = false
	var callback = func() []HourlyPrice {
		called = true
		return nil
	}

	initPriceCache(callback)

	if called == false {
		t.Error("did not call callback for fetching prices")
	}
}

func TestGetHourlyPriceShouldReturnCorrectPrice(t *testing.T) {
	var callback = func() []HourlyPrice {
		return []HourlyPrice{
			{Total: 1.2, StartsAt: buildTime("2022-10-24T22:00:00+02:00", t)},
			{Total: 1.4, StartsAt: buildTime("2022-10-24T23:00:00+02:00", t)},
			{Total: 0.5, StartsAt: buildTime("2022-10-25T00:00:00+02:00", t)},
		}
	}

	var cache = initPriceCache(callback)

	var price, err = cache.getHourlyPrice(buildTime("2022-10-24T20:00:00+00:00", t))

	if err != nil {
		t.Error("Received error from cache: ", err)
	}
	if price.Total != 1.2 {
		t.Errorf("Did not receive correct price from cache %f expected %f", price.Total, 1.2)
	}
}
