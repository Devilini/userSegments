package request

type UserAddSegmentRequest struct {
	UserId         int      `json:"userId" validate:"required"`
	AddSegments    []string `json:"addSegments" validate:"required"`
	DeleteSegments []string `json:"deleteSegments" validate:"required"`
}

type UserCreateRequest struct {
	Name string `json:"name" validate:"required"`
}

type SegmentCreateRequest struct {
	Slug string `json:"slug" validate:"required"`
}
