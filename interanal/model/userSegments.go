package model

import (
	_type "userSegments/interanal/model/type"
)

type UserSegments struct {
	UserId    int             `json:"user_id"`
	SegmentId int             `json:"segment_id"`
	CreatedAt _type.DateTime  `json:"created_at"`
	ExpiredAt *_type.DateTime `json:"expired_at,omitempty"`
}
