-- 创建主从复制的账号 - 用户slave，密码slave
CREATE USER 'slave'@'%' IDENTIFIED BY 'slave';
-- 授权主从复制
GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'slave'@'%';
-- 刷新权限
FLUSH PRIVILEGES;