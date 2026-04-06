-- +goose Up
CREATE TABLE `session` (
	`id` text PRIMARY KEY,
	`project_id` text NOT NULL,
	`slug` text NOT NULL,
	`directory` text NOT NULL,
	`title` text NOT NULL,
	`created_at` integer NOT NULL,
	`updated_at` integer NOT NULL,
	CONSTRAINT `fk_session_project` 
	FOREIGN KEY (`project_id`) 
	REFERENCES `project`(`id`) ON DELETE CASCADE
);

CREATE TABLE `message` (
	`id` text PRIMARY KEY,
	`session_id` text NOT NULL,
	`created_at` integer NOT NULL,
	`updated_at` integer NOT NULL,
	`data` text NOT NULL,
	CONSTRAINT `fk_message_session` 
	FOREIGN KEY (`session_id`) 
	REFERENCES `session`(`id`) ON DELETE CASCADE
);

CREATE TABLE `part` (
	`id` text PRIMARY KEY,
	`message_id` text NOT NULL,
	`session_id` text NOT NULL,
	`created_at` integer NOT NULL,
	`updated_at` integer NOT NULL,
	`data` text NOT NULL,
	CONSTRAINT `fk_part_message` 
	FOREIGN KEY (`message_id`) 
	REFERENCES `message`(`id`) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE `part`;
DROP TABLE `message`;
DROP TABLE `session`;