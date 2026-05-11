-- Add manager_id to orders table
ALTER TABLE orders ADD COLUMN manager_id INT NULL;
ALTER TABLE orders ADD CONSTRAINT fk_orders_manager FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE SET NULL;
