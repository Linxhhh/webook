-- 创建数据库
CREATE DATABASE IF NOT EXISTS webook;

-- 创建复制用户
CREATE USER 'replicate'@'%' IDENTIFIED BY 'replicate_password';

-- 授予复制权限
GRANT REPLICATION SLAVE ON *.* TO 'replicate'@'%';

-- 刷新权限
FLUSH PRIVILEGES;