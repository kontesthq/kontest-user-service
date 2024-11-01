package model

import (
	"github.com/google/uuid"
	"time"
)

type Site struct {
	SiteName                               string `json:"site_name" db:"site_name"`
	IsSiteEnabled                          bool   `json:"is_site_enabled" db:"is_site_enabled"`
	IsAutomaticCalendarNotificationEnabled bool   `json:"is_automatic_calendar_notification_enabled" db:"is_automatic_calendar_notification_enabled"`
	SecondsBeforeWhichAppNotificationToSet []int  `json:"seconds_before_which_app_notification_to_set" db:"seconds_before_which_app_notification_to_set"`
}

type User struct {
	ID                uuid.UUID `json:"id" db:"id"`
	FirstName         string    `json:"first_name" db:"first_name"`
	LastName          string    `json:"last_name" db:"last_name"`
	AccountCreateDate time.Time `json:"account_create_date" db:"account_create_date"`

	LeetcodeUsername   string `json:"leetcode_username" db:"leetcode_username"`
	CodechefUsername   string `json:"codechef_username" db:"codechef_username"`
	CodeforcesUsername string `json:"codeforces_username" db:"codeforces_username"`

	Sites []Site `json:"sites" db:"sites"`

	MinDurationInSecond int `json:"min_duration_in_second" db:"min_duration_in_seconds"`
	MaxDurationInSecond int `json:"max_duration_in_second" db:"max_duration_in_seconds"`

	CollegeName  string `json:"college_name" db:"college_name"`
	CollegeState string `json:"college_state" db:"college_state"`
}
