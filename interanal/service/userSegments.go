package service

import (
	"context"
	"userSegments/interanal/controller/request"
	"userSegments/interanal/model"
	"userSegments/interanal/storage"
)

type UserSegments interface {
	GetUserSegments(ctx context.Context, id int) ([]model.Segment, error)
	ChangeUserSegments(ctx context.Context, req request.UserAddSegmentRequest) (int, error)
	DeleteAllExpired(ctx context.Context) (int64, error)
}

type UserSegmentsService struct {
	userSegmentsStorage   storage.UserSegments
	segmentsStorage       storage.SegmentStorage
	segmentHistoryStorage storage.SegmentHistory
}

func NewUserSegmentsService(userSegmentsStorage storage.UserSegments, segmentsStorage storage.SegmentStorage, segmentHistoryStorage storage.SegmentHistory) *UserSegmentsService {
	return &UserSegmentsService{userSegmentsStorage: userSegmentsStorage, segmentsStorage: segmentsStorage, segmentHistoryStorage: segmentHistoryStorage}
}

func (s *UserSegmentsService) GetUserSegments(ctx context.Context, id int) ([]model.Segment, error) {
	return s.userSegmentsStorage.GetUserSegments(ctx, id)
}

func (s *UserSegmentsService) ChangeUserSegments(ctx context.Context, req request.UserAddSegmentRequest) (int, error) {
	var userSegmentsToInsert []model.UserSegments
	if len(req.AddSegments) > 0 {
		segments, err := s.segmentsStorage.GetSegmentsBySlug(ctx, req.AddSegments)
		if err != nil {
			return 0, err
		}

		for _, segment := range segments {
			userSegmentsToInsert = append(userSegmentsToInsert, model.UserSegments{
				UserId:    req.UserId,
				SegmentId: segment.Id,
				ExpiredAt: req.ExpiredDate,
			})
		}
	}

	var segmentDelIds []int
	if len(req.DeleteSegments) > 0 {
		userSegmentsToDelete, err := s.segmentsStorage.GetSegmentsBySlug(ctx, req.DeleteSegments)
		if err != nil {
			return 0, err
		}

		for _, segment := range userSegmentsToDelete {
			segmentDelIds = append(segmentDelIds, segment.Id)
		}
	}

	return s.userSegmentsStorage.ChangeUserSegments(ctx, userSegmentsToInsert, segmentDelIds, req.UserId)
}

func (s *UserSegmentsService) DeleteAllExpired(ctx context.Context) (int64, error) {
	var countRows int64 = 0
	deletedSegments, err := s.userSegmentsStorage.DeleteAllExpired(ctx)
	if err != nil {
		return countRows, err
	}

	if len(deletedSegments) > 0 {
		var entries []model.SegmentHistory
		for _, item := range deletedSegments {
			entries = append(entries, model.SegmentHistory{UserId: item.UserId, SegmentId: item.SegmentId, Operation: storage.OperationTypeDelete})
		}

		countRows, err = s.segmentHistoryStorage.BulkInsert(ctx, entries)
		if err != nil {
			return countRows, err
		}
	}

	return countRows, nil
}
