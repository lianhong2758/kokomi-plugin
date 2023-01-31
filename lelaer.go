package kokomi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/FloatTech/floatbox/web"
)

const (
	// k_lelaer_damage = "https://api.lelaer.com/ys/getDamageResult.php"
	// k_lelaer_team   = "https://api.lelaer.com/ys/getTeamResult.php"
	kLelaerSum = "https://api.lelaer.com/ys/getSumComment.php"
)

func (ndata Data) GetSumComment(uid string, wife FindMap) ([]byte, error) {
	p, err := ndata.transToTeyvat(uid, wife)
	if err != nil {
		return nil, err
	}
	d, _ := json.Marshal(p)
	if d, err = web.RequestDataWith(web.NewTLS12Client(),
		kLelaerSum,
		"POST",
		"https://servicewechat.com/wx2ac9dce11213c3a8/192/page-frame.html",
		"Mozilla/5.0 (Linux; Android 12; SM-G977N Build/SP1A.210812.016; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/86.0.4240.99 XWEB/4375 MMWEBSDK/20221011 Mobile Safari/537.36 MMWEBID/4357 MicroMessenger/8.0.30.2244(0x28001E44) WeChat/arm64 Weixin GPVersion/1 NetType/WIFI Language/zh_CN ABI/arm64 MiniProgramEnv/android",
		bytes.NewReader(d),
	); err != nil {
		return nil, fmt.Errorf("出现错误捏：%v", err)
	}
	return d, nil
}

// 角色数据转换为 Teyvat Helper 请求格式
type (
	TeyvatDetail struct {
		Name      string `json:"artifacts_name"`
		Type      string `json:"artifacts_type"`
		Level     int    `json:"level"`
		MainTips  string `json:"maintips"`
		MainValue any    `json:"mainvalue"`
		Tips1     string `json:"tips1"`
		Tips2     string `json:"tips2"`
		Tips3     string `json:"tips3"`
		Tips4     string `json:"tips4"`
	}

	TeyvatData struct {
		Server      string         `json:"server"`
		UserLevel   int            `json:"user_level"`
		Uid         string         `json:"uid"`
		Role        string         `json:"role"`
		Cons        int            `json:"role_class"`
		Level       int            `json:"level"`
		Weapon      string         `json:"weapon"`
		WeaponLevel int            `json:"weapon_level"`
		WeaponClass string         `json:"weapon_class"`
		HP          int            `json:"hp"`
		BaseHP      int            `json:"base_hp"`
		Attack      int            `json:"attack"`
		BaseAttack  int            `json:"base_attack"`
		Defend      int            `json:"defend"`
		BaseDefend  int            `json:"base_defend"`
		Element     int            `json:"element"`
		Crit        string         `json:"crit"`
		CritDmg     string         `json:"crit_dmg"`
		Heal        string         `json:"heal"`
		Recharge    string         `json:"recharge"`
		FireDmg     string         `json:"fire_dmg"`
		WaterDmg    string         `json:"water_dmg"`
		ThunderDmg  string         `json:"thunder_dmg"`
		WindDmg     string         `json:"wind_dmg"`
		IceDmg      string         `json:"ice_dmg"`
		RockDmg     string         `json:"rock_dmg"`
		GrassDmg    string         `json:"grass_dmg"`
		PhysicalDmg string         `json:"physical_dmg"`
		Artifacts   string         `json:"artifacts"`
		Fetter      int            `json:"fetter"`
		Ability1    int            `json:"ability1"`
		Ability2    int            `json:"ability2"`
		Ability3    int            `json:"ability3"`
		Detail      []TeyvatDetail `json:"artifacts_detail"`
	}

	Teyvat struct {
		Role []TeyvatData `json:"role_data"`
		Time int64        `json:"timestamp"`
	}
)

var (
	k_error_sys    = errors.New("程序错误")
	k_error_promap = errors.New("获取角色失败")
)

func (ndata Data) transToTeyvat(uid string, wife FindMap) (*Teyvat, error) {
	if wife == nil {
		if wife = GetWifeOrWq("wife"); wife == nil {
			return nil, k_error_sys
		}
	}
	reliquary := GetReliquary()
	if reliquary == nil {
		return nil, k_error_sys
	}
	syw := GetSywName()
	if syw == nil {
		return nil, k_error_sys
	}

	s := getServer(uid)
	res := &Teyvat{Time: time.Now().Unix()}

	for _, v := range ndata.AvatarInfoList {
		name := wife.Idmap(strconv.Itoa(v.AvatarID))
		cons := len(v.TalentIDList)

		n := len(v.EquipList) // 纠正圣遗物空缺报错的无返回情况
		if n == 0 {
			return nil, k_error_sys
		}
		equip_last := v.EquipList[n-1]
		for m := range equip_last.Weapon.AffixMap {
			n = m
		}
		affix := equip_last.Weapon.AffixMap[n] + 1

		// 武器名
		wqname := reliquary.WQ[equip_last.Flat.NameTextHash]
		if wqname == "" {
			return nil, k_error_sys
		}

		teyvat_data := TeyvatData{
			Uid:         uid,
			Server:      s,
			UserLevel:   ndata.PlayerInfo.Level,
			Role:        name,
			Cons:        cons,
			Weapon:      wqname,
			WeaponLevel: equip_last.Weapon.Level,
			WeaponClass: fmt.Sprintf("精炼%d阶", affix),
			Fetter:      v.FetterInfo.ExpLevel,
		}

		for _, item := range ndata.PlayerInfo.ShowAvatarInfoList {
			if item.AvatarID == v.AvatarID {
				teyvat_data.Level = item.Level
				break
			}
		}

		hp := v.FightPropMap.Num2000           //生命
		crit := v.FightPropMap.Num20 * 100     //暴击
		critDmg := v.FightPropMap.Num22 * 100  //爆伤
		recharge := v.FightPropMap.Num23 * 100 //充能

		physicalDmg := v.FightPropMap.Num30 * 100 // 物理加伤
		fireDmg := v.FightPropMap.Num40 * 100     // 火元素加伤
		thunderDmg := v.FightPropMap.Num41 * 100  // 雷元素加伤
		waterDmg := v.FightPropMap.Num42 * 100    // 水元素加伤
		windDmg := v.FightPropMap.Num44 * 100     // 风元素加伤
		rockDmg := v.FightPropMap.Num45 * 100     // 岩元素加伤
		iceDmg := v.FightPropMap.Num46 * 100      // 冰元素加伤
		grassDmg := v.FightPropMap.Num43 * 100    // 草元素加伤

		// dataFix from https://github.com/yoimiya-kokomi/miao-plugin/blob/ac27075276154ef5a87a458697f6e5492bd323bd/components/profile-data/enka-data.js#L186  # noqa: E501
		if name == "雷电将军" {
			thunderDmg = max(0, thunderDmg-(recharge-100)*0.4) // 雷元素伤害加成
		} else if name == "莫娜" {
			waterDmg = max(0, waterDmg-recharge*0.2) // 水元素伤害加成
		} else if name == "妮露" && cons == 6 {
			crit = max(5, crit-min(30, hp*0.6))        // 暴击率
			critDmg = max(50, critDmg-min(60, hp*1.2)) // 暴击伤害
		}
		for _, item := range []string{"息灾", "波乱月白经津", "雾切之回光", "猎人之径"} {
			if item == wqname {
				z := 12 + 12*(float64(affix)-1)/4
				fireDmg = max(0, fireDmg-z)       // 火元素加伤
				thunderDmg = max(0, thunderDmg-z) // 雷元素加伤
				waterDmg = max(0, waterDmg-z)     // 水元素加伤
				windDmg = max(0, windDmg-z)       // 风元素加伤
				rockDmg = max(0, rockDmg-z)       // 岩元素加伤
				iceDmg = max(0, iceDmg-z)         // 冰元素加伤
				grassDmg = max(0, grassDmg-z)     // 草元素加伤
				break
			}
		}

		// fmt.Println(name)

		// 获取角色
		role := GetRole(name)
		if role == nil {
			return nil, k_error_sys
		}
		// 天赋等级
		talentid := role.GetTalentId()
		teyvat_data.Ability1 = v.SkillLevelMap[talentid[0]]
		teyvat_data.Ability2 = v.SkillLevelMap[talentid[1]]
		teyvat_data.Ability3 = v.SkillLevelMap[talentid[2]]
		ming := len(v.TalentIDList) // 命之座
		// 天赋等级修复
		if ming >= role.TalentCons.E {
			teyvat_data.Ability2 += 3
		}
		if ming >= role.TalentCons.Q {
			teyvat_data.Ability3 += 3
		}

		// 圣遗物数据
		var syws []string
		for i, equip := range v.EquipList {
			if equip.Flat.SetNameTextHash == "" {
				continue
			}
			if wqname = reliquary.WQ[equip.Flat.SetNameTextHash]; wqname == "" {
				return nil, k_error_sys
			}
			syws = append(syws, wqname)
			sywallname := syw.Names(wqname)[i] // 圣遗物name
			// fmt.Println(wqname, sywallname)

			var mainValue any
			if s = Stofen(equip.Flat.ReliquaryMainStat.MainPropID); s == "" {
				mainValue = int(0.5 + equip.Flat.ReliquaryMainStat.Value)
			} else {
				mainValue = Ftoone(equip.Flat.ReliquaryMainStat.Value) + s
			}

			detail := TeyvatDetail{
				Name:      sywallname,
				Type:      GetEquipType(equip.Flat.EquipType),
				Level:     equip.Reliquary.Level - 1,
				MainTips:  GetAppendProp(equip.Flat.ReliquaryMainStat.MainPropID),
				MainValue: mainValue,
			}
			for i, stats := range equip.Flat.ReliquarySubStats {
				s = fmt.Sprintf("%s+%v%s", GetAppendProp(stats.SubPropID), stats.Value, Stofen(stats.SubPropID))
				switch i {
				case 0:
					detail.Tips1 = s
				case 1:
					detail.Tips2 = s
				case 2:
					detail.Tips3 = s
				case 3:
					detail.Tips4 = s
				}
			}
			teyvat_data.Detail = append(teyvat_data.Detail, detail)
		}

		teyvat_data.Artifacts = Sywsuit(syws)

		teyvat_data.HP = int(0.5 + hp)
		teyvat_data.BaseHP = int(0.5 + v.FightPropMap.Num1)       // 基础生命值
		teyvat_data.Attack = int(0.5 + v.FightPropMap.Num2001)    // 攻击
		teyvat_data.BaseAttack = int(0.5 + v.FightPropMap.Num4)   // 基础攻击力
		teyvat_data.Defend = int(0.5 + v.FightPropMap.Num2002)    // 防御
		teyvat_data.BaseDefend = int(0.5 + v.FightPropMap.Num7)   // 基础防御力
		teyvat_data.Element = int(0.5 + v.FightPropMap.Num28)     // 元素精通
		teyvat_data.Heal = Ftoone(v.FightPropMap.Num26*100) + "%" // 治疗加成
		teyvat_data.Crit = Ftoone(crit) + "%"
		teyvat_data.CritDmg = Ftoone(critDmg) + "%"
		teyvat_data.Recharge = Ftoone(recharge) + "%"
		teyvat_data.FireDmg = Ftoone(fireDmg) + "%"
		teyvat_data.WaterDmg = Ftoone(waterDmg) + "%"
		teyvat_data.ThunderDmg = Ftoone(thunderDmg) + "%"
		teyvat_data.WindDmg = Ftoone(windDmg) + "%"
		teyvat_data.IceDmg = Ftoone(iceDmg) + "%"
		teyvat_data.RockDmg = Ftoone(rockDmg) + "%"
		teyvat_data.GrassDmg = Ftoone(grassDmg) + "%"
		teyvat_data.PhysicalDmg = Ftoone(physicalDmg) + "%"

		// fmt.Println(teyvat_data)

		res.Role = append(res.Role, teyvat_data) // 单个角色最终结果
	}
	return res, nil
}

// 获取指定 UID 所属服务器
func getServer(uid string) string {
	switch uid[0] {
	case '5':
		return "cn_qd01"
	case '6':
		return "os_usa"
	case '7':
		return "os_euro"
	case '8':
		return "os_asia"
	case '9':
		return "世界树" // os_cht
	}
	return "天空岛" // cn_gf01
}

func max(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}
