ALTER TABLE `urbs_group` ADD COLUMN `status` bigint NOT NULL  DEFAULT 0;
ALTER TABLE `urbs_label` ADD COLUMN `rls` bigint NOT NULL  DEFAULT 0;
ALTER TABLE `urbs_setting` ADD COLUMN `rls` bigint NOT NULL  DEFAULT 0;
ALTER TABLE `user_label` ADD COLUMN `rls` bigint NOT NULL  DEFAULT 0;
ALTER TABLE `user_setting` ADD COLUMN `rls` bigint NOT NULL  DEFAULT 0;
ALTER TABLE `group_label` ADD COLUMN `rls` bigint NOT NULL  DEFAULT 0;
ALTER TABLE `group_setting` ADD COLUMN `rls` bigint NOT NULL  DEFAULT 0;

ALTER TABLE `user_label` ADD INDEX `idx_user_label_label_id` (`label_id`);
ALTER TABLE `user_setting` ADD INDEX `idx_user_setting_setting_id` (`setting_id`);
ALTER TABLE `group_label` ADD INDEX `idx_group_label_label_id` (`label_id`);
ALTER TABLE `group_setting` ADD INDEX `idx_group_setting_setting_id` (`setting_id`);

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
  `name` varchar(255) NOT NULL,
  `value` varchar(8190) NOT NULL DEFAULT '',
  `status` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_urbs_statistic_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `urbs`.`urbs_lock` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `expire_at` datetime(3) NOT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_urbs_statistic_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
