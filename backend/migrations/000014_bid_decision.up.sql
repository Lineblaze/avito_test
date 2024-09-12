CREATE TABLE bid_decision (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
      decision decision,
      username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE
);

CREATE UNIQUE INDEX bid_decision_unique_idx ON bid_decision (bid_id, username);