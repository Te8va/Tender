package domain

import (
	"time"
)

type Employee struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type OrganizationType string

const (
	OrganizationTypeIE  OrganizationType = "IE"
	OrganizationTypeLLC OrganizationType = "LLC"
	OrganizationTypeJSC OrganizationType = "JSC"
)

type Organization struct {
	ID          string           `json:"id" db:"id"`
	Name        string           `json:"name" db:"name"`
	Description string           `json:"description" db:"description"`
	Type        OrganizationType `json:"type" db:"type"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
}

type OrganizationResponsible struct {
	ID             string `json:"id" db:"id"`
	OrganizationID string `json:"organization_id" db:"organization_id"`
	UserID         string `json:"user_id" db:"user_id"`
}

type Tender struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	ServiceType     string    `json:"serviceType"`
	OrganizationId  string    `json:"organizationId"`
	CreatorUsername string    `json:"creatorUsername"`
	Version         int       `json:"version"`
	CreatedAt       time.Time `json:"createdAt"`
}
type CreateTenderRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	OrganizationId  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}

type TenderResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ServiceType string    `json:"serviceType"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
}

type TenderVersion struct {
	ID              int    `json:"id"`
	TenderID        string `json:"tender_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	Status          string `json:"status"`
	OrganizationId  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
	Version         int    `json:"version"`
}

type Bid struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	TenderId        int    `json:"tenderId"`
	OrganizationId  int    `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}

type TenderStatusUpdate struct {
	Status string `json:"status"`
}
