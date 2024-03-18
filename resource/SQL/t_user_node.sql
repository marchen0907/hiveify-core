CREATE TYPE NODE_ENUM AS ENUM (
    'file',
	'folder'
);


CREATE TABLE IF NOT EXISTS t_user_node(
    id SERIAL PRIMARY KEY,
    email VARCHAR(64) NOT NULL,
	type NODE_ENUM NOT NULL DEFAULT 'file',
    sha512 CHAR(128) NOT NULL,
	parent_node CHAR(128) NOT NULL DEFAULT 'root',
    name varchar(256) NOT NULL,
    sync BOOLEAN DEFAULT FALSE,
    create_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    status BOOLEAN DEFAULT TRUE
);


CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_user_node ON t_user_node (email, parent_node, name);
CREATE INDEX  IF NOT EXISTS idx_user_node_status ON t_user_node (email, status);

COMMENT ON INDEX idx_unique_user_node IS '文件（夹）唯一索引，保证一个用户同一文件（夹）只能上传一次';
COMMENT ON INDEX idx_user_node_status IS '用户文件表状态索引';
COMMENT ON TABLE t_user_node IS '用户文件（夹）表';
COMMENT ON COLUMN t_user_node.id IS '主键';
COMMENT ON COLUMN t_user_node.email IS '邮箱与用户表关联';
COMMENT ON COLUMN t_user_node.type IS '节点类型';
COMMENT ON COLUMN t_user_node.sha512 IS 'sha512';
COMMENT ON COLUMN t_user_node.parent_node IS '父节点';
COMMENT ON COLUMN t_user_node.name IS '用户定义的文件（夹）名';
COMMENT ON COLUMN t_user_node.sync IS '是否自动同步';
COMMENT ON COLUMN t_user_node.create_time IS '创建时间';
COMMENT ON COLUMN t_user_node.update_time IS '更新时间';
COMMENT ON COLUMN t_user_node.status IS '状态';

-- 创建更新时间触发器
CREATE TRIGGER t_user_node_update_time_trigger
    BEFORE UPDATE ON t_user_node
    FOR EACH ROW
EXECUTE FUNCTION update_update_time();

