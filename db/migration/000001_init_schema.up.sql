-- Active: 1684737051366@@127.0.0.1@3306@simple_bank
CREATE TABLE `accounts` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `owner` NVARCHAR(255) NOT NULL,
  `balance` DECIMAL(10, 2) NOT NULL,
  `currency` NVARCHAR(255) NOT NULL,
  -- mysql当前时区
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT `check_balance_positive` CHECK (`balance` > 0)
);

CREATE TABLE `entries` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `account_id` BIGINT NOT NULL,
  `amount` DECIMAL(10, 2) NOT NULL COMMENT '可以是正或者负，表示存入或者取出',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `transfers` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `from_account_id` BIGINT NOT NULL,
  `to_account_id` BIGINT NOT NULL COMMENT 'Content of the post',
  `amount` DECIMAL(10, 2) NOT NULL COMMENT '只能是正数',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT `check_amount_positive` CHECK (`amount` > 0)
);

CREATE INDEX `accounts_index_0` ON `accounts` (`owner`);

CREATE INDEX `entries_index_1` ON `entries` (`account_id`);

CREATE INDEX `transfers_index_2` ON `transfers` (`from_account_id`);

CREATE INDEX `transfers_index_3` ON `transfers` (`to_account_id`);

CREATE INDEX `transfers_index_4` ON `transfers` (`from_account_id`, `to_account_id`);
