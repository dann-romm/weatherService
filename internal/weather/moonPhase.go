package weather

import (
	"math"
	"time"
)

func CalcMoonPhase(t time.Time) float64 {
	firstNewMoon := time.Date(2022, time.January, 2, 21, 36, 0, 0, time.UTC)
	const moonMonth float64 = 29.53058812
	days := t.Sub(firstNewMoon).Hours() / 24 / moonMonth
	return days - math.Floor(days)
}

func MoonPhase(t time.Time) string {
	moonPhase := CalcMoonPhase(t)
	if moonPhase < 0.05 {
		return "new moon"
	} else if moonPhase < 0.20 {
		return "waxing crescent"
	} else if moonPhase < 0.30 {
		return "first quarter moon"
	} else if moonPhase < 0.45 {
		return "waxing gibous"
	} else if moonPhase < 0.55 {
		return "full moon"
	} else if moonPhase < 0.70 {
		return "waning gibous"
	} else if moonPhase < 0.80 {
		return "last quarter moon"
	} else if moonPhase < 0.95 {
		return "waning crescent"
	} else {
		return "new moon"
	}
}
