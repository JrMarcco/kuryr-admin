DROP TABLE IF EXISTS sys_user;
CREATE TABLE sys_user (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(128) NOT NULL,
    password VARCHAR(128) NOT NULL,
    real_name VARCHAR(64) NOT NULL,
    user_type user_type_enum NOT NULL,
    biz_id BIGINT NOT NULL DEFAULT 0,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    CONSTRAINT uk_email UNIQUE (email)
);

COMMENT ON TABLE sys_user IS '用户信息表';
COMMENT ON COLUMN sys_user.email IS '邮箱';
COMMENT ON COLUMN sys_user.password IS '密码（请存储哈希后的值）';
COMMENT ON COLUMN sys_user.real_name IS '用户姓名';
COMMENT ON COLUMN sys_user.user_type IS '用户类型: administrator-系统管理员 / operator-业务操作员';
COMMENT ON COLUMN sys_user.biz_id IS '关联的业务方 id, 当 user_type 为 operator 时不应为 0';
COMMENT ON COLUMN sys_user.created_at IS '创建时间戳 ( Unix 毫秒值 )';
COMMENT ON COLUMN sys_user.updated_at IS '更新时间戳 ( Unix 毫秒值 )';

-- 插入初始用户数据
INSERT INTO sys_user (
    email,
    password,
    real_name,
    user_type,
    biz_id,
    created_at,
    updated_at
) VALUES (
    'jrmarcco@gmail.com',
    '$2a$10$besICPqbCRWOocqlsaKXV.rniGRyCNPLHeFT.osXbhgisW4XSW/um',
    'jrmarcco',
    'administrator',
    0,
    EXTRACT(EPOCH FROM NOW()) * 1000,
    EXTRACT(EPOCH FROM NOW()) * 1000
);

DROP TABLE IF EXISTS biz_info;
CREATE TABLE biz_info (
    id BIGSERIAL PRIMARY KEY,
    biz_type biz_type_enum NOT NULL,
    biz_key VARCHAR(64) NOT NULL,
    biz_secret VARCHAR(128) NOT NULL,
    biz_name VARCHAR(128) NOT NULL,
    contact varchar(64) NOT NULL,
    contact_email varchar(128) NOT NULL,
    creator_id BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    CONSTRAINT uk_biz_key UNIQUE (biz_key)
);

-- 添加表注释
COMMENT ON TABLE biz_info IS '业务信息表';
COMMENT ON COLUMN biz_info.biz_type IS '业务类型';
COMMENT ON COLUMN biz_info.biz_key IS '业务 key 用于识别业务方身份';
COMMENT ON COLUMN biz_info.biz_secret IS '业务密钥 用于认证';
COMMENT ON COLUMN biz_info.biz_name IS '业务名';
COMMENT ON COLUMN biz_info.contact IS '业务联系人';
COMMENT ON COLUMN biz_info.contact_email IS '联系人邮箱';
COMMENT ON COLUMN biz_info.creator_id IS '创建人 id';
COMMENT ON COLUMN biz_info.created_at IS '创建时间戳 ( Unix 毫秒值 )';
COMMENT ON COLUMN biz_info.updated_at IS '更新时间戳 ( Unix 毫秒值 )';
