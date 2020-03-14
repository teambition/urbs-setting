ALTER TABLE `urbs_user` MODIFY COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE `urbs_user` MODIFY COLUMN `uid` varchar(63) NOT NULL COLLATE utf8mb4_bin;
ALTER TABLE `urbs_user` MODIFY COLUMN `labels` varchar(8190) NOT NULL COLLATE utf8mb4_bin DEFAULT '';

 COLLATE=utf8mb4_bin

ALTER TABLE `urbs_group` MODIFY COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE `urbs_group` MODIFY COLUMN `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);
ALTER TABLE `urbs_group` MODIFY COLUMN `uid` varchar(63) NOT NULL COLLATE utf8mb4_bin;
ALTER TABLE `urbs_group` MODIFY COLUMN `description` varchar(1022) NOT NULL COLLATE utf8mb4_bin DEFAULT '';
ALTER TABLE `urbs_group` ADD COLUMN `kind` varchar(63) NOT NULL COLLATE utf8mb4_bin DEFAULT '';
ALTER TABLE `urbs_group` DROP INDEX `idx_group_sync_at`;
ALTER TABLE `urbs_group` ADD INDEX `idx_group_kind` (`kind`);

ALTER TABLE `urbs_product` MODIFY COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE `urbs_product` MODIFY COLUMN `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);
ALTER TABLE `urbs_product` MODIFY COLUMN `deleted_at` datetime(3) DEFAULT NULL;
ALTER TABLE `urbs_product` MODIFY COLUMN `offline_at` datetime(3) DEFAULT NULL;
ALTER TABLE `urbs_product` MODIFY COLUMN `name` varchar(63) NOT NULL COLLATE utf8mb4_bin;
ALTER TABLE `urbs_product` MODIFY COLUMN `description` varchar(1022) NOT NULL COLLATE utf8mb4_bin DEFAULT '';

ALTER TABLE `urbs_label` MODIFY COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE `urbs_label` MODIFY COLUMN `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);
ALTER TABLE `urbs_label` MODIFY COLUMN `offline_at` datetime(3) DEFAULT NULL;
ALTER TABLE `urbs_label` MODIFY COLUMN `name` varchar(63) NOT NULL COLLATE utf8mb4_bin;
ALTER TABLE `urbs_label` MODIFY COLUMN `description` varchar(1022) NOT NULL COLLATE utf8mb4_bin DEFAULT '';
ALTER TABLE `urbs_label` MODIFY COLUMN `channels` varchar(255) NOT NULL COLLATE utf8mb4_bin DEFAULT '';
ALTER TABLE `urbs_label` MODIFY COLUMN `clients` varchar(255) NOT NULL COLLATE utf8mb4_bin DEFAULT '';

ALTER TABLE `urbs_module` MODIFY COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE `urbs_module` MODIFY COLUMN `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);
ALTER TABLE `urbs_module` MODIFY COLUMN `offline_at` datetime(3) DEFAULT NULL;
ALTER TABLE `urbs_module` MODIFY COLUMN `name` varchar(63) NOT NULL COLLATE utf8mb4_bin;
ALTER TABLE `urbs_module` MODIFY COLUMN `description` varchar(1022) NOT NULL COLLATE utf8mb4_bin DEFAULT '';

ALTER TABLE `urbs_setting` MODIFY COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE `urbs_setting` MODIFY COLUMN `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);
ALTER TABLE `urbs_setting` MODIFY COLUMN `offline_at` datetime(3) DEFAULT NULL;
ALTER TABLE `urbs_setting` MODIFY COLUMN `name` varchar(63) NOT NULL COLLATE utf8mb4_bin;
ALTER TABLE `urbs_setting` MODIFY COLUMN `description` varchar(1022) NOT NULL COLLATE utf8mb4_bin DEFAULT '';
ALTER TABLE `urbs_setting` MODIFY COLUMN `channels` varchar(255) NOT NULL COLLATE utf8mb4_bin DEFAULT '';
ALTER TABLE `urbs_setting` MODIFY COLUMN `clients` varchar(255) NOT NULL COLLATE utf8mb4_bin DEFAULT '';
ALTER TABLE `urbs_setting` MODIFY COLUMN `vals` varchar(1022) NOT NULL COLLATE utf8mb4_bin DEFAULT '';

ALTER TABLE `user_group` MODIFY COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE `user_group` DROP INDEX `idx_user_group_sync_at`;
ALTER TABLE `user_group` ADD INDEX `idx_user_group_group_id` (`group_id`);

ALTER TABLE `user_label` MODIFY COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

ALTER TABLE `user_setting` MODIFY COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE `user_setting` MODIFY COLUMN `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);
ALTER TABLE `user_setting` MODIFY COLUMN `value` varchar(255) NOT NULL COLLATE utf8mb4_bin DEFAULT '';
ALTER TABLE `user_setting` MODIFY COLUMN `last_value` varchar(255) NOT NULL COLLATE utf8mb4_bin DEFAULT '';

ALTER TABLE `group_label` MODIFY COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

ALTER TABLE `group_setting` MODIFY COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE `group_setting` MODIFY COLUMN `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);
ALTER TABLE `group_setting` MODIFY COLUMN `value` varchar(255) NOT NULL COLLATE utf8mb4_bin DEFAULT '';
ALTER TABLE `group_setting` MODIFY COLUMN `last_value` varchar(255) NOT NULL COLLATE utf8mb4_bin DEFAULT '';
