// Package kokomi_wiki 原神查询功能
package kokomi

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	url1  = "https://raw.githubusercontents.com/CMHopeSunshine/GenshinWikiMap/master/results/character_map/%v.jpg" //角色图鉴
	url2  = "https://raw.githubusercontents.com/CMHopeSunshine/LittlePaimonRes/main/genshin_guide/guide/%v.jpg"    //角色攻略
	url3  = "https://raw.githubusercontents.com/Nwflower/genshin-atlas/master/material%20for%20role/%v.png"        //角色材料
	url4  = "https://raw.githubusercontents.com/CMHopeSunshine/LittlePaimonRes/main/genshin_guide/curve/%v.jpg"    //收益曲线
	url5  = "https://raw.githubusercontents.com/CMHopeSunshine/LittlePaimonRes/main/genshin_guide/panel/%v.jpg"    //参考面板
	url6  = "https://raw.githubusercontents.com/Nwflower/genshin-atlas/master/weapon/%v.png"                       //武器图鉴
	url7  = "https://raw.githubusercontents.com/Nwflower/genshin-atlas/master/artifact/%v.png"                     //圣遗物图鉴
	url8  = "https://raw.githubusercontents.com/CMHopeSunshine/GenshinWikiMap/master/results/monster_map/%v.jpg"   //原魔图鉴
	url9  = "https://raw.githubusercontents.com/Nwflower/genshin-atlas/master/specialty/%v.png"                    //特产图鉴
	url10 = "https://ghproxy.com/https://raw.githubusercontent.com/Nwflower/genshin-atlas/master%v"                //七圣召唤卡
)

func init() { // 主函数
	en := control.Register("kokomi_wiki", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "原神wiki查询",
		Help: "原神wiki查询功能\n" +
			"- -#查卡xxx[七圣召唤查卡]\n" +
			"- xxx\n" +
			"- xxx\n" +
			"- xxx",
	})
	en.OnPrefix("#").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		str := ctx.State["args"].(string) // 获取key
		keys := str[0:6]                  //匹配种类
		word := str[6:]                   //关键字
		var url, k string                 //匹配链接
		var paths Wikimap
		t, err := os.ReadFile("plugin/kokomi/data/json/path.json") //获取文件
		if err != nil {
			ctx.SendChain(message.Text("获取路径文件失败", err))
			return
		}
		_ = json.Unmarshal(t, &paths)
		switch keys {
		case "查卡":
			url = url10
			k = paths.Card[word]
		}
		if k == "" {
			ctx.SendChain(message.Text("未找到信息呜"))
			return
		}
		data, err := web.GetData(fmt.Sprintf(url, k))
		if err != nil {
			ctx.SendChain(message.Text("获取图片失败惹", err))
			return
		}
		ctx.SendChain(message.ImageBytes(data))
	})
}
