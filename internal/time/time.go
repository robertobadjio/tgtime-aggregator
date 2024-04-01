package time

import (
	"database/sql"
	"tgtime-aggregator/internal/db"
	"time"
)

// AggregateDayTotalTime Подсчет общего количества секунд
func AggregateDayTotalTime(times []*TimeUser) int64 {
	var sum int64
	for i, t := range times {
		if i == 0 {
			continue
		}
		delta := t.Seconds - times[i-1].Seconds
		// Не учитываем перерывы меньше 15 минут
		if delta <= 15*60 {
			sum += delta
		}
	}

	return sum
}

func GetAllBreaksByTimesOld(times []*TimeUser) []*Break {
	breaks := make([]*Break, 0)
	for i, time := range times {
		if i == 0 {
			continue
		}

		breakStruct := new(Break)

		delta := time.Seconds - times[i-1].Seconds
		if delta <= 33 {
			continue
		} else if delta <= (10 * 60) { // TODO: в параметры
			continue
		} else {
			breakStruct.BeginTime = times[i-1].Seconds
			breakStruct.EndTime = time.Seconds
			breaks = append(breaks, breakStruct)
		}
	}

	return breaks
}

func GetAllByDate(macAddress string, date time.Time, routerId int) []*TimeUser {
	var args []interface{}
	args = append(args, macAddress)
	args = append(args, GetSecondsByBeginDate(date.Format("2006-01-02")))
	args = append(args, GetSecondsByEndDate(date.Format("2006-01-02")))

	var rQuery string
	if routerId != 0 {
		rQuery = " AND t.router_id = $4"
		args = append(args, routerId)
	}

	rows, err := db.GetDB().Query("SELECT t.mac_address, t.second FROM time t WHERE t.mac_address = $1 AND t.second BETWEEN $2 AND $3"+rQuery+" ORDER BY t.second", args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	times := make([]*TimeUser, 0)
	for rows.Next() {
		time := new(TimeUser)
		err = rows.Scan(&time.MacAddress, &time.Seconds)
		if err != nil {
			panic(err)
		}

		times = append(times, time)
	}

	return times
}

func GetSecondsByBeginDate(date string) int64 {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	t, _ := time.ParseInLocation("2006-01-02", date, moscowLocation)

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, moscowLocation).Unix()
}

func GetSecondsByEndDate(date string) int64 {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	t, _ := time.ParseInLocation("2006-01-02", date, moscowLocation)

	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()).Unix()
}

func GetDayTimeFromTimeTable(macAddress string, date time.Time, sort string) int64 {
	var beginSecond int64
	row := db.GetDB().QueryRow("SELECT t.second FROM time t WHERE t.mac_address = $1 AND t.second BETWEEN $2 AND $3 ORDER BY t.second "+sort+" LIMIT 1", macAddress, GetSecondsByBeginDate(date.Format("2006-01-02")), GetSecondsByEndDate(date.Format("2006-01-02")))
	err := row.Scan(&beginSecond)

	if err == sql.ErrNoRows {
		return 0
	}

	if err != nil {
		panic(err)
	}

	return beginSecond
}
