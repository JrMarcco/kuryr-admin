DROP TYPE IF EXISTS user_type_enum;
CREATE TYPE user_type_enum AS ENUM ('administrator', 'operator');

DROP TYPE IF EXISTS biz_type_enum;
CREATE TYPE biz_type_enum AS ENUM ('individual', 'organization');
