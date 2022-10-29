package helper

import "time"

//go:generate mockery --name=Timer --output=./mocks
type Timer interface {
	NowInUTC() time.Time
}

type TimerImplementation struct{}

func (t *TimerImplementation) NowInUTC() time.Time {
	return time.Now().UTC()
}
