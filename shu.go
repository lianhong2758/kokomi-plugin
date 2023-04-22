package kokomi // Package kokomi

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/FloatTech/gg"
)

// Fff 圣遗物武器名匹配
type Fff struct {
	WQ map[string]string `json:"zh-CN"`
}

// 评分权重结构
type wifequan struct {
	Hp       int //生命
	Atk      int //攻击力
	Def      int //防御力
	Cpct     int //暴击率
	Cdmg     int //暴击伤害
	Mastery  int //元素精通
	Dmg      int //元素伤害
	Phy      int //物理伤害
	Recharge int //元素充能
	Heal     int //治疗加成
}

// config内容
type config struct {
	Apis    []string `json:"apis"`
	Apiid   int      `json:"api_id"`
	Postfix string   `json:"postfix"`
	Edition string   `json:"edition"`
}

// Wikimap wiki查询地址结构解析
type Wikimap struct {
	Card      map[string]string `json:"card"`
	Matera    map[string]string `json:"material for role"`
	Specialty map[string]string `json:"specialty"`
	Weapon    map[string]string `json:"weapon"`
}

// Role 角色信息json解析
type Role struct {
	TalentID   map[string]string `json:"talentId"`
	Elem       string            `json:"elem"`
	TalentCons struct {
		E int `json:"e"`
		Q int `json:"q"`
	} `json:"talentCons"`
}

// Dam 角色伤害解析
type Dam struct {
	Result []struct {
		DamageResultArr []struct {
			Title  string      `json:"title"`
			Value  interface{} `json:"value"`
			Expect string      `json:"expect"`
		} `json:"damage_result_arr"`
	} `json:"result"`
}
type Damgroup struct {
	Result struct {
		ZdlResult any        `json:"zdl_result"`
		ZdlTips0  string     `json:"zdl_tips0"`
		ZdlTips1  string     `json:"zdl_tips1"`
		ZdlTips2  string     `json:"zdl_tips2"`
		ZdlTips3  string     `json:"zdl_tips3"`
		ChartData []struct { //统计图数据
			Name  string  `json:"name"`
			Ename string  `json:"ename"`
			Value float64 `json:"value"`
			Label struct {
				Color string `json:"color"`
			} `json:"label"`
		} `json:"chart_data"`
		ComboIntro   string `json:"combo_intro"`
		RechargeInfo []struct {
			Ename    string  `json:"ename"`
			Energy   float64 `json:"energy"`
			Rate     string  `json:"rate"`
			Height   int     `json:"height"`
			Color    string  `json:"color"`
			Recharge string  `json:"recharge"`
		} `json:"recharge_info"`
	} `json:"result"`
}

// Getuid qquid->uid
func Getuid(sqquid string) (uid int) { // 获取对应游戏uid
	// 获取本地缓存数据
	txt, err := os.ReadFile("plugin/kokomi/data/uid/" + sqquid + ".kokomi")
	if err != nil {
		return 0
	}
	uid, _ = strconv.Atoi(string(txt))
	return
}

// StoS 圣遗物词条简单描述
func StoS(val string) string {
	switch val {
	case "FIGHT_PROP_HP":
		return "小生命"
	case "FIGHT_PROP_HP_PERCENT":
		return "大生命"
	case "FIGHT_PROP_ATTACK":
		return "小攻击"
	case "FIGHT_PROP_ATTACK_PERCENT":
		return "大攻击"
	case "FIGHT_PROP_DEFENSE":
		return "小防御"
	case "FIGHT_PROP_DEFENSE_PERCENT":
		return "大防御"
	case "FIGHT_PROP_CRITICAL":
		return "暴击率"
	case "FIGHT_PROP_CRITICAL_HURT":
		return "暴击伤害"
	case "FIGHT_PROP_CHARGE_EFFICIENCY":
		return "元素充能"
	case "FIGHT_PROP_HEAL_ADD":
		return "治疗加成"
	case "FIGHT_PROP_ELEMENT_MASTERY":
		return "元素精通"
	case "FIGHT_PROP_PHYSICAL_ADD_HURT":
		return "物理加伤"
	case "FIGHT_PROP_FIRE_ADD_HURT":
		return "火元素加伤"
	case "FIGHT_PROP_ELEC_ADD_HURT":
		return "雷元素加伤"
	case "FIGHT_PROP_WATER_ADD_HURT":
		return "水元素加伤"
	case "FIGHT_PROP_GRASS_ADD_HURT":
		return "草元素加伤"
	case "FIGHT_PROP_WIND_ADD_HURT":
		return "风元素加伤"
	case "FIGHT_PROP_ROCK_ADD_HURT":
		return "岩元素加伤"
	case "FIGHT_PROP_ICE_ADD_HURT":
		return "冰元素加伤"
	}
	return ""
}

func GetAppendProp(v string) string {
	switch v {
	case "FIGHT_PROP_HP", "FIGHT_PROP_HP_PERCENT", "小生命", "大生命":
		return "生命值"
	case "FIGHT_PROP_ATTACK", "FIGHT_PROP_ATTACK_PERCENT", "小攻击", "大攻击":
		return "攻击力"
	case "FIGHT_PROP_DEFENSE", "FIGHT_PROP_DEFENSE_PERCENT", "小防御", "大防御":
		return "防御力"
	case "FIGHT_PROP_CRITICAL", "暴击率":
		return "暴击率"
	case "FIGHT_PROP_CRITICAL_HURT", "暴击伤害":
		return "暴击伤害"
	case "FIGHT_PROP_CHARGE_EFFICIENCY", "元素充能":
		return "元素充能效率"
	case "FIGHT_PROP_HEAL_ADD", "治疗加成":
		return "治疗加成"
	case "FIGHT_PROP_ELEMENT_MASTERY", "元素精通":
		return "元素精通"
	case "FIGHT_PROP_PHYSICAL_ADD_HURT", "物理加伤":
		return "物理伤害加成"
	case "FIGHT_PROP_FIRE_ADD_HURT", "火元素加伤":
		return "火元素伤害加成"
	case "FIGHT_PROP_ELEC_ADD_HURT", "雷元素加伤":
		return "雷元素伤害加成"
	case "FIGHT_PROP_WATER_ADD_HURT", "水元素加伤":
		return "水元素伤害加成"
	case "FIGHT_PROP_GRASS_ADD_HURT", "草元素加伤":
		return "草元素伤害加成"
	case "FIGHT_PROP_WIND_ADD_HURT", "风元素加伤":
		return "风元素伤害加成"
	case "FIGHT_PROP_ROCK_ADD_HURT", "岩元素加伤":
		return "岩元素伤害加成"
	case "FIGHT_PROP_ICE_ADD_HURT", "冰元素加伤":
		return "冰元素伤害加成"
	}
	return ""
}

func GetEquipType(v string) string {
	switch v {
	case "EQUIP_BRACER", "0":
		return "生之花"
	case "EQUIP_NECKLACE", "1":
		return "死之羽"
	case "EQUIP_SHOES", "2":
		return "时之沙"
	case "EQUIP_RING", "3":
		return "空之杯"
	case "EQUIP_DRESS", "4":
		return "理之冠"
	}
	return ""
}

// Stofen 判断词条分号
func Stofen(val string) string {
	switch val {
	case "FIGHT_PROP_HP", "FIGHT_PROP_ATTACK", "FIGHT_PROP_DEFENSE", "FIGHT_PROP_ELEMENT_MASTERY", "小生命", "小攻击", "小防御", "元素精通":
		return ""
		/*
			case "FIGHT_PROP_HP_PERCENT":
			case "FIGHT_PROP_ATTACK_PERCENT":
			case "FIGHT_PROP_DEFENSE_PERCENT":
			case "FIGHT_PROP_CRITICAL":
			case "FIGHT_PROP_CRITICAL_HURT":
			case "FIGHT_PROP_CHARGE_EFFICIENCY":
			case "FIGHT_PROP_HEAL_ADD":
			case "FIGHT_PROP_PHYSICAL_ADD_HURT":
			case "FIGHT_PROP_FIRE_ADD_HURT":
			case "FIGHT_PROP_ELEC_ADD_HURT":
			case "FIGHT_PROP_WATER_ADD_HURT":
			case "FIGHT_PROP_GRASS_ADD_HURT":
			case "FIGHT_PROP_WIND_ADD_HURT":
			case "FIGHT_PROP_ROCK_ADD_HURT":
			case "FIGHT_PROP_ICE_ADD_HURT":
		*/
	}
	return "%"
}

// Countcitiao 计算圣遗物单词条分
func Countcitiao(wifename, funame string, figure float64) float64 {
	ti := Wifequanmap[wifename]
	switch funame {
	case "大生命":
		return figure * 1.33 * float64(ti.Hp) / 100
	case "大攻击":
		return figure * 1.33 * float64(ti.Atk) / 100
	case "大防御":
		return figure * 1.33 * float64(ti.Def) / 100
	case "暴击率":
		return figure * 2.0 * float64(ti.Cpct) / 100
	case "暴击伤害":
		return figure * 1.0 * float64(ti.Cdmg) / 100
	case "元素精通":
		return figure * 0.33 * float64(ti.Mastery) / 100
	case "雷元素加伤", "水元素加伤", "火元素加伤", "风元素加伤", "草元素加伤", "岩元素加伤", "冰元素加伤", "元素加伤":
		return figure * 1.33 * float64(ti.Dmg) / 100
	case "物理加伤":
		return figure * 1.33 * float64(ti.Phy) / 100
	case "元素充能":
		return figure * 1.2 * float64(ti.Recharge) / 100
	case "治疗加成":
		return figure * 1.73 * float64(ti.Heal) / 100
	}
	return 0
}

// Pingji 词条评级
func Pingji(val float64) string {
	switch {
	case val < 10:
		return "D"
	case val < 16.5:
		return "C"
	case val < 23.1:
		return "B"
	case val < 29.7:
		return "A"
	case val < 36.3:
		return "S"
	case val < 42.9:
		return "SS"
	case val < 49.5:
		return "SSS"
	case val < 56.1:
		return "ACE"
	}
	return "ACES"
}

// Ftoone 保留一位小数并转化string
func Ftoone(f float64) string {
	// return strconv.FormatFloat(f, 'f', 1, 64)
	if f == 0 {
		return "0"
	}
	return strconv.FormatFloat(f, 'f', 1, 64)
}

// FindMap 各种简称map查询
type FindMap map[string][]string

func GetWifeOrWq(val string) FindMap {
	var txt []byte
	switch val {
	case "wife":
		txt, _ = os.ReadFile("plugin/kokomi/data/json/wife_list.json")
	case "wq":
		txt, _ = os.ReadFile("plugin/kokomi/data/json/wq.json")
	}
	var m FindMap = make(map[string][]string)
	if nil == json.Unmarshal(txt, &m) {
		return m
	}
	return nil
}

// Findnames 遍历寻找匹配昵称
func (m FindMap) Findnames(val string) string {
	for k, v := range m {
		for _, vv := range v {
			if vv == val {
				return k
			}
		}
	}
	return ""
}

// Idmap wifeid->wifename
func (m FindMap) Idmap(val string) string {
	f, b := m[val]
	if !b {
		return ""
	}
	return f[0]
}

// StringStrip 字符串删空格
func StringStrip(input string) string {
	if input == "" {
		return ""
	}
	reg := regexp.MustCompile(`[\s\p{Zs}]{1,}`)
	return reg.ReplaceAllString(input, "")
}

// GetReliquary 读取圣遗物信息
func GetReliquary() *Fff {
	txt, err := os.ReadFile("plugin/kokomi/data/json/loc.json")
	if err != nil {
		return nil
	}
	var p Fff
	if nil == json.Unmarshal(txt, &p) {
		return &p
	}
	return nil
}

// Findwq Findwq圣遗物,武器名匹配
func (m *Fff) Findwq(a string) string {
	return m.WQ[a]
}

// GetRole 角色信息
func GetRole(str string) (*Role, error) {
	txt, err := os.ReadFile("plugin/kokomi/data/character/" + str + "/data.json")
	if err != nil {
		return nil, err
	}
	var p Role
	if err = json.Unmarshal(txt, &p); err == nil {
		return &p, nil
	}
	return nil, err
}

// GetTalentId 天赋列表
func (m *Role) GetTalentId() []int {
	f := make([]int, 3)
	for k, v := range m.TalentID {
		switch v {
		case "a":
			f[0], _ = strconv.Atoi(k)
		case "e":
			f[1], _ = strconv.Atoi(k)
		case "q":
			f[2], _ = strconv.Atoi(k)
		}
	}
	return f
}

// Syws 圣遗物列表名解析
type Syws map[string]struct {
	Name string `json:"name"`
	Sets struct {
		Num1 struct {
			Name string `json:"name"`
		} `json:"1"`
		Num2 struct {
			Name string `json:"name"`
		} `json:"2"`
		Num3 struct {
			Name string `json:"name"`
		} `json:"3"`
		Num4 struct {
			Name string `json:"name"`
		} `json:"4"`
		Num5 struct {
			Name string `json:"name"`
		} `json:"5"`
	} `json:"sets"`
}

func GetSywName() Syws {
	data, err := os.ReadFile("plugin/kokomi/data/json/sywname_list.json")
	if err != nil {
		return nil
	}
	var p Syws
	if nil == json.Unmarshal(data, &p) {
		return p
	}
	return nil
}

// Names 圣遗物名列表
func (m Syws) Names(syw string) []string {
	for _, v := range m {
		if v.Name == syw {
			return []string{
				v.Sets.Num1.Name,
				v.Sets.Num2.Name,
				v.Sets.Num3.Name,
				v.Sets.Num4.Name,
				v.Sets.Num5.Name,
			}
		}
	}
	return nil
}

// Sywsuit 圣遗物套装判断
func Sywsuit(syws []string) string {
	sywMap := make(map[string]int)
	var c0, c1 string
	for _, v := range syws {
		i := sywMap[v]
		sywMap[v] = i + 1
	}
	sywMap[""] = 0
	for k, v := range sywMap {
		if v >= 4 {
			return k + "4"
		}
		if v >= 2 {
			if c0 == "" {
				c0 = k
			} else {
				c1 = k
			}
		}
	}
	if c0 != "" {
		if c1 != "" {
			return c0 + "2+" + c1 + "2"
		}
		return c0 + "2"
	}
	return "+"
}

// 解析为本地结构
func (n Data) ConvertData() (Thisdata, error) {
	t := new(Thisdata)
	t.Chars = make(map[int]CharRole)
	wife := GetWifeOrWq("wife")
	t.UID = n.UID
	t.Nickname = n.PlayerInfo.Nickname
	t.Level = n.PlayerInfo.Level
	for k, v := range n.AvatarInfoList {
		name := wife.Idmap(strconv.Itoa(v.AvatarID))
		//数据处理区
		adds, addf := "元素加伤:", 0.0
		if v.FightPropMap.Num30*100 > addf {
			adds = "物理加伤:"
			addf = v.FightPropMap.Num30 * 100
		}
		if v.FightPropMap.Num40*100 > addf {
			adds = "火元素加伤:"
			addf = v.FightPropMap.Num40 * 100
		}
		if v.FightPropMap.Num41*100 > addf {
			adds = "雷元素加伤:"
			addf = v.FightPropMap.Num41 * 100
		}
		if v.FightPropMap.Num42*100 > addf {
			adds = "水元素加伤:"
			addf = v.FightPropMap.Num42 * 100
		}
		if v.FightPropMap.Num44*100 > addf {
			adds = "风元素加伤:"
			addf = v.FightPropMap.Num44 * 100
		}
		if v.FightPropMap.Num45*100 > addf {
			adds = "岩元素加伤:"
			addf = v.FightPropMap.Num45 * 100
		}
		if v.FightPropMap.Num46*100 > addf {
			adds = "冰元素加伤:"
			addf = v.FightPropMap.Num46 * 100
		}
		if v.FightPropMap.Num43*100 > addf {
			adds = "草元素加伤:"
			addf = v.FightPropMap.Num43 * 100
		}
		l := len(v.EquipList)
		reliquary := GetReliquary()
		if reliquary == nil {
			return *t, errors.New("1")
		}
		wq := reliquary.WQ[v.EquipList[l-1].Flat.NameTextHash]
		if wq == "" {
			return *t, errors.New("2")
		}
		var wqjl = 0
		for m := range v.EquipList[l-1].Weapon.AffixMap {
			wqjl = m
		}
		role, err := GetRole(wife.Idmap(strconv.Itoa(v.AvatarID)))
		if err != nil {
			return *t, err
		}
		talentId := role.GetTalentId()
		syw := GetSywName()
		var sywhua, sywyu, sywsha, sywbei, sywguan sywm
		for i := 0; i < l-1; i++ {
			switch v.EquipList[i].Flat.EquipType {
			case "EQUIP_BRACER":
				sywhua = sywm{
					Set:   reliquary.WQ[v.EquipList[i].Flat.SetNameTextHash],
					Name:  syw.Names(reliquary.WQ[v.EquipList[i].Flat.SetNameTextHash])[0],
					Level: v.EquipList[i].Reliquary.Level,
					Main: attrs{
						Title: StoS(v.EquipList[i].Flat.ReliquaryMainStat.MainPropID),
						Value: v.EquipList[i].Flat.ReliquaryMainStat.Value,
					},
				}
				for _, stats := range v.EquipList[i].Flat.ReliquarySubStats {
					sywhua.Attrs = append(sywhua.Attrs, attrs{
						Title: StoS(stats.SubPropID),
						Value: stats.Value,
					})
				}
			case "EQUIP_NECKLACE":
				sywyu = sywm{
					Set:   reliquary.WQ[v.EquipList[i].Flat.SetNameTextHash],
					Name:  syw.Names(reliquary.WQ[v.EquipList[i].Flat.SetNameTextHash])[1],
					Level: v.EquipList[i].Reliquary.Level,
					Main: attrs{
						Title: StoS(v.EquipList[i].Flat.ReliquaryMainStat.MainPropID),
						Value: v.EquipList[i].Flat.ReliquaryMainStat.Value,
					},
				}
				for _, stats := range v.EquipList[i].Flat.ReliquarySubStats {
					sywyu.Attrs = append(sywyu.Attrs, attrs{
						Title: StoS(stats.SubPropID),
						Value: stats.Value,
					})
				}
			case "EQUIP_SHOES":
				sywsha = sywm{
					Set:   reliquary.WQ[v.EquipList[i].Flat.SetNameTextHash],
					Name:  syw.Names(reliquary.WQ[v.EquipList[i].Flat.SetNameTextHash])[2],
					Level: v.EquipList[i].Reliquary.Level,
					Main: attrs{
						Title: StoS(v.EquipList[i].Flat.ReliquaryMainStat.MainPropID),
						Value: v.EquipList[i].Flat.ReliquaryMainStat.Value,
					},
				}
				for _, stats := range v.EquipList[i].Flat.ReliquarySubStats {
					sywsha.Attrs = append(sywsha.Attrs, attrs{
						Title: StoS(stats.SubPropID),
						Value: stats.Value,
					})
				}
			case "EQUIP_RING":
				sywbei = sywm{
					Set:   reliquary.WQ[v.EquipList[i].Flat.SetNameTextHash],
					Name:  syw.Names(reliquary.WQ[v.EquipList[i].Flat.SetNameTextHash])[3],
					Level: v.EquipList[i].Reliquary.Level,
					Main: attrs{
						Title: StoS(v.EquipList[i].Flat.ReliquaryMainStat.MainPropID),
						Value: v.EquipList[i].Flat.ReliquaryMainStat.Value,
					},
				}
				for _, stats := range v.EquipList[i].Flat.ReliquarySubStats {
					sywbei.Attrs = append(sywbei.Attrs, attrs{
						Title: StoS(stats.SubPropID),
						Value: stats.Value,
					})
				}
			case "EQUIP_DRESS":
				sywguan = sywm{
					Set:   reliquary.WQ[v.EquipList[i].Flat.SetNameTextHash],
					Name:  syw.Names(reliquary.WQ[v.EquipList[i].Flat.SetNameTextHash])[4],
					Level: v.EquipList[i].Reliquary.Level,
					Main: attrs{
						Title: StoS(v.EquipList[i].Flat.ReliquaryMainStat.MainPropID),
						Value: v.EquipList[i].Flat.ReliquaryMainStat.Value,
					},
				}
				for _, stats := range v.EquipList[i].Flat.ReliquarySubStats {
					sywguan.Attrs = append(sywguan.Attrs, attrs{
						Title: StoS(stats.SubPropID),
						Value: stats.Value,
					})
				}
			}
		}
		//导入
		t.Chars[k] = CharRole{
			ID:     v.AvatarID,
			Name:   name,
			Level:  v.PropMap.Num4001.Val,
			Fetter: v.FetterInfo.ExpLevel,
			Cons:   len(v.TalentIDList),
			Attr: attr{
				Atk:      v.FightPropMap.Num4*(1+(v.FightPropMap.Num6)) + v.FightPropMap.Num5,
				AtkBase:  v.FightPropMap.Num4,
				Def:      v.FightPropMap.Num2002,
				DefBase:  v.FightPropMap.Num7,
				Hp:       v.FightPropMap.Num2000,
				HpBase:   v.FightPropMap.Num1,
				Mastery:  v.FightPropMap.Num28,
				Recharge: v.FightPropMap.Num23 * 100,
				Heal:     v.FightPropMap.Num26 * 100,
				Cpct:     v.FightPropMap.Num20 * 100,
				Cdmg:     v.FightPropMap.Num22 * 100,
				Dmg:      addf,
				DmgName:  adds,
				Phy:      v.FightPropMap.Num30 * 100,
			},
			Weapon: weapon{
				Name:  wq,
				Star:  v.EquipList[l-1].Flat.RankLevel,
				Level: v.EquipList[l-1].Weapon.Level,
				Affix: v.EquipList[l-1].Weapon.AffixMap[wqjl] + 1,
				Atk:   v.EquipList[l-1].Flat.WeaponStat[0].Value,
			},
			Talent: talent{
				A: v.SkillLevelMap[talentId[0]],
				E: v.SkillLevelMap[talentId[1]],
				Q: v.SkillLevelMap[talentId[2]],
			},
			Artis: artis{
				Hua:  sywhua,
				Yu:   sywyu,
				Sha:  sywsha,
				Bei:  sywbei,
				Guan: sywguan,
			},
			DataSource: "Enka.Network",
		}
	}
	return *t, nil
}

// 合并映射
func (t *Thisdata) MergeFile(suid string) {
	tx, err := os.ReadFile("plugin/kokomi/data/js/" + suid + ".kokomi")
	if err != nil {
		return
	}
	// 解析
	var alldata Thisdata
	err = json.Unmarshal(tx, &alldata)
	if err != nil {
		return
	}
OuterLoop:
	for i := 0; i < len(alldata.Chars); i++ {
		for l := 0; l < len(t.Chars); l++ {
			if alldata.Chars[i].Name == t.Chars[l].Name {
				if i == len(alldata.Chars)-1 {
					return
				} else {
					continue OuterLoop
				}
			}
		}
		//未找到相同
		t.Chars[len(t.Chars)] = alldata.Chars[i]
	}
	return
}

// 字符串分行
func truncation(canvas *gg.Context, text string, width int) (buff []string) {
	buff = make([]string, 0, 32)
	s := bufio.NewScanner(strings.NewReader(text))
	line := strings.Builder{}
	for s.Scan() {
		for _, v := range s.Text() {
			length, _ := canvas.MeasureString(line.String())
			if int(length) <= width {
				line.WriteRune(v)
			} else {
				buff = append(buff, line.String())
				line.Reset()
				line.WriteRune(v)
			}
		}
		buff = append(buff, line.String())
		line.Reset()
	}
	return
}

// cmd后台执行
func RunCmd(path, order string) (output []byte, err error) {
	var cmd *exec.Cmd
	cmd = exec.Command("bash", "-c", "cd "+path+" ; "+order)
	output, err = cmd.CombinedOutput()
	return
}
