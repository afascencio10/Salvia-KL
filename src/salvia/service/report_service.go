package service

import "bitsflow/internal/repository"

type ReportService interface{}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}
