BEGIN;

CREATE TABLE IF NOT EXISTS tender_versions (
    id SERIAL PRIMARY KEY,  
    tender_id UUID NOT NULL REFERENCES tender(id) ON DELETE CASCADE,
    version INT NOT NULL, 
    name VARCHAR(255),  
    description TEXT,  
    service_type VARCHAR(50),  
    status VARCHAR(10) CHECK (status IN ('CREATED', 'PUBLISHED', 'CLOSED', 'OPEN')), 
    organization_id UUID NOT NULL, 
    created_by_user VARCHAR(255) NOT NULL 
    -- created_at TIMESTAMP DEFAULT NOW()  
);

COMMIT;