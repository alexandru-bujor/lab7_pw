-- Migration to drop the serial_number column from the products table, per user request.

ALTER TABLE products DROP COLUMN IF EXISTS serial_number;
