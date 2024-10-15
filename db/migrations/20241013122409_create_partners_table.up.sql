CREATE TABLE IF NOT EXISTS partners (
    id CHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    company_id CHAR(26) NOT NULL,
    branch_id CHAR(26) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    address TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (company_id) REFERENCES companies (id),
    FOREIGN KEY (branch_id) REFERENCES branches (id),
    CONSTRAINT partners_email_unique UNIQUE (company_id, email)
);