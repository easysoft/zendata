DROP TABLE IF EXISTS `by_field_type`;

CREATE TABLE `by_field_type` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT,
`f_bit` bit(1) NOT NULL,
`f_tinyint` tinyint DEFAULT NULL,
`f_smallint` smallint DEFAULT NULL,
`f_mediumint` mediumint DEFAULT NULL,
`f_int` int DEFAULT NULL,
`f_bigint` bigint DEFAULT NULL,
`f_float` float DEFAULT NULL,
`f_double` double DEFAULT NULL,
`f_decimal` decimal(10,2) DEFAULT NULL,
`f_char` char(10) DEFAULT NULL,
`f_tinytext` tinytext,
`f_text` text,
`f_mediumtext` mediumtext,
`f_longtext` longtext,
`f_tinyblob` tinyblob,
`f_blob` blob,
`f_mediumblob` mediumblob,
`f_longblob` longblob,
`f_binary` binary(10) DEFAULT NULL,
`f_varbinary` varbinary(1000) DEFAULT NULL,
`f_date` date DEFAULT NULL,
`f_time` time DEFAULT NULL,
`f_year` year DEFAULT NULL,
`f_datetime` datetime DEFAULT NULL,
`f_timestamp` timestamp(6) NULL DEFAULT NULL,
`f_enum` enum('a','b','c') DEFAULT NULL,
`f_set` set('a','b','c') DEFAULT NULL,
PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
