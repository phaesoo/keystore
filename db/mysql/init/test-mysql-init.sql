CREATE DATABASE IF NOT EXISTS shield_test;
GRANT ALL PRIVILEGES ON `shield_test`.* TO `shield-user`@`%`;


USE shield_test;


-- CREATE TABLE `auth_key` (
--     `id` INT PRIMARY KEY AUTO_INCREMENT,
--     `access_key` VARCHAR(63) UNIQUE,
--     `secret_key` VARCHAR(63) UNIQUE,
--     `user_uuid` VARCHAR(36) UNIQUE
-- ) ENGINE=INNODB;


-- CREATE TABLE `path_permission` (
--     `id` INT PRIMARY KEY AUTO_INCREMENT,
--     `path_pattern` VARCHAR(63) UNIQUE
-- ) ENGINE=INNODB;

-- CREATE TABLE `auth_key_path_permissions` (
--     `key_id` INT,
--     `perm_id` INT,
--     FOREIGN KEY(`key_id`) REFERENCES `auth_key`(`id`),
--     FOREIGN KEY(`perm_id`) REFERENCES `path_permission`(`id`)
-- ) ENGINE=INNODB;

