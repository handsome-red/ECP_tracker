package model

import (
	"ecp_app/internal/constants"
	"time"

	"github.com/google/uuid"
)

type Ecp struct {
	ID         uuid.UUID
	StartDate  time.Time
	EndDate    time.Time
	Department constants.Department
	FullName   string
}

func NewEcp(fullName string, department constants.Department, startDate, endDate time.Time) *Ecp {
	return &Ecp{
		ID:         uuid.New(),
		StartDate:  startDate,
		EndDate:    endDate,
		Department: department,
		FullName:   fullName,
	}
}

func (e *Ecp) IsExpiringNextMonth() bool {
	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)

	firstDayNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDayNextMonth := firstDayNextMonth.AddDate(0, 1, -1)

	return !e.EndDate.Before(firstDayNextMonth) && !e.EndDate.After(lastDayNextMonth)
}
