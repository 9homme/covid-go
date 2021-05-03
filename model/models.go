package model

import (
	"github.com/globalsign/mgo/bson"
	"strings"
	"time"
)

type User struct {
	ID           bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username     string        `bson:"username" json:"username"`
	PasswordHash string        `bson:"password_hash" json:"password_hash"`
}

type CovidSummary struct {
	Data []CovidByProvince
}

type CovidByProvince struct {
	Province string
	Count    int
	LastDate time.Time
}

type CovidData struct {
	Data []CovidCase `json:"Data"`
}

type CovidCase struct {
	ConfirmDate    CustomTime `json:"ConfirmDate"`
	No             string     `json:"No"`
	Age            float32    `json:"Age"`
	Gender         string     `json:"Gender"`
	GenderEn       string     `json:"GenderEn"`
	Nation         string     `json:"Nation"`
	NationEn       string     `json:"NationEn"`
	Province       string     `json:"Province"`
	ProvinceId     int        `json:"ProvinceId"`
	District       string     `json:"District"`
	ProvinceEn     string     `json:"ProvinceEn"`
	Detail         string     `json:"Detail"`
	StatQuarantine int        `json:"StatQuarantine"`
}

type CustomTime struct {
	time.Time
}

const ctLayout = "2006-01-02 15:04:05"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(ctLayout, s)
	return
}
