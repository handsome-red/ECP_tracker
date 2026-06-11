package service

import (
	"ecp_app/internal/model"
	"ecp_app/internal/repository"
)

type EcpServiceInterface interface {
	Add(ecp model.Ecp) error
	FindExpiring() ([]model.Ecp, error)
}

type EcpService struct {
	repo repository.RepositoryInterface
}

func NewEcpService(repo repository.RepositoryInterface) *EcpService {
	return &EcpService{repo: repo}
}

func (s *EcpService) Add(ecp model.Ecp) error {
	return s.repo.Add(ecp)
}

func (s *EcpService) FindExpiring() ([]model.Ecp, error) {
	return s.repo.FindExpiring()
}
