
CREATE TABLE `biz_project`   (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `name` varchar,
    `desc` longtext,
    `is_default` longtext,
    `disabled` tinyint(1) DEFAULT NULL,
    `disabled_at` datetime(3) DEFAULT NULL,
    `path` longtext,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `biz_task` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `disabled` tinyint(1) DEFAULT '0',
    `name` longtext,
    `project_id` bigint(20) unsigned NOT NULL,
    `project_name` longtext,
    `disabled_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY  `fk_project_id_idx`  (`project_id`),
    CONSTRAINT `fk_project_id` FOREIGN KEY (`project_id`) REFERENCES `biz_project` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
