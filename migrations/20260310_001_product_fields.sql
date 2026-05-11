-- Migration to add new product fields (Model, Code, Serial, Battery, Labels)
ALTER TABLE products 
    ADD COLUMN model VARCHAR(255) DEFAULT '',
    ADD COLUMN product_code VARCHAR(255) DEFAULT '',
    ADD COLUMN serial_number VARCHAR(255) DEFAULT '',
    ADD COLUMN battery_health VARCHAR(50) DEFAULT '',
    ADD COLUMN is_recommended BOOLEAN DEFAULT FALSE,
    ADD COLUMN is_new_arrival BOOLEAN DEFAULT FALSE;
