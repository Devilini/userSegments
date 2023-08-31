package request

type UserAddSegmentRequest struct {
	UserId         int      `json:"userId" validate:"required"`
	AddSegments    []string `json:"addSegments" validate:"required"`
	DeleteSegments []string `json:"deleteSegments" validate:"required"`
}
