CREATE DATABASE IF NOT EXISTS shield;
GRANT ALL PRIVILEGES ON `shield`.* TO `shield-user`@`%`;


USE shield;


CREATE TABLE `key_user` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `access_key` VARCHAR(63) UNIQUE,
    `secret_key` VARCHAR(63) UNIQUE,
    `user_uuid` VARCHAR(32) UNIQUE,
) ENGINE=INNODB;


CREATE TABLE `permission` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `path` VARCHAR(63) UNIQUE,
)

CREATE TABLE `key_permission` (
    `key_id` INT
    `perm_id` INT
    FOREIGN KEY(`key_id`) REFERENCES `key_user`(`id`)
    FOREIGN KEY(`perm_id`) REFERENCES `permission`(`id`)
) ENGINE=INNODB;


INSERT INTO key_user (id, access_key, secret_key, user_uuid)
VALUES
(0, "123", "456", "user1"),
(1, "789", "012", "user2");


INSERT INTO permission (id, path)
VALUES
(0, "/markets/all");


INSERT INTO key_permission (key_id, perm_id)
VALUES
(0, 0);