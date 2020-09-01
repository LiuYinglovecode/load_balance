CREATE TABLE `lb_pool` (
  `lb_pool_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `type` tinyint(1) NOT NULL DEFAULT '0',
  `ip` varchar(511) DEFAULT NULL,
  `start_port` int(11) DEFAULT NULL,
  `end_port` int(11) DEFAULT NULL,
  `domain_regex` varchar(511) DEFAULT NULL,
  `deleted` tinyint(1) NOT NULL DEFAULT '0',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `modify_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`lb_pool_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='LoadBalancer Pool';

CREATE TABLE `lb_agent` (
  `lb_agent_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `lb_pool_id` bigint(20) NOT NULL,
  `type` tinyint(1) NOT NULL DEFAULT '0',
  `description` varchar(511) DEFAULT NULL,
  `deleted` tinyint(1) NOT NULL DEFAULT '0',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `modify_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`lb_agent_id`),
  KEY `lb_pool_fk` (`lb_pool_id`),
  CONSTRAINT `lb_pool_fk` FOREIGN KEY (`lb_pool_id`) REFERENCES `lb_pool` (`lb_pool_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='LoadBalancer Agent';

CREATE TABLE `lb_record` (
  `lb_record_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `type` tinyint(1) NOT NULL DEFAULT '0',
  `owner` varchar(255) NOT NULL,
  `name` varchar(511) DEFAULT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `ip` varchar(511) DEFAULT NULL,
  `port` int(11) DEFAULT NULL,
  `domain` varchar(511) DEFAULT NULL,
  `deleted` tinyint(1) NOT NULL DEFAULT '0',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `modify_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`lb_record_id`),
  KEY `owner_idx` (`owner`),
  KEY `create_time_idx` (`create_time`),
  KEY `modify_time_idx` (`modify_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='LoadBalancer Record';

CREATE TABLE `lb_request` (
  `lb_request_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `owner` varchar(255) NOT NULL,
  `service` varchar(511) NOT NULL DEFAULT '-',
  `status` varchar(255) DEFAULT NULL,
  `action` tinyint(1) NOT NULL,
  `policy` text COMMENT 'LBPolicy的json格式',
  `finish_time` datetime DEFAULT NULL COMMENT '请求处理结束时间',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` datetime DEFAULT NULL COMMENT '删除时间',
  `note` varchar(511) DEFAULT NULL COMMENT '如果处理失败记录失败原因,或者记录其他需要特别存储的信息',
  PRIMARY KEY (`lb_request_id`),
  KEY `owner_idx` (`owner`),
  KEY `create_time_idx` (`create_time`),
  KEY `update_time` (`update_time`),
  FULLTEXT KEY `note_idx` (`note`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='LoadBalancer Request';
