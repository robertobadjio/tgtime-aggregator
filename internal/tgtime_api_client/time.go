package tgtime_api_client

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Time struct {
	MacAddress string
	Second     int64
	RouterId   int
}

type TimesPeriodResponse struct {
	Period int            `json:"period"`
	Times  []TimeResponse `json:"time"`
}

type TimeResponse struct {
	Date      string          `json:"date"`
	Total     int             `json:"total"`
	BeginTime int64           `json:"beginTime"`
	EndTime   int64           `json:"endTime"`
	Breaks    []BreakResponse `json:"breaks,omitempty"`
}

type BreakResponse struct {
	BeginTime int64 `json:"beginTime"`
	EndTime   int64 `json:"endTime"`
}

type Period struct {
	Id        int
	Name      string
	Year      int
	BeginDate string
	EndDate   string
}

type StatByWorkingPeriod struct {
	StartWorkingDate  string `json:"start_working_date"`
	EndWorkingDate    string `json:"end_working_date"`
	WorkingHours      int    `json:"working_hours"`
	TotalWorkingHours int    `json:"total_working_hours"`
}

type Periods struct {
	Periods []Period `json:"periods"`
}

func (otc *tgTimeClient) InsertTime(macAddress string, seconds int64, routerId int) (*Routers, error) {
	timeStruct := new(Time)
	timeStruct.MacAddress = macAddress
	timeStruct.Second = seconds
	timeStruct.RouterId = routerId
	params, _ := json.Marshal(timeStruct)
	payload := strings.NewReader(string(params))

	request, err := http.NewRequest(http.MethodPost, otc.baseURL+"/time", payload)
	if err != nil {
		return nil, err
	}

	routers := Routers{}
	if err := otc.sendRequest(request, &routers); err != nil {
		return nil, err
	}

	return &routers, nil
}

func (otc *tgTimeClient) GetTimesByPeriod(userId, periodId int) (*TimesPeriodResponse, error) {
	request, err := http.NewRequest(http.MethodGet, otc.baseURL+"/time/"+strconv.Itoa(userId)+"/period/"+strconv.Itoa(periodId), nil)
	if err != nil {
		return nil, err
	}

	timePeriod := TimesPeriodResponse{}
	if err := otc.sendRequest(request, &timePeriod); err != nil {
		return nil, err
	}

	return &timePeriod, nil
}

func (otc *tgTimeClient) GetTimesByDate(userId int, date time.Time) (*TimeResponse, error) {
	request, err := http.NewRequest(http.MethodGet, otc.baseURL+"/time/"+strconv.Itoa(userId)+"/day/"+date.Format("2006-01-02"), nil)
	if err != nil {
		return nil, err
	}

	timeStruct := TimeResponse{}
	if err := otc.sendRequest(request, &timeStruct); err != nil {
		return nil, err
	}

	return &timeStruct, nil
}

func (otc *tgTimeClient) GetStatByWorkingPeriod(userId, periodId int) (*StatByWorkingPeriod, error) {
	request, err := http.NewRequest("GET", otc.baseURL+"/stat/working-period/"+strconv.Itoa(userId)+"/period/"+strconv.Itoa(periodId), nil)
	if err != nil {
		return nil, err
	}

	stat := StatByWorkingPeriod{}
	if err := otc.sendRequest(request, &stat); err != nil {
		return nil, err
	}

	return &stat, nil
}
