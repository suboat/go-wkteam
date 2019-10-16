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
	//
	testGroupID     = "18217585821@chatroom" // 测试群号
	testFriendUID   = "好友微信号"                // 测试好友微信号
	testFriendAlias = "这是一个足够长的测试备注名"        // 测试好友微信号
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
	api := &WkTeam{}
	d, err := api.GetAgent()
	as.Nil(err)
	t.Log(PubJSON(d))
}

// 获取群列表
func Test_GetGroups(t *testing.T) {
	testConfigRead()
	as := require.New(t)
	api := NewWkTeam(nil)
	d, err := api.GetGroups(nil)
	as.Nil(err)
	t.Log(PubJSON(d))
}

// 取群消息
func Test_GetMsgGroup(t *testing.T) {
	testConfigRead()
	as := require.New(t)
	api := NewWkTeam(nil)
	query := &Query{Limit: 30, Page: 0}
	d, err := api.GetMsgGroup(testGroupID, query)
	as.Nil(err)
	// 调试信息
	for _i, _d := range d {
		api.Log.Infof("#%d %s -> %s(%s): %s", _i+1, _d.Time, _d.Uid, _d.GetName(), _d.Content)
		if _i == len(d)-1 {
			api.Log.Info(PubJSON(_d))
		}
	}
	api.Log.Infof("获取到群消息 %d/%d", len(d), query.Total)
}

// 同意好友添加申请
func Test_PassAddFriend(t *testing.T) {
	testConfigRead()
	as := require.New(t)
	api := NewWkTeam(nil)
	d, err := api.PassAddFriend(testFriendUID)
	as.Nil(err)
	t.Log(d)
	return
}

// 备注好友
func Test_RemarkFriend(t *testing.T) {
	testConfigRead()
	as := require.New(t)
	api := NewWkTeam(nil)
	d, err := api.RemarkFriend(testFriendUID, testFriendAlias)
	as.Nil(err)
	t.Log(d)
	return
}
