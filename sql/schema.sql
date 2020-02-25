-- Keywords and Reserved Words https://dev.mysql.com/doc/refman/5.7/en/keywords.html name status channel value values group user
CREATE DATABASE IF NOT EXISTS `urbs` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `urbs`.`user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `active_at` bigint NOT NULL DEFAULT 0,
  `uid` varchar(63) NOT NULL,
  `labels` varchar(8191) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_uid` (`uid`),
  KEY `idx_user_active_at` (`active_at`)
) ENGINE=InnoDB;

 CREATE TABLE IF NOT EXISTS `urbs`.`group` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `sync_at` bigint NOT NULL DEFAULT 0,
  `uid` varchar(63) NOT NULL,
  `desc` varchar(1023) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_uid` (`uid`),
  KEY `idx_group_sync_at` (`sync_at`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `urbs`.`product` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  `offline_at` datetime DEFAULT NULL,
  `name` varchar(63) NOT NULL,
  `desc` varchar(1023) NOT NULL DEFAULT '',
  `status` bigint NOT NULL  DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_name` (`name`),
  KEY `idx_product_created_at` (`created_at`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `urbs`.`label` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `offline_at` datetime DEFAULT NULL,
  `product_id` bigint NOT NULL,
  `name` varchar(63) NOT NULL,
  `desc` varchar(1023) NOT NULL DEFAULT '',
  `channels` varchar(255) NOT NULL DEFAULT '', -- split by comma
  `clients` varchar(255) NOT NULL DEFAULT '', -- split by comma
  `status` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_label_product_id_name` (`product_id`,`name`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `urbs`.`module` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `offline_at` datetime DEFAULT NULL,
  `product_id` bigint NOT NULL,
  `name` varchar(63) NOT NULL,
  `desc` varchar(1023) NOT NULL DEFAULT '',
  `status` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_module_product_id_name` (`product_id`,`name`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `urbs`.`setting` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `offline_at` datetime DEFAULT NULL,
  `module_id` bigint NOT NULL,
  `name` varchar(63) NOT NULL,
  `desc` varchar(1023) NOT NULL DEFAULT '',
  `channels` varchar(255) NOT NULL DEFAULT '', -- split by comma
  `clients` varchar(255) NOT NULL DEFAULT '', -- split by comma
  `values` varchar(1023) NOT NULL DEFAULT '', -- split by comma
  `status` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_setting_module_id_name` (`module_id`,`name`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `urbs`.`user_group` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `sync_at` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `group_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_group_user_id_group_id` (`user_id`,`group_id`),
  KEY `idx_user_group_sync_at` (`sync_at`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `urbs`.`user_label` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_id` bigint NOT NULL,
  `label_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_label_user_id` (`user_id`),
  KEY `idx_user_label_label_id` (`label_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `urbs`.`user_setting` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_id` bigint NOT NULL,
  `setting_id` bigint NOT NULL,
  `value` varchar(255) NOT NULL DEFAULT '',
  `last_value` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_user_setting_user_id` (`user_id`),
  KEY `idx_user_setting_setting_id` (`setting_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `urbs`.`group_label` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `group_id` bigint NOT NULL,
  `label_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_id_label_id` (`group_id`,`label_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `urbs`.`group_setting` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `group_id` bigint NOT NULL,
  `setting_id` bigint NOT NULL,
  `value` varchar(255) NOT NULL DEFAULT '',
  `last_value` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_id_setting_id` (`group_id`,`setting_id`)
) ENGINE=InnoDB;
