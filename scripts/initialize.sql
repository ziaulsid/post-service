-- clean up all the data
DROP DATABASE IF EXISTS post_db;

-- recreate database
CREATE DATABASE IF NOT EXISTS post_db;

USE post_db;

CREATE TABLE `posts`
(
    `id`                  bigint unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `caption`             varchar(255),
    `original_image`      longtext,
    `original_image_name` varchar(255),
    `resized_image`       longtext,
    `created_by`          bigint unsigned NOT NULL,
    `created_at`          timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at`          timestamp       NULL     DEFAULT NULL,
     KEY `created_by_idx` (`created_by`)
) DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ENGINE = INNODB;

CREATE TABLE `comments`
(
    `id`         bigint unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `content`    varchar(255),
    `created_by` bigint unsigned NOT NULL,
    `post_id`    bigint unsigned NOT NULL,
    `created_at` timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` timestamp       NULL     DEFAULT NULL,
     KEY `post_id_idx` (`post_id`)
) DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ENGINE = INNODB;

CREATE TABLE `users`
(
    `id`         bigint unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `username`   varchar(255)    NOT NULL,
    `name`       varchar(255)    NOT NULL,
    `created_at` timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` timestamp       NULL     DEFAULT NULL
) DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ENGINE = INNODB;
