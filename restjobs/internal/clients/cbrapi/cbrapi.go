package cbrapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	cbr_url = "https://www.cbr-xml-daily.ru/daily_json.js"
)

type CbrClient struct {
}

type cbr_response struct {
	Date   string
	Valute valutes
}

type valutes struct {
	USD valute
	EUR valute
}

type valute struct {
	CharCode string
	Value    json.Number `json:"Value"`
}

func (c *CbrClient) GetRate(valute string) string {

	res := new(cbr_response)
	err := getJson(cbr_url, res)
	if err != nil {
		errText := "[restjobs] getRate: " + err.Error()
		log.Println(errText)
		return errText
	}
	switch valute {
	case "USD":
		return fmt.Sprintf("%s %s <b>%s</b>", res.Date, res.Valute.USD.CharCode, string(res.Valute.USD.Value))
	case "EUR":
		return fmt.Sprintf("%s %s <b>%s</b>", res.Date, res.Valute.EUR.CharCode, string(res.Valute.EUR.Value))
	}

	return ""
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
