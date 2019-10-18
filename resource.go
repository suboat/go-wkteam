package wkteam

import (
	"crypto/sha1"
	"fmt"
	"time"
)

// 微控返回数据

const (
	// 1文字、2图片、3表情、 4语音、5视频、6文件、10系统消息
	MsgCategoryTxt = "txt" // 文字
	MsgCategoryImg = "img" // 图片
	MsgCategoryGif = "gif" // 表情
	MsgCategoryWav = "wav" // 语音
	MsgCategoryMp4 = "mp4" // 视频
	MsgCategoryDoc = "doc" // 文件
	MsgCategorySys = "sys" // 系统消息
)

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
	Account   string `json:"-"`                   // 获取消息的微信号
	Gid       string `json:"gid,omitempty"`       // 群唯一ID
	Uid       string `json:"uid,omitempty"`       // 消息发送者唯一ID
	Name      string `json:"name,omitempty"`      // 消息发送者昵称
	NameAlias string `json:"nameAlias,omitempty"` // 消息发送者备注
	GroupName string `json:"groupName,omitempty"` // 群名称
	// 消息相关
	Category string    `json:"category,omitempty"` // 消息类型, 见顶部定义
	Content  string    `json:"content,omitempty"`  // 消息内容
	Time     time.Time `json:"time,omitempty"`     // 消息创建时间
}

// 用户消息
type MsgUser struct {
	Account       string `json:"-"`                       // 获取消息的微信号
	FromUid       string `json:"fromUid,omitempty"`       // 发送者唯一ID
	FromName      string `json:"fromName,omitempty"`      // 发送者昵称
	FromNameAlias string `json:"fromNameAlias,omitempty"` // 发送者备注
	ToUid         string `json:"toUid,omitempty"`         // 接受者唯一ID
	ToName        string `json:"toName,omitempty"`        // 接受者昵称
	ToNameAlias   string `json:"toNameAlias,omitempty"`   // 接受者备注
	IsMe          bool   `json:"isMe,omitempty"`          // true: 这条消息是我发送的
	// 消息相关
	Category string    `json:"category,omitempty"` // 消息类型, 见顶部定义
	Content  string    `json:"content,omitempty"`  // 消息内容
	Time     time.Time `json:"time,omitempty"`     // 消息创建时间
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

// 群一条群聊消息的哈希, 去重判断用
func (d *MsgGroup) GetHash() (ret string) {
	if d != nil {
		ret = d.Uid + d.Category + d.Content + d.Time.String()
	}
	return fmt.Sprintf(`%x`, sha1.Sum([]byte(ret)))
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

// 群一条单聊消息的哈希, 去重判断用
func (d *MsgUser) GetHash() (ret string) {
	if d != nil {
		ret = d.FromUid + d.ToUid + d.Category + d.Content + d.Time.String()
	}
	return fmt.Sprintf(`%x`, sha1.Sum([]byte(ret)))
}

// 消息发出者名字
func (d *MsgUser) GetName() string {
	if d == nil {
		return ""
	} else if d.IsMe {
		return d.GetFromName()
	} else {
		return d.GetToName()
	}
}

func (d *MsgUser) GetFromName() (ret string) {
	if d == nil {
		return
	}
	if len(d.FromNameAlias) > 0 {
		ret = d.FromNameAlias
	} else if len(d.FromName) > 0 {
		ret = d.FromName
	} else {
		ret = d.FromUid
	}
	return
}

func (d *MsgUser) GetToName() (ret string) {
	if d == nil {
		return
	}
	if len(d.ToNameAlias) > 0 {
		ret = d.ToNameAlias
	} else if len(d.ToName) > 0 {
		ret = d.ToName
	} else {
		ret = d.ToUid
	}
	return
}
