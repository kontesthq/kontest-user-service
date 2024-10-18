package model

import "time"

type Site struct {
	isSiteEnabled                          bool
	isAutomaticCalendarNotificationEnabled bool
	secondsBeforeWhichNotificationToSet    []int
}

type User struct {
	FirstName         string
	LastName          string
	AccountCreateDate time.Time

	LeetcodeUsername   string
	CodechefUsername   string
	CodeforcesUsername string

	sites []Site

	minDurationInSecond int
	maxDurationInSecond int
}
