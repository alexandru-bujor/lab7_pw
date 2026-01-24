-- Migration: Set all existing products to 'on_book'
-- This query sets all products that don't have accounting_type set to 'on_book'

UPDATE products 
SET accounting_type = 'on_book' 
WHERE accounting_type IS NULL OR accounting_type = '';

-- Verify the update
SELECT COUNT(*) as total_products, 
       COUNT(CASE WHEN accounting_type = 'on_book' THEN 1 END) as on_book_products,
       COUNT(CASE WHEN accounting_type = 'off_book' THEN 1 END) as off_book_products,
       COUNT(CASE WHEN accounting_type IS NULL OR accounting_type = '' THEN 1 END) as null_products
FROM products;

