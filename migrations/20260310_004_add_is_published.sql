-- Migration: Add is_published column to products
ALTER TABLE products ADD COLUMN is_published BOOLEAN DEFAULT TRUE;
