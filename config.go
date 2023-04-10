package kokomi

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	Config  config
	url     string
	edition string
	Postfix string
)

func init() {
	cache, err := os.ReadFile("plugin/kokomi/config.json")
	if err != nil {
		logrus.Errorln("获取kokomi配置文件错误,请重新配置config.josn文件,若无法解决,请加群解决,678586912")
		os.Exit(1)
	}
	err = json.Unmarshal(cache, &Config)
	if err != nil {
		logrus.Errorln("解析kokomi配置文件错误,请重新配置config.josn文件,若无法解决,请加群解决,678586912")
		os.Exit(1)
	}
	fmt.Print(
		"==========[ZeroBot-Plugin & ", Config.Edition, "]==================",
		"\n\n插件配置加载完成,本插件完全免费,作者尽力维护",
		"\n若出现问题请加群解决,678586912",
		"\n\n============================================================\n\n",
	)
	url = Config.Apis[Config.Apiid]
	edition = Config.Edition
	Postfix = Config.Postfix
}
