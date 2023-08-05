ALTER TABLE `accounts`
DROP FOREIGN KEY `fk_accounts_owner_id_users_id`;

ALTER TABLE `accounts`
DROP INDEX `accounts_index_0`;

ALTER TABLE `accounts`
CHANGE `owner_id` `owner` NVARCHAR(255) NOT NULL;

ALTER TABLE `accounts`
ADD INDEX `accounts_index_0` (`owner`);

DROP TABLE IF EXISTS `users`;