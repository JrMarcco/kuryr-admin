DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `username` VARCHAR(64) NOT NULL COMMENT '用户名',
    `password` VARCHAR(128) NOT NULL COMMENT '密码（请存储哈希后的值）',
    `email` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '邮箱',
    `user_type` VARCHAR(16) NOT NULL COMMENT '用户类型：administrator-系统管理员 operator-业务操作员',
    `biz_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '关联的业务方 id，当 user_type 为 operator 时不应为 0',
    `created_at` BIGINT UNSIGNED NOT NULL COMMENT '创建时间戳（Unix 毫秒值）',
    `updated_at` BIGINT UNSIGNED NOT NULL COMMENT '更新时间戳（Unix 毫秒值）',
    UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息表';

DROP TABLE IF EXISTS `biz_info`;
CREATE TABLE `biz_info` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `biz_key` VARCHAR(64)  NOT NULL COMMENT '业务 key，用于识别业务方身份',
    `biz_secret` VARCHAR(128) NOT NULL COMMENT '业务密钥，用于认证',
    `biz_name` VARCHAR(128) NOT NULL COMMENT '业务名',
    `created_at` BIGINT UNSIGNED NOT NULL COMMENT '创建时间戳（Unix 毫秒值）',
    `updated_at` BIGINT UNSIGNED NOT NULL COMMENT '更新时间戳（Unix 毫秒值）',
    UNIQUE KEY `uk_biz_key` (`biz_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='业务方信息表';
