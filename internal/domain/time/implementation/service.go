package implementation

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	time2 "github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
)

// TimeService Сервис для работы с временем сотрудника
type TimeService struct {
	repository time2.Repository
	logger     log.Logger
}

// NewTimeService Конструктор сервиса
func NewTimeService(rep time2.Repository, logger log.Logger) *TimeService {
	return &TimeService{
		repository: rep,
		logger:     logger,
	}
}

// CreateTime Добавить время пребывания сотрудника на работе / в офисе
func (s *TimeService) CreateTime(ctx context.Context, t *time2.Time) error {
	if err := s.repository.CreateTime(ctx, t); err != nil {
		_ = s.logger.Log("msg", err.Error())
		return err // TODO: !
	}

	return nil
}

// GetByFilters ???
func (s *TimeService) GetByFilters(
	ctx context.Context,
	macAddress string,
	date time.Time,
	routerID int,
) ([]*time2.Time, error) {
	secondsStart := getSecondsByBeginDate(date.Format("2006-01-02"))
	secondsEnd := getSecondsByEndDate(date.Format("2006-01-02"))
	q := time2.Query{MacAddress: macAddress, SecondsStart: secondsStart, SecondsEnd: secondsEnd, RouterID: routerID}
	users, err := s.repository.GetByFilters(ctx, q)
	if err != nil {
		_ = s.logger.Log("msg", err.Error())
		return []*time2.Time{}, err
	}

	return users, nil
}

// GetStartSecondDayByDate ???
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
		_ = s.logger.Log("msg", err.Error())
		return seconds, err
	}

	return seconds, nil
}

// GetEndSecondDayByDate ???
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
		_ = s.logger.Log("msg", err.Error())
		return seconds, err
	}

	return seconds, nil
}

// AggregateDayTotalTime Подсчет общего количества секунд
func (s *TimeService) AggregateDayTotalTime(times []*time2.Time) (int64, error) {
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

// GetAllBreaksByTimes Подсчет перерывов в работе
func (s *TimeService) GetAllBreaksByTimes(times []*time2.Time) ([]*time_summary.Break, error) {
	breaks := make([]*time_summary.Break, 0) // TODO: len
	for i, t := range times {
		if i == 0 {
			continue
		}

		breakStruct := new(time_summary.Break)

		delta := t.Seconds - times[i-1].Seconds
		if delta <= 33 {
			continue
		}

		if delta <= (10 * 60) { // TODO: в параметры
			continue
		}

		breakStruct.SecondsStart = times[i-1].Seconds
		breakStruct.SecondsEnd = t.Seconds
		breaks = append(breaks, breakStruct)
	}

	return breaks, nil
}

// GetMacAddresses Получение mac-адресов сотрудников присутствовавших в офисе за период
func (s *TimeService) GetMacAddresses(ctx context.Context, date time.Time) ([]string, error) {
	secondsStart := getSecondsByBeginDate(date.Format("2006-01-02"))
	secondsEnd := getSecondsByEndDate(date.Format("2006-01-02"))
	q := time2.Query{SecondsStart: secondsStart, SecondsEnd: secondsEnd}
	macAddresses, err := s.repository.GetMacAddresses(ctx, q)
	if err != nil {
		_ = s.logger.Log("msg", err.Error())
		return []string{}, err
	}

	return macAddresses, nil
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
