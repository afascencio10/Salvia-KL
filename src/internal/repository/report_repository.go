package repository

import "gorm.io/gorm"

type ReportRepository interface{}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}
