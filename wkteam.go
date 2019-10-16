package wkteam

import (
	"github.com/suboat/go-contrib"
	"github.com/suboat/go-contrib/log"

	"fmt"
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

// Config 系统配置参数
type Config struct {
	contrib.Config `yaml:"-"`
	//
	ApiHost  string // 微控api入口
	Phone    string // 开发者手机号
	Password string // 开发者密码
	Secret   string // 开发者密钥
	Account  string // 默认要管理的微信号
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

// NewWkTeam 新建一个微控对象
func NewWkTeam(s *WkTeam) (ret *WkTeam) {
	if s != nil {
		ret = s
	} else {
		ret = new(WkTeam)
	}
	// defaults
	if ret.Log == nil {
		ret.Log = log.Log
	}
	if len(ret.ApiHost) == 0 {
		ret.ApiHost = Settings.ApiHost
	}
	if len(ret.Phone) == 0 {
		ret.Phone = Settings.Phone
	}
	if len(ret.Password) == 0 {
		ret.Password = Settings.Password
	}
	if len(ret.Secret) == 0 {
		ret.Secret = Settings.Secret
	}
	return
}

//
func init() {
	// 默认设置
	Settings = &Config{
		ApiHost: "http://admin.wkgjhome.com",
	}
	// 配置注释
	_ = Settings.SetComments(map[string]string{
		"ApiHost":  "微控api入口",
		"Phone":    "开发者手机号",
		"Password": "开发者密码",
		"Secret":   "开发者密钥",
		"Account":  "默认要管理的微信号",
	})
	// version
	if len(GitCommit) > 0 {
		Version.Commit = &GitCommit
	}
}
