CREATE TABLE IF NOT EXISTS delivery_templates (
    id CHAR(26) PRIMARY KEY,
    company_id CHAR(26) NOT NULL,
    branch_id CHAR(26),
    partner_id CHAR(26) NOT NULL,
    courier_id CHAR(26) NOT NULL,
    time_start TIME WITHOUT TIME ZONE NOT NULL,
    time_end TIME WITHOUT TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (branch_id) REFERENCES branches(id),
    FOREIGN KEY (courier_id) REFERENCES users(id),
    FOREIGN KEY (partner_id) REFERENCES partners(id)
);