-- Fix category_id / id types for products/categories foreign key
SET FOREIGN_KEY_CHECKS=0;

ALTER TABLE products DROP FOREIGN KEY products_ibfk_1;

ALTER TABLE categories MODIFY COLUMN id BIGINT NOT NULL AUTO_INCREMENT;
ALTER TABLE products MODIFY COLUMN category_id BIGINT NOT NULL;

ALTER TABLE products
    ADD CONSTRAINT products_ibfk_1
    FOREIGN KEY (category_id) REFERENCES categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE;

SET FOREIGN_KEY_CHECKS=1;
