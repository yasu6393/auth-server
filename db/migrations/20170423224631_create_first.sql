
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+09:00";

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) UNSIGNED NOT NULL COMMENT '�V�[�P���X',
  `user` varchar(32) NOT NULL COMMENT '���[�U�[ID',
  `password` varchar(1024) NOT NULL COMMENT '�p�X���[�h',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='���[�U�[���';

DROP TABLE IF EXISTS `client`;
CREATE TABLE `client` (
  `id` bigint(20) UNSIGNED NOT NULL COMMENT '�V�[�P���X',
  `client_id` varchar(32) NOT NULL COMMENT '�N���C�A���gID',
  `client_secret` varchar(1024) NOT NULL COMMENT '�N���C�A���g�V�[�N���b�g',
  `redirect_uri` varchar(1024) NOT NULL COMMENT '���_�C���N�gURL',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='���[�U�[���';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS `user`;