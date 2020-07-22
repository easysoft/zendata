-- sqlite
CREATE TABLE user (
    field_common VARCHAR (500)
);

-- DROP TABLE IF EXISTS `zt_action`;
CREATE TABLE IF NOT EXISTS `zt_action` (
  `id` mediumint(8) unsigned NOT NULL auto_increment,
  `objectType` varchar(30) NOT NULL default '',
  `objectID` mediumint(8) unsigned NOT NULL default '0',
  `product` varchar(255) NOT NULL,
  `project` mediumint(9) NOT NULL,
  `actor` varchar(100) NOT NULL default '',
  `action` varchar(30) NOT NULL default '',
  `date` datetime NOT NULL,
  `comment` text NOT NULL,
  `extra` text NOT NULL,
  `read` enum('0','1') NOT NULL default '0',
  PRIMARY KEY (`id`),
  KEY `date` (`date`),
  KEY `actor` (`actor`),
  KEY `project` (`project`),
  KEY `objectID` (`objectID`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
