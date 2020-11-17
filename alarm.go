package alarm

import (
	"context"
	"fmt"
	"time"
)

// Alarm ...
type Alarm struct {
	A chan time.Time
}

// New возвращает структуру содержащую только одно поле A, которое является каналом и
// в который в заданое время отправится значение текущего времени. Значения в канал будет посылаться
// каждые сутки пока через контекст не будет вызвана отмена. Точность срабатывани до секунды.
// На вход функция принимает контекст, время срабатывания и не обязательный параметр
// временной зоны по которой необходимо срабатывание.
// Если временная зона не задана, по умолчанию принимается локальная зона.
// Время задается в формате 15:04:05, зона задается в формате IANA.
func New(ctx context.Context, timeAlarm time.Time, location ...string) (*Alarm, error) {

	timeWithDate, err := time.Parse("2006-01-02 15:04:05",
		fmt.Sprintf("%s %s",
			time.Now().Format("2006-01-02"),
			timeAlarm.Format("15:04:05")))
	if err != nil {
		return nil, err
	}

	l := time.Now().Local().Location()
	if len(location) > 0 {
		l, err = time.LoadLocation(location[0])
		if err != nil {
			return nil, err
		}
	}

	timeAlarmWithLoc, err := time.ParseInLocation("2006-01-02 15:04:05",
		timeWithDate.Format("2006-01-02 15:04:05"), l)
	if err != nil {
		return nil, err
	}
	local := time.Now().Local().Location()
	shedule := timeAlarmWithLoc.In(local)

	ta, err := time.ParseInLocation("15:04:05", shedule.Format("15:04:05"), local)
	if err != nil {
		return nil, err
	}

	a := &Alarm{A: make(chan time.Time, 1)}
	go func(ctx context.Context, tAlarm time.Time, sigAlarm chan time.Time) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				tnow, _ := time.ParseInLocation("15:04:05",
					time.Now().Format("15:04:05"), time.Now().Local().Location())
				if tnow.Equal(tAlarm) {
					sigAlarm <- time.Now()
				}
				time.Sleep(1 * time.Second)
			}
		}
	}(ctx, ta, a.A)

	return a, nil
}
