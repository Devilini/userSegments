package service

import (
	"context"
	"fmt"
	"userSegments/interanal/controller/request"
	"userSegments/interanal/model"
	"userSegments/interanal/storage"
)

type Segment interface {
	GetSegmentById(ctx context.Context, id int) (model.Segment, error)
	CreateSegment(ctx context.Context, req request.SegmentCreateRequest) (int, error)
	DeleteSegmentBySlug(ctx context.Context, slug string) error
}

type SegmentService struct {
	segmentStorage        storage.Segment
	userSegmentsStorage   storage.UserSegments
	segmentHistoryStorage storage.SegmentHistory
}

func NewSegmentService(storage storage.Segment, userSegmentsStorage storage.UserSegments, segmentHistoryStorage storage.SegmentHistory) *SegmentService {
	return &SegmentService{segmentStorage: storage, userSegmentsStorage: userSegmentsStorage, segmentHistoryStorage: segmentHistoryStorage}
}

func (s *SegmentService) GetSegmentById(ctx context.Context, id int) (model.Segment, error) {
	segment, err := s.segmentStorage.GetSegmentById(ctx, id)
	if err != nil {
		return segment, err
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

	userIds, err := s.userSegmentsStorage.DeleteSlugForUsers(ctx, segment.Id)
	if err != nil {
		return err
	}

	var entries []model.SegmentHistory
	for _, userId := range userIds {
		entries = append(entries, model.SegmentHistory{UserId: userId, SegmentId: segment.Id, Operation: storage.OperationTypeDelete})
	}

	_, err = s.segmentHistoryStorage.BulkInsert(ctx, entries)
	if err != nil {
		return err
	}

	return nil
}

func (s *SegmentService) CreateSegment(ctx context.Context, req request.SegmentCreateRequest) (int, error) {
	segment, _ := s.segmentStorage.GetSegmentBySlug(ctx, req.Slug)
	if segment.Id != 0 {
		return 0, fmt.Errorf("segment already exists")
	}

	segmentId, err := s.segmentStorage.CreateSegment(ctx, req.Slug, req.Percent)
	if err != nil {
		return 0, err
	}

	if req.Percent != nil {
		err := s.userSegmentsStorage.CreateByPercent(ctx, *req.Percent, segmentId)
		if err != nil {
			return 0, err
		}
	}

	return segmentId, nil
}
