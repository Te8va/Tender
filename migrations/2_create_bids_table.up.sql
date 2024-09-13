BEGIN;

CREATE TABLE IF NOT EXISTS bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(10) CHECK (status IN ('CREATED', 'SUBMITTED', 'PUBLISHED', 'CANCELED')) NOT NULL DEFAULT 'CREATED',
    tender_id UUID NOT NULL REFERENCES tender(id) ON DELETE CASCADE,
    -- organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    -- created_by_user_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    version INT DEFAULT 1 NOT NULL
);


COMMIT;