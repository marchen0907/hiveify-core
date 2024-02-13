-- noinspection SqlResolveForFile @ language/"plpgsql"

CREATE TABLE t_file (
    id SERIAL PRIMARY KEY,
    sha512 CHAR(128) NOT NULL UNIQUE,
    name VARCHAR(256) NOT NULL,
    size BIGINT DEFAULT 0,
    path VARCHAR(1024),
    create_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    status BOOLEAN DEFAULT true
);

COMMENT ON TABLE t_file IS '文件表';
COMMENT ON COLUMN t_file.id IS '主键，自增长';
COMMENT ON COLUMN t_file.sha512 IS '文件sha512';
COMMENT ON COLUMN t_file.name IS '文件名';
COMMENT ON COLUMN t_file.size IS '文件大小';
COMMENT ON COLUMN t_file.path IS '文件路径';
COMMENT ON COLUMN t_file.create_time IS '文件创建时间';
COMMENT ON COLUMN t_file.update_time IS '文件更新时间';
COMMENT ON COLUMN t_file.status IS '是否删除，true为没有删除';


CREATE INDEX file_status_index ON t_file (status);

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_update_time()
RETURNS TRIGGER AS $$
BEGIN
    NEW.update_time = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建更新时间触发器
CREATE TRIGGER t_file_update_time_trigger
BEFORE UPDATE ON t_file
FOR EACH ROW
EXECUTE FUNCTION update_update_time();

