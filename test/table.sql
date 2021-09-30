CREATE TABLE `ddl_table1` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `enum1` tinyint(4) NOT NULL DEFAULT 1 COMMENT '枚举1。枚举（1：枚举值1，2：枚举值2，100：枚举值100）',
  `field1` varchar(32) NOT NULL DEFAULT '' COMMENT '字段1',
  `field2` varbinary(64) NOT NULL DEFAULT '' COMMENT '字段2',
  `ctime` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
  `mtime` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1281 DEFAULT CHARSET=utf8mb4 COMMENT='表1';

CREATE TABLE `ddl_table2` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `field1` varchar(255) NOT NULL DEFAULT '' COMMENT '字段1',
  `field2` varchar(255) NOT NULL DEFAULT '' COMMENT '字段2',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='表2'