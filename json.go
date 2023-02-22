package kokomi // Package kokomi

import ()

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

	RankLevel int    `json:"rankLevel"` // 3, 4 or 5
	ItemType  string `json:"itemType"`  // ITEM_WEAPON or ITEM_RELIQUARY
	Icon      string `json:"icon"`      // You can get the icon from https://enka.network/ui/{Icon}.png
}

// Stat ...  属性对
type Stat struct {
	MainPropID string  `json:"mainPropId,omitempty"`
	SubPropID  string  `json:"appendPropId,omitempty"`
	Value      float64 `json:"statValue"`
}

// 本地数据
type Thisdata struct {
	UID      string           `json:"uid"`
	Nickname string           `json:"nickname"`
	Level    int              `json:"level"`
	Chars    map[int]CharRole `json:"chars"`
}

type CharRole struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Level      string `json:"level"`
	Cons       int    `json:"cons"`   //命之座
	Fetter     int    `json:"fetter"` //好感
	Attr       attr   `json:"attr"`
	Weapon     weapon `json:"weapon"`
	Talent     talent `json:"talent"`
	Artis      artis  `json:"artis"`
	DataSource string `json:"dataSource"`
}

type attr struct {
	Atk      float64 `json:"atk"`
	AtkBase  float64 `json:"atkBase"`
	Def      float64 `json:"def"`
	DefBase  float64 `json:"defBase"`
	Hp       float64 `json:"hp"`
	HpBase   float64 `json:"hpBase"`
	Mastery  float64 `json:"mastery"`
	Recharge float64 `json:"recharge"`
	Heal     float64 `json:"heal"`
	Cpct     float64 `json:"cpct"`
	Cdmg     float64 `json:"cdmg"`
	Dmg      float64 `json:"dmg"`
	DmgName  string  `json:"dmgname"`
	Phy      float64 `json:"phy"`
}

type weapon struct {
	Name  string  `json:"name"`
	Star  int     `json:"star"`
	Level int     `json:"level"`
	Affix int     `json:"affix"`
	Atk   float64 `json:"akt"` //攻击力,可能不存在
}
type talent struct {
	A int `json:"a"`
	E int `json:"e"`
	Q int `json:"q"`
}

type sywm struct {
	Name  string  `json:"name"`
	Set   string  `json:"set"`
	Level int     `json:"level"`
	Main  attrs   `json:"main"`
	Attrs []attrs `json:"attrs"`
}

type attrs struct {
	Title string  `json:"title"`
	Value float64 `json:"value"`
}
type artis struct {
	Hua  sywm `json:"1"`
	Yu   sywm `json:"2"`
	Sha  sywm `json:"3"`
	Bei  sywm `json:"4"`
	Guan sywm `json:"5"`
}
