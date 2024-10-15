CREATE TABLE IF NOT EXISTS products (
    id CHAR(26) PRIMARY KEY,
    parent_id CHAR(26),
    company_id CHAR(26) NOT NULL,
    branch_id CHAR(26) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(19, 4) NOT NULL,
    stock INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (parent_id) REFERENCES products (id),
    FOREIGN KEY (company_id) REFERENCES companies (id),
    FOREIGN KEY (branch_id) REFERENCES branches (id)
)