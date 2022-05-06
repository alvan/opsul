package model

import "time"

type Note struct {
	Name string    `json:"name"`
	Text string    `json:"text"`
	Time time.Time `json:"time"`
}
