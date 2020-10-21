ALTER TABLE `urbs_group` DROP index `uk_group_uid`;
ALTER TABLE `urbs_group` ADD unique `uk_group_uid_kind` (`uid`,`kind`);