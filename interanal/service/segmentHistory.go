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
		fmt.Println(err)
		return "", err
	}
	fmt.Println(history)
	fileName := time.Now().Format("20060102150405") + ".csv"
	path := fmt.Sprintf("public/reports/segments_%s", fileName)
	// creating csv file
	csvFile, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	// creating writer object for csv file
	csvWriter := csv.NewWriter(csvFile)
	// Specify the delimiter
	csvWriter.Comma = ','

	fmt.Println(history)

	// Adding all the students in CSV file using csv writer object
	row := []string{"UserId", "Segment", "Operation", "Date"}
	err = csvWriter.Write(row)
	if err != nil {
		fmt.Println("Error:", err)
		return "", nil
	}
	for _, item := range history {
		row := []string{item.UserId, item.Segment, item.Operation, item.CreatedAt.Format("2006-01-02 15:04:05")}
		err := csvWriter.Write(row)
		if err != nil {
			fmt.Println("Error:", err)
			return "", nil
		}
	}
	csvWriter.Flush()

	return path, nil
}
