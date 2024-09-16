package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Te8va/Tender/internal/tender/domain"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

var (
	_ domain.TenderRepositoryPingProvider = (*pingProvider)(nil)
)

type pingProvider struct {
	db *postgres
}

func NewPingProvider(pg *postgres) *pingProvider {
	return &pingProvider{db: pg}
}

func (r *pingProvider) Ping(ctx context.Context) error {
	err := r.db.Ping(ctx)
	if err != nil {
		return fmt.Errorf("repository.Ping: %w", err)
	}

	return nil
}

var (
	_ domain.TenderRepository = (*TenderService)(nil)
)

type TenderService struct {
	pool *pgxpool.Pool
}

func NewTenderService(pool *pgxpool.Pool) *TenderService {
	return &TenderService{pool: pool}
}

func (t *TenderService) ListTender(ctx context.Context, limit, offset int, serviceTypes []string) ([]domain.Tender, error) {
	query := `SELECT id, name, description, service_type, status, version, created_at
              FROM tender`
	var args []interface{}
	argIndex := 1

	if len(serviceTypes) > 0 {
		query += ` WHERE service_type = ANY($` + strconv.Itoa(argIndex) + `)`
		args = append(args, pq.Array(serviceTypes))
		argIndex++
	}

	query += ` LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := t.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAllTenders: %w", err)
	}
	defer rows.Close()

	var tenders []domain.Tender
	for rows.Next() {
		var tender domain.Tender
		if err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType,
			&tender.Status, &tender.Version, &tender.CreatedAt); err != nil {
			return nil, fmt.Errorf("repository.GetAllTenders: %w", err)
		}
		tenders = append(tenders, tender)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository.GetAllTenders: %w", err)
	}

	return tenders, nil
}

func (r *TenderService) CreateTender(ctx context.Context, tender domain.Tender) (domain.Tender, error) {
	isAuthorized, err := r.IsUserAuthorizedForOrganization(ctx, tender.CreatorUsername, tender.OrganizationId)
	if err != nil {
		if err.Error() == "user does not exist" {
			return domain.Tender{}, fmt.Errorf("user does not exist")
		}
		return domain.Tender{}, fmt.Errorf("repository.CreateTender: %w", err)
	}

	if !isAuthorized {
		return domain.Tender{}, fmt.Errorf("user is not authorized to create tender for this organization")
	}

	query := `INSERT INTO tender (id, name, description, service_type, status, organization_id, created_by_user, version, created_at)
			  VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING id`

	var tenderID string
	err = r.pool.QueryRow(ctx, query, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId, tender.CreatorUsername, tender.Version).Scan(&tenderID)
	if err != nil {
		return domain.Tender{}, fmt.Errorf("repository.CreateTender: %w", err)
	}

	var createdTender domain.Tender
	err = r.pool.QueryRow(ctx, `
		SELECT id, name, description, service_type, status, organization_id, created_by_user, version, created_at
		FROM tender
		WHERE id = $1
	`, tenderID).Scan(
		&createdTender.ID,
		&createdTender.Name,
		&createdTender.Description,
		&createdTender.ServiceType,
		&createdTender.Status,
		&createdTender.OrganizationId,
		&createdTender.CreatorUsername,
		&createdTender.Version,
		&createdTender.CreatedAt,
	)
	if err != nil {
		return domain.Tender{}, fmt.Errorf("repository.GetTenderByID: %w", err)
	}

	return createdTender, nil
}

func (r *TenderService) GetUserTenders(ctx context.Context, limit, offset int, username string) ([]domain.Tender, error) {
	if exists, err := r.UserExists(ctx, username); err != nil || !exists {
		return nil, fmt.Errorf("repository.UpdateTenderStatus: %w", err)
	}

	query := `SELECT id, name, description, status, service_type, created_at, version
	          FROM tenders WHERE created_by_user = $1 LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, limit, offset, username)
	if err != nil {
		return nil, fmt.Errorf("repository.GetUserTenders: %w", err)
	}
	defer rows.Close()

	tenders := []domain.Tender{}
	for rows.Next() {
		var tender domain.Tender
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.Status, &tender.ServiceType, &tender.CreatedAt, &tender.Version)
		if err != nil {
			return nil, fmt.Errorf("repository.GetUserTenders: error scanning row: %w", err)
		}
		tenders = append(tenders, tender)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository.GetUserTenders: error iterating rows: %w", err)
	}

	return tenders, nil
}

func (r *TenderService) IsUserAuthorizedForOrganization(ctx context.Context, username, organizationId string) (bool, error) {
	var userID string
	err := r.pool.QueryRow(ctx, `SELECT id FROM employee WHERE username = $1`, username).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, fmt.Errorf("user does not exist")
		}
		return false, fmt.Errorf("user does not exist")
	}

	var isAuthorized bool
	err = r.pool.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1 
            FROM organization_responsible 
            WHERE user_id = $1 AND organization_id = $2
        )
    `, userID, organizationId).Scan(&isAuthorized)
	if err != nil {
		return false, fmt.Errorf("repository.IsUserAuthorizedForOrganization: %w", err)
	}

	return isAuthorized, nil
}

func (r *TenderService) GetTenderStatus(ctx context.Context, tenderID string, username string) (string, error) {
	if exists, err := r.UserExists(ctx, username); err != nil || !exists {
		return "", fmt.Errorf("repository.GetTenderStatus: %w", err)
	}

	var status string
	err := r.pool.QueryRow(ctx, `SELECT status FROM tender WHERE id = $1`, tenderID).Scan(&status)
	if err != nil {
		return "", fmt.Errorf("repository.GetTenderStatus: %w", err)
	}

	return status, nil
}

func (r *TenderService) UpdateTenderStatus(ctx context.Context, tenderID string, status string, username string) (domain.Tender, error) {
	if exists, err := r.UserExists(ctx, username); err != nil || !exists {
		return domain.Tender{}, fmt.Errorf("repository.UpdateTenderStatus: %w", err)
	}

	updateQuery := `UPDATE tender SET status = $1 WHERE id = $2`
	tag, err := r.pool.Exec(ctx, updateQuery, status, tenderID)
	if err != nil {
		return domain.Tender{}, fmt.Errorf("repository.UpdateTenderStatus: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.Tender{}, fmt.Errorf("no rows updated; check the ID")
	}

	updatedTender, err := r.GetTenderByID(ctx, tenderID)
	if err != nil {
		return domain.Tender{}, fmt.Errorf("repository.UpdateTenderStatus: failed to retrieve updated tender: %w", err)
	}

	return updatedTender, nil
}

func (r *TenderService) UserExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM employee WHERE username = $1)`
	err := r.pool.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("repository.UserExists: %w", err)
	}
	return exists, nil
}

func (r *TenderService) GetTenderByID(ctx context.Context, tenderID string) (domain.Tender, error) {
	var tender domain.Tender

	query := `SELECT id, name, description, service_type, status, organization_id, created_by_user, version, created_at FROM tender WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, tenderID).Scan(
		&tender.ID,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.CreatorUsername,
		&tender.Version,
		&tender.CreatedAt,
	)

	if err != nil {
		return domain.Tender{}, fmt.Errorf("repository.GetTenderByID: %w", err)
	}

	return tender, nil
}

func (r *TenderService) UpdatePartTender(ctx context.Context, id string, updates map[string]interface{}, username string) (domain.Tender, error) {
	if exists, err := r.UserExists(ctx, username); err != nil || !exists {
		return domain.Tender{}, fmt.Errorf("repository.UpdatePartTender: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Tender{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var currentVersion int
	err = tx.QueryRow(ctx, `
		SELECT version
		FROM tender
		WHERE id = $1
	`, id).Scan(&currentVersion)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Tender{}, fmt.Errorf("tender not found: %w", err)
		}
		return domain.Tender{}, fmt.Errorf("error fetching current version: %w", err)
	}

	query := "UPDATE tender SET "
	values := []interface{}{}
	i := 1
	if name, ok := updates["name"].(string); ok && name != "" {
		query += fmt.Sprintf("name = $%d, ", i)
		values = append(values, name)
		i++
	}
	if description, ok := updates["description"].(string); ok && description != "" {
		query += fmt.Sprintf("description = $%d, ", i)
		values = append(values, description)
		i++
	}
	if serviceType, ok := updates["serviceType"].(string); ok && serviceType != "" {
		query += fmt.Sprintf("service_type = $%d, ", i)
		values = append(values, serviceType)
		i++
	}

	query += fmt.Sprintf("version = version + 1 ")
	query += fmt.Sprintf("WHERE id = $%d", i)
	values = append(values, id)

	if len(values) == 0 {
		return domain.Tender{}, fmt.Errorf("no fields to update")
	}

	_, err = tx.Exec(ctx, query, values...)
	if err != nil {
		return domain.Tender{}, fmt.Errorf("failed to execute update query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Tender{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	var updatedTender domain.Tender
	err = r.pool.QueryRow(ctx, `
		SELECT id, name, description, service_type, status, organization_id, created_by_user, version
		FROM tender
		WHERE id = $1
	`, id).Scan(
		&updatedTender.ID,
		&updatedTender.Name,
		&updatedTender.Description,
		&updatedTender.ServiceType,
		&updatedTender.Status,
		&updatedTender.OrganizationId,
		&updatedTender.CreatorUsername,
		&updatedTender.Version,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Tender{}, fmt.Errorf("tender not found: %w", err)
		}
		return domain.Tender{}, fmt.Errorf("error fetching updated tender: %w", err)
	}

	if err := r.SaveTenderVersion(ctx, updatedTender); err != nil {
		return domain.Tender{}, fmt.Errorf("failed to save tender version: %w", err)
	}

	return updatedTender, nil
}

func (r *TenderService) SaveTenderVersion(ctx context.Context, tender domain.Tender) error {
	query := `
        INSERT INTO tender_versions (tender_id, version, name, description, service_type, status, organization_id, created_by_user)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err := r.pool.Exec(ctx, query, tender.ID, tender.Version, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId, tender.CreatorUsername)
	if err != nil {
		return fmt.Errorf("failed to save tender version: %w", err)
	}
	return nil
}

func (r *TenderService) RollbackTenderVersion(ctx context.Context, id string, targetVersion int, username string) (domain.Tender, error) {
	if exists, err := r.UserExists(ctx, username); err != nil || !exists {
		return domain.Tender{}, fmt.Errorf("repository.UpdatePartTender: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Tender{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var targetTender domain.Tender
	err = tx.QueryRow(ctx, `
        SELECT name, description, service_type, status, organization_id, created_by_user
        FROM tender_versions
        WHERE tender_id = $1 AND version = $2
    `, id, targetVersion).Scan(
		&targetTender.Name,
		&targetTender.Description,
		&targetTender.ServiceType,
		&targetTender.Status,
		&targetTender.OrganizationId,
		&targetTender.CreatorUsername,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Tender{}, fmt.Errorf("target version not found: %w", err)
		}
		return domain.Tender{}, fmt.Errorf("error fetching target version: %w", err)
	}

	var maxVersion int
	err = tx.QueryRow(ctx, `
        SELECT COALESCE(MAX(version), 0)
        FROM tender_versions
        WHERE tender_id = $1
    `, id).Scan(&maxVersion)
	if err != nil {
		return domain.Tender{}, fmt.Errorf("failed to get max version: %w", err)
	}

	newVersion := maxVersion + 1

	_, err = tx.Exec(ctx, `
        UPDATE tender
        SET name = $1, description = $2, service_type = $3, status = $4, version = $5
        WHERE id = $6
    `, targetTender.Name, targetTender.Description, targetTender.ServiceType, targetTender.Status, newVersion, id)
	if err != nil {
		return domain.Tender{}, fmt.Errorf("failed to update tender: %w", err)
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO tender_versions (tender_id, version, name, description, service_type, status, organization_id, created_by_user)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, id, newVersion, targetTender.Name, targetTender.Description, targetTender.ServiceType, targetTender.Status, targetTender.OrganizationId, targetTender.CreatorUsername)
	if err != nil {
		return domain.Tender{}, fmt.Errorf("failed to save new version: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Tender{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	var updatedTender domain.Tender
	err = r.pool.QueryRow(ctx, `
        SELECT id, name, description, service_type, status, organization_id, created_by_user, version
        FROM tender
        WHERE id = $1
    `, id).Scan(
		&updatedTender.ID,
		&updatedTender.Name,
		&updatedTender.Description,
		&updatedTender.ServiceType,
		&updatedTender.Status,
		&updatedTender.OrganizationId,
		&updatedTender.CreatorUsername,
		&updatedTender.Version,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Tender{}, fmt.Errorf("tender not found: %w", err)
		}
		return domain.Tender{}, fmt.Errorf("error fetching updated tender: %w", err)
	}

	return updatedTender, nil
}
