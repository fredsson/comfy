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
	var pricesCache = initPricesCache(tibberProxy.FetchPricesToday)

	for {
		var currentPrice, priceErr = pricesCache.getHourlyPrice(time.Now())

		if priceErr != nil {
			log.Fatal(priceErr)
		}

		var shouldSmartModeBeEnabledNow = (currentPrice.Total <= 1.0000)
		if shouldSmartModeBeEnabledNow != SmartModeEnabledInSensibo {
			if shouldSmartModeBeEnabledNow {
				SmartModeEnabledInSensibo = true
				sensiboProxy.EnableSmartMode(pod)
			} else {
				SmartModeEnabledInSensibo = false
				sensiboProxy.DisableSmartMode(pod)
			}
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
