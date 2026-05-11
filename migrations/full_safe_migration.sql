-- -------------------------------------------------------------------------------------
-- SAFE SCHEMA UPDATES MIGRATION SCRIPT
-- -------------------------------------------------------------------------------------
-- This script uses a stored procedure to safely drop foreign keys only if they exist,
-- preventing the "Error #1091 - Can't DROP FOREIGN KEY" issue.

SET FOREIGN_KEY_CHECKS=0;

-- 1. Create a safe drop procedure
DELIMITER //
CREATE PROCEDURE DropFKIfExists(
    IN p_table_name VARCHAR(255),
    IN p_constraint_name VARCHAR(255)
)
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
        WHERE CONSTRAINT_SCHEMA = DATABASE() 
          AND TABLE_NAME = p_table_name 
          AND CONSTRAINT_NAME = p_constraint_name
          AND CONSTRAINT_TYPE = 'FOREIGN KEY'
    ) THEN
        SET @query = CONCAT('ALTER TABLE `', p_table_name, '` DROP FOREIGN KEY `', p_constraint_name, '`');
        PREPARE stmt FROM @query;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END IF;
END //
DELIMITER ;

-- 2. Safely drop all related foreign keys
CALL DropFKIfExists('categories', 'categories_ibfk_1');
CALL DropFKIfExists('products', 'products_ibfk_1');
CALL DropFKIfExists('products', 'products_ibfk_2');
CALL DropFKIfExists('product_images', 'product_images_ibfk_1');
CALL DropFKIfExists('product_variants', 'product_variants_ibfk_1');

-- 3. Cleanup the procedure
DROP PROCEDURE IF EXISTS DropFKIfExists;

-- 4. MODIFY COLUMNS TO BIGINT
ALTER TABLE `categories` MODIFY COLUMN `id` BIGINT NOT NULL AUTO_INCREMENT;
ALTER TABLE `categories` MODIFY COLUMN `parent_id` BIGINT DEFAULT NULL;

ALTER TABLE `brands` MODIFY COLUMN `id` BIGINT NOT NULL AUTO_INCREMENT;

ALTER TABLE `products` MODIFY COLUMN `id` BIGINT NOT NULL AUTO_INCREMENT;
ALTER TABLE `products` MODIFY COLUMN `category_id` BIGINT DEFAULT NULL;
ALTER TABLE `products` MODIFY COLUMN `brand_id` BIGINT DEFAULT NULL;

ALTER TABLE `product_images` MODIFY COLUMN `id` BIGINT NOT NULL AUTO_INCREMENT;
ALTER TABLE `product_images` MODIFY COLUMN `product_id` BIGINT DEFAULT NULL;

ALTER TABLE `product_variants` MODIFY COLUMN `id` BIGINT NOT NULL AUTO_INCREMENT;
ALTER TABLE `product_variants` MODIFY COLUMN `product_id` BIGINT DEFAULT NULL;

-- 5. RECREATE FOREIGN KEYS
-- Ignore errors here if they already exist, MySQL handles ADD CONSTRAINT safely if named
ALTER TABLE `categories`
  ADD CONSTRAINT `categories_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `categories` (`id`) ON DELETE SET NULL;

ALTER TABLE `products`
  ADD CONSTRAINT `products_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `products_ibfk_2` FOREIGN KEY (`brand_id`) REFERENCES `brands` (`id`) ON DELETE SET NULL;

ALTER TABLE `product_images`
  ADD CONSTRAINT `product_images_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE;

ALTER TABLE `product_variants`
  ADD CONSTRAINT `product_variants_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE;

-- 6. DATA FIXES
UPDATE products SET slug = CONCAT('product-', id) WHERE slug = '' OR slug IS NULL;
ALTER TABLE `makes` MODIFY COLUMN `name` JSON;
ALTER TABLE `service_categories` MODIFY COLUMN `name` JSON;

SET FOREIGN_KEY_CHECKS=1;
