// Package kokomi  原神面板v2.1
package kokomi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Coloured-glaze/gg"
	"github.com/FloatTech/floatbox/img/writer"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	//"github.com/golang/freetype"
	"golang.org/x/image/webp"

	//"github.com/FloatTech/zbputils/img"
	"github.com/nfnt/resize"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	url      = "https://enka.minigg.cn/u/%v/__data.json"
	Postfix  = "~" //语句后缀
	edition  = "Created By ZeroBot-Plugin v1.6.1 & kokomi v2.1"
	tu       = "https://api.yimian.xyz/img?type=moe&size=1920x1080"
	NameFont = "plugin/kokomi/data/font/NZBZ.ttf"        // 名字字体
	FontFile = "plugin/kokomi/data/font/HYWH-65W.ttf"    // 汉字字体
	FiFile   = "plugin/kokomi/data/font/tttgbnumber.ttf" // 其余字体(数字英文)
	BaFile   = "plugin/kokomi/data/font/STLITI.TTF"      // 华文隶书版本版本号字体

)

func init() { // 主函数
	en := control.Register("kokomi", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "原神面板查询",
		Help: "- kokomi菜单\n" +
			"- 绑定......(uid)\n" +
			"- 更新面板\n" +
			"- 全部面板\n" +
			"- XX面板\n" +
			"- 删除账号[@xx]",
	})
	en.OnSuffix("面板").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		str := ctx.State["args"].(string) // 获取key
		var wifeid int64
		var allfen = 0.00
		qquid := ctx.Event.UserID
		// 获取uid
		uid := Getuid(qquid)
		suid := strconv.Itoa(uid)
		if uid == 0 {
			ctx.SendChain(message.Text("-未绑定uid" + Postfix))
			return
		}
		//############################################################判断数据更新,逻辑原因不能合并进switch
		if str == "更新" || str == "#更新" {
			es, err := web.GetData(fmt.Sprintf(url, uid)) // 网站返回结果
			if err != nil {
				ctx.SendChain(message.Text("网站获取信息失败", err))
				return
			}
			// 创建存储文件,路径plugin/kokomi/data/js
			file, _ := os.OpenFile("plugin/kokomi/data/js/"+suid+".kokomi", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			_, _ = file.Write(es)
			ctx.SendChain(message.Text("-获取角色面板成功" + Postfix + "\n-请发送 全部面板 查看已展示角色" + Postfix))
			file.Close()
			return
		}
		//############################################################
		// 获取本地缓存数据
		txt, err := os.ReadFile("plugin/kokomi/data/js/" + suid + ".kokomi")
		if err != nil {
			ctx.SendChain(message.Text("-本地未找到账号信息, 请更新面板" + Postfix))
			return
		}

		// 解析
		var alldata Data
		err = json.Unmarshal(txt, &alldata)
		if err != nil {
			ctx.SendChain(message.Text("出现错误捏：", err))
			return
		}
		if len(alldata.PlayerInfo.ShowAvatarInfoList) == 0 {
			ctx.SendChain(message.Text("-请在游戏中打开角色展柜,并将想查询的角色进行展示" + "\n-完成上述操作并等待5分钟后,请使用 更新面板 获取账号信息" + Postfix))
			return
		}
		switch str {
		case "全部", "全部角色", "#全部":
			var msg strings.Builder
			msg.WriteString("您的展示角色为:\n")
			for i := 0; i < len(alldata.PlayerInfo.ShowAvatarInfoList); i++ {
				mmm := Idmap(strconv.Itoa(alldata.PlayerInfo.ShowAvatarInfoList[i].AvatarID), "wife")
				if mmm == "" {
					ctx.SendChain(message.Text("Idmap数据缺失"))
					return
				}
				msg.WriteString(mmm)
				if i < len(alldata.PlayerInfo.ShowAvatarInfoList)-1 {
					msg.WriteByte('\n')
				}
			}
			ctx.SendChain(message.Text(msg.String()))
			return
		default: // 角色名解析为id
			//排除#
			if str[0:1] == "#" {
				str = str[1:]
			}
			//匹配简称/外号
			//str = FindName(str)
			swifeid := Findnames(str, "wife")
			if swifeid == "" {
				ctx.SendChain(message.Text("-请输入角色全名" + Postfix))
				return
			}
			wifeid, _ = strconv.ParseInt(swifeid, 10, 64)
			str = Idmap(swifeid, "wife")
			if str == "" {
				ctx.SendChain(message.Text("Idmap数据缺失"))
				return
			}
		}
		var t = -1
		// 匹配角色
		for i := 0; i < len(alldata.PlayerInfo.ShowAvatarInfoList); i++ {
			if wifeid == int64(alldata.PlayerInfo.ShowAvatarInfoList[i].AvatarID) {
				t = i
			}
		}
		if t == -1 { // 在返回数据中未找到想要的角色
			ctx.SendChain(message.Text("-该角色未展示" + Postfix))
			return
		}

		// 画图
		var height int = 2400
		dc := gg.NewContext(1080, height) // 画布大小
		dc.SetHexColor("#98F5FF")
		dc.Clear() // 背景
		pro, flg := Promap[wifeid]
		if !flg {
			ctx.SendChain(message.Text("匹配角色元素失败"))
			return
		}
		beijing, err := gg.LoadImage("plugin/kokomi/data/pro/" + pro + ".jpg")
		if err != nil {
			ctx.SendChain(message.Text("获取背景失败", err))
			return
		}
		dc.Scale(5/3.0, 5/3.0)
		dc.DrawImage(beijing, -792, 0)
		dc.Scale(3/5.0, 3/5.0)
		dc.SetRGB(1, 1, 1) // 换白色

		//武器图层
		//新建图层,实现阴影
		yinwq := Yinying(340, 180, 16, 0.6)
		// 字图层
		two := gg.NewContext(340, 180)
		if err := two.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		two.SetRGB(1, 1, 1) //白色
		//武器名
		//纠正圣遗物空缺报错的无返回情况
		l := len(alldata.AvatarInfoList[t].EquipList)

		wq := IdforNamemap[alldata.AvatarInfoList[t].EquipList[l-1].Flat.NameTextHash]
		two.DrawString(wq, 150, 50)

		//详细
		two.DrawString("攻击力:", 150, 130)
		two.DrawString("精炼:", 245, 90)
		if err := two.LoadFontFace(FiFile, 30); err != nil { // 字体大小
			panic(err)
		}
		two.DrawString(strconv.FormatFloat(alldata.AvatarInfoList[t].EquipList[l-1].Flat.WeaponStat[0].Value, 'f', 1, 32), 250, 130)
		//Lv,精炼
		var wqjl int
		for m := range alldata.AvatarInfoList[t].EquipList[l-1].Weapon.AffixMap {
			wqjl = m
		}
		two.DrawString("Lv."+strconv.Itoa(alldata.AvatarInfoList[t].EquipList[l-1].Weapon.Level), 150, 90)
		two.DrawString(strconv.Itoa(alldata.AvatarInfoList[t].EquipList[l-1].Weapon.AffixMap[wqjl]+1), 316, 90)
		/*副词条,放不下
		fucitiao, _ := IdforNamemap[alldata.AvatarInfoList[t].EquipList[5].Flat.WeaponStat[1].SubPropId] //名称
		var baifen = "%"
		if fucitiao == "元素精通" {
			baifen = ""
		}
		dc.DrawString(fucitiao+":"+strconv.Itoa(int(alldata.AvatarInfoList[t].EquipList[5].Flat.WeaponStat[1].Value))+baifen, 820, 270)
		*/
		//图片
		tuwq, err := gg.LoadPNG("plugin/kokomi/data/wq/" + wq + ".png")
		if err != nil {
			ctx.SendChain(message.Text("获取武器图标失败", err))
			return
		}
		tuwq = resize.Resize(130, 0, tuwq, resize.Bilinear)
		two.DrawImage(tuwq, 10, 10)
		dc.DrawImage(yinwq, 20, 920)
		dc.DrawImage(two.Image(), 20, 920)

		//圣遗物
		yinsyw := Yinying(340, 350, 16, 0.6)
		for i := 0; i < l-1; i++ {
			// 字图层
			three := gg.NewContext(340, 350)
			if err := three.LoadFontFace(FontFile, 30); err != nil {
				panic(err)
			}
			//字号30,间距50
			three.SetRGB(1, 1, 1) //白色
			sywname := IdforNamemap[alldata.AvatarInfoList[t].EquipList[i].Flat.SetNameTextHash]
			tusyw, err := gg.LoadImage("plugin/kokomi/data/syw/" + sywname + "/" + strconv.Itoa(i+1) + ".webp")
			if err != nil {
				ctx.SendChain(message.Text("获取圣遗物图标失败", err))
				return
			}
			tusyw = resize.Resize(80, 0, tusyw, resize.Bilinear) //缩小
			three.DrawImage(tusyw, 15, 15)
			//圣遗物name
			sywallname := SywNamemap[sywname]
			three.DrawString(sywallname[i], 110, 50)
			//圣遗物属性
			zhuci := StoS(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquaryMainStat.MainPropID) //主词条
			zhucitiao := Ftoone(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquaryMainStat.Value)
			//间隔45,初始145
			var xx, yy, pingfeng float64 //xx,yy词条相对位置,x,y文本框在全图位置
			var x, y int
			xx = 15
			yy = 145
			pingfeng = 0
			//主词条
			three.DrawString("主:"+zhuci, xx, yy)
			if err := three.LoadFontFace(FiFile, 30); err != nil {
				panic(err)
			} //主词条名字
			three.DrawString("+"+zhucitiao+Stofen(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquaryMainStat.MainPropID), 200, yy) //主词条属性
			//算分
			if i > 1 { //不算前两主词条属性
				pingfeng += Countcitiao(str, zhuci, alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquaryMainStat.Value/4)
			}
			//副词条
			p := len(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquarySubStats)
			for k := 0; k < p; k++ {
				switch k {
				case 0:
					yy = 190
				case 1:
					yy = 235
				case 2:
					yy = 280
				case 3:
					yy = 325
				}
				if err := three.LoadFontFace(FontFile, 30); err != nil {
					panic(err)
				}
				three.DrawString(StoS(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquarySubStats[k].SubPropID), xx, yy)
				if err := three.LoadFontFace(FiFile, 30); err != nil {
					panic(err)
				}
				three.DrawString("+"+strconv.FormatFloat(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquarySubStats[k].Value, 'f', 1, 64)+Stofen(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquarySubStats[k].SubPropID), 200, yy)
				var fuciname = StoS(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquarySubStats[k].SubPropID)
				var fufigure = alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquarySubStats[k].Value
				switch fuciname {
				case "小攻击":
					fufigure = fufigure / alldata.AvatarInfoList[t].FightPropMap.Num4
					fuciname = "大攻击"
				case "小防御":
					fufigure = fufigure / alldata.AvatarInfoList[t].FightPropMap.Num7
					fuciname = "大防御"
				case "小生命":
					fufigure = fufigure / alldata.AvatarInfoList[t].FightPropMap.Num1
					fuciname = "大生命"
				default:
				}
				pingfeng += Countcitiao(str, fuciname, fufigure)
			}
			//评分处理,对齐
			if i == 2 {
				pingfeng *= 0.90
			} else if i > 2 {
				pingfeng *= 0.85
			}
			allfen += pingfeng

			//圣遗物单个评分
			if err := three.LoadFontFace(FiFile, 30); err != nil {
				panic(err)
			}
			three.DrawString(Ftoone(pingfeng), 110, 85)
			three.DrawString("-"+Pingji(pingfeng), 222, 85) //评级
			if err := three.LoadFontFace(FontFile, 30); err != nil {
				panic(err)
			}
			three.DrawString("分", 175, 85)

			switch i {
			case 0:
				x = 370
				y = 920
			case 1:
				x = 720
				y = 920
			case 2:
				x = 20
				y = 1280
			case 3:
				x = 370
				y = 1280
			case 4:
				x = 720
				y = 1280
			}
			dc.DrawImage(yinsyw, x, y)
			dc.DrawImage(three.Image(), x, y)
		}

		//总评分框
		yinping := Yinying(340, 160, 16, 0.6)
		// 字图层
		four := gg.NewContext(340, 160)
		if err := four.LoadFontFace(FontFile, 25); err != nil {
			panic(err)
		}
		four.SetRGB(1, 1, 1) //白色
		four.DrawString("评分规则:通用评分规则", 50, 35)

		if err := four.LoadFontFace(FiFile, 50); err != nil {
			panic(err)
		}
		four.DrawString(Ftoone(allfen), 50, 100)
		four.DrawString("-"+Pingji(allfen/5), 200, 100)
		if err := four.LoadFontFace(FontFile, 25); err != nil {
			panic(err)
		}
		four.DrawString("圣遗物总分", 50, 150)
		four.DrawString("评级", 230, 150)
		dc.DrawImage(yinping, 20, 1110)
		dc.DrawImage(four.Image(), 20, 1110)

		//伤害显示区,暂时展示图片
		pic, err := web.GetData(tu)
		var dst image.Image
		if err != nil {
			dst, err = gg.LoadJPG("plugin/kokomi/data/tietu/tietie.jpg")
			if err != nil {
				ctx.SendChain(message.Text("获取本地插图失败", err))
				return
			}
		} else {
			dst, _, err = image.Decode(bytes.NewReader(pic))
			if err != nil {
				ctx.SendChain(message.Text("插图解析失败", err))
				return
			}
		}
		sx := float64(1080) / float64(dst.Bounds().Size().X) // 计算缩放倍率（宽）
		dc.Scale(sx, sx)                                     // 使画笔按倍率缩放
		dc.DrawImage(dst, 0, int(1700*(1/sx)))               // 贴图（会受上述缩放倍率影响）
		dc.Scale(1/sx, 1/sx)
		//************************************************************************************
		//部分数据提前计算获取
		//命之座
		ming := len(alldata.AvatarInfoList[t].TalentIDList)
		//天赋等级
		talentid := IdtoTalent[wifeid]
		lin1 := alldata.AvatarInfoList[t].SkillLevelMap[talentid[0]]
		lin2 := alldata.AvatarInfoList[t].SkillLevelMap[talentid[1]]
		lin3 := alldata.AvatarInfoList[t].SkillLevelMap[talentid[2]]
		// 角色立绘
		var lihuifile *os.File
		var lihui image.Image
		if allfen/5 > 49.5 || ming > 4 || (lin1 == 10 && lin2 == 10 && lin3 == 10) { //第二立绘判定条件
			lihui, err = gg.LoadPNG("plugin/kokomi/data/lihui_two/" + str + ".png")
			if err != nil { //失败使用第一立绘
				lihui, err = gg.LoadPNG("plugin/kokomi/data/lihui_one/" + str + ".png")
				if err != nil { //失败使用默认立绘
					lihuifile, err = os.Open("plugin/kokomi/data/character/" + str + "/imgs/splash.webp")
					defer lihuifile.Close() // 关闭文件
					if err != nil {
						ctx.SendChain(message.Text("获取立绘失败", err))
						return
					}
					lihui, err = webp.Decode(lihuifile)
					if err != nil {
						ctx.SendChain(message.Text("解析立绘失败", err))
						return
					}
				}
			}
		} else { //第一立绘
			lihui, err = gg.LoadPNG("plugin/kokomi/data/lihui_one/" + str + ".png")
			if err != nil { //失败使用默认立绘
				lihuifile, err = os.Open("plugin/kokomi/data/character/" + str + "/imgs/splash.webp")
				defer lihuifile.Close() // 关闭文件
				if err != nil {
					ctx.SendChain(message.Text("获取立绘失败", err))
					return
				}
				lihui, err = webp.Decode(lihuifile)
				if err != nil {
					ctx.SendChain(message.Text("解析立绘失败", err))
					return
				}
			}
		}
		//立绘参数
		//syy := lihui.Bounds().Size().Y
		lihui = resize.Resize(0, 880, lihui, resize.Bilinear)
		sxx := lihui.Bounds().Size().X
		dc.DrawImage(lihui, int(260-float64(sxx)/2), 0)

		//角色名字
		if err := dc.LoadFontFace(NameFont, 80); err != nil {
			panic(err)
		}
		namelen := utf8.RuneCountInString(str)
		dc.DrawString(str, float64(1050-namelen*90), float64(130))
		// 好感度,uid
		if err := dc.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}

		//好感度位置
		dc.DrawString("好感度"+strconv.Itoa(alldata.AvatarInfoList[t].FetterInfo.ExpLevel), 20, 910)
		dc.DrawString("昵称:"+alldata.PlayerInfo.Nickname, 700, 40)
		if err := dc.LoadFontFace(FiFile, 30); err != nil {
			panic(err)
		}
		//计算宽度
		b, _ := dc.MeasureString("UID:" + suid + "---LV" + strconv.Itoa(alldata.PlayerInfo.ShowAvatarInfoList[t].Level) + "---" + strconv.Itoa(ming))
		dc.DrawString("UID:"+suid+"---LV"+strconv.Itoa(alldata.PlayerInfo.ShowAvatarInfoList[t].Level)+"---"+strconv.Itoa(ming), 976-b, 180)
		if err := dc.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		dc.DrawString("命", 976, 180)
		// 角色等级,命之座(合并上程序)
		//dc.DrawString("LV"+strconv.Itoa(alldata.PlayerInfo.ShowAvatarInfoList[t].Level), 630, 130) // 角色等级
		//dc.DrawString(strconv.Itoa(ming)+"命", 765, 130)

		//新建图层,实现阴影
		bg := Yinying(540, 470, 16, 0.6)
		//字图层
		one := gg.NewContext(540, 470)
		if err := one.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		// 属性540*460,字30,间距15,60
		one.SetRGB(1, 1, 1) //白色
		one.DrawString("生命值:", 70, 40)
		one.DrawString("攻击力:", 70, 100)
		one.DrawString("防御力:", 70, 160)
		one.DrawString("元素精通:", 70, 220)
		one.DrawString("暴击率:", 70, 280)
		one.DrawString("暴击伤害:", 70, 340)
		one.DrawString("元素充能:", 70, 400)
		// 元素加伤判断
		adds, addf := "元素加伤", 0.0
		if alldata.AvatarInfoList[t].FightPropMap.Num30*100 > addf {
			adds = "物理加伤:"
			addf = alldata.AvatarInfoList[t].FightPropMap.Num30 * 100
		}
		if alldata.AvatarInfoList[t].FightPropMap.Num40*100 > addf {
			adds = "火元素加伤:"
			addf = alldata.AvatarInfoList[t].FightPropMap.Num40 * 100
		}
		if alldata.AvatarInfoList[t].FightPropMap.Num41*100 > addf {
			adds = "雷元素加伤:"
			addf = alldata.AvatarInfoList[t].FightPropMap.Num41 * 100
		}
		if alldata.AvatarInfoList[t].FightPropMap.Num42*100 > addf {
			adds = "水元素加伤:"
			addf = alldata.AvatarInfoList[t].FightPropMap.Num42 * 100
		}
		if alldata.AvatarInfoList[t].FightPropMap.Num44*100 > addf {
			adds = "风元素加伤:"
			addf = alldata.AvatarInfoList[t].FightPropMap.Num44 * 100
		}
		if alldata.AvatarInfoList[t].FightPropMap.Num45*100 > addf {
			adds = "岩元素加伤:"
			addf = alldata.AvatarInfoList[t].FightPropMap.Num45 * 100
		}
		if alldata.AvatarInfoList[t].FightPropMap.Num46*100 > addf {
			adds = "冰元素加伤:"
			addf = alldata.AvatarInfoList[t].FightPropMap.Num46 * 100
		}
		if alldata.AvatarInfoList[t].FightPropMap.Num43*100 > addf {
			adds = "草元素加伤:"
			addf = alldata.AvatarInfoList[t].FightPropMap.Num43 * 100
		}
		one.DrawString(adds, 70, 460)

		//值,一一对应
		if err := one.LoadFontFace(FiFile, 30); err != nil {
			panic(err)
		}
		// 属性540*460,字30,间距15,60
		one.SetRGB(1, 1, 1)                                                                    //白色
		one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num2000), 335, 40)        //生命
		one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num2001), 335, 100)       //攻击
		one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num2002), 335, 160)       //防御
		one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num28), 335, 220)         //精通
		one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num20*100)+"%", 335, 280) //暴击
		one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num22*100)+"%", 335, 340) //爆伤
		one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num23*100)+"%", 335, 400) //充能
		one.DrawString(Ftoone(addf)+"%", 335, 460)

		dc.DrawImage(bg, 505, 420)
		dc.DrawImage(one.Image(), 505, 420)

		// 天赋等级
		if err := dc.LoadFontFace(FiFile, 30); err != nil { // 字体大小
			panic(err)
		}
		//贴图
		tulin1, err := gg.LoadImage("plugin/kokomi/data/character/" + str + "/icons/talent-a.webp")
		tulin1 = resize.Resize(80, 0, tulin1, resize.Bilinear)
		if err != nil {
			ctx.SendChain(message.Text("获取天赋图标失败", err))
			return
		}
		tulin2, err := gg.LoadImage("plugin/kokomi/data/character/" + str + "/icons/talent-e.webp")
		tulin2 = resize.Resize(80, 0, tulin2, resize.Bilinear)
		if err != nil {
			ctx.SendChain(message.Text("获取天赋图标失败", err))
			return
		}
		tulin3, err := gg.LoadImage("plugin/kokomi/data/character/" + str + "/icons/talent-q.webp")
		tulin3 = resize.Resize(80, 0, tulin3, resize.Bilinear)
		if err != nil {
			ctx.SendChain(message.Text("获取天赋图标失败", err))
			return
		}
		//边框间隔180
		kuang, err := gg.LoadPNG("plugin/kokomi/data/pro/" + pro + ".png")
		if err != nil {
			ctx.SendChain(message.Text("获取天赋边框失败", err))
			return
		}
		dc.DrawImage(kuang, 520, 220)
		dc.DrawImage(kuang, 700, 220)
		dc.DrawImage(kuang, 880, 220)

		//贴图间隔214
		dc.DrawImage(tulin1, 550, 260)
		//纠正素材问题
		bb := Tianfujiuzhen(str)
		dc.DrawImage(tulin2, 733, bb)
		dc.DrawImage(tulin3, 910, 260)

		//Lv背景
		talentying := gg.NewContext(40, 35)
		talentying.SetRGB(1, 1, 1) //白色
		talentying.DrawRoundedRectangle(0, 0, 40, 35, 5)
		talentying.Fill()
		talenty := AdjustOpacity(talentying.Image(), 0.9)
		dc.DrawImage(talenty, 570, 350)
		dc.DrawImage(talenty, 750, 350)
		dc.DrawImage(talenty, 930, 350)

		//Lv间隔180
		dc.SetRGB(0, 0, 0) // 换黑色
		dc.DrawString(strconv.Itoa(lin1), float64(580-lin1/10*8), 380)
		dc.DrawString(strconv.Itoa(lin2), float64(760-lin2/10*8), 380)
		dc.DrawString(strconv.Itoa(lin3), float64(940-lin3/10*8), 380)
		dc.SetRGB(1, 1, 1) // 换白色
		//皇冠
		tuguan, err := gg.LoadImage("plugin/kokomi/data/zawu/crown.png")
		if err != nil {
			ctx.SendChain(message.Text("获取皇冠图标失败", err))
			return
		}
		tuguan = resize.Resize(0, 55, tuguan, resize.Bilinear)
		if lin1 == 10 {
			dc.DrawImage(tuguan, 568, 215)
		}
		if lin2 == 10 {
			dc.DrawImage(tuguan, 748, 215)
		}
		if lin3 == 10 {
			dc.DrawImage(tuguan, 928, 215)
		}

		//命之座
		for m, mm := 1, 1; m < 7; m++ {
			tuming, err := gg.LoadImage("plugin/kokomi/data/character/" + str + "/icons/cons-" + strconv.Itoa(m) + ".webp")
			tuming = resize.Resize(40, 40, tuming, resize.Bilinear)
			if err != nil {
				ctx.SendChain(message.Text("获取命之座图标失败", err))
				return
			}
			kuang = resize.Resize(80, 0, kuang, resize.Bilinear)
			dc.DrawImage(kuang, -50+m*70, 800)
			if mm > ming {
				tuming = AdjustOpacity(tuming, 0.5)
			}
			dc.DrawImage(tuming, -30+m*70, 825)
			mm++
		}
		//**************************************************************************************************
		// 版本号
		if err := dc.LoadFontFace(BaFile, 30); err != nil {
			panic(err)
		}
		dc.DrawString(edition, 180, float64(height)-20)

		// 输出图片
		ff, cl := writer.ToBytes(dc.Image())  // 图片放入缓存
		ctx.SendChain(message.ImageBytes(ff)) // 输出
		cl()
	})

	// 绑定uid
	en.OnRegex(`^(#)?绑定\s*(uid)?(\d+)?`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		suid := ctx.State["regex_matched"].([]string)[3] // 获取uid
		int64uid, err := strconv.ParseInt(suid, 10, 64)
		if suid == "" || int64uid < 100000000 || int64uid > 1000000000 || err != nil {
			//ctx.SendChain(message.Text("-请输入正确的uid"))
			return
		}
		sqquid := strconv.Itoa(int(ctx.Event.UserID))
		file, _ := os.OpenFile("plugin/kokomi/data/uid/"+sqquid+".kokomi", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		_, _ = file.Write([]byte(suid))
		file.Close()
		ctx.SendChain(message.Text("-绑定uid" + suid + "成功" + "\n-尝试获取角色面板信息" + Postfix))

		//更新面板程序
		es, err := web.GetData(fmt.Sprintf(url, suid)) // 网站返回结果
		if err != nil {
			ctx.SendChain(message.Text("-网站获取角色信息失败"+Postfix, err))
			return
		}
		// 创建存储文件,路径plugin/kokomi/data/js
		file1, _ := os.OpenFile("plugin/kokomi/data/js/"+suid+".kokomi", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		_, _ = file1.Write(es)
		ctx.SendChain(message.Text("-获取角色面板成功" + "\n-请发送 全部面板 查看已展示角色" + Postfix))
		file1.Close()
	})
	//菜单命令
	en.OnFullMatchGroup([]string{"原神菜单", "kokomi菜单"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		menu, err := gg.LoadPNG("plugin/kokomi/data/zawu/menu.png")
		if err != nil {
			ctx.SendChain(message.Text("-获取菜单图片失败"+Postfix, err))
			return
		}
		ff, cl := writer.ToBytes(menu)
		ctx.SendChain(message.ImageBytes(ff))
		cl()
	})

	//删除账号信息,限制群内,权限管理员+可以删除别人账号信息
	en.OnRegex(`^删除账号\s*(\[CQ:at,qq=)?(\d+)?`, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		var sqquid = ""
		if ctx.State["regex_matched"].([]string)[2] != "" {
			if zero.AdminPermission(ctx) {
				sqquid = ctx.State["regex_matched"].([]string)[2] // 获取qquid
			} else {
				ctx.SendChain(message.Text("-您的权限不足" + Postfix))
			}
		}
		if sqquid == "" { // user
			sqquid = strconv.FormatInt(ctx.Event.UserID, 10)
		}
		err := os.Remove("plugin/kokomi/data/uid/" + sqquid + ".kokomi")
		if err != nil {
			//如果删除失败则输出 file remove Error!
			ctx.SendChain(message.Text("-未找到该账号信息" + Postfix))
		} else {
			//如果删除成功则输出 file remove OK!
			ctx.SendChain(message.Text("-删除成功" + Postfix))
		}
	})

	//上传立绘,限制群内,权限管理员+
	en.OnRegex(`^上传第(1|2|一|二)立绘\s*(.*)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		z := ctx.State["regex_matched"].([]string)[1] // 获取编号
		wifename := ctx.State["regex_matched"].([]string)[2]
		var pathw string
		swifeid := Findnames(wifename, "wife")
		if swifeid == "" {
			ctx.SendChain(message.Text("-请输入角色全名" + Postfix))
			return
		}
		wifename = Idmap(swifeid, "wife")
		if wifename == "" {
			ctx.SendChain(message.Text("Idmap数据缺失"))
			return
		}
		switch z {
		case "1", "一":
			pathw = "plugin/kokomi/data/lihui_one/" + wifename + ".png"
		case "2", "二":
			pathw = "plugin/kokomi/data/lihui_two/" + wifename + ".png"
		}
		next := zero.NewFutureEvent("message", 999, false, zero.OnlyGroup, ctx.CheckSession())
		recv, stop := next.Repeat()
		defer stop()
		ctx.SendChain(message.Text("-请发送面板图" + Postfix))
		var step int
		var origin string
		for {
			select {
			case <-time.After(time.Second * 120):
				ctx.SendChain(message.Text("-时间太久啦！摆烂惹!"))
				return
			case c := <-recv:
				switch step {
				case 0:
					re := regexp.MustCompile(`https:(.*)is_origin=(0|1)`)
					origin = re.FindString(c.Event.RawMessage)
					ctx.SendChain(message.Text("-请输入\"确定\"或者\"取消\"来决定是否上传" + Postfix))
					step++
				case 1:
					msg := c.Event.Message.ExtractPlainText()
					if msg != "确定" && msg != "取消" {
						ctx.SendChain(message.Text("-请输入\"确定\"或者\"取消\"" + Postfix))
						continue
					}
					if msg == "确定" {
						ctx.SendChain(message.Text("-正在上传..." + Postfix))
						pic, err := web.GetData(origin)
						if err != nil {
							ctx.SendChain(message.Text("-获取插图失败"+Postfix, err))
							return
						}
						dst, _, err := image.Decode(bytes.NewReader(pic))
						if err != nil {
							ctx.SendChain(message.Text("-插图解析失败"+Postfix, err))
							return
						}
						err = gg.SavePNG(pathw, dst)
						if err != nil {
							ctx.SendChain(message.Text("-上传失败惹~", err))
							return
						}
						ctx.SendChain(message.Text("-上传成功了" + Postfix))
						return
					}
					ctx.SendChain(message.Text("-已经取消上传了" + Postfix))
					return
				}
			}
		}
	})
	//删除立绘图,权限同上
	en.OnRegex(`^删除第(1|2|一|二)立绘\s*(.*)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		z := ctx.State["regex_matched"].([]string)[1] // 获取编号
		wifename := ctx.State["regex_matched"].([]string)[2]
		var pathw string
		swifeid := Findnames(wifename, "wife")
		if swifeid == "" {
			ctx.SendChain(message.Text("-请输入角色全名" + Postfix))
			return
		}
		wifename = Idmap(swifeid, "wife")
		if wifename == "" {
			ctx.SendChain(message.Text("Idmap数据缺失"))
			return
		}
		switch z {
		case "1", "一":
			pathw = "plugin/kokomi/data/lihui_one/" + wifename + ".png"
		case "2", "二":
			pathw = "plugin/kokomi/data/lihui_two/" + wifename + ".png"
		}
		next := zero.NewFutureEvent("message", 999, false, zero.OnlyGroup, ctx.CheckSession())
		recv, stop := next.Repeat()
		defer stop()
		ctx.SendChain(message.Text("-请输入\"确定\"或者\"取消\"来决定是否删除" + Postfix))
		var origin string
		for {
			select {
			case <-time.After(time.Second * 120):
				ctx.SendChain(message.Text("-时间太久啦！摆烂惹!"))
				return
			case c := <-recv:
				origin = c.Event.Message.ExtractPlainText()
				if origin != "确定" && origin != "取消" {
					ctx.SendChain(message.Text("-请输入\"确定\"或者\"取消\"" + Postfix))
					continue
				}
				if origin == "确定" {
					err := os.Remove(pathw)
					if err != nil {
						//如果删除失败则输出 file remove Error!
						ctx.SendChain(message.Text("-未找到该面板图" + Postfix))
					} else {
						//如果删除成功则输出 file remove OK!
						ctx.SendChain(message.Text("-删除成功" + Postfix))
					}
					return
				}
				ctx.SendChain(message.Text("-已经取消删除了" + Postfix))
				return
			}
		}
	})
}
