package wkteam

import (
	"time"
)

// 微控返回数据

// query参数
type Query struct {
	Account string                 // 要管理的微信号
	Skip    int                    // 分页参数
	Page    int                    // 分页参数,填写会覆盖Skip信息 0起始
	Limit   int                    // 分页参数,默认10
	Params  map[string]interface{} // 其它参数
	// 回应参数
	Total int // 总数据量
}

// 群信息
type Group struct {
	Gid  string `json:"number"` // 群唯一ID
	Name string `json:"name"`   // 群名
	// 别名字段
}

// 群消息
type MsgGroup struct {
	Gid       string    `json:"gid"`            // 群唯一ID
	Uid       string    `json:"wac_account"`    // 消息发送者唯一ID
	Name      string    `json:"wac_name"`       // 消息发送者昵称
	NameAlias string    `json:"form_name"`      // 消息发送者备注
	Content   string    `json:"content"`        // 消息内容
	Time      time.Time `json:"time,omitempty"` // 消息创建时间
	//
	ID         int   `json:"id,omitempty"`          //
	GroId      int   `json:"gro_id,omitempty"`      //
	CreateTime int64 `json:"create_time,omitempty"` //
	inited     bool  //
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
	inited      bool  //
}

// 整理为go友好数据格式
func (d *MsgGroup) init() (err error) {
	if d.inited {
		return
	}
	d.Time = time.Unix(d.CreateTime, 0)
	d.inited = true
	return
}

//
func (d *MsgGroup) GetName() string {
	if d == nil {
		return ""
	} else if len(d.NameAlias) > 0 {
		return d.NameAlias
	} else if len(d.NameAlias) > 0 {
		return d.Name
	} else {
		return d.Uid
	}
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
