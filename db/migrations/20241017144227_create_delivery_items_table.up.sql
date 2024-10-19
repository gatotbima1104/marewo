CREATE TABLE IF NOT EXISTS delivery_items (
    id CHAR(26) PRIMARY KEY,
    delivery_schedule_id CHAR(26) NOT NULL,
    product_id CHAR(26) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(19, 4) NOT NULL,
    total DECIMAL(19, 4) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (delivery_schedule_id) REFERENCES delivery_schedules(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
);