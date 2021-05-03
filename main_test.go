package main

import (
	"encoding/json"
	"example.com/covid-go/model"
	"example.com/covid-go/repository"
	mock_repository "example.com/covid-go/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func TestGetCovidDataWithCorrectUserCredential(t *testing.T) {
	date := time.Now().Round(time.Minute)
	setupApiMock(t, date)
	setupDbMock(t)
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/covid19", nil)
	req.Header.Add("Authorization", "Basic dXNlcjp1c2VyMTIz")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	expected := model.CovidSummary{
		Data: []model.CovidByProvince{
			{
				Province: "Nonthaburi",
				Count:    2,
				LastDate: date.Add(2 * time.Hour),
			},
			{
				Province: "Bangkok",
				Count:    3,
				LastDate: date.Add(time.Hour),
			},
			{
				Province: "Ayuthaya",
				Count:    1,
				LastDate: date,
			},
		},
	}
	var actual model.CovidSummary
	err := json.NewDecoder(w.Body).Decode(&actual)
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, actual)

}

func TestGetCovidDataWithoutCredentialShouldFailed(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/covid19", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "{\"error\":\"Unauthorized\"}", w.Body.String())

}

func TestGetCovidDataWithInvalidCredentialShouldFailed(t *testing.T) {
	setupDbMock(t)
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/covid19", nil)
	req.Header.Add("Authorization", "Basic dXNlcjp1c2VyMTIzNDU2")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "{\"error\":\"Unauthorized\"}", w.Body.String())

}

func setupApiMock(t *testing.T, date time.Time) {
	ctrl := gomock.NewController(t)
	m := mock_repository.NewMockApiInvokerInterface(ctrl)
	m.EXPECT().GetCovid19Data().Return(getMockCovidData(date), nil)

	repository.Api = m

}

func setupDbMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mock_repository.NewMockDataSource(ctrl)
	db.EXPECT().GetUserByUsername(gomock.Any()).Return(model.User{
		Username:     "user",
		PasswordHash: "2dff86d9d96e4422da72b62a9966961a3c55d185a6f955589f1aacc59cdd198ef1ffea84c1e6033af39299151a12bb825ff3eb3780fd73c15b6d68d412f852d9",
	}, nil)

	repository.DB = db
}

func getMockCovidData(date time.Time) model.CovidData {
	return model.CovidData{
		Data: []model.CovidCase{
			{
				ConfirmDate: model.CustomTime{Time: date},
				Province:    "Bangkok",
			},
			{
				ConfirmDate: model.CustomTime{Time: date},
				Province:    "Nonthaburi",
			},
			{
				ConfirmDate: model.CustomTime{Time: date.Add(time.Hour)},
				Province:    "Bangkok",
			},
			{
				ConfirmDate: model.CustomTime{Time: date},
				Province:    "Ayuthaya",
			},
			{
				ConfirmDate: model.CustomTime{Time: date},
				Province:    "Bangkok",
			},
			{
				ConfirmDate: model.CustomTime{Time: date.Add(2 * time.Hour)},
				Province:    "Nonthaburi",
			},
		},
	}
}
