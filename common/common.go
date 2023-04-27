package common

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func GetYamlCfg(fileName string, index ...string) map[string]any {
	// 读取文件所有内容装到 []byte 中
	bytes, err := os.ReadFile(CONFIG_PATH + `/` + fileName + ".yaml")
	if err != nil {
		Logger.errorLog("GET_YAML_CFG_ERROR:" + fileName)
	}
	// 创建配置文件的结构体
	var conf map[string]any
	// 调用 Unmarshall 去解码文件内容
	// 注意要穿配置结构体的指针进去
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		Logger.errorLog("YAML_UNMARSHAL_ERROR:" + fmt.Sprintln(err))
	}

	for _, i := range index {
		var value any
		var ok bool
		if value, ok = conf[i]; !ok {
			Logger.errorLog("GET_YAML_CFG_INDEX_ERROR:" + fmt.Sprintln(err) + "INDEX:" + i)
		}

		conf = value.(map[string]any)
	}

	return conf
}
