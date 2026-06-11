package repository

import (
	"ecp_app/internal/model"
)

type RepositoryInterface interface {
	Add(ecp model.Ecp) error
	FindExpiring() ([]model.Ecp, error)
}
