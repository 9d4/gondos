CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `crypted_password` varchar(255) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime(3) NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `lists` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) DEFAULT NULL,
  `title` varchar(100) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT current_timestamp(),
  `updated_at` datetime(3) DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `idx_lists_user_id` (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `list_items` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `list_id` bigint(20) NOT NULL,
  `body` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT current_timestamp(),
  `updated_at` datetime(3) DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `idx_list_items_list_id` (`list_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
