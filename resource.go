package wkteam

import (
	"github.com/suboat/go-contrib"

	"encoding/json"
	"time"
)

// 微控返回数据

// query参数
type Query struct {
	Account string                 // 要管理的微信号
	Skip    int                    // 分页参数
	Page    int                    // 分页参数,填写会覆盖Skip信息
	Limit   int                    // 分页参数,默认10
	Params  map[string]interface{} // 其它参数
}

// 群信息
type Group struct {
	Gid  string `json:"number"` // 群唯一ID
	Name string `json:"name"`   // 群名
}

// 用户信息
type User struct {
	Account      string `json:"account"`       // 微信号
	AccountAlias string `json:"account_alias"` // 微信唯一id(原始的微信号)
	Name         string `json:"name"`          // 微信昵称
	Sex          int    `json:"sex"`           // 性别
	Area         string `json:"area"`          // 所在地
	Description  string `json:"description"`   // 签名
}

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

// 整理为go友好数据格式
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
