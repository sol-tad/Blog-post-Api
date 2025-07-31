package domian

import "time"

type Blog struct {
	Title   string
	Content string
	Tags    string
	Date    time.Time
	User 	User
}