package domain

import "time"

type LinkVisit struct {
	Id        int64
	LinkId    int64
	Ip        string
	UserAgent string
	Referer   string
	Status    int
	CreatedAt time.Time
}

type LinkVisitVO struct {
	LinkId    int64
	Ip        string
	UserAgent string
	Referer   string
	Status    int
}
