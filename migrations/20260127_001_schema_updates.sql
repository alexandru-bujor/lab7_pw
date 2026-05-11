-- Migration: Schema updates for categories, products, settings, and secondary/accounting fields
-- Run in order. If a column/table already exists (e.g. from GORM AutoMigrate), skip that statement or ignore the error.

-- ---------------------------------------------------------------------------
-- Categories: parent_id for subcategories
-- ---------------------------------------------------------------------------
ALTER TABLE categories ADD COLUMN parent_id INT NULL DEFAULT NULL AFTER slug;
CREATE INDEX idx_categories_parent_id ON categories(parent_id);

-- ---------------------------------------------------------------------------
-- Products: old_price for sale/discount display
-- ---------------------------------------------------------------------------
ALTER TABLE products ADD COLUMN old_price DECIMAL(12,2) NOT NULL DEFAULT 0 AFTER price;

-- ---------------------------------------------------------------------------
-- Settings table (for secondary panel password, etc.)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS settings (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `key` VARCHAR(255) NOT NULL,
    value TEXT,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    UNIQUE KEY uk_settings_key (`key`)
);

-- ---------------------------------------------------------------------------
-- Category templates (specs per category)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS category_templates (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    category_id INT UNSIGNED NOT NULL,
    name VARCHAR(255) DEFAULT '',
    specs JSON NULL,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    INDEX idx_category_templates_category_id (category_id)
);

-- ---------------------------------------------------------------------------
-- Orders: accounting_type and is_visible_in_secondary
-- ---------------------------------------------------------------------------
ALTER TABLE orders ADD COLUMN accounting_type VARCHAR(32) DEFAULT 'on_book' AFTER order_source;
ALTER TABLE orders ADD COLUMN is_visible_in_secondary TINYINT(1) DEFAULT 0 AFTER accounting_type;

-- ---------------------------------------------------------------------------
-- Inventory items: accounting_type and is_visible_in_secondary
-- ---------------------------------------------------------------------------
ALTER TABLE inventory_items ADD COLUMN accounting_type VARCHAR(32) DEFAULT 'on_book' AFTER notes;
ALTER TABLE inventory_items ADD COLUMN is_visible_in_secondary TINYINT(1) DEFAULT 0 AFTER accounting_type;

-- ---------------------------------------------------------------------------
-- Clients: is_visible_in_secondary
-- ---------------------------------------------------------------------------
ALTER TABLE clients ADD COLUMN is_visible_in_secondary TINYINT(1) DEFAULT 0 AFTER last_order_date;

-- ---------------------------------------------------------------------------
-- Products: accounting_type and is_visible_in_secondary
-- ---------------------------------------------------------------------------
ALTER TABLE products ADD COLUMN accounting_type VARCHAR(32) DEFAULT 'on_book' AFTER variant_ram;
ALTER TABLE products ADD COLUMN is_visible_in_secondary TINYINT(1) DEFAULT 0 AFTER accounting_type;

-- ---------------------------------------------------------------------------
-- Services: prices, photo_url, is_visible_in_secondary (skip if column exists)
-- ---------------------------------------------------------------------------
ALTER TABLE services ADD COLUMN prices JSON NULL;
ALTER TABLE services ADD COLUMN photo_url VARCHAR(512) DEFAULT '';
ALTER TABLE services ADD COLUMN is_visible_in_secondary TINYINT(1) DEFAULT 0;

-- ---------------------------------------------------------------------------
-- Lombard requests: is_visible_in_secondary
-- ---------------------------------------------------------------------------
ALTER TABLE lombard_requests ADD COLUMN is_visible_in_secondary TINYINT(1) DEFAULT 0 AFTER status;

-- ---------------------------------------------------------------------------
-- Users: permissions (JSON)
-- ---------------------------------------------------------------------------
ALTER TABLE users ADD COLUMN permissions JSON NULL AFTER role;
