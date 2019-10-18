package wkteam

import (
	"github.com/suboat/go-contrib"
	"github.com/suboat/go-contrib/log"

	"fmt"
	"strings"
	"sync"
)

// 文档地址 https://wkteam.gitbook.io/doc/api-wen-dang

var (
	// 版本
	Version = &contrib.Version{Major: 0, Minor: 0, Patch: 1, Model: "wkteam"}
	// GitCommit 当前代码对应的git-commit记录
	GitCommit = ""
	// Settings 系统设置
	Settings *Config
)

// 通用hook
var (
	DefaultHookHookMsgGroup = func(msg *MsgGroup) error {
		log.Debugf(`[default] HookHookMsgGroup %s <- %s`, msg.Account, PubJSON(msg))
		return nil
	}
)

// Config 系统配置参数
type Config struct {
	contrib.Config `yaml:"-"`
	//
	ApiHost        string // 微控api入口
	Phone          string // 开发者手机号
	Password       string // 开发者密码
	Secret         string // 开发者密钥
	Account        string // 默认要管理的微信号
	CallbackLocal  string // 实际监听回调地址
	CallbackPublic string // 映射到公网的回调地址
}

//
type WkTeam struct {
	// 账号信息
	Phone    string // 开发者手机号
	Password string // 开发者密码
	Secret   string // 开发者密钥
	Account  string // 要管理的微信号
	//
	ApiHost string         // 微控api入口
	Log     contrib.Logger //
	// hooks
	HookMsgGroup func(group *MsgGroup) error
	//
	lock   sync.RWMutex //
	apiKey string       //
	//
	inited bool
}

//
func (cfg *Config) Valid() (err error) {
	if len(cfg.Phone) == 0 {
		return fmt.Errorf(`phone undefined`)
	} else if len(cfg.Password) == 0 {
		return fmt.Errorf(`password undefined`)
	} else if len(cfg.Secret) == 0 {
		return fmt.Errorf(`secret undefined`)
	}
	return
}

//
func (api *WkTeam) init() (err error) {
	if api.inited {
		return
	}
	api.lock.Lock()
	defer api.lock.Unlock()
	if api.inited {
		return
	}
	// defaults
	if api.Log == nil {
		api.Log = log.Log
	}
	if len(api.ApiHost) == 0 {
		api.ApiHost = Settings.ApiHost
	}
	if len(api.Phone) == 0 {
		api.Phone = Settings.Phone
	}
	if len(api.Password) == 0 {
		api.Password = Settings.Password
	}
	if len(api.Secret) == 0 {
		api.Secret = Settings.Secret
	}
	api.ApiHost = strings.TrimRight(api.ApiHost, "/")
	api.inited = true
	return
}

// NewWkTeam 新建一个微控对象
func NewWkTeam(s *WkTeam) (ret *WkTeam) {
	if s != nil {
		ret = s
	} else {
		ret = new(WkTeam)
	}
	_ = ret.init()
	return
}

//
func init() {
	// 默认设置
	Settings = &Config{
		ApiHost:        "http://admin.wkgjhome.com",
		CallbackLocal:  "http://127.0.0.1:8080/v1/wkteam/",
		CallbackPublic: "https://yourhost/api/wkteam/",
	}
	// 配置注释
	_ = Settings.SetComments(map[string]string{
		"ApiHost":        "微控api入口",
		"Phone":          "开发者手机号",
		"Password":       "开发者密码",
		"Secret":         "开发者密钥",
		"Account":        "默认要管理的微信号",
		"CallbackLocal":  "实际监听回调地址",
		"CallbackPublic": "映射到公网的回调地址",
	})
	// version
	if len(GitCommit) > 0 {
		Version.Commit = &GitCommit
	}
}
