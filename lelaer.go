package kokomi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/FloatTech/floatbox/web"
)

func (ndata Thisdata) GetSumComment(uid string) (data []byte, err error) {
	var teyvat *Teyvat
	if teyvat, err = ndata.transToTeyvat(uid); err == nil {
		data, _ = json.Marshal(teyvat)
		data, err = web.RequestDataWith(web.NewTLS12Client(),
			"https://api.lelaer.com/ys/getSumComment.php",
			"POST",
			"https://servicewechat.com/wx2ac9dce11213c3a8/192/page-frame.html",
			"Mozilla/5.0 (Linux; Android 12; SM-G977N Build/SP1A.210812.016; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/86.0.4240.99 XWEB/4375 MMWEBSDK/20221011 Mobile Safari/537.36 MMWEBID/4357 MicroMessenger/8.0.30.2244(0x28001E44) WeChat/arm64 Weixin GPVersion/1 NetType/WIFI Language/zh_CN ABI/arm64 MiniProgramEnv/android",
			bytes.NewReader(data),
		)
	}
	return
}
func (ndata Thisdata) Getgroupdata(uid string, is [4]int) (data []byte, err error) {
	var teyvat *Teyvat
	if teyvat, err = ndata.transTogroup(uid, is); err == nil {
		data, _ = json.Marshal(teyvat)
		data, err = web.RequestDataWith(web.NewTLS12Client(),
			"https://api.lelaer.com/ys/getTeamResult.php",
			"POST",
			"https://servicewechat.com/wx2ac9dce11213c3a8/211/page-frame.html",
			"Mozilla/5.0 (Linux; Android 12; SM-G977N Build/SP1A.210812.016; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/86.0.4240.99 XWEB/4375 MMWEBSDK/20221011 Mobile Safari/537.36 MMWEBID/4357 MicroMessenger/8.0.30.2244(0x28001E44) WeChat/arm64 Weixin GPVersion/1 NetType/WIFI Language/zh_CN ABI/arm64 MiniProgramEnv/android",
			bytes.NewReader(data),
		)
	}
	return
}

// 角色数据转换为 Teyvat Helper 请求格式
type (
	TeyvatDetail struct {
		Name      string `json:"artifacts_name,omitempty"`
		Type      string `json:"artifacts_type,omitempty"`
		Level     int    `json:"level,omitempty"`
		MainTips  string `json:"maintips,omitempty"`
		MainValue any    `json:"mainvalue,omitempty"`
		Tips1     string `json:"tips1,omitempty"`
		Tips2     string `json:"tips2,omitempty"`
		Tips3     string `json:"tips3,omitempty"`
		Tips4     string `json:"tips4,omitempty"`
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
		Uid    string       `json:"uid"`
		Server string       `json:"server"`
		Data   []TeyvatData `json:"role_data"`
		Time   int64        `json:"timestamp"`
	}
)

var lelaerErrorSYS = errors.New("程序错误")

func (ndata Thisdata) transToTeyvat(uid string) (*Teyvat, error) {
	res := &Teyvat{Time: time.Now().Unix()}
	for l := 0; l < len(ndata.Chars); l++ {
		_, err := res.getroeldata(l, uid, ndata)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (ndata Thisdata) transTogroup(uid string, is [4]int) (*Teyvat, error) {
	res := &Teyvat{Time: time.Now().Unix()}
	res.Uid = uid
	res.Server = "cn_gf01"
	for _, l := range is {
		_, err := res.getroeldata(l, uid, ndata)
		if err != nil {
			return nil, err
		}
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
	default:
		return "天空岛" // cn_gf01
	}
}
func (res *Teyvat) getroeldata(l int, uid string, ndata Thisdata) (*Teyvat, error) {
	s := getServer(uid)
	v := ndata.Chars[l]
	name := v.Name
	role, err := GetRole(name) // 获取角色
	if err != nil {
		return nil, lelaerErrorSYS
	}
	affix := v.Weapon.Affix
	// 武器名
	wqname := v.Weapon.Name
	if wqname == "" {
		return nil, lelaerErrorSYS
	}

	cons := v.Cons // 命之座
	rolelv, _ := strconv.Atoi(v.Level)
	teyvatData := TeyvatData{
		Uid:         uid,
		Server:      s,
		UserLevel:   ndata.Level,
		Role:        name,
		Cons:        cons,
		Weapon:      wqname,
		WeaponLevel: v.Weapon.Level,
		WeaponClass: fmt.Sprintf("精炼%d阶", affix),
		Fetter:      v.Fetter,
		Ability1:    v.Talent.A,
		Ability2:    v.Talent.E,
		Ability3:    v.Talent.Q,
		Level:       rolelv,
	}

	hp := v.Attr.Hp             //生命
	crit := v.Attr.Cpct         //暴击
	critDmg := v.Attr.Cdmg      //爆伤
	recharge := v.Attr.Recharge //充能

	physicalDmg := v.Attr.Phy // 物理加伤
	var fireDmg, thunderDmg, waterDmg, windDmg, rockDmg, iceDmg, grassDmg float64
	switch v.Attr.DmgName {
	case "火元素加伤:":
		fireDmg = v.Attr.Dmg
	case "雷元素加伤:":
		thunderDmg = v.Attr.Dmg
	case "水元素加伤:":
		waterDmg = v.Attr.Dmg
	case "风元素加伤:":
		windDmg = v.Attr.Dmg
	case "岩元素加伤:":
		rockDmg = v.Attr.Dmg
	case "冰元素加伤:":
		iceDmg = v.Attr.Dmg
	case "草元素加伤:":
		grassDmg = v.Attr.Dmg
	}
	// 天赋等级修复
	if cons >= role.TalentCons.E {
		teyvatData.Ability2 += 3
	}
	if cons >= role.TalentCons.Q {
		teyvatData.Ability3 += 3
	}

	// dataFix from https://github.com/yoimiya-kokomi/miao-plugin/blob/ac27075276154ef5a87a458697f6e5492bd323bd/components/profile-data/enka-data.js#L186  # noqa: E501
	switch name {
	case "雷电将军":
		thunderDmg = math.Max(0, thunderDmg-(recharge-100)*0.4) // 雷元素伤害加成
	case "莫娜":
		waterDmg = math.Max(0, waterDmg-recharge*0.2) // 水元素伤害加成
	case "妮露":
		if cons == 6 {
			crit = math.Max(5, crit-math.Min(30, hp*0.6))        // 暴击率
			critDmg = math.Max(50, critDmg-math.Min(60, hp*1.2)) // 暴击伤害
		}
	case "达达利亚":
		teyvatData.Ability1 += 1

	default:

	}
	for _, item := range []string{"息灾", "波乱月白经津", "雾切之回光", "猎人之径"} {
		if item == wqname {
			z := 12 + 12*(float64(affix)-1)/4
			fireDmg = math.Max(0, fireDmg-z)       // 火元素加伤
			thunderDmg = math.Max(0, thunderDmg-z) // 雷元素加伤
			waterDmg = math.Max(0, waterDmg-z)     // 水元素加伤
			windDmg = math.Max(0, windDmg-z)       // 风元素加伤
			rockDmg = math.Max(0, rockDmg-z)       // 岩元素加伤
			iceDmg = math.Max(0, iceDmg-z)         // 冰元素加伤
			grassDmg = math.Max(0, grassDmg-z)     // 草元素加伤
			break
		}
	}

	// 圣遗物数据
	var syws []string = []string{v.Artis.Hua.Set, v.Artis.Yu.Set, v.Artis.Sha.Set, v.Artis.Bei.Set, v.Artis.Guan.Set}
	teyvatData.Detail = []TeyvatDetail{}
	for i := 0; i < 5; i++ {
		var equip sywm
		switch i {
		case 0:
			equip = v.Artis.Hua
		case 1:
			equip = v.Artis.Yu
		case 2:
			equip = v.Artis.Sha
		case 3:
			equip = v.Artis.Bei
		case 4:
			equip = v.Artis.Guan
		}
		if equip.Name == "" {
			continue
		}
		// 圣遗物name
		detail := TeyvatDetail{
			Name:     equip.Name,
			Type:     GetEquipType(strconv.Itoa(i)),
			Level:    equip.Level - 1,
			MainTips: GetAppendProp(equip.Main.Title),
		}

		if s = Stofen(equip.Main.Title); s == "" {
			detail.MainValue = int(0.5 + equip.Main.Value)
		} else {
			detail.MainValue = Ftoone(equip.Main.Value) + s
		}

		for i, stats := range equip.Attrs {
			s = fmt.Sprintf("%s+%v%s", GetAppendProp(stats.Title), stats.Value, Stofen(stats.Title))
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
		teyvatData.Detail = append(teyvatData.Detail, detail)
	}

	teyvatData.Artifacts = Sywsuit(syws)

	teyvatData.HP = int(0.5 + hp)
	teyvatData.BaseHP = int(0.5 + v.Attr.HpBase)      // 基础生命值
	teyvatData.Attack = int(0.5 + v.Attr.Atk)         // 攻击
	teyvatData.BaseAttack = int(0.5 + v.Attr.AtkBase) // 基础攻击力
	teyvatData.Defend = int(0.5 + v.Attr.Def)         // 防御
	teyvatData.BaseDefend = int(0.5 + v.Attr.DefBase) // 基础防御力
	teyvatData.Element = int(0.5 + v.Attr.Mastery)    // 元素精通
	teyvatData.Heal = Ftoone(v.Attr.Heal) + "%"       // 治疗加成
	teyvatData.Crit = Ftoone(crit) + "%"
	teyvatData.CritDmg = Ftoone(critDmg) + "%"
	teyvatData.Recharge = Ftoone(recharge) + "%"
	teyvatData.FireDmg = Ftoone(fireDmg) + "%"
	teyvatData.WaterDmg = Ftoone(waterDmg) + "%"
	teyvatData.ThunderDmg = Ftoone(thunderDmg) + "%"
	teyvatData.WindDmg = Ftoone(windDmg) + "%"
	teyvatData.IceDmg = Ftoone(iceDmg) + "%"
	teyvatData.RockDmg = Ftoone(rockDmg) + "%"
	teyvatData.GrassDmg = Ftoone(grassDmg) + "%"
	teyvatData.PhysicalDmg = Ftoone(physicalDmg) + "%"

	res.Data = append(res.Data, teyvatData) // 单个角色最终结果
	return nil, nil
}
