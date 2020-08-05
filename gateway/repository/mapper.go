package repository

import "time"

type MapperList struct {
}

type MapperVo struct {
	ID      int64
	Deleted *bool
	Created time.Time
	Updated time.Time
	URL     string
	Type    int16
	PSM     string
	Domain  string
	Base    string
	Method  string
}
