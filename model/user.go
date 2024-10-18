package model

import (
	"github.com/google/uuid"
	"time"
)

type Site struct {
	IsSiteEnabled                          bool  `json:"is_site_enabled"`
	IsAutomaticCalendarNotificationEnabled bool  `json:"is_automatic_calendar_notification_enabled"`
	SecondsBeforeWhichNotificationToSet    []int `json:"seconds_before_which_notification_to_set"`
}

type User struct {
	ID                uuid.UUID `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	AccountCreateDate time.Time `json:"account_create_date"`

	LeetcodeUsername   string `json:"leetcode_username"`
	CodechefUsername   string `json:"codechef_username"`
	CodeforcesUsername string `json:"codeforces_username"`

	Sites []Site `json:"sites"`

	MinDurationInSecond int `json:"min_duration_in_second"`
	MaxDurationInSecond int `json:"max_duration_in_second"`
}
