package service

import (
	"context"
	"fmt"

	"git.codenrock.com/cnrprod1725727333-user-88349/zadanie-6105/internal/tender/domain"
)

type pingProvider struct {
	repo domain.TenderRepositoryPingProvider
}

func NewPingProvider(repo domain.TenderRepositoryPingProvider) *pingProvider {
	return &pingProvider{repo: repo}
}

func (s *pingProvider) Ping(ctx context.Context) error {
	err := s.repo.Ping(ctx)
	if err != nil {
		return fmt.Errorf("service.Ping: %w", err)
	}

	return nil
}

type Tender struct {
	repo domain.TenderRepository
}

func NewTender(repo domain.TenderRepository) *Tender {
	return &Tender{repo: repo}
}

func (t *Tender) ListTender(ctx context.Context, limit, offset int, serviceTypes []string) ([]domain.Tender, error) {
	tenders, err := t.repo.ListTender(ctx, limit, offset, serviceTypes)
	if err != nil {
		return nil, fmt.Errorf("service.ListBanners: %w", err)
	}

	return tenders, nil
}

func (s *Tender) CreateTender(ctx context.Context, tender domain.Tender) (domain.Tender, error) {

	createdTender, err := s.repo.CreateTender(ctx, tender)
	if err != nil {
		return domain.Tender{}, fmt.Errorf("service.CreateTender: %w", err)
	}

	return createdTender, nil
}

func (s *Tender) GetUserTenders(ctx context.Context, limit, offset int, username string) ([]domain.Tender, error) {
	tenders, err := s.repo.GetUserTenders(ctx, limit, offset, username)
	if err != nil {
		return nil, fmt.Errorf("service.CreateTender: %w", err)
	}

	return tenders, nil
}

func (s *Tender) UpdateTenderStatus(ctx context.Context, tenderID string, status string, username string) (domain.Tender, error) {
	updateTender,err := s.repo.UpdateTenderStatus(ctx, tenderID, status, username)
	if err != nil {
		return domain.Tender{},fmt.Errorf("service.UpdateTenderStatus: %w", err)
	}

	return updateTender, nil
}

func (s *Tender) GetTenderStatus(ctx context.Context, tenderID string, username string) (string, error) {
	status, err := s.repo.GetTenderStatus(ctx, tenderID, username)
	if err != nil {
		return "", fmt.Errorf("service.GetTenderByID: %w", err)
	}

	return status, nil
}

func (s *Tender) UpdatePartTender(ctx context.Context, id string, updates map[string]interface{}, username string) (domain.Tender, error) {
	updatedTender, err := s.repo.UpdatePartTender(ctx, id, updates, username)
	if err != nil {
		return domain.Tender{}, fmt.Errorf("failed to update tender in repository: %w", err)
	}

	return updatedTender, nil
}

func (s *Tender) RollbackTenderVersion(ctx context.Context, id string, version int, username string) (domain.Tender, error) {
	tender, err := s.repo.RollbackTenderVersion(ctx, id, version, username)
	if err != nil {
		return domain.Tender{}, fmt.Errorf("service.RollbackTenderVersion: %w", err)
	}
	return tender, nil
}
