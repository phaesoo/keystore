CREATE DATABASE IF NOT EXISTS shield;
GRANT ALL PRIVILEGES ON `shield`.* TO `shield-user`@`%`;


USE shield;


CREATE TABLE `auth_key` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `access_key` VARCHAR(63) UNIQUE,
    `secret_key` VARCHAR(63) UNIQUE,
    `user_uuid` VARCHAR(32) UNIQUE
) ENGINE=INNODB;


CREATE TABLE `path_permission` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `path_pattern` VARCHAR(63) UNIQUE
) ENGINE=INNODB;

CREATE TABLE `auth_key_path_permissions` (
    `key_id` INT,
    `perm_id` INT,
    FOREIGN KEY(`key_id`) REFERENCES `auth_key`(`id`),
    FOREIGN KEY(`perm_id`) REFERENCES `path_permission`(`id`)
) ENGINE=INNODB;


INSERT INTO auth_key (access_key, secret_key, user_uuid)
VALUES
("123", "456", "user1"),
("789", "012", "user2");


INSERT INTO path_permission (path_pattern)
VALUES
("/markets/all");


INSERT INTO auth_key_path_permissions (key_id, perm_id)
VALUES
(1, 1);