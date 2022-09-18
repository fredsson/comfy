package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type HourlyPrice struct {
	Total    float32
	StartsAt time.Time
}

type PricesTodayResponse struct {
	Data struct {
		Viewer struct {
			Homes []struct {
				CurrentSubscription struct {
					PriceInfo struct {
						Today []HourlyPrice
					}
				}
			}
		}
	}
}

const pricesTodayQuery = `{"query":"{\nviewer{\nhomes{\ncurrentSubscription{\npriceInfo{\ntoday{\n total\n startsAt\n}\n\n}\n}\n}\n}\n}\n"}`

type TibberProxy struct {
	apiKey string
}

func (p TibberProxy) FetchPricesToday() []HourlyPrice {
	var url = "https://api.tibber.com/v1-beta/gql"
	var mapper = func(body io.ReadCloser) interface{} {
		return mapToPricesTodayResponse(body)
	}

	var response = Post(url, p.apiKey, mapper).(*PricesTodayResponse)
	return response.Data.Viewer.Homes[0].CurrentSubscription.PriceInfo.Today
}

func mapToPricesTodayResponse(body io.ReadCloser) *PricesTodayResponse {
	response := new(PricesTodayResponse)
	err := decodeFromJson(body, response)
	if err != nil {
		log.Fatal("Could not decode Prices today response", err)
	}
	return response
}

func Post(url string, apiKey string, m mapper) interface{} {
	var myReader = strings.NewReader(pricesTodayQuery)
	var req, err = http.NewRequest("POST", url, myReader)
	if err != nil {
		log.Fatal("woopsi ", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	log.Println("Sending Post request to " + url)
	var resp, reqErr = http.DefaultClient.Do(req)
	if reqErr != nil {
		log.Fatal("woopsi ", reqErr)
	}
	defer resp.Body.Close()

	return m(resp.Body)
}
