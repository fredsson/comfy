package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type HeaderDefinition struct {
	name  string
	value string
}

type HourlyPrice struct {
	Total    float32
	StartsAt time.Time
}

type PricesResponse struct {
	Data struct {
		Viewer struct {
			Homes []struct {
				CurrentSubscription struct {
					PriceInfo struct {
						Today    []HourlyPrice
						Tomorrow []HourlyPrice
					}
				}
			}
		}
	}
}

const pricesQuery = `{"query":"{\nviewer{\nhomes{\ncurrentSubscription{\npriceInfo{\ntoday{\ntotal\nstartsAt\n}\ntomorrow{\ntotal\nstartsAt\n}\n}\n}\n}\n}\n}"}`

type TibberProxy struct {
	apiKey string
}

func (p TibberProxy) FetchPrices() []HourlyPrice {
	var url = "https://api.tibber.com/v1-beta/gql"
	var mapper = func(body io.ReadCloser) interface{} {
		return mapToPricesResponse(body)
	}

	var body = strings.NewReader(pricesQuery)
	var headers = []HeaderDefinition{
		{"Content-Type", "application/json"},
		{"Authorization", "Bearer " + p.apiKey},
	}
	var response = Post(url, body, headers, mapper).(*PricesResponse)
	var today = response.Data.Viewer.Homes[0].CurrentSubscription.PriceInfo.Today
	var tomorrow = response.Data.Viewer.Homes[0].CurrentSubscription.PriceInfo.Tomorrow
	return append(today, tomorrow...)
}

func mapToPricesResponse(body io.ReadCloser) *PricesResponse {
	response := new(PricesResponse)
	err := decodeFromJson(body, response)
	if err != nil {
		log.Fatal("Could not decode Prices response", err)
	}
	return response
}

func Post(url string, body io.Reader, headers []HeaderDefinition, m mapper) interface{} {
	var req, err = http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatal("woopsi ", err)
	}

	for _, h := range headers {
		req.Header.Add(h.name, h.value)
	}

	log.Println("Sending Post request to " + url)
	var resp, reqErr = http.DefaultClient.Do(req)
	if reqErr != nil {
		log.Fatal("woopsi ", reqErr)
	}
	defer resp.Body.Close()

	return m(resp.Body)
}
