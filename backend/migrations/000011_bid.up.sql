CREATE TABLE bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status bid_status,
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    author_type author_type,
    author_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    version INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

