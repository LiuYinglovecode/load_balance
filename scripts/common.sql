-- 测试环境mysql地址
user: alb_rw
password: 72r8pmqSViA=
host: 10.153.51.186
port: 32168
-- 181上连接
mysql -h 10.153.51.186 -u alb_rw --port 32168 -p72r8pmqSViA= alb

-- 插入Lbpool
INSERT INTO `alb`.`lb_pool`
(`type`,
`ip`,
`start_port`,
`end_port`,
`domain_regex`)
VALUES
( 1, '106.74.152.34',10000, 20000, '[0-9a-z]{1,2}.htres.cn')

-- 插入lbr
-- type: 0是ip 1是domain
INSERT INTO `alb`.`lb_record`
(`type`,
`owner`,
`name`,
`status`,
`ip`,
`port`,
`domain`)
VALUES
(0, '50000003590000', '-', 0, '106.74.152.34', 10118, null);

-- 插入lbagent
INSERT INTO `alb`.`lb_agent`
(`lb_pool_id`,
`type`,
`description`)
VALUES
(1, 0, "-");

-- 查看etcd
-- ETCDCTL_API=3 etcdctl --endpoints=https://10.153.51.208:2379 --cacert=ca.pem  --cert=kubernetes.pem --key=kubernetes-key.pem get / --prefix --keys-only
-- 在负载均衡机器上手动启动docker
-- docker run -d --name adc-haproxy-106.74.152.45_80 --network host -v "$(pwd)"/106.74.152.45_80-haproxy.cfg:/usr/local/etc/haproxy/haproxy.cfg hub.htres.cn/pub/haproxy:1.8
