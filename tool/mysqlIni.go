package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var userName = "root"
var password = "&2xs5CzX00ymRmEl"
var host = "127.0.0.1"
var port = "63763"

var dsn string = userName + ":" + password + "@tcp(" + host + ":" + port + ")/?charset=utf8mb4"

func main() {
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("error")
		os.Exit(1)
	}
	var dbName string
	for i := 51; i <= 100; i++ {
		dbName = fmt.Sprintf("%s%02d ", "hdyx_game_", i)

		// 创建数据库
		sql := "CREATE DATABASE IF NOT EXISTS " + dbName + " CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"
		fmt.Println(db.Exec(sql))

		// 切换到目标数据库
		fmt.Println(db.Exec("USE " + dbName + ";"))

		sql = `CREATE TABLE IF NOT EXISTS user(
				uid BIGINT NOT NULL COMMENT 'uid',
				sex int(10) DEFAULT '0' COMMENT '男=1，女=2',
				fortune int(10) DEFAULT '0' COMMENT '财富值',
				isGame bool DEFAULT '0' COMMENT '是否在游戏中',
				loginIp varchar(16) DEFAULT '0' COMMENT '登录IP',
				playTime  BIGINT DEFAULT '0' COMMENT '累计玩游戏时间',
				lastLogin int(12) DEFAULT '0' COMMENT '最后一次登录时间'
			)`

		fmt.Println(db.Exec(sql))

		sql = `CREATE TABLE IF NOT EXISTS act(
					uid bigint(64) DEFAULT '0',
					actId int(10) DEFAULT '0' COMMENT '活动类型编号',
					tJson longtext COMMENT '活动数据',
					createTime datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
					updateTime datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
					UNIQUE KEY idx_uid (uid,actid)
				) ENGINE=InnoDB DEFAULT CHARSET=utf8;`

		fmt.Println(db.Exec(sql))
	}
}
