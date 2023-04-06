DROP TABLE IF EXISTS `by_records`;

CREATE TABLE `by_records` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `f1` varchar(100) DEFAULT NULL,
  `f2` varchar(100) DEFAULT NULL,
  `f3` varchar(100) DEFAULT NULL,
  `f4` varchar(100) DEFAULT NULL,
  `f5` varchar(100) DEFAULT NULL,
  `f6` varchar(100) DEFAULT NULL,
  `f7` varchar(100) DEFAULT NULL,
  `f8` varchar(100) DEFAULT NULL,
  `f9` varchar(100) DEFAULT NULL,
  `f10` varchar(100) DEFAULT NULL,
  `f11` varchar(100) DEFAULT NULL,
  `f12` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
