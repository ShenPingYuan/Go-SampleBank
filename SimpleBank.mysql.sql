-- Active: 1684737051366@@127.0.0.1@3306@simple_bank
CREATE TABLE `accounts` (
  `id` BIGINT PRIMARY KEY,
  `owner` NVARCHAR(255) NOT NULL,
  `balance` DECIMAL(10, 2) NOT NULL,
  `currency` NVARCHAR(255) NOT NULL,
  --mysql当前时区
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `entries` (
  `id` BIGINT PRIMARY KEY,
  `account_id` BIGINT NOT NULL,
  `amount` BIGINT NOT NULL COMMENT '可以是正或者负，表示存入或者取出',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `transfers` (
  `id` BIGINT PRIMARY KEY,
  `from_account_id` BIGINT NOT NULL,
  `to_account_id` BIGINT NOT NULL COMMENT 'Content of the post',
  `amount` BIGINT NOT NULL COMMENT '只能是正数',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX `accounts_index_0` ON `accounts` (`owner`);

CREATE INDEX `entries_index_1` ON `entries` (`account_id`);

CREATE INDEX `transfers_index_2` ON `transfers` (`from_account_id`);

CREATE INDEX `transfers_index_3` ON `transfers` (`to_account_id`);

CREATE INDEX `transfers_index_4` ON `transfers` (`from_account_id`, `to_account_id`);

ALTER TABLE `entries` ADD FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`);

ALTER TABLE `transfers` ADD FOREIGN KEY (`from_account_id`) REFERENCES `accounts` (`id`);

ALTER TABLE `transfers` ADD FOREIGN KEY (`to_account_id`) REFERENCES `accounts` (`id`);