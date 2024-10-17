CREATE TABLE IF NOT EXISTS partners (
    id CHAR(26) PRIMARY KEY,
    company_id CHAR(26) NOT NULL,
    branch_id CHAR(26),
    name VARCHAR(255) NOT NULL,
    address TEXT,
    phone_country_code VARCHAR(5),
    phone VARCHAR(255) NOT NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (company_id) REFERENCES companies (id),
    FOREIGN KEY (branch_id) REFERENCES branches (id)
);