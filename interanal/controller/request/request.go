package request

import _type "userSegments/interanal/model/type"

type UserAddSegmentRequest struct {
	UserId         int             `json:"userId" validate:"required"`
	AddSegments    []string        `json:"addSegments" validate:"required"`
	DeleteSegments []string        `json:"deleteSegments" validate:"required"`
	ExpiredDate    *_type.DateTime `json:"expiredDate,omitempty"`
}

type UserCreateRequest struct {
	Name string `json:"name" validate:"required"`
}

type SegmentCreateRequest struct {
	Slug    string `json:"slug" validate:"required"`
	Percent *int   `json:"percent,omitempty" validate:"omitempty,min=1,max=100"`
}
