CREATE TABLE organization_responsible (
     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
     organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
     user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);
