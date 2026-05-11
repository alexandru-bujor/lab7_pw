-- -------------------------------------------------------------------------------------
-- Schema Updates Migration Script (Since Yesterday)
-- -------------------------------------------------------------------------------------

-- 1. Fix duplicated or empty `slug` values to allow the UNIQUE constraint `uni_products_slug`
UPDATE products SET slug = CONCAT('product-', id) WHERE slug = '' OR slug IS NULL;

-- 2. Modify `name` columns to JSON for translations (RO, RU, EN)
ALTER TABLE `makes` MODIFY COLUMN `name` JSON;
ALTER TABLE `service_categories` MODIFY COLUMN `name` JSON;

-- 3. Correct products Category foreign key and IDs (fixing Error 1832 constraints)
SET FOREIGN_KEY_CHECKS=0;
-- We use ignore errors if constraint doesn't exist by dropping it first, assuming it is products_ibfk_1
-- (If it fails on your system because it's named differently, you can skip the drop line or adjust it)
ALTER TABLE products DROP FOREIGN KEY IF EXISTS products_ibfk_1;

ALTER TABLE categories MODIFY COLUMN id BIGINT NOT NULL AUTO_INCREMENT;     
ALTER TABLE products MODIFY COLUMN category_id BIGINT NOT NULL;

ALTER TABLE products 
    ADD CONSTRAINT products_ibfk_1 
    FOREIGN KEY (category_id) REFERENCES categories(id) 
    ON DELETE RESTRICT ON UPDATE CASCADE;
SET FOREIGN_KEY_CHECKS=1;

-- 4. Add new columns to `inventory_items` if they don't already exist 
-- Note: GORM's AutoMigrate usually adds these automatically, but we include them here.
-- MySQL automatically skips adding them if you run GORM, but this is the SQL equivalent:

-- ALTER TABLE `inventory_items` ADD COLUMN `is_visible_in_secondary` BOOLEAN DEFAULT FALSE;
-- ALTER TABLE `inventory_items` ADD COLUMN `accounting_type` VARCHAR(50) DEFAULT 'on_book';
