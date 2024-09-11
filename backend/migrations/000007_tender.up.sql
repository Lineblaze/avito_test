CREATE TABLE tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    service_type tender_service_type,
    status tender_status,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    version INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
