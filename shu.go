package kokomi // 导入yuan-shen模块

import (
	"encoding/json"
	"os"
	"regexp"
	"strconv"
)

// 圣遗物武器名匹配
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
	Apis     []string `json:"apis"`
	Apiid    int      `json:"api_id"`
	Postfix  string   `json:"postfix"`
	Datafrom string   `json:"from"`
	Edition  string   `json:"edition"`
}

// wiki查询地址结构解析
type Wikimap struct {
	Card      map[string]string `json:"card"`
	Matera    map[string]string `json:"material for role"`
	Specialty map[string]string `json:"specialty"`
	Weapon    map[string]string `json:"weapon"`
}

// 角色信息json解析
type Role struct {
	TalentID   map[int]int    `json:"talentId"`
	TalentKey  map[int]string `json:"talentKey"`
	Elem       string         `json:"elem"`
	TalentCons struct {
		E int `json:"e"`
		Q int `json:"q"`
	} `json:"talentCons"`
}

// 角色伤害解析
type Dam struct {
	Result []struct {
		DamageResultArr []struct {
			Title  string      `json:"title"`
			Value  interface{} `json:"value"`
			Expect string      `json:"expect"`
		} `json:"damage_result_arr"`
	} `json:"result"`
}

// Data 从网站获取的数据
type Data struct {
	PlayerInfo struct {
		Nickname             string `json:"nickname"`
		Level                int    `json:"level"`
		Signature            string `json:"signature"`
		WorldLevel           int    `json:"worldLevel"`
		NameCardID           int    `json:"nameCardId"`
		FinishAchievementNum int    `json:"finishAchievementNum"`
		TowerFloorIndex      int    `json:"towerFloorIndex"`
		TowerLevelIndex      int    `json:"towerLevelIndex"`
		ShowAvatarInfoList   []struct {
			AvatarID  int `json:"avatarId"`
			Level     int `json:"level"`
			CostumeID int `json:"costumeId,omitempty"`
		} `json:"showAvatarInfoList"`
		ShowNameCardIDList []int `json:"showNameCardIdList"`
		ProfilePicture     struct {
			AvatarID int `json:"avatarId"`
		} `json:"profilePicture"`
	} `json:"playerInfo"`
	AvatarInfoList []struct {
		AvatarID int `json:"avatarId"`
		PropMap  struct {
			Num1001 struct {
				Type int    `json:"type"`
				Ival string `json:"ival"`
			} `json:"1001"`
			Num1002 struct {
				Type int    `json:"type"`
				Ival string `json:"ival"`
				Val  string `json:"val"`
			} `json:"1002"`
			Num1003 struct {
				Type int    `json:"type"`
				Ival string `json:"ival"`
			} `json:"1003"`
			Num1004 struct {
				Type int    `json:"type"`
				Ival string `json:"ival"`
			} `json:"1004"`
			Num4001 struct {
				Type int    `json:"type"`
				Ival string `json:"ival"`
				Val  string `json:"val"`
			} `json:"4001"`
			Num10010 struct {
				Type int    `json:"type"`
				Ival string `json:"ival"`
				Val  string `json:"val"`
			} `json:"10010"`
		} `json:"propMap"`
		FightPropMap struct {
			Num1    float64 `json:"1"`
			Num2    float64 `json:"2"`
			Num3    float64 `json:"3"`
			Num4    float64 `json:"4"`
			Num5    float64 `json:"5"`
			Num6    float64 `json:"6"`
			Num7    float64 `json:"7"`
			Num8    float64 `json:"8"`
			Num20   float64 `json:"20"`
			Num21   float64 `json:"21"`
			Num22   float64 `json:"22"`
			Num23   float64 `json:"23"`
			Num26   float64 `json:"26"`
			Num27   float64 `json:"27"`
			Num28   float64 `json:"28"`
			Num29   float64 `json:"29"`
			Num30   float64 `json:"30"`
			Num40   float64 `json:"40"`
			Num41   float64 `json:"41"`
			Num42   float64 `json:"42"`
			Num43   float64 `json:"43"`
			Num44   float64 `json:"44"`
			Num45   float64 `json:"45"`
			Num46   float64 `json:"46"`
			Num50   float64 `json:"50"`
			Num51   float64 `json:"51"`
			Num52   float64 `json:"52"`
			Num53   float64 `json:"53"`
			Num54   float64 `json:"54"`
			Num55   float64 `json:"55"`
			Num56   float64 `json:"56"`
			Num70   float64 `json:"70"`
			Num80   float64 `json:"80"`
			Num1000 float64 `json:"1000"`
			Num1010 float64 `json:"1010"`
			Num2000 float64 `json:"2000"`
			Num2001 float64 `json:"2001"`
			Num2002 float64 `json:"2002"`
			Num2003 float64 `json:"2003"`
			Num3007 float64 `json:"3007"`
			Num3008 float64 `json:"3008"`
			Num3015 float64 `json:"3015"`
			Num3016 float64 `json:"3016"`
			Num3017 float64 `json:"3017"`
			Num3018 float64 `json:"3018"`
			Num3019 float64 `json:"3019"`
			Num3020 float64 `json:"3020"`
			Num3021 float64 `json:"3021"`
			Num3022 float64 `json:"3022"`
			Num3045 float64 `json:"3045"`
			Num3046 float64 `json:"3046"`
		} `json:"fightPropMap"`
		SkillDepotID           int         `json:"skillDepotId"`
		InherentProudSkillList []int       `json:"inherentProudSkillList"`
		SkillLevelMap          map[int]int `json:"skillLevelMap"`
		EquipList              []struct {
			ItemID    int `json:"itemId"`
			Reliquary struct {
				Level            int   `json:"level"`
				MainPropID       int   `json:"mainPropId"`
				AppendPropIDList []int `json:"appendPropIdList"`
			} `json:"reliquary,omitempty"`
			Flat   Flat `json:"flat"` //标记
			Weapon struct {
				Level        int         `json:"level"`
				PromoteLevel int         `json:"promoteLevel"`
				AffixMap     map[int]int `json:"affixMap"`
			} `json:"weapon,omitempty"`
		} `json:"equipList"`
		FetterInfo struct {
			ExpLevel int `json:"expLevel"`
		} `json:"fetterInfo"`
		TalentIDList            []int `json:"talentIdList,omitempty"`
		ProudSkillExtraLevelMap struct {
			Num4239 int `json:"4239"`
		} `json:"proudSkillExtraLevelMap,omitempty"`
		CostumeID int `json:"costumeId,omitempty"`
	} `json:"avatarInfoList"`
	TTL int    `json:"ttl"`
	UID string `json:"uid"`
}

// Flat ... 详细数据
type Flat struct {
	// l10n
	NameTextHash    string `json:"nameTextMapHash"`
	SetNameTextHash string `json:"setNameTextMapHash,omitempty"`

	// artifact
	ReliquaryMainStat Stat   `json:"reliquaryMainstat,omitempty"`
	ReliquarySubStats []Stat `json:"reliquarySubstats,omitempty"`
	EquipType         string `json:"equipType,omitempty"`

	// weapon
	WeaponStat []Stat `json:"weaponStats,omitempty"`

	RankLevel uint8  `json:"rankLevel"` // 3, 4 or 5
	ItemType  string `json:"itemType"`  // ITEM_WEAPON or ITEM_RELIQUARY
	Icon      string `json:"icon"`      // You can get the icon from https://enka.network/ui/{Icon}.png
}

// Stat ...  属性对
type Stat struct {
	MainPropID string  `json:"mainPropId,omitempty"`
	SubPropID  string  `json:"appendPropId,omitempty"`
	Value      float64 `json:"statValue"`
}

// Getuid qquid->uid
func Getuid(qquid int64) (uid int) { // 获取对应游戏uid
	sqquid := strconv.Itoa(int(qquid))
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
	case "FIGHT_PROP_HP", "FIGHT_PROP_HP_PERCENT":
		return "生命值"
	case "FIGHT_PROP_ATTACK", "FIGHT_PROP_ATTACK_PERCENT":
		return "攻击力"
	case "FIGHT_PROP_DEFENSE", "FIGHT_PROP_DEFENSE_PERCENT":
		return "防御力"
	case "FIGHT_PROP_CRITICAL":
		return "暴击率"
	case "FIGHT_PROP_CRITICAL_HURT":
		return "暴击伤害"
	case "FIGHT_PROP_CHARGE_EFFICIENCY":
		return "元素充能效率"
	case "FIGHT_PROP_HEAL_ADD":
		return "治疗加成"
	case "FIGHT_PROP_ELEMENT_MASTERY":
		return "元素精通"
	case "FIGHT_PROP_PHYSICAL_ADD_HURT":
		return "物理伤害加成"
	case "FIGHT_PROP_FIRE_ADD_HURT":
		return "火元素伤害加成"
	case "FIGHT_PROP_ELEC_ADD_HURT":
		return "雷元素伤害加成"
	case "FIGHT_PROP_WATER_ADD_HURT":
		return "水元素伤害加成"
	case "FIGHT_PROP_GRASS_ADD_HURT":
		return "草元素伤害加成"
	case "FIGHT_PROP_WIND_ADD_HURT":
		return "风元素伤害加成"
	case "FIGHT_PROP_ROCK_ADD_HURT":
		return "岩元素伤害加成"
	case "FIGHT_PROP_ICE_ADD_HURT":
		return "冰元素伤害加成"
	}
	return ""
}

func GetEquipType(v string) string {
	switch v {
	case "EQUIP_BRACER":
		return "生之花"
	case "EQUIP_NECKLACE":
		return "死之羽"
	case "EQUIP_SHOES":
		return "时之沙"
	case "EQUIP_RING":
		return "空之杯"
	case "EQUIP_DRESS":
		return "理之冠"
	}
	return ""
}

// Stofen 判断词条分号
func Stofen(val string) string {
	switch val {
	case "FIGHT_PROP_HP", "FIGHT_PROP_ATTACK", "FIGHT_PROP_DEFENSE", "FIGHT_PROP_ELEMENT_MASTERY":
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

// Tianfujiuzhen 修复部分贴图大小错误
func Tianfujiuzhen(val string) int {
	switch val {
	case "芭芭拉", "北斗", "多莉", "甘雨", "胡桃", "科莱", "雷电将军", "罗莎莉亚", "凝光", "赛诺", "魈", "行秋", "烟绯", "夜兰", "早柚":
		return 280
	}
	return 257
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
	case "雷元素加伤", "水元素加伤", "火元素加伤", "风元素加伤", "草元素加伤", "岩元素加伤", "冰元素加伤":
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
	case val < 18:
		return "C"
	case val < 24:
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

// 各种简称map查询
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
	for k, v := range m {
		if k == val {
			return v[0]
		}
	}
	return ""
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

// Findwq圣遗物,武器名匹配
func (m *Fff) Findwq(a string) string {
	return m.WQ[a]
}

// GetRole 角色信息
func GetRole(str string) *Role {
	txt, err := os.ReadFile("plugin/kokomi/data/character/" + str + "/data.json")
	if err != nil {
		return nil
	}
	var p Role
	if nil == json.Unmarshal(txt, &p) {
		return &p
	}
	return nil
}

// GetTalentId 天赋列表
func (m *Role) GetTalentId() []int {
	var a, e, q int
	for k, v := range m.TalentKey {
		switch v {
		case "a":
			a = k
		case "e":
			e = k
		case "q":
			q = k
		}
	}
	f := make([]int, 3)
	for k, v := range m.TalentID {
		switch v {
		case a:
			f[0] = k
		case e:
			f[1] = k
		case q:
			f[2] = k
		}
	}
	return f
}

// 圣遗物列表名解析
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

// 圣遗物名列表
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

// 圣遗物套装判断
func Sywsuit(syws []string) string {
	syw_map := make(map[string]int)
	var c0, c1 string
	for _, v := range syws {
		i := syw_map[v]
		syw_map[v] = i + 1
	}
	syw_map[""] = 0
	for k, v := range syw_map {
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
