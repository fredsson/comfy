package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	var envErr error = godotenv.Load()
	if envErr != nil {
		log.Fatalf("Some error occured. Err: %s", envErr)
	}

	var sensiboProxy = SensiboProxy{
		apiKey: os.Getenv("SENSIBO_API_KEY"),
	}
	var pods = sensiboProxy.FetchPods()
	var pod = pods[0]
	var smartMode = sensiboProxy.FetchSmartModeForPod(pod)
	var SmartModeEnabledInSensibo = smartMode.Enabled

	var tibberProxy = TibberProxy{
		apiKey: os.Getenv("TIBBER_API_KEY"),
	}
	var priceCache = initPriceCache(tibberProxy.FetchPrices)

	for {
		var currentPrice, priceErr = priceCache.getHourlyPrice(time.Now())
		var nextHoursPrice, priceErrNextHour = priceCache.getHourlyPrice(time.Now().Add(time.Hour))

		if priceErr != nil {
			log.Fatal(priceErr)
		}
		if priceErrNextHour != nil {
			log.Fatal(priceErrNextHour)
		}

		var shouldSmartModeBeEnabledNow = !priceExceedsTreshold(currentPrice)
		if shouldSmartModeBeEnabledNow != SmartModeEnabledInSensibo {
			if shouldSmartModeBeEnabledNow {
				SmartModeEnabledInSensibo = true
				sensiboProxy.EnableSmartMode(pod)
			} else {
				SmartModeEnabledInSensibo = false
				sensiboProxy.DisableSmartMode(pod)
			}
		}
		if shouldSmartModeBeEnabledNow && priceExceedsTreshold(nextHoursPrice) {
			sensiboProxy.enableAc(pod)
		}

		var timeForNextRun = getTimeForNextRun()
		var now = time.Now().UTC()
		var durationToNextRun = timeForNextRun.Sub(now)
		log.Println("Waiting", durationToNextRun, "for next update (", timeForNextRun, ")")
		time.Sleep(durationToNextRun)
	}
}

func getTimeForNextRun() time.Time {
	return time.Now().UTC().Add(time.Hour).Truncate(time.Hour)
}

func priceExceedsTreshold(price *HourlyPrice) bool {
	return price.Total > 1.0000
}
