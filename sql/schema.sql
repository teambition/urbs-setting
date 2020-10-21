-- Keywords and Reserved Words https://dev.mysql.com/doc/refman/5.7/en/keywords.html name status channel value values group user
CREATE DATABASE IF NOT EXISTS `urbs` CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`urbs_user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `active_at` bigint NOT NULL DEFAULT 0,
  `uid` varchar(63) NOT NULL,
  `labels` varchar(8190) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_uid` (`uid`),
  KEY `idx_user_active_at` (`active_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

 CREATE TABLE IF NOT EXISTS `urbs`.`urbs_group` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `sync_at` bigint NOT NULL DEFAULT 0,
  `uid` varchar(63) NOT NULL,
  `kind` varchar(63) NOT NULL DEFAULT '',
  `description` varchar(1022) NOT NULL DEFAULT '',
  `status` bigint NOT NULL  DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_uid_kind` (`uid`,`kind`),
  KEY `idx_group_kind` (`kind`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`urbs_product` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `offline_at` datetime(3) DEFAULT NULL,
  `name` varchar(63) NOT NULL,
  `description` varchar(1022) NOT NULL DEFAULT '',
  `status` bigint NOT NULL  DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_name` (`name`),
  KEY `idx_product_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`urbs_label` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `offline_at` datetime(3) DEFAULT NULL,
  `product_id` bigint NOT NULL,
  `name` varchar(63) NOT NULL,
  `description` varchar(1022) NOT NULL DEFAULT '',
  `channels` varchar(255) NOT NULL DEFAULT '', -- split by comma
  `clients` varchar(255) NOT NULL DEFAULT '', -- split by comma
  `status` bigint NOT NULL DEFAULT 0,
  `rls` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_label_product_id_name` (`product_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`urbs_module` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `offline_at` datetime(3) DEFAULT NULL,
  `product_id` bigint NOT NULL,
  `name` varchar(63) NOT NULL,
  `description` varchar(1022) NOT NULL DEFAULT '',
  `status` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_module_product_id_name` (`product_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`urbs_setting` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `offline_at` datetime(3) DEFAULT NULL,
  `module_id` bigint NOT NULL,
  `name` varchar(63) NOT NULL,
  `description` varchar(1022) NOT NULL DEFAULT '',
  `channels` varchar(255) NOT NULL DEFAULT '', -- split by comma
  `clients` varchar(255) NOT NULL DEFAULT '', -- split by comma
  `vals` varchar(1022) NOT NULL DEFAULT '', -- split by comma
  `status` bigint NOT NULL DEFAULT 0,
  `rls` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_setting_module_id_name` (`module_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`user_group` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `sync_at` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `group_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_group_user_id_group_id` (`user_id`,`group_id`),
  KEY `idx_user_group_group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`user_label` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `user_id` bigint NOT NULL,
  `label_id` bigint NOT NULL,
  `rls` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_label_user_id_label_id` (`user_id`,`label_id`),
  KEY `idx_user_label_label_id` (`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`user_setting` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `user_id` bigint NOT NULL,
  `setting_id` bigint NOT NULL,
  `value` varchar(255) NOT NULL DEFAULT '',
  `last_value` varchar(255) NOT NULL DEFAULT '',
  `rls` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_setting_user_id_setting_id` (`user_id`,`setting_id`),
  KEY `idx_user_setting_setting_id` (`setting_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`group_label` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `group_id` bigint NOT NULL,
  `label_id` bigint NOT NULL,
  `rls` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_label_group_id_label_id` (`group_id`,`label_id`),
  KEY `idx_group_label_label_id` (`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`group_setting` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `group_id` bigint NOT NULL,
  `setting_id` bigint NOT NULL,
  `value` varchar(255) NOT NULL DEFAULT '',
  `last_value` varchar(255) NOT NULL DEFAULT '',
  `rls` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_setting_group_id_setting_id` (`group_id`,`setting_id`),
  KEY `idx_group_setting_setting_id` (`setting_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`label_rule` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `product_id` bigint NOT NULL,
  `label_id` bigint NOT NULL,
  `kind` varchar(63) NOT NULL,
  `rule` varchar(1022) NOT NULL DEFAULT '',
  `rls` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_label_rule_label_id_kind` (`label_id`,`kind`),
  KEY `idx_label_rule_product_id` (`product_id`),
  KEY `idx_label_rule_label_id` (`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`setting_rule` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `product_id` bigint NOT NULL,
  `setting_id` bigint NOT NULL,
  `kind` varchar(63) NOT NULL,
  `rule` varchar(1022) NOT NULL DEFAULT '',
  `value` varchar(255) NOT NULL DEFAULT '',
  `rls` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_setting_rule_setting_id_kind` (`setting_id`,`kind`),
  KEY `idx_setting_rule_product_id` (`product_id`),
  KEY `idx_setting_rule_setting_id` (`setting_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`urbs_statistic` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `name` varchar(127) NOT NULL,
  `value` varchar(8190) NOT NULL DEFAULT '',
  `status` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_urbs_statistic_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`urbs_lock` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `expire_at` datetime(3) NOT NULL,
  `name` varchar(127) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_urbs_lock_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
