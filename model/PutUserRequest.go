package model

type PutUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`

	LeetcodeUsername   *string `json:"leetcode_username,omitempty"`
	CodechefUsername   *string `json:"codechef_username,omitempty"`
	CodeforcesUsername *string `json:"codeforces_username,omitempty"`

	Sites []Site `json:"sites,omitempty"`

	MinDurationInSecond *int `json:"min_duration_in_second,omitempty"`
	MaxDurationInSecond *int `json:"max_duration_in_second,omitempty"`

	CollegeName  *string `json:"college_name,omitempty"`
	CollegeState *string `json:"college_state,omitempty"`
}
