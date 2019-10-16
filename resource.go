package wkteam

import (
	"github.com/suboat/go-contrib"

	"encoding/json"
	"time"
)

// 微控返回数据

// 开发者信息
type Agent struct {
	Uid        int       `json:"uid"`                  // 开发者uid
	Phone      string    `json:"phone"`                // 开发者手机号
	Name       string    `json:"name"`                 // 开发者名称
	Sex        string    `json:"sex,omitempty"`        // 性别
	TimeExpire time.Time `json:"timeExpire,omitempty"` // 账户过期时间
	TimeLogin  time.Time `json:"timeLogin,omitempty"`  // 上次登录时间
	//
	OverdueTime int64 `json:"overdue_time,omitempty"` //
	LastTime    int64 `json:"last_time,omitempty"`    //
	//
	inited bool //
}

func (d *Agent) init() (err error) {
	if d.inited {
		return
	}
	d.TimeExpire = time.Unix(d.OverdueTime, 0)
	d.TimeLogin = time.Unix(d.LastTime, 0)
	d.inited = true
	return
}

// 解析返回
func (d *Agent) Parse(resp string) (err error) {
	if d == nil {
		return contrib.ErrUndefined
	} else if err = json.Unmarshal([]byte(resp), d); err != nil {
		println(err.Error())
		return
	}
	return d.init()
}
