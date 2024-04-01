package aggregator

import (
	"database/sql"
	"encoding/json"
	"tgtime-aggregator/internal/tgtime_api_client"
	timeModel "tgtime-aggregator/internal/time"
	"time"
)

var Db *sql.DB

func AggregateTime() {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	date := time.Now().AddDate(0, 0, -1).In(moscowLocation)

	apiClient := tgtime_api_client.NewTgTimeClient()
	users, _ := apiClient.GetAllUsers()

	for _, user := range users.Users {
		times := timeModel.GetAllByDate(user.MacAddress, date, 0)
		seconds := timeModel.AggregateDayTotalTime(times)
		breaks := timeModel.GetAllBreaksByTimesOld(times)
		breaksJson, err := json.Marshal(breaks)
		begin := timeModel.GetDayTimeFromTimeTable(user.MacAddress, date, "ASC")
		end := timeModel.GetDayTimeFromTimeTable(user.MacAddress, date, "DESC")

		_, err = Db.Exec("INSERT INTO time_summary (mac_address, date, seconds, breaks, seconds_begin, seconds_end) VALUES ($1, $2, $3, $4, $5, $6)", user.MacAddress, date, seconds, breaksJson, begin, end)
		if err != nil {
			panic(err)
		}
	}
}
