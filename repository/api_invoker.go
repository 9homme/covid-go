package repository

import (
	"encoding/json"
	"example.com/covid-go/model"
	"net/http"
)

var Api ApiInvokerInterface = ApiInvoker{}

//go:generate mockgen --source api_invoker.go --destination mock/mock_api.go
type ApiInvokerInterface interface {
	GetCovid19Data() (model.CovidData, error)
}

type ApiInvoker struct{}

func (ApiInvoker) GetCovid19Data() (model.CovidData, error) {
	covidData := model.CovidData{}
	err := getJson("https://covid19.th-stat.com/api/open/cases", &covidData)
	return covidData, err
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
