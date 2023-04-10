// Package kokomi_wiki 原神查询功能
package kokomi

import (
	"encoding/json"
	"fmt"
	www "net/url"
	"os"

	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const ( //~标记已实现
	url1 = "https://ghproxy.com/https://raw.githubusercontent.com/CMHopeSunshine/GenshinWikiMap/master/results/character_map/%v.jpg" //~角色图鉴
	url2 = "https://ghproxy.com/https://raw.githubusercontent.com/CMHopeSunshine/LittlePaimonRes/main/genshin_guide/guide/%v.jpg"    //~角色攻略
	url3 = "https://ghproxy.com/https://raw.githubusercontent.com/Nwflower/genshin-atlas/master%v"                                   //~角色材料~七圣召唤卡~特产图鉴[已经替换]~武器图鉴
	url4 = "https://ghproxy.com/https://raw.githubusercontent.com/CMHopeSunshine/LittlePaimonRes/main/genshin_guide/curve/%v.jpg"    //~收益曲线
	url5 = "https://ghproxy.com/https://raw.githubusercontent.com/CMHopeSunshine/LittlePaimonRes/main/genshin_guide/panel/%v.jpg"    //~参考面板
	url6 = "https://map.minigg.cn/map/get_map?resource_name=%v&is_cluster=false"                                                     //地图资源截图
	url7 = "https://ghproxy.com/https://raw.githubusercontent.com/Nwflower/genshin-atlas/master/artifact/%v.png"                     //圣遗物图鉴
	url8 = "https://ghproxy.com/https://raw.githubusercontent.com/CMHopeSunshine/GenshinWikiMap/master/results/monster_map/%v.jpg"   //~原魔图鉴
)

func init() { // 主函数
	en := control.Register("kokomi_wiki", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "原神wiki查询",
		Help: "原神wiki查询功能\n" +
			"- #xxx材料/培养[角色培养材料查询]\n" +
			"- #xxx特产/位置[地区特产查询]\n" +
			"- #xxx收益[角色收益曲线查询]\n" +
			"- #xxx参考[角色参考面板查询]\n" +
			"- #xxx查卡[七圣召唤查卡]\n" +
			"- #xxx攻略[角色攻略查询]\n" +
			"- #xxx原魔[原魔图鉴查询]\n" +
			"- #xxx武器[武器图鉴查询]\n" +
			"- #xxx图鉴[角色图鉴查询]",
	})
	en.OnRegex(`^(?:#|＃)\s?(\D+)(查卡|七圣|培养|材料|特产|位置|武器|图鉴|收益|参考|攻略|原魔)\s?(\d)?`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		keys := ctx.State["regex_matched"].([]string)[2] //匹配种类
		word := ctx.State["regex_matched"].([]string)[1] //关键字
		//num := ctx.State["regex_matched"].([]string)[3]  //可能存在的选择数
		var url, k string //匹配链接
		var paths Wikimap
		t, err := os.ReadFile("plugin/kokomi/data/json/path.json") //获取文件
		if err != nil {
			ctx.SendChain(message.Text("获取路径文件失败", err))
			return
		}
		_ = json.Unmarshal(t, &paths)
		switch keys {
		case "查卡", "七圣": //七圣召唤
			url = url3
			k = paths.Card[word]
		case "培养", "材料": //角色素材
			url = url3
			wife := GetWifeOrWq("wife")
			swifeid := wife.Findnames(word)
			if swifeid == "" {
				ctx.SendChain(message.Text("-请输入角色全名" + Postfix))
				return
			}
			word = wife.Idmap(swifeid)
			if word == "" {
				ctx.SendChain(message.Text("Idmap数据缺失"))
				return
			}
			k = paths.Matera[word]
		case "特产", "位置": //区域特产
			url = url6
			k = www.QueryEscape(word)
		case "武器": //武器图鉴
			url = url3
			wq := GetWifeOrWq("wq")
			word = wq.Findnames(word)
			if word == "" {
				ctx.SendChain(message.Text("-请输入武器全名" + Postfix))
				return
			}
			k = paths.Weapon[word]
		case "收益": //收益曲线
			url = url4
			wife := GetWifeOrWq("wife")
			swifeid := wife.Findnames(word)
			if swifeid == "" {
				ctx.SendChain(message.Text("-请输入角色全名" + Postfix))
				return
			}
			k = wife.Idmap(swifeid)
			if k == "" {
				ctx.SendChain(message.Text("Idmap数据缺失"))
				return
			}
		case "参考": //参考面板
			url = url5
			wife := GetWifeOrWq("wife")
			swifeid := wife.Findnames(word)
			if swifeid == "" {
				ctx.SendChain(message.Text("-请输入角色全名" + Postfix))
				return
			}
			k = wife.Idmap(swifeid)
			if k == "" {
				ctx.SendChain(message.Text("Idmap数据缺失"))
				return
			}
		case "攻略": //角色攻略
			url = url2
			wife := GetWifeOrWq("wife")
			swifeid := wife.Findnames(word)
			if swifeid == "" {
				ctx.SendChain(message.Text("-请输入角色全名" + Postfix))
				return
			}
			k = wife.Idmap(swifeid)
			if k == "" {
				ctx.SendChain(message.Text("Idmap数据缺失"))
				return
			}
		case "原魔": //原魔图鉴
			url = url8
			k = word
		case "图鉴": //角色/武器图鉴
			url = url1
			wife := GetWifeOrWq("wife")
			swifeid := wife.Findnames(word)
			if swifeid != "" {
				k = wife.Idmap(swifeid)
				if k == "" {
					ctx.SendChain(message.Text("Idmap数据缺失"))
					return
				}
			} else { //未找到角色,开始匹配武器
				url = url3
				wq := GetWifeOrWq("wq")
				word = wq.Findnames(word)
				if word == "" {
					ctx.SendChain(message.Text("-未找到信息" + Postfix))
					return
				}
				k = paths.Weapon[word]
			}
		}
		if k == "" {
			ctx.SendChain(message.Text("-未找到信息" + Postfix))
			return
		}

		data, err := web.GetData(fmt.Sprintf(url, k))
		if err != nil {
			ctx.SendChain(message.Text("-获取图片失败惹", err))
			return
		}
		ctx.SendChain(message.ImageBytes(data))
	})
}
