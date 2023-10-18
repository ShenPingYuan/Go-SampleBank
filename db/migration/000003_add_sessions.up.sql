CREATE TABLE `sessions` (
  `id` CHAR(36) PRIMARY KEY DEFAULT (UUID()),
  `user_id` BIGINT NOT NULL,
  `refresh_token` varchar(255) NOT NULL,
  `user_agent` varchar(255) NOT NULL,
  `client_ip` varchar(255) NOT NULL,
  `is_blocked` boolean NOT NULL,
  `expire_time` DATETIME NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);
ALTER TABLE `sessions`
ADD CONSTRAINT `fk_sessions_user_id_users_id` 
FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);