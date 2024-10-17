CREATE TABLE IF NOT EXISTS delivery_schedules (
    ID CHAR(26) PRIMARY KEY,
    company_id CHAR(26) NOT NULL,
    branch_id CHAR(26),
    courier_id CHAR(26) NOT NULL,
    partner_id CHAR(26) NOT NULL,
    status VARCHAR(255) NOT NULL, -- scheduled, on_delivery, delivered, canceled
    total DECIMAL(19, 4) NOT NULL DEFAULT 0,
    delivery_at TIMESTAMP WITH TIME ZONE NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (branch_id) REFERENCES branches(id),
    FOREIGN KEY (courier_id) REFERENCES users(id),
    FOREIGN KEY (partner_id) REFERENCES partners(id)
)