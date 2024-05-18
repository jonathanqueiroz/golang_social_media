package models

type Profile struct {
	User           User   `json:"user"`
	Posts          []Post `json:"posts"`
	FollowersCount int    `json:"followers_count"`
	FollowingCount int    `json:"following_count"`
}
