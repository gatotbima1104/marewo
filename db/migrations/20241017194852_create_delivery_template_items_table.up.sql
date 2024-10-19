CREATE TABLE IF NOT EXISTS delivery_template_items (
    id CHAR(26) PRIMARY KEY,
    delivery_template_id CHAR(26) NOT NULL,
    product_id CHAR(26) NOT NULL,
    quantity INT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (delivery_template_id) REFERENCES delivery_templates(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
);