package main

import (
	"example.com/covid-go/model"
	"example.com/covid-go/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func Covid19Handler(c *gin.Context) {
	covidData, err := repository.Api.GetCovid19Data()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, handleCovidData(covidData))
	}
}

func GetUsersHandler(c *gin.Context) {
	users, err := repository.DB.GetUsers()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, users)
	}
}

func handleCovidData(data model.CovidData) model.CovidSummary {
	groupedSlices := make(map[string]model.CovidByProvince)
	for _, covidCase := range data.Data {
		if province, ok := groupedSlices[covidCase.Province]; ok {
			province.Count = province.Count + 1
			if covidCase.ConfirmDate.After(province.LastDate) {
				province.LastDate = covidCase.ConfirmDate.Time
			}
			groupedSlices[covidCase.Province] = province
		} else {
			province := model.CovidByProvince{
				Province: covidCase.Province,
				Count:    1,
				LastDate: covidCase.ConfirmDate.Time,
			}
			groupedSlices[covidCase.Province] = province
		}
	}
	var summaryData []model.CovidByProvince
	for _, v := range groupedSlices {
		summaryData = append(summaryData, v)
	}
	sort.Slice(summaryData, func(i, j int) bool {
		return summaryData[i].LastDate.After(summaryData[j].LastDate)
	})
	return model.CovidSummary{
		Data: summaryData,
	}
}
