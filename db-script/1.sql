-- 可用的ip池
INSERT INTO `lb_pool` (`type`,`ip`,`start_port`,`end_port`,`domain_regex`) VALUES ( 1, '106.74.152.45',1, 20000, '[0-9a-z]{1,2}.htres.cn')
INSERT INTO `lb_pool` (`type`,`ip`,`start_port`,`end_port`,`domain_regex`) VALUES ( 1, '106.74.152.44',1, 20000, '[0-9a-z]{1,2}.htres.cn')
INSERT INTO `lb_pool` (`type`,`ip`,`start_port`,`end_port`,`domain_regex`) VALUES ( 1, '106.74.152.39',1, 20000, '[0-9a-z]{1,2}.htres.cn')
INSERT INTO `lb_pool` (`type`,`ip`,`start_port`,`end_port`,`domain_regex`) VALUES ( 1, '106.74.152.35',1, 20000, '[0-9a-z]{1,2}.htres.cn')
INSERT INTO `lb_pool` (`type`,`ip`,`start_port`,`end_port`,`domain_regex`) VALUES ( 1, '106.74.152.33',1, 20000, '[0-9a-z]{1,2}.htres.cn')
INSERT INTO `lb_pool` (`type`,`ip`,`start_port`,`end_port`,`domain_regex`) VALUES ( 1, '106.74.152.32',1, 20000, '[0-9a-z]{1,2}.htres.cn')

-- 负载均衡机器
INSERT INTO `lb_agent`(`lb_pool_id`,`type`,`description`) VALUES (1, 0, "-");
INSERT INTO `lb_agent`(`lb_pool_id`,`type`,`description`) VALUES (2, 0, "-");
INSERT INTO `lb_agent`(`lb_pool_id`,`type`,`description`) VALUES (3, 0, "-");
INSERT INTO `lb_agent`(`lb_pool_id`,`type`,`description`) VALUES (4, 0, "-");
INSERT INTO `lb_agent`(`lb_pool_id`,`type`,`description`) VALUES (5, 0, "-");
INSERT INTO `lb_agent`(`lb_pool_id`,`type`,`description`) VALUES (6, 0, "-");

-- default命名空间目前正在使用的端口
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES(0, 'default', 'chart仓库', 0, '106.74.152.45', 8000, null);

INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES(0, 'default', 'moli日志处理kibana', 0, '106.74.152.45', 12901, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES(0, 'default', 'moli日志处理kafka-manager', 0, '106.74.152.45', 12902, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES(0, 'default', 'moli日志处理openrest-1', 0, '106.74.152.45', 12080, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES(0, 'default', 'moli日志处理openrest-2', 0, '106.74.152.45', 9312, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES(0, 'default', 'moli日志处理openrest-3', 0, '106.74.152.45', 9315, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES(0, 'default', 'moli日志处理openrest-4', 0, '106.74.152.45', 9900, null);

INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES(0, 'default', 'springboot-svc', 0, '106.74.152.45', 10008, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES(0, 'default', 'nic-80', 0, '106.74.152.45', 80, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES(0, 'default', 'nic-443', 0, '106.74.152.45', 443, null);

-- uc项目迁移正在使用的端口
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'uc-user', 0, '106.74.152.45', 18996, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'uc-org', 0, '106.74.152.45', 18995, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'uc-common', 0, '106.74.152.45', 18994, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'uc-auth', 0, '106.74.152.45', 18993, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'uc-token', 0, '106.74.152.45', 18992, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'uc-oauth', 0, '106.74.152.45', 18991, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'uc-sub-system', 0, '106.74.152.45', 18990, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'binding', 0, '106.74.152.45', 18989, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'sync', 0, '106.74.152.45', 18988, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'appclient', 0, '106.74.152.45', 18987, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'appgroup', 0, '106.74.152.45', 18986, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'application', 0, '106.74.152.45', 18985, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'uct01', 0, '106.74.152.45', 18984, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'appapi', 0, '106.74.152.45', 18983, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'authentication', 0, '106.74.152.45', 18982, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'authorization', 0, '106.74.152.45', 18981, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'orglevel', 0, '106.74.152.45', 18980, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'permission', 0, '106.74.152.45', 18979, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cast', 0, '106.74.152.45', 18978, null);

-- 云端业务工作室项目迁移使用的端口
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18800, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18801, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18802, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18803, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18804, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18805, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18806, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18807, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18808, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18809, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18810, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18811, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18812, null);
INSERT INTO `lb_record` (`type`,`owner`,`name`,`status`,`ip`,`port`,`domain`) VALUES (0, '10000077913330', 'cm-', 0, '106.74.152.45', 18813, null);