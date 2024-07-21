change master to 
    master_host='mysql-master',
    master_port=3306,
    master_user='slave',
    master_password='123456', 
    master_log_file='mysql-bin.000003',
    master_log_pos= 154,
    get_master_public_key=1,
    master_connect_retry=30;

start slave;

show slave status \G;

-- master_log_file 和 master_log_pos 是不确定的
-- 需要在主数据库中通过 show master status; 命令查看