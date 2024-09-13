package domain

import (
	"context"
)

type TenderServicePingProvider interface {
	Ping(ctx context.Context) error
}

type TenderRepositoryPingProvider interface {
	Ping(ctx context.Context) error
}

type TenderService interface {
	ListTender(ctx context.Context, limit int, offset int, serviceTypes []string) ([]Tender, error)
	CreateTender(ctx context.Context, tender Tender) (Tender, error)
	GetUserTenders(ctx context.Context, limit int, offset int, username string) ([]Tender, error)
	UpdateTenderStatus(ctx context.Context, tenderID string, status string, username string) (Tender, error)
	GetTenderStatus(ctx context.Context, tenderID string, username string) (string, error)
	UpdatePartTender(ctx context.Context, id string, updates map[string]interface{}, username string) (Tender, error)
	RollbackTenderVersion(ctx context.Context, tenderID string, version int, username string) (Tender, error)
}

//go:generate mockgen -destination=mocks/repo_mock.gen.go -package=mocks . TenderRepositoryGetter
type TenderRepository interface {
	ListTender(ctx context.Context, limit int, offset int, serviceTypes []string) ([]Tender, error)
	CreateTender(ctx context.Context, tender Tender) (Tender, error)
	GetUserTenders(ctx context.Context, limit int, offset int, username string) ([]Tender, error)
	UpdateTenderStatus(ctx context.Context, tenderID string, status string, username string) (Tender, error)
	GetTenderStatus(ctx context.Context, tenderID string, username string) (string, error)
	UpdatePartTender(ctx context.Context, id string, updates map[string]interface{}, username string) (Tender, error)
	RollbackTenderVersion(ctx context.Context, tenderID string, version int, username string) (Tender, error)
}
