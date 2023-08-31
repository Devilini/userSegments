package service

import (
	"context"
	"fmt"
	"userSegments/interanal/model"
	"userSegments/interanal/storage"
)

type Segment interface {
	GetSegmentById(ctx context.Context, id int) (model.Segment, error)
	CreateSegment(ctx context.Context, slug string) (int, error)
	DeleteSegmentBySlug(ctx context.Context, slug string) error
}

type SegmentService struct {
	segmentStorage      storage.Segment
	userSegmentsStorage storage.UserSegments
}

func NewSegmentService(storage storage.Segment, userSegmentsStorage storage.UserSegments) *SegmentService {
	return &SegmentService{segmentStorage: storage, userSegmentsStorage: userSegmentsStorage}
}

func (s *SegmentService) GetSegmentById(ctx context.Context, id int) (model.Segment, error) {
	segment, err := s.segmentStorage.GetSegmentById(ctx, id)
	if segment.Id == 0 {
		return segment, fmt.Errorf("segment does not exists") // todo error
	}
	return segment, err
}

func (s *SegmentService) DeleteSegmentBySlug(ctx context.Context, slug string) error {
	segment, _ := s.segmentStorage.GetSegmentBySlug(ctx, slug)
	if segment.Id == 0 {
		return fmt.Errorf("segment does not exists")
	}

	_, err := s.segmentStorage.DeleteSegmentBySlug(ctx, slug)
	if err != nil {
		return err
	}

	err = s.userSegmentsStorage.DeleteSlugForUsers(ctx, segment.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *SegmentService) CreateSegment(ctx context.Context, slug string) (int, error) {
	segment, _ := s.segmentStorage.GetSegmentBySlug(ctx, slug)
	if segment.Id != 0 {
		return 0, fmt.Errorf("segment already exists") // todo error
	}

	return s.segmentStorage.CreateSegment(ctx, slug)
}
