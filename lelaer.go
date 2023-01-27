package kokomi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/FloatTech/floatbox/web"
)

const (
	// k_lelaer_damage = "https://api.lelaer.com/ys/getDamageResult.php"
	// k_lelaer_team   = "https://api.lelaer.com/ys/getTeamResult.php"
	k_lelaer_sum = "https://api.lelaer.com/ys/getSumComment.php"
)

type LelaerApi struct {
	ndata Data

	reliquary *Fff
	syw       Syws
}

func (ndata Data) GetSumComment(uid string, wife FindMap) ([]byte, error) {
	p, err := ndata.transToTeyvat(uid, wife)
	if err != nil {
		return nil, err
	}
	d, _ := json.Marshal(p)
	if d, err = web.RequestDataWith(web.NewTLS12Client(),
		k_lelaer_sum,
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
	TeyvatHelperDetail struct {
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

	TeyvatHelperData struct {
		Server      string               `json:"server"`
		UserLevel   int                  `json:"user_level"`
		Uid         string               `json:"uid"`
		Role        string               `json:"role"`
		Cons        int                  `json:"role_class"`
		Level       int                  `json:"level"`
		Weapon      string               `json:"weapon"`
		WeaponLevel int                  `json:"weapon_level"`
		WeaponClass string               `json:"weapon_class"`
		HP          int                  `json:"hp"`
		BaseHP      int                  `json:"base_hp"`
		Attack      int                  `json:"attack"`
		BaseAttack  int                  `json:"base_attack"`
		Defend      int                  `json:"defend"`
		BaseDefend  int                  `json:"base_defend"`
		Element     int                  `json:"element"`
		Crit        string               `json:"crit"`
		CritDmg     string               `json:"crit_dmg"`
		Heal        string               `json:"heal"`
		Recharge    string               `json:"recharge"`
		FireDmg     string               `json:"fire_dmg"`
		WaterDmg    string               `json:"water_dmg"`
		ThunderDmg  string               `json:"thunder_dmg"`
		WindDmg     string               `json:"wind_dmg"`
		IceDmg      string               `json:"ice_dmg"`
		RockDmg     string               `json:"rock_dmg"`
		GrassDmg    string               `json:"grass_dmg"`
		PhysicalDmg string               `json:"physical_dmg"`
		Artifacts   string               `json:"artifacts"`
		Fetter      int                  `json:"fetter"`
		Ability1    int                  `json:"ability1"`
		Ability2    int                  `json:"ability2"`
		Ability3    int                  `json:"ability3"`
		Detail      []TeyvatHelperDetail `json:"artifacts_detail"`
	}

	TeyvatHelper struct {
		Role []TeyvatHelperData `json:"role_data"`
		Time int64              `json:"timestamp"`
	}
)

func max[T int | int32 | int64 | float64](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func min[T int | int32 | int64 | float64](x, y T) T {
	if x > y {
		return y
	}
	return x
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

var (
	k_error_sys    = errors.New("程序错误")
	k_error_promap = errors.New("获取角色失败")
)

func (ndata Data) transToTeyvat(uid string, wife FindMap) (*TeyvatHelper, error) {
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

	server := getServer(uid)
	res := &TeyvatHelper{Time: 0}

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
		equip_affix := equip_last.Weapon.AffixMap[n] + 1

		// 武器名
		wqname := reliquary.WQ[equip_last.Flat.NameTextHash]
		if wqname == "" {
			return nil, k_error_sys
		}

		roleData := TeyvatHelperData{
			Uid:         uid,
			Server:      server,
			UserLevel:   ndata.PlayerInfo.Level,
			Role:        name,
			Cons:        cons,
			Weapon:      wqname,
			WeaponLevel: equip_last.Weapon.Level,
			WeaponClass: fmt.Sprintf("精炼%d阶", equip_affix),
			Fetter:      v.FetterInfo.ExpLevel,
		}

		for _, item := range ndata.PlayerInfo.ShowAvatarInfoList {
			if item.AvatarID == v.AvatarID {
				roleData.Level = item.Level
				break
			}
		}

		hp := v.FightPropMap.Num2000           //生命
		crit := v.FightPropMap.Num20 * 100     //暴击
		crit_dmg := v.FightPropMap.Num22 * 100 //爆伤
		recharge := v.FightPropMap.Num23 * 100 //充能

		physical_dmg := v.FightPropMap.Num30 * 100 // 物理加伤
		fire_dmg := v.FightPropMap.Num40 * 100     // 火元素加伤
		thunder_dmg := v.FightPropMap.Num41 * 100  // 雷元素加伤
		water_dmg := v.FightPropMap.Num42 * 100    // 水元素加伤
		wind_dmg := v.FightPropMap.Num44 * 100     // 风元素加伤
		rock_dmg := v.FightPropMap.Num45 * 100     // 岩元素加伤
		ice_dmg := v.FightPropMap.Num46 * 100      // 冰元素加伤
		grass_dmg := v.FightPropMap.Num43 * 100    // 草元素加伤

		// dataFix from https://github.com/yoimiya-kokomi/miao-plugin/blob/ac27075276154ef5a87a458697f6e5492bd323bd/components/profile-data/enka-data.js#L186  # noqa: E501
		if name == "雷电将军" {
			thunder_dmg = max(0, thunder_dmg-(recharge-100)*0.4) // 雷元素伤害加成
		} else if name == "莫娜" {
			water_dmg = max(0, water_dmg-recharge*0.2) // 水元素伤害加成
		} else if name == "妮露" && cons == 6 {
			crit = max(5, crit-min(30, hp*0.6))          // 暴击率
			crit_dmg = max(50, crit_dmg-min(60, hp*1.2)) // 暴击伤害
		}
		for _, item := range []string{"息灾", "波乱月白经津", "雾切之回光", "猎人之径"} {
			if item == wqname {
				z := 12 + 12*(float64(equip_affix)-1)/4
				fire_dmg = max(0, fire_dmg-z)       // 火元素加伤
				thunder_dmg = max(0, thunder_dmg-z) // 雷元素加伤
				water_dmg = max(0, water_dmg-z)     // 水元素加伤
				wind_dmg = max(0, wind_dmg-z)       // 风元素加伤
				rock_dmg = max(0, rock_dmg-z)       // 岩元素加伤
				ice_dmg = max(0, ice_dmg-z)         // 冰元素加伤
				grass_dmg = max(0, grass_dmg-z)     // 草元素加伤
				break
			}
		}

		// fmt.Println(name)

		// 天赋等级修复
		role := GetRole(name)
		if role == nil {
			return nil, k_error_promap
		}
		talentid := role.GetTalentId()
		roleData.Ability1 = v.SkillLevelMap[talentid[0]]
		roleData.Ability2 = v.SkillLevelMap[talentid[1]]
		roleData.Ability3 = v.SkillLevelMap[talentid[2]]
		ming := len(v.TalentIDList) // 命之座

		if ming >= role.TalentCons.E {
			roleData.Ability2 += 3
		}
		if ming >= role.TalentCons.Q {
			roleData.Ability3 += 3
		}

		// 圣遗物数据
		for i, equip := range v.EquipList {
			if equip.Flat.SetNameTextHash == "" {
				return nil, k_error_sys
			}
			roleData.Artifacts = reliquary.WQ[equip.Flat.SetNameTextHash]
			if roleData.Artifacts == "" {
				return nil, k_error_sys
			}

			sywallname := syw.Names(roleData.Artifacts)[i] // 圣遗物name
			// fmt.Println(roleData.Artifacts, sywallname)

			var main_value any
			if s := Stofen(equip.Flat.ReliquaryMainStat.MainPropID); s == "" {
				main_value = int(0.5 + equip.Flat.ReliquaryMainStat.Value)
			} else {
				main_value = Ftoone(equip.Flat.ReliquaryMainStat.Value) + s
			}

			detail := TeyvatHelperDetail{
				Name:      sywallname,
				Type:      GetEquipType(equip.Flat.EquipType),
				Level:     equip.Reliquary.Level - 1,
				MainTips:  GetAppendProp(equip.Flat.ReliquaryMainStat.MainPropID),
				MainValue: main_value,
			}
			for i, stats := range equip.Flat.ReliquarySubStats {
				s := fmt.Sprintf("%s+%v%s", GetAppendProp(stats.SubPropID), stats.Value, Stofen(stats.SubPropID))
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
			roleData.Detail = append(roleData.Detail, detail)
		}

		if roleData.Artifacts == "" {
			roleData.Artifacts = "+"
		} else {
			roleData.Artifacts += "4"
		}

		roleData.HP = int(0.5 + hp)
		roleData.BaseHP = int(0.5 + v.FightPropMap.Num1)     // 基础生命值
		roleData.Attack = int(0.5 + v.FightPropMap.Num2001)  // 攻击
		roleData.BaseAttack = int(0.5 + v.FightPropMap.Num4) // 基础攻击力
		roleData.Defend = int(0.5 + v.FightPropMap.Num2002)  // 防御
		roleData.BaseDefend = int(0.5 + v.FightPropMap.Num7) // 基础防御力
		roleData.Element = int(0.5 + v.FightPropMap.Num28)   // 精通
		roleData.Heal = Ftoone(v.FightPropMap.Num26) + "%"
		roleData.Crit = Ftoone(crit) + "%"
		roleData.CritDmg = Ftoone(crit_dmg) + "%"
		roleData.Recharge = Ftoone(recharge) + "%"
		roleData.FireDmg = Ftoone(fire_dmg) + "%"
		roleData.WaterDmg = Ftoone(water_dmg) + "%"
		roleData.ThunderDmg = Ftoone(thunder_dmg) + "%"
		roleData.WindDmg = Ftoone(wind_dmg) + "%"
		roleData.IceDmg = Ftoone(ice_dmg) + "%"
		roleData.RockDmg = Ftoone(rock_dmg) + "%"
		roleData.GrassDmg = Ftoone(grass_dmg) + "%"
		roleData.PhysicalDmg = Ftoone(physical_dmg) + "%"

		// fmt.Println(roleData)

		res.Role = append(res.Role, roleData) // 单个角色最终结果
	}
	return res, nil
}
