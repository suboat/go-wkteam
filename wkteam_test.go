package wkteam

import (
	"github.com/stretchr/testify/require"
	"github.com/suboat/go-contrib"

	"fmt"
	"testing"
)

var (
	// 单元测试读取的配置信息, 账号信息在内
	testConfig = "config.test.yaml"
)

// 读取测试配置文件
func testConfigRead() {
	if len(testConfig) == 0 {
		return
	}
	var err error
	if err = contrib.PubConfigRead(testConfig, Settings, Settings); err != nil {
		panic(fmt.Errorf(`配置文件读取失败:%s, err:%v`, testConfig, err))
	} else if err = Settings.Valid(); err != nil {
		panic(fmt.Errorf(`配置文件非法:%s, err:%v`, testConfig, err))
	}
	if err = Settings.Save(nil); err != nil {
		panic(fmt.Errorf(`配置文件更新失败:%s, err:%v`, testConfig, err))
	}
	testConfig = ""
	return
}

// 获取账号信息
func Test_GetAgent(t *testing.T) {
	testConfigRead()
	as := require.New(t)
	api := NewWkTeam(nil)
	d, err := api.GetAgent()
	as.Nil(err)
	t.Log(PubJSON(d))
}
