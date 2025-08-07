-- 系统用户信息表
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

COMMENT ON TABLE sys_user IS '系统用户信息表';
COMMENT ON COLUMN sys_user.id IS 'id';
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
