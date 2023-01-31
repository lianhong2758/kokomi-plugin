package kokomi

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
)

var Config config

func init() {
	cache, err := os.ReadFile("plugin/kokomi/config.json")
	if err != nil {
		logrus.Errorln("获取kokomi配置文件错误")
		os.Exit(1)
	}
	err = json.Unmarshal(cache, &Config)
	if err != nil {
		logrus.Errorln("解析kokomi配置文件错误")
		os.Exit(1)
	}
}
