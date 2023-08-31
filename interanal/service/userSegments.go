package service

import (
	"context"
	"time"
	"userSegments/interanal/controller/request"
	"userSegments/interanal/model"
	"userSegments/interanal/storage"
	"userSegments/pkg/helper"
)

type UserSegments interface {
	GetUserSegments(ctx context.Context, id int) ([]model.Segment, error)
	ChangeUserSegments(ctx context.Context, req request.UserAddSegmentRequest) (int, error)
}

type UserSegmentsService struct {
	userSegmentsStorage storage.UserSegments
	segmentsStorage     storage.SegmentStorage
}

func NewUserSegmentsService(userSegmentsStorage storage.UserSegments, segmentsStorage storage.SegmentStorage) *UserSegmentsService {
	return &UserSegmentsService{userSegmentsStorage: userSegmentsStorage, segmentsStorage: segmentsStorage}
}

func (s *UserSegmentsService) GetUserSegments(ctx context.Context, id int) ([]model.Segment, error) {
	return s.userSegmentsStorage.GetUserSegments(ctx, id)
}

func (s *UserSegmentsService) ChangeUserSegments(ctx context.Context, req request.UserAddSegmentRequest) (int, error) {
	userActiveSegments, err := s.userSegmentsStorage.GetUserSegments(ctx, req.UserId)
	if err != nil {
		return 0, err
	}
	var segmentSlugs []string
	for _, segment := range userActiveSegments {
		segmentSlugs = append(segmentSlugs, segment.Slug)
	}

	var userSegmentsToInsert []model.UserSegments
	if len(req.AddSegments) > 0 {
		toAdd := helper.DifferenceSlices(req.AddSegments, segmentSlugs)
		if len(toAdd) > 0 {
			for _, slug := range toAdd {
				segment, err := s.segmentsStorage.GetSegmentBySlug(ctx, slug)
				if segment.Id == 0 {
					continue
				}
				if err != nil {
					return 0, err
				}
				userSegmentsToInsert = append(userSegmentsToInsert, model.UserSegments{
					UserId:    req.UserId,
					SegmentId: segment.Id,
					CreatedAt: time.Now(),
				})
			}
		}
	}

	var segmentIds []int
	if len(segmentSlugs) > 0 && len(req.DeleteSegments) > 0 {
		segmentsToDel := helper.IntersectionSlices(req.DeleteSegments, segmentSlugs)
		existingsSegmentsForDel, err := s.userSegmentsStorage.GetSegmentsBySlugs(ctx, segmentsToDel, req.UserId)
		if err != nil {
			return 0, err
		}

		for _, segment := range existingsSegmentsForDel {
			segmentIds = append(segmentIds, segment.Id)
		}
	}

	return s.userSegmentsStorage.ChangeUserSegments(ctx, userSegmentsToInsert, segmentIds, req.UserId)
}
