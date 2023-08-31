package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"
	"userSegments/interanal/storage"
)

type SegmentHistory interface {
	GetSegmentsReport(ctx context.Context, date string) (string, error)
}

type SegmentHistoryService struct {
	segmentHistoryStorage storage.SegmentHistory
}

func NewSegmentsHistoryService(segmentHistoryStorage storage.SegmentHistory) *SegmentHistoryService {
	return &SegmentHistoryService{segmentHistoryStorage: segmentHistoryStorage}
}

func (s *SegmentHistoryService) GetSegmentsReport(ctx context.Context, date string) (string, error) {
	history, err := s.segmentHistoryStorage.GetSegmentsHistory(ctx, date)
	if err != nil {
		return "", err
	}
	fileName := time.Now().Format("20060102150405") + ".csv"
	path := fmt.Sprintf("public/reports/segments_%s", fileName)
	csvFile, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	csvWriter.Comma = ','

	cols := []string{"UserId", "Segment", "Operation", "Date"}
	err = csvWriter.Write(cols)
	if err != nil {
		return "", err
	}
	for _, item := range history {
		row := []string{item.UserId, item.Segment, item.Operation, item.CreatedAt.Format("2006-01-02 15:04:05")}
		err := csvWriter.Write(row)
		if err != nil {
			return "", err
		}
	}
	csvWriter.Flush()

	return fileName, nil
}
