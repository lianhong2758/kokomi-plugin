// Package kokomi 原神面板v2.4.2
package kokomi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	//"unicode/utf8"
	kanban "github.com/FloatTech/ZeroBot-Plugin/kanban/banner"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/imgfactory"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	//"github.com/fogleman/gg"//原版gg
	"github.com/FloatTech/gg"
	//"github.com/golang/freetype"
	"github.com/nfnt/resize"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"golang.org/x/image/webp"
)

const (
	//tu       = "https://api.yimian.xyz/img?type=moe&size=1920x1080"
	NameFont = "plugin/kokomi/data/font/NZBZ.ttf"                    // 名字字体
	FontFile = "plugin/kokomi/data/font/SourceHanMonoSC-HeavyIt.ttf" // 汉字字体
	FiFile   = "plugin/kokomi/data/font/tttgbnumber.ttf"             // 其余字体(数字英文)
	BaFile   = "plugin/kokomi/data/font/STLITI.TTF"                  // 华文隶书版本版本号字体
)

func init() { // 主函数
	en := control.Register("kokomi", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "原神面板查询",
		Help: "- kokomi菜单([]里面为可选项)\n" +
			"- 绑定......(uid)\n" +
			"- 更新面板[@xx]\n" +
			"- 全部面板[@xx]\n" +
			"- XX面板[@xx]\n" +
			"- 删除账号[@xx]\n" +
			"- [@xx]队伍伤害[xx xx xx xx]\n" +
			"- 管理员专属指令:\n" +
			"- (上传|删除)第(1|2)立绘 XX\n",
	})
	en.OnRegex(`^(?:#|＃)(.*)面板\s*(?:\[CQ:at,qq=)?(\d+)?(.*)`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		var allfen = 0.00
		sqquid := ctx.State["regex_matched"].([]string)[2] // 获取第三者qquid
		if sqquid == "" {
			sqquid = strconv.FormatInt(ctx.Event.UserID, 10)
		}
		str := ctx.State["regex_matched"].([]string)[1] // 获取key
		if str == "" {
			str = ctx.State["regex_matched"].([]string)[3]
		}
		//尝试匹配角色名或者指令
		wife := GetWifeOrWq("wife")
		if str != "更新" && str != "全部" && str != "全部角色" {
			swifeid := wife.Findnames(str)
			if swifeid == "" {
				//ctx.SendChain(message.Text("-请输入角色全名" + Postfix))
				return
			}
			str = wife.Idmap(swifeid)
			if str == "" {
				ctx.SendChain(message.Text("Idmap数据缺失"))
				return
			}
		}
		if str == "空" || str == "荧" || str == "旅行者" {
			ctx.SendChain(message.Text("-暂不支持查看该角色" + Postfix))
			return
		}
		// 获取uid
		uid := Getuid(sqquid)
		suid := strconv.Itoa(uid)
		if uid == 0 {
			ctx.SendChain(message.Text("-未绑定uid\n-第一次使用请发送\"绑定xxx\"" + Postfix))
			return
		}
		//############################################################判断数据更新,逻辑原因不能合并进switch
		if str == "更新" {
			es, err := getEnka(ctx, suid)
			if err != nil || len(es) == 0 {
				ctx.SendChain(message.Text("-网站获取角色信息失败"+Postfix, err))
				return
			}
			//解析
			var dam_a []byte
			var msg strings.Builder
			if Config.Apiid < 3 {
				var ndata Data
				err = json.Unmarshal(es, &ndata)
				if err != nil {
					ctx.SendChain(message.Text("出现错误捏：", err))
					return
				}
				if len(ndata.PlayerInfo.ShowAvatarInfoList) == 0 || len(ndata.AvatarInfoList) == 0 {
					ctx.SendChain(message.Text("-请在游戏中打开角色展柜,并将想查询的角色进行展示" + "\n-完成上述操作并等待5分钟后,请使用 更新面板 获取账号信息" + Postfix))
					return
				}
				//映射
				thisdata, err := ndata.ConvertData()
				if err != nil {
					ctx.SendChain(message.Text("数据映射错误捏：", err))
					return
				}
				//合并映射
				thisdata.MergeFile(suid)
				es, err = json.Marshal(&thisdata)
				if err != nil {
					ctx.SendChain(message.Text("数据反解析错误捏：", err))
					return
				}
				msg.WriteString("-获取角色面板成功\n")
				msg.WriteString("-您的展示角色为:\n")
				for i := 0; i < len(thisdata.Chars); i++ {
					msg.WriteString(" ")
					msg.WriteString(thisdata.Chars[i].Name)
					if i < len(thisdata.Chars)-1 {
						msg.WriteByte('\n')
					}
				}
				dam_a, err = thisdata.GetSumComment(suid)
				if err != nil {
					ctx.SendChain(message.Text("-获取伤害数据失败\n-请尝试联系维护者或者重新绑定uid\n-若多次尝试无效,请等待修复"+Postfix, err))
				}
			}
			//存储伤害计算返回值
			file2, _ := os.OpenFile("plugin/kokomi/data/damage/"+suid+".kokomi", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			_, _ = file2.Write(dam_a)
			file2.Close()
			// 创建存储文件,路径plugin/kokomi/data/js
			file, _ := os.OpenFile("plugin/kokomi/data/js/"+suid+".kokomi", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			_, _ = file.Write(es)
			ctx.SendChain(message.Text(msg.String()))
			file.Close()
			return
		}
		//############################################################
		// 获取本地缓存数据
		txt, err := os.ReadFile("plugin/kokomi/data/js/" + suid + ".kokomi")
		if err != nil {
			ctx.SendChain(message.Text("-本地未找到账号信息, 请\"更新面板\"" + Postfix))
			return
		}

		// 解析
		var alldata Thisdata
		err = json.Unmarshal(txt, &alldata)
		if err != nil {
			ctx.SendChain(message.Text("出现错误捏：", err))
			return
		}
		if len(alldata.Chars) == 0 {
			ctx.SendChain(message.Text("-请在游戏中打开角色展柜,并将想查询的角色进行展示" + "\n-完成上述操作并等待5分钟后,请使用\"更新面板\"获取账号信息" + Postfix))
			return
		}

		switch str {
		case "全部", "全部角色":
			var msg strings.Builder
			msg.WriteString("-您的展示角色为:\n")
			//	for i := 0; i < len(alldata.Chars); i++ {
			i := 0
			for _, v := range alldata.Chars {
				msg.WriteString(" ")
				msg.WriteString(v.Name)
				if i < len(alldata.Chars)-1 {
					msg.WriteByte('\n')
				}
				i++
			}
			ctx.SendChain(message.Text(msg.String()))
			return
		}
		var t = -1
		// 匹配角色
		for i, v := range alldata.Chars {
			if str == v.Name {
				t = i
			}
		}

		if t == -1 { // 在返回数据中未找到想要的角色
			ctx.SendChain(message.Text("-该角色未展示" + Postfix))
			return
		}

		// 画图
		const height = 2400 - 360
		dc := gg.NewContext(1080, height) // 画布大小
		//*******************************************************
		role, err := GetRole(str)
		if err != nil {
			ctx.SendChain(message.Text("获取角色失败", err))
			return
		}
		//*******************************************************
		pro := role.Elem
		beijing, err := gg.LoadImage("plugin/kokomi/data/pro/" + pro + ".jpg")
		if err != nil {
			ctx.SendChain(message.Text("获取背景失败", err))
			return
		}
		beijing = resize.Resize(0, 2040, beijing, resize.Bilinear) //改比例
		dc.DrawImage(beijing, -690, 0)
		dc.SetRGB(1, 1, 1) // 换白色

		//武器图层
		// 字图层
		two := gg.NewContext(340, 180)
		if err := two.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		two.SetRGB(1, 1, 1) //白色

		//武器名
		//纠正圣遗物空缺报错的无返回情况
		wq := alldata.Chars[t].Weapon.Name
		if wq == "" {
			ctx.SendChain(message.Text("获取武器名称失败"))
		}
		two.DrawString(wq, 150, 50)
		//星级
		two.DrawImage(resize.Resize(0, 30, Drawstars("#FFCC00", "#FFE43A", alldata.Chars[t].Weapon.Star), resize.Bilinear), 150, 60)
		//详细
		if alldata.Chars[t].Weapon.Atk != 0.0 {
			two.DrawString("攻击力:", 145, 160)
		}
		two.DrawString("精炼:", 240, 120)
		if err := two.LoadFontFace(FiFile, 30); err != nil { // 字体大小
			panic(err)
		}
		if alldata.Chars[t].Weapon.Atk != 0.0 {
			two.DrawString(strconv.FormatFloat(alldata.Chars[t].Weapon.Atk, 'f', 1, 32), 250, 160) //攻击力
		}
		//Lv,精炼
		two.DrawString("Lv."+strconv.Itoa(alldata.Chars[t].Weapon.Level), 150, 120)
		two.DrawString(strconv.Itoa(alldata.Chars[t].Weapon.Affix), 316, 120)
		//图片
		tuwq, err := gg.LoadPNG("plugin/kokomi/data/wq/" + wq + ".png")
		if err != nil {
			ctx.SendChain(message.Text("获取武器图标失败", err))
			return
		}
		tuwq = resize.Resize(130, 0, tuwq, resize.Bilinear)
		// int(213*0.6)
		yinyinBlack127 := color.NRGBA{R: 0, G: 0, B: 0, A: 127}

		two.DrawImage(tuwq, 10, 10)
		dc.DrawImage(Yinying(340, 180, 16, yinyinBlack127), 20, 920) // 背景
		dc.DrawImage(two.Image(), 20, 920)

		//圣遗物
		yinsyw := Yinying(340, 350, 16, yinyinBlack127)
		var syw sywm
		for i := 0; i < 5; i++ {
			switch i {
			case 0:
				syw = alldata.Chars[t].Artis.Hua
			case 1:
				syw = alldata.Chars[t].Artis.Yu
			case 2:
				syw = alldata.Chars[t].Artis.Sha
			case 3:
				syw = alldata.Chars[t].Artis.Bei
			case 4:
				syw = alldata.Chars[t].Artis.Guan
			}
			if syw.Name == "" {
				continue
			}
			// 字图层
			three := gg.NewContext(340, 350)
			if err := three.LoadFontFace(FontFile, 30); err != nil {
				panic(err)
			}
			//字号30,间距50
			three.SetRGB(1, 1, 1) //白色
			//画线
			for c := 0; c < 4; c++ {
				three.DrawLine(0, 157+float64(c)*45, 350, 157+float64(c)*45) //横线条分割
			}
			three.Stroke()
			sywname := syw.Set
			tusyw, err := gg.LoadImage("plugin/kokomi/data/syw/" + sywname + "/" + strconv.Itoa(i+1) + ".webp")
			if err != nil {
				ctx.SendChain(message.Text("获取圣遗物图标失败", err))
				return
			}
			tusyw = resize.Resize(80, 0, tusyw, resize.Bilinear) //缩小
			three.DrawImage(tusyw, 15, 15)
			//圣遗物name
			three.DrawString(syw.Name, 110, 50)
			//圣遗物属性 主词条
			//间隔45,初始145
			var xx, yy, pingfeng float64 //xx,yy词条相对位置,x,y文本框在全图位置
			var x, y int
			xx = 15
			yy = 145
			pingfeng = 0
			//主词条
			three.DrawString(syw.Main.Title, xx, yy) //"主:"
			if err := three.LoadFontFace(FiFile, 30); err != nil {
				panic(err)
			}
			//主词条属性
			//+对齐three.DrawString("+"+zhucitiao+Stofen(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquaryMainStat.MainPropID), 200, yy)
			three.DrawStringAnchored("+"+Ftoone(syw.Main.Value)+Stofen(syw.Main.Title), 325, yy, 1, 0) //主词条属性
			//算分
			if i > 1 { //不算前两主词条属性
				pingfeng += Countcitiao(str, syw.Main.Title, syw.Main.Value/4)
			}
			//副词条
			three.SetHexColor("#98F5FF") //蓝色
			p := len(syw.Attrs)
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
				var fuciname = syw.Attrs[k].Title
				var fuciname2 = syw.Attrs[k].Title
				var fufigure = syw.Attrs[k].Value
				var fufigure2 = syw.Attrs[k].Value
				switch fuciname2 {
				case "小攻击":
					fufigure2 = fufigure / alldata.Chars[t].Attr.AtkBase * 100
					fuciname2 = "大攻击"
				case "小防御":
					fufigure2 = fufigure / alldata.Chars[t].Attr.DefBase * 100
					fuciname2 = "大防御"
				case "小生命":
					fufigure2 = fufigure / alldata.Chars[t].Attr.HpBase * 100
					fuciname2 = "大生命"
				}
				pingfeng += Countcitiao(str, fuciname2, fufigure2) //单个圣遗物分数合计
				if Countcitiao(str, fuciname2, fufigure2) == 0.0 {
					three.SetHexColor("#BEBEBE") //灰色#BEBEBE,浅灰色#D3D3D3
				}
				if err := three.LoadFontFace(FontFile, 30); err != nil {
					panic(err)
				}
				three.DrawString((fuciname), xx, yy)
				if err := three.LoadFontFace(FiFile, 30); err != nil {
					panic(err)
				}
				three.DrawStringAnchored("+"+Ftoone(fufigure)+Stofen(fuciname), 325, yy, 1, 0)
				three.SetHexColor("#98F5FF") //蓝色
			}
			//评分处理,对齐
			if i == 2 {
				pingfeng *= 0.90
			} else if i > 2 {
				pingfeng *= 0.85
			}
			allfen += pingfeng

			three.SetRGB(1, 1, 1)
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
				x, y = 370, 920
			case 1:
				x, y = 720, 920
			case 2:
				x, y = 20, 1280
			case 3:
				x, y = 370, 1280
			case 4:
				x, y = 720, 1280
			}
			dc.DrawImage(yinsyw, x, y)
			dc.DrawImage(three.Image(), x, y)
		}

		//总评分框
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
		four.DrawStringAnchored("-"+Pingji(allfen/5), 255, 100, 0.5, 0)
		if err := four.LoadFontFace(FontFile, 25); err != nil {
			panic(err)
		}
		four.DrawString("圣遗物总分", 50, 150)
		four.DrawString("评级", 230, 150)
		dc.DrawImage(Yinying(340, 160, 16, yinyinBlack127), 20, 1110) // 背景
		dc.DrawImage(four.Image(), 20, 1110)

		var ok = -1
		damfile, err := os.ReadFile("plugin/kokomi/data/damage/" + suid + ".kokomi")
		if err != nil {
			ok = 0
		}
		var roleDam Dam
		if ok != 0 {
			err = json.Unmarshal(damfile, &roleDam)
			if err != nil {
				ok = 1
			}
		}
		//绘图区
		six := gg.NewContext(1040, 325)
		six.SetRGB(1, 1, 1) //白色
		//汉字描述
		if err := six.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		for c := 1; c <= 4; c++ {
			six.DrawLine(0, 65*float64(c), 1040, 65*float64(c)) //横线条分割
		}
		for c := 1; c < 3; c++ {
			six.DrawLine(346*float64(c), 65, 346*float64(c), 325) //竖线条分割
		}
		six.DrawString("伤害计算[结果仅供参考,以实际为准]", 50, 40)
		six.DrawStringAnchored("伤害类型", 290, 105, 1, 0)
		six.DrawStringAnchored("暴击伤害/治疗/护盾", 520, 105, 0.5, 0)
		six.DrawStringAnchored("期望伤害(EX)", 867, 105, 0.5, 0)
		switch ok {
		case -1:
			for c := 1; c <= 3 && c <= len(roleDam.Result[t].DamageResultArr); c++ {
				six.DrawStringAnchored(roleDam.Result[t].DamageResultArr[c-1].Title, 290, 105+65*float64(c), 1, 0)
			}
			if len(roleDam.Result[t].DamageResultArr) < 3 {
				six.DrawStringAnchored("暂无数据", 290, 300, 1, 0)
			}
			if err := six.LoadFontFace(FiFile, 30); err != nil {
				panic(err)
			}
			for c := 1; c <= 3 && c <= len(roleDam.Result[t].DamageResultArr); c++ {
				six.DrawStringAnchored(fmt.Sprint(roleDam.Result[t].DamageResultArr[c-1].Value), 520, 105+65*float64(c), 0.5, 0)
				if roleDam.Result[t].DamageResultArr[c-1].Expect != "" {
					six.DrawStringAnchored(roleDam.Result[t].DamageResultArr[c-1].Expect[6:], 867, 105+65*float64(c), 0.5, 0)
				} else {
					six.DrawLine(692, 65*float64(c+1), 1040, 65*float64(c+2))
				}
			}
		case 0:
			six.DrawStringAnchored("暂无数据", 290, 170, 1, 0)
			six.DrawString("请\"更新面板\"", 360, 170)
		case 1:
			six.DrawStringAnchored("数据错误", 290, 170, 1, 0)
			six.DrawString("请联系维护人员", 360, 170)
		}
		six.Stroke()
		dc.DrawImage(Yinying(1040, 325, 16, yinyinBlack127), 20, 1660) // 背景
		dc.DrawImage(six.Image(), 20, 1660)

		//************************************************************************************
		//部分数据提前计算获取
		//命之座
		ming := alldata.Chars[t].Cons
		//天赋等级
		lin1 := alldata.Chars[t].Talent.A
		lin2 := alldata.Chars[t].Talent.E
		lin3 := alldata.Chars[t].Talent.Q
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
		dc.DrawImage(lihui, int(300-float64(sxx)/2), 0)

		// 好感度,uid
		if err := dc.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		//好感度位置,数据来源
		dc.DrawString("好感度"+strconv.Itoa(alldata.Chars[t].Fetter), 20, 910)
		dc.DrawStringAnchored("Data From "+alldata.Chars[t].DataSource, 1045, 910, 1, 0)
		//昵称图框
		five := gg.NewContext(540, 200)
		five.SetRGB(1, 1, 1) //白色
		//角色名字
		if err := five.LoadFontFace(NameFont, 80); err != nil {
			panic(err)
		}
		five.DrawStringAnchored(str, 505, 130, 1, 0)
		if err := five.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		five.DrawStringAnchored("昵称:"+alldata.Nickname, 505, 40, 1, 0)
		five.DrawString("命", 470, 180)
		if err := five.LoadFontFace(FiFile, 30); err != nil {
			panic(err)
		}
		five.DrawStringAnchored("UID:"+suid+"--LV"+strconv.Itoa(alldata.Level)+"--"+strconv.Itoa(ming), 470, 180, 1, 0)
		// 角色等级,命之座(合并上程序)
		//dc.DrawString("LV"+strconv.Itoa(alldata.PlayerInfo.ShowAvatarInfoList[t].Level), 630, 130) // 角色等级
		//dc.DrawString(strconv.Itoa(ming)+"命", 765, 130)
		// 透明度 int(213*0.5)
		newying := Yinying(540, 200, 16, color.NRGBA{R: 0, G: 0, B: 0, A: 106})
		dc.DrawImage(newying, 505, 20)
		dc.DrawImage(five.Image(), 505, 20)

		//字图层
		one := gg.NewContext(540, 470)
		if err := one.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		// 属性540*460,字30,间距15,60
		one.SetRGB(1, 1, 1) //白色
		one.DrawString("角色等级:", 70, 45)
		one.DrawString("生命值:", 70, 96.25)
		one.DrawString("攻击力:", 70, 147.5)
		one.DrawString("防御力:", 70, 198.75)
		one.DrawString("元素精通:", 70, 250)
		one.DrawString("暴击率:", 70, 301.25)
		one.DrawString("暴击伤害:", 70, 352.5)
		one.DrawString("元素充能:", 70, 403.75)
		// 元素加伤判断
		adds, addf := alldata.Chars[t].Attr.DmgName, alldata.Chars[t].Attr.Dmg
		if adds == "" {
			adds = "元素加伤:"
		}
		one.DrawString(adds, 70, 455)

		//值,一一对应
		if err := one.LoadFontFace(FiFile, 30); err != nil {
			panic(err)
		}
		// 属性540*460,字30,间距15,60
		one.SetRGB(1, 1, 1)                                                                   //白色
		one.DrawStringAnchored("Lv"+alldata.Chars[t].Level, 470, 45, 1, 0)                    //Lv
		one.DrawStringAnchored(Ftoone(alldata.Chars[t].Attr.Hp), 470, 96.25, 1, 0)            //生命
		one.DrawStringAnchored(Ftoone(alldata.Chars[t].Attr.Atk), 470, 147.5, 1, 0)           //攻击
		one.DrawStringAnchored(Ftoone(alldata.Chars[t].Attr.Def), 470, 198.75, 1, 0)          //防御
		one.DrawStringAnchored(Ftoone(alldata.Chars[t].Attr.Mastery), 470, 250, 1, 0)         //精通
		one.DrawStringAnchored(Ftoone(alldata.Chars[t].Attr.Cpct)+"%", 470, 301.25, 1, 0)     //暴击
		one.DrawStringAnchored(Ftoone(alldata.Chars[t].Attr.Cdmg)+"%", 470, 352.5, 1, 0)      //爆伤
		one.DrawStringAnchored(Ftoone(alldata.Chars[t].Attr.Recharge)+"%", 470, 403.75, 1, 0) //充能
		one.DrawStringAnchored(Ftoone(addf)+"%", 470, 455, 1, 0)                              //元素加伤
		//评分权重
		ttt := []string{
			strconv.Itoa(Wifequanmap[str].Hp),
			strconv.Itoa(Wifequanmap[str].Atk),
			strconv.Itoa(Wifequanmap[str].Def),
			strconv.Itoa(Wifequanmap[str].Mastery),
			strconv.Itoa(Wifequanmap[str].Cpct),
			strconv.Itoa(Wifequanmap[str].Cdmg),
			strconv.Itoa(Wifequanmap[str].Recharge),
			strconv.Itoa(Wifequanmap[str].Dmg),
		}
		one.SetHexColor("#98F5FF")
		for i, v := range ttt {
			one.DrawString(v, 8, 96+51.25*float64(i))
		}

		dc.DrawImage(Yinying(540, 470, 16, yinyinBlack127), 505, 410) // 背景
		dc.DrawImage(one.Image(), 505, 410)

		// 天赋等级
		seven := gg.NewContext(540, 190)
		if err := seven.LoadFontFace(FiFile, 30); err != nil { // 字体大小
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
		seven.DrawImage(kuang, 15, 10)
		seven.DrawImage(kuang, 195, 10)
		seven.DrawImage(kuang, 375, 10)
		//贴图
		seven.DrawImageAnchored(tulin1, 85, 92, 0.5, 0.5)
		seven.DrawImageAnchored(tulin2, 268, 90, 0.5, 0.5)
		seven.DrawImageAnchored(tulin3, 445, 92, 0.5, 0.5)

		// Lv背景 透明度 int(255*0.9)
		talenty := Yinying(40, 35, 5, color.NRGBA{R: 255, G: 255, B: 255, A: 226})
		seven.DrawImageAnchored(talenty, 85, 145, 0.5, 0.5)
		seven.DrawImageAnchored(talenty, 265, 145, 0.5, 0.5)
		seven.DrawImageAnchored(talenty, 445, 145, 0.5, 0.5)

		//皇冠
		tuguan, err := gg.LoadImage("plugin/kokomi/data/zawu/crown.png")
		if err != nil {
			ctx.SendChain(message.Text("获取皇冠图标失败", err))
			return
		}
		tuguan = resize.Resize(0, 55, tuguan, resize.Bilinear)
		if lin1 == 10 {
			seven.DrawImageAnchored(tuguan, 90, 30, 0.5, 0.5)
		}
		if lin2 == 10 {
			seven.DrawImageAnchored(tuguan, 270, 30, 0.5, 0.5)
		}
		if lin3 == 10 {
			seven.DrawImageAnchored(tuguan, 450, 30, 0.5, 0.5)
		}

		//Lv-天赋等级修复
		if ming >= role.TalentCons.E {
			lin2 += 3
		}
		if ming >= role.TalentCons.Q {
			lin3 += 3
		}
		if str == "达达利亚" {
			lin1++
		}
		//Lv间隔180
		seven.SetRGB(0, 0, 0) // 换黑色
		seven.DrawStringAnchored(strconv.Itoa(lin1), 85, 145, 0.5, 0.5)
		seven.DrawStringAnchored(strconv.Itoa(lin2), 265, 145, 0.5, 0.5)
		seven.DrawStringAnchored(strconv.Itoa(lin3), 445, 145, 0.5, 0.5)
		dc.DrawImage(seven.Image(), 505, 220)

		// 命之座
		kuang = resize.Resize(80, 0, kuang, resize.Bilinear)
		kuangblack := AdjustOpacity(kuang, 0.5)
		for m, mm := 1, 1; m < 7; m++ {
			tuming, err := gg.LoadImage("plugin/kokomi/data/character/" + str + "/icons/cons-" + strconv.Itoa(m) + ".webp")
			if err != nil {
				ctx.SendChain(message.Text("获取命之座图标失败", err))
				return
			}
			tuming = resize.Resize(40, 0, tuming, resize.Bilinear)
			if mm > ming {
				dc.DrawImage(kuangblack, -50+m*70, 800)
				tuming = AdjustOpacity(tuming, 0.5)
			} else {
				dc.DrawImage(kuang, -50+m*70, 800)
			}
			dc.DrawImageAnchored(tuming, -30+m*70, 845, 0, 0.5)
			mm++
		}
		//**************************************************************************************************
		// 版本号
		if err := dc.LoadFontFace(BaFile, 30); err != nil {
			panic(err)
		}
		dc.DrawStringAnchored("Created By ZeroBot-Plugin "+kanban.Version+" & "+edition, 540, float64(height)-30, 0.5, 0.5)
		// 输出图片
		ff, err := imgfactory.ToBytes(dc.Image()) // 图片放入缓存
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		ctx.SendChain(message.ImageBytes(ff)) // 输出
	})

	// 绑定uid
	en.OnRegex(`^(?:#|＃)\s*绑定+?\s*(?:uid|UID|Uid)?\s*(\d+)?`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		suid := ctx.State["regex_matched"].([]string)[1] // 获取uid
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
		//es, err := web.GetData(fmt.Sprintf(url, suid)) // 网站返回结果
		es, err := getEnka(ctx, suid)
		if err != nil || len(es) == 0 {
			ctx.SendChain(message.Text("-网站获取角色信息失败"+Postfix, err))
			return
		}
		var dam_a []byte
		var msg strings.Builder
		//解析
		if Config.Apiid < 3 {
			var ndata Data
			err = json.Unmarshal(es, &ndata)
			if err != nil {
				ctx.SendChain(message.Text("出现错误捏：", err))
				return
			}
			if len(ndata.PlayerInfo.ShowAvatarInfoList) == 0 || len(ndata.AvatarInfoList) == 0 {
				ctx.SendChain(message.Text("-请在游戏中打开角色展柜,并将想查询的角色进行展示" + "\n-完成上述操作并等待5分钟后,请使用 更新面板 获取账号信息" + Postfix))
				return
			}
			//映射
			thisdata, err := ndata.ConvertData()
			if err != nil {
				ctx.SendChain(message.Text("数据映射错误捏：", err))
				return
			}
			es, err = json.Marshal(&thisdata)
			if err != nil {
				ctx.SendChain(message.Text("数据反解析错误捏：", err))
				return
			}
			wife := GetWifeOrWq("wife")
			msg.WriteString("-获取角色面板成功\n")
			msg.WriteString("-您的展示角色为:\n")
			for i := 0; i < len(thisdata.Chars); i++ {
				mmm := wife.Idmap(strconv.Itoa(thisdata.Chars[i].ID))
				if mmm == "" {
					ctx.SendChain(message.Text("Idmap数据缺失"))
					return
				}
				msg.WriteString(" ")
				msg.WriteString(mmm)
				if i < len(thisdata.Chars)-1 {
					msg.WriteByte('\n')
				}
			}
			dam_a, err = thisdata.GetSumComment(suid)
			if err != nil {
				ctx.SendChain(message.Text("-获取伤害数据失败"+Postfix, err))
			}
		}
		//存储伤害计算返回值
		file2, _ := os.OpenFile("plugin/kokomi/data/damage/"+suid+".kokomi", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		_, _ = file2.Write(dam_a)
		file2.Close()
		// 创建存储文件,路径plugin/kokomi/data/js
		file1, _ := os.OpenFile("plugin/kokomi/data/js/"+suid+".kokomi", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		_, _ = file1.Write(es)
		ctx.SendChain(message.Text(msg.String()))
		file1.Close()
	})
	//菜单命令
	en.OnFullMatchGroup([]string{"原神菜单", "kokomi菜单", "菜单"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		menu, err := gg.LoadPNG("plugin/kokomi/data/zawu/menu.png")
		if err != nil {
			ctx.SendChain(message.Text("-获取菜单图片失败"+Postfix, err))
			return
		}
		ff, err := imgfactory.ToBytes(menu) // 图片放入缓存
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		ctx.SendChain(message.ImageBytes(ff))
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
	en.OnRegex(`^上传第(1|2|一|二)立绘\s*(.*)`, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		z := ctx.State["regex_matched"].([]string)[1] // 获取编号
		wifename := ctx.State["regex_matched"].([]string)[2]
		var pathw string
		wife := GetWifeOrWq("wife")
		swifeid := wife.Findnames(wifename)
		if swifeid == "" {
			ctx.SendChain(message.Text("-请输入角色全名" + Postfix))
			return
		}
		wifename = wife.Idmap(swifeid)
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
	en.OnRegex(`^删除第(1|2|一|二)立绘\s*(.*)`, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		z := ctx.State["regex_matched"].([]string)[1] // 获取编号
		wifename := ctx.State["regex_matched"].([]string)[2]
		var pathw string
		wife := GetWifeOrWq("wife")
		swifeid := wife.Findnames(wifename)
		if swifeid == "" {
			ctx.SendChain(message.Text("-请输入角色全名" + Postfix))
			return
		}
		wifename = wife.Idmap(swifeid)
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
				ctx.SendChain(message.Text("-已经删除了" + Postfix))
				return
			}
		}
	})
	//切换api
	en.OnRegex(`^切换api(\d)?`, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		z := ctx.State["regex_matched"].([]string)[1] // 获取编号
		if z != "" {
			zz, _ := strconv.Atoi(z)
			if zz >= len(Config.Apis) {
				ctx.SendChain(message.Text("-api不存在" + Postfix))
				return
			}
			url = Config.Apis[zz]
			goto success
		}
		if Config.Apiid+1 < len(Config.Apis) {
			url = Config.Apis[Config.Apiid+1]
			Config.Apiid++
			goto success
		}
		url = Config.Apis[0]
		Config.Apiid = 0
		goto success
	success:
		ctx.SendChain(message.Text("切换api成功" + Postfix))
	})

	//队伍伤害
	en.OnRegex(`^(?:\[CQ:at,qq=)?(\d+)?\]?\s*(?:#|＃)队伍伤害\s*((\D+)\s(\D+)\s(\D+)\s(\D+))?`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		var alldata Thisdata
		is := [4]int{}
		sqquid := ctx.State["regex_matched"].([]string)[1] // 获取第三者qquid
		if sqquid == "" {
			sqquid = strconv.FormatInt(ctx.Event.UserID, 10)
		}
		// 获取uid
		uid := Getuid(sqquid)
		suid := strconv.Itoa(uid)
		if uid == 0 {
			ctx.SendChain(message.Text("-未绑定uid\n-第一次使用请发送\"绑定xxx\"" + Postfix))
			return
		}

		{
			txt, err := os.ReadFile("plugin/kokomi/data/js/" + suid + ".kokomi")
			if err != nil {
				ctx.SendChain(message.Text("-本地未找到账号信息, 请\"更新面板\"" + Postfix))
				return
			}
			err = json.Unmarshal(txt, &alldata)
			if err != nil {
				ctx.SendChain(message.Text("出现错误捏：", err))
				return
			}
			if len(alldata.Chars) == 0 {
				ctx.SendChain(message.Text("-请在游戏中打开角色展柜,并将想查询的角色进行展示" + "\n-完成上述操作并等待5分钟后,请使用\"更新面板\"获取账号信息" + Postfix))
				return
			} else if len(alldata.Chars) < 4 {
				ctx.SendChain(message.Text("-展示的角色数量不足4" + Postfix))
				return
			}
		}

		names := []string{ctx.State["regex_matched"].([]string)[3], ctx.State["regex_matched"].([]string)[4], ctx.State["regex_matched"].([]string)[5], ctx.State["regex_matched"].([]string)[6]} // 获取key
		if names[0] == "" {
			is = [4]int{0, 1, 2, 3}
			for i := 0; i < 4; i++ {
				names[i] = alldata.Chars[i].Name
			}
		} else {
			wife := GetWifeOrWq("wife")
			for i := 0; i < 4; i++ {
				swifeid := wife.Findnames(StringStrip(names[i]))
				if swifeid == "" {
					ctx.SendChain(message.Text("-未找到角色" + names[i] + Postfix))
					return
				}
				names[i] = wife.Idmap(swifeid)
				if names[i] == "" {
					ctx.SendChain(message.Text("Idmap数据缺失"))
					return
				} else if names[i] == "空" || names[i] == "荧" || names[i] == "旅行者" {
					ctx.SendChain(message.Text("-暂不支持空/荧伤害数据" + Postfix))
					return
				}
				var t = -1
				// 匹配角色
				for p, v := range alldata.Chars {
					if names[i] == v.Name {
						t = p
					}
				}
				if t == -1 { // 在返回数据中未找到想要的角色
					ctx.SendChain(message.Text("-角色", names[i], "未展示", Postfix))
					return
				} else {
					is[i] = t //匹配成功
				}
			}
		}

		ctx.SendChain(message.Text("-伤害计算中...\n-队伍配置", fmt.Sprint(names)))
		da, err := alldata.Getgroupdata("123456789", is)
		if err != nil {
			ctx.SendChain(message.Text("Error:", err))
			return
		}
		//解析
		var gdate Damgroup
		err = json.Unmarshal(da, &gdate)
		if err != nil {
			ctx.SendChain(message.Text("解析伤害数据错误捏：", err))
			return
		}
		//绘图
		{
			dc := gg.NewContext(1080, 1620)
			dc.SetRGB(1, 1, 1)
			window, err := gg.LoadImage("plugin/kokomi/data/zawu/window.jpg")
			if err != nil {
				ctx.SendChain(message.Text("获取背景失败", err))
				return
			}
			window = resize.Resize(1080, 0, window, resize.Bilinear) //改比例
			dc.DrawImage(window, 0, -100)
			black127 := color.NRGBA{R: 0, G: 0, B: 0, A: 127}
			yingone := Yinying(410, 570, 16, black127)
			yingtwo := Yinying(460, 600, 16, black127)
			yingthree := Yinying(500, 230, 16, black127)
			yingfour := Yinying(500, 400, 16, black127)
			yingzero := Yinying(410, 160, 16, black127)
			one := gg.NewContext(410, 570)
			two := gg.NewContext(460, 600)
			three := gg.NewContext(500, 230)
			four := gg.NewContext(500, 400)
			zero := gg.NewContext(410, 160)
			//图层零,用户信息
			{
				zero.SetRGB(1, 1, 1) //白色
				if err := zero.LoadFontFace(NameFont, 80); err != nil {
					panic(err)
				}
				zero.DrawStringAnchored(string([]rune(names[0])[:1])+
					string([]rune(names[1])[:1])+
					string([]rune(names[2])[:1])+
					string([]rune(names[3])[:1]), 390, 130, 1, 0)
				if err := zero.LoadFontFace(FontFile, 30); err != nil {
					panic(err)
				}
				zero.DrawStringAnchored("昵称:"+alldata.Nickname, 390, 40, 1, 0)
				dc.DrawImage(yingzero, 630, 20)
				dc.DrawImage(zero.Image(), 630, 20)
			}
			//图层1,角色配置
			{
				if err := one.LoadFontFace(FontFile, 40); err != nil {
					panic(err)
				}
				one.SetRGB(1, 1, 1) //白色
				one.DrawString("队伍配置", 30, 50)
				for i := 0; i < 4; i++ {
					turole, err := gg.LoadImage("plugin/kokomi/data/character/" + names[i] + "/imgs/face.webp")
					if err != nil {
						ctx.SendChain(message.Text("获取角色头像图失败", err))
						return
					}
					turole = resize.Resize(100, 0, turole, resize.Bilinear) //缩小
					tuwq, err := gg.LoadPNG("plugin/kokomi/data/wq/" + alldata.Chars[is[i]].Weapon.Name + ".png")
					if err != nil {
						ctx.SendChain(message.Text("获取武器图片失败", err))
						return
					}
					tuwq = resize.Resize(50, 0, tuwq, resize.Bilinear) //缩小
					one.DrawImage(turole, 10+i%2*200, 70+i/2*250)
					one.DrawImage(tuwq, 105+i%2*195, 120+i/2*250)
					if err := one.LoadFontFace(FontFile, 20); err != nil {
						panic(err)
					}
					one.DrawStringWrapped(alldata.Chars[is[i]].Level+"\n"+
						strconv.Itoa(alldata.Chars[is[i]].Talent.A)+"-"+
						strconv.Itoa(alldata.Chars[is[i]].Talent.E)+"-"+
						strconv.Itoa(alldata.Chars[is[i]].Talent.Q)+"\n"+
						Ftoone(alldata.Chars[is[i]].Attr.Cpct)+"%/"+
						Ftoone(alldata.Chars[is[i]].Attr.Cdmg)+"%\n生命"+
						Ftoone(alldata.Chars[is[i]].Attr.Hp)+"\n"+
						strings.ReplaceAll(Sywsuit([]string{
							alldata.Chars[is[i]].Artis.Hua.Set,
							alldata.Chars[is[i]].Artis.Yu.Set,
							alldata.Chars[is[i]].Artis.Sha.Set,
							alldata.Chars[is[i]].Artis.Bei.Set,
							alldata.Chars[is[i]].Artis.Guan.Set}), "+", "\n"),
						float64(20+i%2*200), float64(180+i/2*255), 0, 0, 120, 1.5, gg.AlignLeft)
				}
				dc.DrawImage(yingone, 630, 200)
				dc.DrawImage(one.Image(), 630, 200)
			}
			//图层2,伤害图
			{
				two.SetRGB(1, 1, 1)             //白色
				two.DrawLine(20, 40, 20, 515)   //|
				two.DrawLine(20, 515, 440, 515) //--
				two.SetLineWidth(6)
				//two.Stroke()
				two.StrokePreserve()
				var numchart float64
				if err := two.LoadFontFace(FontFile, 20); err != nil {
					panic(err)
				}
				for _, v := range gdate.Result.ChartData {
					numchart += v.Value
				}
				nn := 440 / (len(gdate.Result.ChartData) + 1)
				for k, v := range gdate.Result.ChartData {
					ss := strings.Split(v.Name, "\n")
					y := -v.Value / numchart * 527
					two.SetRGB(1, 1, 1)
					buff := truncation(two, ss[0], 40) //宽减20
					for i, v := range buff {
						if v != "" {
							two.DrawStringAnchored(v, float64(nn*(k+1))+10, float64(535+i*30), 0.5, 0) // name
						}
					}
					two.DrawStringAnchored(ss[1], float64(nn*(k+1)-8), 515+y, 0.1, -0.5) // 8=0.2*40
					two.DrawRectangle(float64(nn*(k+1)), 515, 20, y)
					two.SetHexColor(v.Label.Color) // 设置画笔颜色为绿色
					two.Fill()                     // 使用当前颜色（绿）填充满当前路径（矩形）所闭合出的区域
				}
				dc.DrawImage(yingtwo, 580, 800)
				dc.DrawImage(two.Image(), 580, 800)
			}
			//图层3,总伤害
			{
				if err := three.LoadFontFace(FontFile, 25); err != nil {
					panic(err)
				}
				three.SetRGB(1, 1, 1) //白色
				strArr := strings.Split(gdate.Result.ZdlTips0, "，")
				three.DrawStringAnchored(strArr[0], 250, 30, 0.5, 0)
				three.DrawStringAnchored(strArr[1], 250, 60, 0.5, 0)
				if err := three.LoadFontFace(FiFile, 120); err != nil {
					panic(err)
				}
				three.SetHexColor("#98F5FF")
				three.DrawStringAnchored(fmt.Sprint(gdate.Result.ZdlResult), 250, 180, 0.5, 0)
				dc.DrawImage(yingthree, 40, 750)
				dc.DrawImage(three.Image(), 40, 750)
			}
			//图层4,手法展示
			{
				if err := four.LoadFontFace(FontFile, 40); err != nil {
					panic(err)
				}
				four.SetRGB(1, 1, 1) //白色
				four.DrawString("操作手法", 30, 40)
				/*	buff := truncation(four, gdate.Result.ComboIntro, 440) //宽减20
					for i, v := range buff {
						if v != "" {
							four.DrawString(v, 30, float64(100+i*35))
						}
					}*/
				strArr := strings.Split(gdate.Result.ComboIntro, ",")
				if err := four.LoadFontFace(FontFile, 30); err != nil {
					panic(err)
				}
				var ws, hs float64
				var c [3]int
				var a = regexp.MustCompile("^[\u4e00-\u9fa5]$")
				for _, v := range strArr {
					if a.MatchString(string([]rune(v)[0:1])) {
						c = randfill()
					}
					four.SetRGB255(c[0], c[1], c[2])
					if ws >= 440 {
						ws = 0
						hs += 50
					}
					ws += DrawStringRec(four, v, "#FFFFFF", ws+5, 50+hs) + 15
				}
				dc.DrawImage(yingfour, 40, 1000)
				dc.DrawImage(four.Image(), 40, 1000)
			}
			// 版本号
			if err := dc.LoadFontFace(BaFile, 30); err != nil {
				panic(err)
			}
			dc.DrawStringAnchored("Created By ZeroBot-Plugin "+kanban.Version+" & "+edition, 540, 1620-30, 0.5, 0.5)
			// 输出图片
			ff, err := imgfactory.ToBytes(dc.Image()) // 图片放入缓存
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			ctx.SendChain(message.ImageBytes(ff)) // 输出
		}
	})
	en.OnRegex(`^更新kokomi$`, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		path := "plugin/kokomi"
		output, err := RunCmd(path, "git pull")
		if err != nil {
			ctx.SendChain(message.Text("更新失败", Postfix, "\n", string(output)))
			return
		}
		ctx.SendChain(message.Text("更新成功", Postfix, "\n", string(output)))
	})
}

// 获取enka数据
func getEnka(ctx *zero.Ctx, uid string) (es []byte, err error) {
	es, err = web.RequestDataWith(web.NewDefaultClient(),
		fmt.Sprintf(url, uid), "GET", "",
		"zerobot-plugin-kokomi",
		nil,
	)
	if err != nil || len(es) == 0 {
		time.Sleep(500 * time.Microsecond) //0.5s
		es, err = web.RequestDataWith(web.NewDefaultClient(),
			fmt.Sprintf(url, uid), "GET", "",
			"zerobot-plugin-kokomi",
			nil,
		)
	}
	return
}
