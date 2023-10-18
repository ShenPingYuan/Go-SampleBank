ALTER TABLE `sessions`
DROP FOREIGN KEY `fk_sessions_user_id_users_id`;

DROP TABLE IF EXISTS `sessions`;