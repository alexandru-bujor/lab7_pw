-- -------------------------------------------------------------------------------------
-- Schema Updates Migration Script (Based on Exact DB Dump)
-- -------------------------------------------------------------------------------------
-- The issue is that GORM wants to upgrade `int(11)` to `BIGINT` because `int` in Go is 64-bit.
-- MySQL blocks this because these columns are tied together with Foreign Keys.
-- We must drop ALL related foreign keys across the product/category domain, upgrade the types,
-- and then recreate them all at once.

SET FOREIGN_KEY_CHECKS=0;

-- 1. DROP ALL FOREIGN KEYS INVOLVING INT COLUMNS
ALTER TABLE `categories` DROP FOREIGN KEY `categories_ibfk_1`;
ALTER TABLE `products` DROP FOREIGN KEY `products_ibfk_1`;
ALTER TABLE `products` DROP FOREIGN KEY `products_ibfk_2`;
ALTER TABLE `product_images` DROP FOREIGN KEY `product_images_ibfk_1`;
ALTER TABLE `product_variants` DROP FOREIGN KEY `product_variants_ibfk_1`;

-- 2. MODIFY COLUMNS TO BIGINT TO MATCH GORM'S EXPECTATIONS (int in Go = bigint in DB)
-- This makes all IDs compatible with each other and GORM.
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

-- 3. RECREATE FOREIGN KEYS WITH THE NEW BIGINT TYPES
ALTER TABLE `categories`
  ADD CONSTRAINT `categories_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `categories` (`id`) ON DELETE SET NULL;

ALTER TABLE `products`
  ADD CONSTRAINT `products_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `products_ibfk_2` FOREIGN KEY (`brand_id`) REFERENCES `brands` (`id`) ON DELETE SET NULL;

ALTER TABLE `product_images`
  ADD CONSTRAINT `product_images_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE;

ALTER TABLE `product_variants`
  ADD CONSTRAINT `product_variants_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE;

-- 4. FIX UNIQUE SLUG CONSTRAINT AND JSON DATA TYPES
UPDATE products SET slug = CONCAT('product-', id) WHERE slug = '' OR slug IS NULL;
ALTER TABLE `makes` MODIFY COLUMN `name` JSON;
ALTER TABLE `service_categories` MODIFY COLUMN `name` JSON;

-- 5. RE-ENABLE FOREIGN KEYS
SET FOREIGN_KEY_CHECKS=1;
