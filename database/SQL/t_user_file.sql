CREATE TABLE IF NOT EXISTS t_user_file(
    id SERIAL PRIMARY KEY,
    email VARCHAR(64) NOT NULL,
    sha512 CHAR(128) NOT NULL,
    file_name varchar(256) NOT NULL,
    create_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    status BOOLEAN DEFAULT TRUE
);


CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_user_file ON t_user_file (email, sha512, file_name);
CREATE INDEX  IF NOT EXISTS idx_user_file_status ON t_user_file (email, status);

COMMENT ON INDEX idx_unique_user_file IS '用户文件表唯一索引，保证一个用户同一文件只能上传一次';
COMMENT ON INDEX idx_user_file_status IS '用户文件表状态索引';

COMMENT ON TABLE t_user_file IS '用户文件表';
COMMENT ON COLUMN t_user_file.id IS '主键';
COMMENT ON COLUMN t_user_file.email IS '邮箱与用户表关联';
COMMENT ON COLUMN t_user_file.sha512 IS '文件sha512';
COMMENT ON COLUMN t_user_file.file_name IS '用户定义的文件名';
COMMENT ON COLUMN t_user_file.create_time IS '创建时间';
COMMENT ON COLUMN t_user_file.update_time IS '更新时间';
COMMENT ON COLUMN t_user_file.status IS '状态';

-- 创建更新时间触发器
CREATE TRIGGER t_user_file_update_time_trigger
    BEFORE UPDATE ON t_user_file
    FOR EACH ROW
EXECUTE FUNCTION update_update_time();

