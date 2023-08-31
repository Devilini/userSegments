package model

import "time"

type UserSegments struct {
	UserId    int        `json:"user_id"`
	SegmentId int        `json:"segment_id"`
	Percent   *int       `json:"percent,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
}
