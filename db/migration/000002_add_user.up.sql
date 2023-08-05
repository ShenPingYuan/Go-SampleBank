-- Active: 1685868021088@@localhost@3306@simple_bank

CREATE TABLE `users` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `username` nvarchar(255) UNIQUE NOT NULL,
  `hashed_password` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `full_name` varchar(255) NOT NULL,
  `password_changed_at` DATETIME NOT NULL DEFAULT '0001-01-01 00:00:00',
  `created_at` DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

ALTER TABLE `accounts`
DROP INDEX `accounts_index_0`;

ALTER TABLE `accounts`
CHANGE `owner` `owner_id` bigint UNIQUE NOT NULL;
CREATE INDEX `users_index_0` ON `users` (`username`);
CREATE INDEX `accounts_index_0` ON `accounts` (`owner_id`);

ALTER TABLE `accounts`
ADD CONSTRAINT `fk_accounts_owner_id_users_id`
FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`);