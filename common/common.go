package common

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type DbConfig struct {
	Common map[string]string   `yaml:"common"`
	Db     []map[string]string `yaml:"db"`
}

func GetYamlMapCfg(fileName string, index ...string) any {
	//	返回值
	var ret any

	// 读取文件所有内容装到 []byte 中
	cfgFile := CONFIG_PATH + `/` + fileName + ".yaml"
	bytes, err := os.ReadFile(cfgFile)
	if err != nil {
		Logger.ErrorLog("GET_YAML_CFG_ERROR:" + cfgFile)
	}
	// 创建配置文件的结构体
	var conf map[string]any
	// 调用 Unmarshall 去解码文件内容
	// 注意要穿配置结构体的指针进去
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		Logger.ErrorLog("YAML_UNMARSHAL_ERROR:" + fmt.Sprintln(err))
	}

	for _, i := range index {
		var value any
		var ok bool

		if value, ok = conf[i]; !ok {
			Logger.ErrorLog("GET_YAML_CFG_INDEX_ERROR:" + fmt.Sprintln(err, "INDEX:", i))
		}

		// 是否读完
		if conf, ok = value.(map[string]any); ok {
			ret = conf
		} else {
			ret = value
		}
	}

	return ret
}

func GetDb(sevId int) {
	dbId := fmt.Sprintln("db", sevId)
	GetYamlMapCfg("mysqlCfg", "mysql", "game", dbId)
}
