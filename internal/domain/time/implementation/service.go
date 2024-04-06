package implementation

import (
	"context"
	"github.com/go-kit/kit/log"
	time2 "tgtime-aggregator/internal/domain/time"
	"time"
)

type TimeService struct {
	repository time2.Repository
	logger     log.Logger
}

func NewTimeService(rep time2.Repository, logger log.Logger) *TimeService {
	return &TimeService{
		repository: rep,
		logger:     logger,
	}
}

func (s *TimeService) CreateTime(ctx context.Context, t *time2.TimeUser) error {
	if err := s.repository.CreateTime(ctx, t); err != nil {
		s.logger.Log("msg", err.Error())
		return err // TODO: !
	}

	return nil
}

func (s *TimeService) GetByFilters(
	ctx context.Context,
	macAddress string,
	date time.Time,
	routerId int,
) ([]*time2.TimeUser, error) {
	secondsStart := getSecondsByBeginDate(date.Format("2006-01-02"))
	secondsEnd := getSecondsByEndDate(date.Format("2006-01-02"))
	q := time2.Query{MacAddress: macAddress, SecondsStart: secondsStart, SecondsEnd: secondsEnd, RouterId: routerId}
	users, err := s.repository.GetByFilters(ctx, q)
	if err != nil {
		s.logger.Log("msg", err.Error())
		return []*time2.TimeUser{}, err
	}

	return users, nil
}

func (s *TimeService) GetStartSecondDayByDate(
	ctx context.Context,
	macAddress string,
	date time.Time,
) (int64, error) {
	secondsStart := getSecondsByBeginDate(date.Format("2006-01-02"))
	secondsEnd := getSecondsByEndDate(date.Format("2006-01-02"))
	q := time2.Query{MacAddress: macAddress, SecondsStart: secondsStart, SecondsEnd: secondsEnd}

	seconds, err := s.repository.GetSecondsDayByDate(ctx, q, "ASC")
	if err != nil {
		s.logger.Log("msg", err.Error())
		return seconds, err
	}

	return seconds, nil
}

func (s *TimeService) GetEndSecondDayByDate(
	ctx context.Context,
	macAddress string,
	date time.Time,
) (int64, error) {
	secondsStart := getSecondsByBeginDate(date.Format("2006-01-02"))
	secondsEnd := getSecondsByEndDate(date.Format("2006-01-02"))
	q := time2.Query{MacAddress: macAddress, SecondsStart: secondsStart, SecondsEnd: secondsEnd}

	seconds, err := s.repository.GetSecondsDayByDate(ctx, q, "DESC")
	if err != nil {
		s.logger.Log("msg", err.Error())
		return seconds, err
	}

	return seconds, nil
}

// AggregateDayTotalTime Подсчет общего количества секунд
func (s *TimeService) AggregateDayTotalTime(times []*time2.TimeUser) (int64, error) {
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

	return sum, nil
}

func (s *TimeService) GetAllBreaksByTimesOld(times []*time2.TimeUser) ([]*time2.Break, error) {
	breaks := make([]*time2.Break, 0)
	for i, t := range times {
		if i == 0 {
			continue
		}

		breakStruct := new(time2.Break)

		delta := t.Seconds - times[i-1].Seconds
		if delta <= 33 {
			continue
		} else if delta <= (10 * 60) { // TODO: в параметры
			continue
		} else {
			breakStruct.BeginTime = times[i-1].Seconds
			breakStruct.EndTime = t.Seconds
			breaks = append(breaks, breakStruct)
		}
	}

	return breaks, nil
}

func getSecondsByBeginDate(date string) int64 {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	t, _ := time.ParseInLocation("2006-01-02", date, moscowLocation)

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, moscowLocation).Unix()
}

func getSecondsByEndDate(date string) int64 {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	t, _ := time.ParseInLocation("2006-01-02", date, moscowLocation)

	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()).Unix()
}
