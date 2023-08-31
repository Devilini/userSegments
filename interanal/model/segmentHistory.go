package model

import "time"

type SegmentHistory struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	SegmentId int       `json:"segment_id"`
	Operation string    `json:"operation"`
	CreatedAt time.Time `json:"created_at"`
}

type SegmentHistoryReport struct {
	Id        int       `json:"id"`
	UserId    string    `json:"user_id"`
	Segment   string    `json:"segment_id"`
	Operation string    `json:"operation"`
	CreatedAt time.Time `json:"created_at"`
}
