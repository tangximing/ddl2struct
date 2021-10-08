# Welcome to ddl2struct !!!
> Generate Golang Struct File From DDL.

## Install
```
$ go get github.com/tangximing/ddl2struct
```

## Usage
```
$ ddl2struct -h
generate golang struct file from ddl

Usage:
 ddl2struct [flags]

Flags:
  -d, --dir string       golang dir to generate
  -h, --help             help for ddl2struct
  -p, --package string   golang package to generate
  -s, --sql string       ddl sql file path
```

## Example
```
$ cat ./test/table.sql
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

$ ddl2struct -s ./test/table.sql -d ./test/model -p model

$ ll ./test/
drwxr-xr-x   4 ximing  staff  128 Oct  8 21:30 ./
drwxr-xr-x  14 ximing  staff  448 Oct  8 21:30 ../
drwxr-xr-x   4 ximing  staff  128 Oct  8 21:30 model/
-rw-r--r--   1 ximing  staff  950 Oct  8 21:16 table.sql

$ cat ./test/model/ddl_table_1.go
package model

import (
        "time"
)

type Enum1Enum int8

const (
        // 1 - 枚举值1
        Enum11 Enum1Enum = 1
        // 2 - 枚举值2
        Enum12 Enum1Enum = 2
        // 100 - 枚举值100
        Enum1100 Enum1Enum = 100
)

type DdlTable1 struct {
        // 主键ID
        Id int `gorm:"column:id" json:"id"`
        // 枚举1
        Enum1 Enum1Enum `gorm:"column:enum1" json:"enum_1"`
        // 字段1
        Field1 string `gorm:"column:field1" json:"field_1"`
        // 字段2
        Field2 string `gorm:"column:field2" json:"field_2"`
        // 创建时间
        Ctime time.Time `gorm:"column:ctime" json:"ctime"`
        // 修改时间
        Mtime time.Time `gorm:"column:mtime" json:"mtime"`
}
```
