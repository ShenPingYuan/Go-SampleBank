-- Active: 1684737051366@@127.0.0.1@3306@simple_bank
SELECT VERSION()

SELECT * FROM accounts;


TRUNCATE TABLE entries;
SET FOREIGN_KEY_CHECKS = 1;
TRUNCATE TABLE accounts;

select * from accounts;


--日志监控
SHOW VARIABLES LIKE "general_log%";
SET GLOBAL general_log = 'on'; --开启日志监控。

SHOW VARIABLES LIKE 'general_log_file';

SET GLOBAL general_log_file = '/var/log/mysql/general.log';

