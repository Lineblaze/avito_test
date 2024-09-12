CREATE TABLE bid_feedback (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    feedback VARCHAR(1000) NOT NULL,
    username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX bid_feedback_unique_idx ON bid_feedback (bid_id, username);