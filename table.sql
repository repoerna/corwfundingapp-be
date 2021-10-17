CREATE TABLE `users` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255),
  `occupation` varchar(255),
  `email` varchar(255),
  `password_hash` varchar(255),
  `avatar_file_name` varchar(255),
  `role` varchar(255),
  `token` varchar(255),
  `created_at` datetime,
  `updated_at` datetime
);

CREATE TABLE `campaigns` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int,
  `name` varchar(255),
  `short_description` varchar(255),
  `description` text,
  `goal_amount` int,
  `current_amount` int,
  `perks` text,
  `backer_count` int,
  `slug` varchar(255),
  `created_at` datetime,
  `updated_at` datetime
);

CREATE TABLE `campaign_images` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `campaign_id` int,
  `file_name` varchar(255),
  `is_primary` boolean,
  `created_at` datetime,
  `updated_at` datetime
);

CREATE TABLE `transactions` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int,
  `campaign_id` int,
  `amount` int,
  `status` varchar(255),
  `code` varchar(255),
  `created_at` datetime,
  `updated_at` datetime
);

ALTER TABLE `campaigns` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `campaign_images` ADD FOREIGN KEY (`campaign_id`) REFERENCES `campaigns` (`id`);

ALTER TABLE `transactions` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `transactions` ADD FOREIGN KEY (`campaign_id`) REFERENCES `campaigns` (`id`);