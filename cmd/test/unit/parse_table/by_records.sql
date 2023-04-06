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

INSERT INTO `by_records` VALUES (1,'462826@qq.com','5139288802098206','3D:F2:C9:A6:B3:4F','a987fbc9-4bed-3078-cf07-9141ba07c9f3','81e36610a43a3cc930c73d6d5ae2cee1','18626203288','0512-65918565','522601198205230053','https://baidu.com','{\"key\":\"value\"}',NULL,NULL);
