-- Migration: Create service_categories table for multilingual service categories
CREATE TABLE IF NOT EXISTS service_categories (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name JSON NOT NULL
);
