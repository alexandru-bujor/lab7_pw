ALTER TABLE `categories` MODIFY COLUMN `name` JSON;
UPDATE `categories` SET `name` = JSON_OBJECT('ro', name, 'en', name, 'ru', name) WHERE JSON_VALID(name) = 0;
