-- noinspection SqlResolveForFile @ language/"plpgsql"

CREATE TABLE t_user(
	id SERIAL PRIMARY KEY,
	name VARCHAR(64) NOT NULL,
	password CHAR(128) NOT NULL,
	email VARCHAR(64) UNIQUE NOT NULL,
	phone VARCHAR(64) UNIQUE,
	email_validated BOOLEAN DEFAULT false,
	phone_validated BOOLEAN DEFAULT false,
	create_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	last_active TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	profile TEXT,
	status SMALLINT
);

COMMENT ON TABLE t_user IS '用户表';
COMMENT ON COLUMN t_user.name IS '用户名';
COMMENT ON COLUMN t_user.password IS '密码';
COMMENT ON COLUMN t_user.email IS '邮箱';
COMMENT ON COLUMN t_user.phone IS '手机';
COMMENT ON COLUMN t_user.email_validated IS '邮箱是否验证';
COMMENT ON COLUMN t_user.phone_validated IS '手机是否验证';
COMMENT ON COLUMN t_user.create_time IS '创建时间';
COMMENT ON COLUMN t_user.last_active IS '上次活跃时间';
COMMENT ON COLUMN t_user.profile IS '用户属性';
COMMENT ON COLUMN t_user.status IS '账户状态';

CREATE INDEX usr_status_index ON t_user (status);


-- 创建更新时间触发器
CREATE TRIGGER t_user_update_time_trigger
BEFORE UPDATE ON t_user
FOR EACH ROW
EXECUTE FUNCTION update_update_time();

-- CREATE TYPE user_status AS ENUM('active', 'disabled','locked','deleted');