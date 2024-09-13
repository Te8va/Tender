BEGIN;

CREATE TABLE IF NOT EXISTS tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    service_type VARCHAR(50) NOT NULL,
    status VARCHAR(10) CHECK (status IN ('CREATED', 'PUBLISHED', 'CLOSED', 'OPEN')) NOT NULL DEFAULT 'CREATED',
    version INT DEFAULT 1 NOT NULL,
    -- organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    -- created_by_user_id UUID REFERENCES employee(id) ON DELETE CASCADE
    organization_id UUID NOT NULL,
    created_by_user VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
       CONSTRAINT fk_organization
        FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT fk_created_by_user
        FOREIGN KEY (created_by_user) REFERENCES employee(username) ON DELETE SET NULL
);

COMMIT;