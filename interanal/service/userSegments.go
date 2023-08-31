package service

import (
	"context"
	"log"
	"time"
	"userSegments/interanal/controller/request"
	"userSegments/interanal/model"
	"userSegments/interanal/storage"
	"userSegments/pkg/helper"
)

type UserSegments interface {
	GetUserSegments(ctx context.Context, id int) ([]model.Segment, error)
	AddUserToSegment(ctx context.Context, req request.UserAddSegmentRequest) (int, error)
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

func (s *UserSegmentsService) AddUserToSegment(ctx context.Context, req request.UserAddSegmentRequest) (int, error) {
	userActiveSegments, err := s.userSegmentsStorage.GetUserSegments(ctx, req.UserId)
	log.Print("userActiveSegments", userActiveSegments)
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
		log.Print("toAd", toAdd)
		if len(toAdd) > 0 {
			for _, slug := range toAdd {
				segment, err := s.segmentsStorage.GetSegmentBySlug(ctx, slug) //todo запрос в цикле
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

	var segmentsToDel []string
	if len(segmentSlugs) > 0 && len(req.DeleteSegments) > 0 {
		segmentsToDel = helper.IntersectionSlices(req.DeleteSegments, segmentSlugs)
		log.Print("toDel", segmentsToDel)
	}

	return s.userSegmentsStorage.AddUserToSegment(ctx, userSegmentsToInsert, segmentsToDel, req.UserId)
}
