package server

import (
	"time"

	"paas/icsoc-monitor/constants"
)

type Date struct {
	time time.Time
}

func (d Date) FormatString() string{
	return d.time.Format(constants.DATE_FORMATE)
}

func (d Date) IsSameDay(date Date) bool{
	return d.time.Year() == date.time.Year() &&
		d.time.Month() == date.time.Month() &&
			d.time.Day() == date.time.Day()
}

func (d Date) Sub(date Date) int {
	subYear := d.time.Year() - date.time.Year()
	if subYear != 0{
		return subYear
	}else {
		subMonth := d.time.Month() - date.time.Month()
		if subMonth != 0 {
			return int(subMonth)
		}else {
			return d.time.Day() - date.time.Day()
		}
	}
}

func (d Date) ChangeDay(day int) Date{
	return Date{d.time.AddDate(0, 0, day)}
}

func SplitDate(start, end, split Date) ([]Date, []Date){
	if start.Sub(split) <= 0 && end.Sub(split) >= 0 {
		if start.IsSameDay(split) && end.IsSameDay(split) {
			return nil, []Date{start, end}
		}else if start.IsSameDay(split) {
			return nil, []Date{start, end}
		}else if end.IsSameDay(split) {
			return []Date{start, end.ChangeDay(-1)}, []Date{end, end}
		}


	}else if start.Sub(split) > 0 {
		return nil, nil
	}else if end.Sub(split) < 0 {
		return nil, nil
	}
	return nil, nil
}
