
-- workflow models
CREATE TABLE `wf_workflows` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` int(11) unsigned NOT NULL default '0',
  `name` varchar(80) NOT NULL default '',
  `desc` text,
  `start_job_id` int(11) unsigned NOT NULL default '0' COMMENT '工作流起始job',
  `is_valid` char(1) NOT NULL default 'N' COMMENT '工作流校验状态',
  `errors` text COMMENT '工作流校验错误信息',
  `start_at` DATETIME default NULL COMMENT '工作流开始生效时间',
  `end_at` DATETIME default NULL COMMENT '工作流结束生效时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` VARCHAR(32) NOT NULL DEFAULT '' ,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` VARCHAR(32) NOT NULL DEFAULT '' ,
 PRIMARY KEY  (`id`),
 KEY `start_job_idx` (`app_id`, `start_job_id`) USING BTREE
) ENGINE=InnoDB COMMENT = '工作流表';

CREATE TABLE `wf_places` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` int(11) unsigned NOT NULL default '0',
  `workflow_id` int(11) unsigned NOT NULL default '0',
  `type` char(1) NOT NULL default '5' COMMENT '类型：1 开始类型，5 中间类型，9 结束类型',
  `name` varchar(80) NOT NULL default '',
  `desc` text,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY  (`id`),
  KEY `workflow_idx` (`app_id`, `workflow_id`) USING BTREE
) ENGINE=InnoDB;

CREATE TABLE `wf_transitions` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` int(11) unsigned NOT NULL default '0',
  `workflow_id` int(11) unsigned NOT NULL default '0',
  `name` varchar(80) NOT NULL default '转换器名称',
  `desc` text,
  `trigger` varchar(4) NOT NULL default 'USER' COMMENT '触发类型：USER（用户）, AUTO（自动）,MSG（消息、事件）,TIME（定时器）',
  `time_limit` int(11) unsigned default NULL COMMENT '转换器有效期 单位：分钟',
  `job_id` int(11) unsigned NOT NULL default '0',
  `role_id` int(11) unsigned NOT NULL default '0',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY  (`id`),
  KEY `workflow_idx` (`app_id`, `workflow_id`) USING BTREE
) ENGINE=InnoDB;

CREATE TABLE `wf_arcs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` int(11) unsigned NOT NULL default '0',
  `workflow_id` int(11) unsigned NOT NULL default '0',
  `transition_id` int(11) unsigned NOT NULL default '0',
  `place_id` int(11) unsigned NOT NULL default '0',
  `direction` char(3) NOT NULL default '' COMMENT '弧方向：IN，OUT',
  `type` varchar(10) NOT NULL default 'SEQ' COMMENT '弧类型：SEQ，Explicit OR split，Implicit OR split，OR join，AND split，AND join',
  `condition` text COMMENT "弧附加条件，仅当arc_type='Explicit OR split' 时才有效",
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY  (`id`),
  KEY `direction_idx` (`app_id`, `workflow_id`,`direction`) USING BTREE,
  KEY `transition_dix` (`app_id`, `workflow_id`,`transition_id`,`direction`) USING BTREE
) ENGINE=InnoDB;

CREATE TABLE `wf_cases` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` int(11) unsigned NOT NULL default '0',
  `workflow_id` int(11) unsigned NOT NULL default '0',
  `context` text,
  `status` char(2) NOT NULL default 'OP' COMMENT 'case状态：OP（open），CL（Closed），SU （Suspended），CA（Cancelled）',
  `start_at` datetime NOT NULL default '0000-00-00 00:00:00',
  `end_at` datetime default NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` VARCHAR(32) NOT NULL DEFAULT '' ,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` VARCHAR(32) NOT NULL DEFAULT '' ,
  PRIMARY KEY  (`id`),
  KEY `workflow_idx` (`app_id`, `workflow_id`) USING BTREE
) ENGINE=InnoDB;

CREATE TABLE `wf_tokens` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` int(11) unsigned NOT NULL default '0',
  `workflow_id` int(11) unsigned NOT NULL default '0',
  `case_id` int(11) unsigned NOT NULL default '0',
  `place_id` int(11) unsigned NOT NULL default '0',
  `context` varchar(255) NOT NULL default '',
  `status` varchar(4) NOT NULL default 'FREE',
  `enabled_date` datetime NOT NULL default '0000-00-00 00:00:00',
  `cancelled_date` datetime default NULL,
  `consumed_date` datetime default NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY  (`id`),
  KEY `place_idx` (`app_id`,`workflow_id`,`case_id`,`place_id`) USING BTREE
) ENGINE=InnoDB;

CREATE TABLE `wf_workitems` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` int(11) unsigned NOT NULL default '0',
  `workflow_id` int(11) unsigned NOT NULL default '0',
  `case_id` int(11) unsigned NOT NULL default '0',
  `transition_id` int(11) unsigned NOT NULL default '0',
  `transition_trigger` varchar(4) NOT NULL default 'USER' COMMENT '触发类型：USER（用户）, AUTO（自动）,MSG（消息、事件）,TIME（定时器）',
  `job_id` varchar(40) NOT NULL default '',
  `context` text,
  `status` char(2) NOT NULL default 'EN' COMMENT '状态：EN（Enabled）， IP（In Progress），CA（Cancelled），FI（Finished）',
  `enabled_date` datetime default NULL,
  `cancelled_date` datetime default NULL,
  `finished_date` datetime default NULL,
  `deadline` datetime default NULL,
  `role_id` int(11) unsigned NOT NULL default '0',
  `user` VARCHAR(32) NOT NULL DEFAULT '' ,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY  (`id`),
  KEY `transition_idx` (`app_id`, `workflow_id`,`transition_id`) USING BTREE,
  KEY `job_idx` (`app_id`, `job_id`) USING BTREE
) ENGINE=InnoDB;
