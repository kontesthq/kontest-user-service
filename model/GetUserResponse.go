package model

type GetUserResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Email string `json:"email"`

	LeetcodeUsername   string `json:"leetcode_username"`
	CodechefUsername   string `json:"codechef_username"`
	CodeforcesUsername string `json:"codeforces_username"`

	Sites []Site `json:"sites"`

	MinDurationInSecond int `json:"min_duration_in_second"`
	MaxDurationInSecond int `json:"max_duration_in_second"`

	CollegeName  string `json:"college_name"`
	CollegeState string `json:"college_state"`
}
