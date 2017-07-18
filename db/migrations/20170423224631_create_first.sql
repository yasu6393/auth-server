
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+09:00";

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'シーケンス',
  `user` varchar(32) NOT NULL COMMENT 'ユーザーID',
  `password` varchar(1024) NOT NULL COMMENT 'パスワード',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(`id`),
  UNIQUE `user_UNIQUE` (`user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='ユーザー情報';

DROP TABLE IF EXISTS `client`;
CREATE TABLE `client` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'シーケンス',
  `client_id` varchar(32) NOT NULL COMMENT 'クライアントID',
  `client_secret` varchar(1024) NOT NULL COMMENT 'クライアントシークレット',
  `redirect_uri` varchar(1024) NOT NULL COMMENT 'リダイレクトURL',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(`id`),
  UNIQUE `client_id_UNIQUE` (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='クライアント情報';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `user` DROP `INDEX user_UNIQUE`;
ALTER TABLE `client` DROP INDEX `client_id_UNIQUE`;

DROP TABLE IF EXISTS `client`;
DROP TABLE IF EXISTS `user`;
