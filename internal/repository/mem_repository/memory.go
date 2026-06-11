package memrepository

import (
	"ecp_app/internal/model"

	"github.com/google/uuid"
)

type MemRepository struct {
	vault map[uuid.UUID]model.Ecp
}

func NewMemRepository() (*MemRepository, error) {
	return &MemRepository{vault: make(map[uuid.UUID]model.Ecp)}, nil
}

func (r *MemRepository) Add(ecp model.Ecp) error {
	r.vault[ecp.ID] = ecp

	return nil
}

func (r *MemRepository) FindExpiring() ([]model.Ecp, error) {
	ecps := make([]model.Ecp, 0)

	for _, ecp := range r.vault {
		if ecp.IsExpiringNextMonth() {
			ecps = append(ecps, ecp)
		}
	}

	return ecps, nil
}
