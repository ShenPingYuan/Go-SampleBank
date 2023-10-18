-- Active: 1697358669446@@127.0.0.1@3306

SELECT VERSION();
CREATE DATABASE IF NOT EXISTS `simple_bank`;
use simple_bank;
show TABLES;
SELECT * from schema_migrations;
SELECT * FROM accounts;


TRUNCATE TABLE entries;
SET FOREIGN_KEY_CHECKS = 1;
TRUNCATE TABLE accounts;

select * from accounts;
SELECT * FROM users;

--日志监控
SHOW VARIABLES LIKE "general_log%";
SET GLOBAL general_log = 'on'; --开启日志监控。

SHOW VARIABLES LIKE 'general_log_file';

SET GLOBAL general_log_file = '/var/log/mysql/general.log';