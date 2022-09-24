package main

import (
	"errors"
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

	var prices = tibberProxy.FetchPricesToday()
	for {
		var currentPrice, priceErr = findHourlyPriceNow(prices)
		if priceErr != nil {
			log.Fatal(priceErr)
		}

		var shouldSmartModeBeEnabledNow = (currentPrice.Total <= 1.0000)
		if shouldSmartModeBeEnabledNow != SmartModeEnabledInSensibo {
			if shouldSmartModeBeEnabledNow {
				sensiboProxy.EnableSmartMode(pod)
			} else {
				sensiboProxy.DisableSmartMode(pod)
			}
		}
		log.Println("Waiting an hour for next update")
		time.Sleep(time.Hour) // TODO: calculate how much to sleep
	}
}

func findHourlyPriceNow(array []HourlyPrice) (*HourlyPrice, error) {
	var currentTime time.Time = time.Now().UTC()
	for _, v := range array {
		var difference = v.StartsAt.UTC().Sub(currentTime)
		if difference > -time.Hour && difference <= 0 {
			return &v, nil
		}
	}

	return nil, errors.New("could not find Hourly Price")
}
