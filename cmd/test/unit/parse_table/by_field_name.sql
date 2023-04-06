DROP TABLE IF EXISTS `by_field_name`;

CREATE TABLE `by_field_name` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(45) DEFAULT NULL,
  `telphone` varchar(45) DEFAULT NULL,
  `mobilephone` varchar(45) DEFAULT NULL,
  `email` varchar(45) DEFAULT NULL,
  `url` varchar(45) DEFAULT NULL,
  `ip` varchar(45) DEFAULT NULL,
  `macaddress` varchar(45) DEFAULT NULL,
  `creditcard` varchar(45) DEFAULT NULL,
  `idcard` varchar(45) DEFAULT NULL,
  `token` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
