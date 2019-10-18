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
}

// 群消息
type MsgGroup struct {
	Account string    `json:"-"`                 // 收到消息的微信号
	Gid     string    `json:"gid,omitempty"`     // 群唯一ID
	Uid     string    `json:"uid,omitempty"`     // 消息发送者唯一ID
	Name    string    `json:"name,omitempty"`    // 消息发送者昵称
	Content string    `json:"content,omitempty"` // 消息内容
	Time    time.Time `json:"time,omitempty"`    // 消息创建时间
	// optional
	NameAlias string `json:"nameAlias,omitempty"` // 消息发送者备注
	GroupName string `json:"groupName,omitempty"` // 群名称
}

type MsgUser struct {
	// 开发者信息
	Uid string `json:"my_account"` //
	// 好友信息
	ToUid   string `json:"to_account"` // 好友唯一ID
	ToName  string `json:"to_name"`    // 昵称
	ToAlias string `json:"form_name"`  // 备注名
	// 消息相关
	Type        int       `json:"type"`           // 类型：1自己发的、2好友发的
	ContentType int       `json:"content_type"`   // 消息类型 消息类型：1文字、2图片、3表情、4语音、5视频、6文件、10系统消息
	Content     string    `json:"content"`        // 消息内容
	Time        time.Time `json:"time,omitempty"` // 消息创建时间
	// 要消化的字段
	ID         int   `json:"id,omitempty"`          // 消息ID
	CreateTime int64 `json:"create_time,omitempty"` // 消息创建时间 (时间戳)
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
	ID         int       `json:"id"`                   // 开发者id
	Phone      string    `json:"phone"`                // 开发者手机号
	Name       string    `json:"name"`                 // 开发者名称
	Sex        string    `json:"sex,omitempty"`        // 性别
	TimeExpire time.Time `json:"timeExpire,omitempty"` // 账户过期时间
	TimeLogin  time.Time `json:"timeLogin,omitempty"`  // 上次登录时间
}

// 整理为go友好数据格式
func (d *MsgUser) init() (err error) {
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
	} else if len(d.Name) > 0 {
		return d.Name
	} else {
		return d.Uid
	}
}

//
func (d *MsgUser) GetName() string {
	if d == nil {
		return ""
	} else if len(d.ToAlias) > 0 {
		return d.ToAlias
	} else if len(d.ToName) > 0 {
		return d.ToName
	} else {
		return d.ToUid
	}
}
