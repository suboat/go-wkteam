package wkteam

import (
	"github.com/stretchr/testify/require"

	"testing"
)

// 消息去重监测
func Test_MsgHash(t *testing.T) {
	as := require.New(t)
	as.Equal("76bc15ccac73d27c76d16bcf72e57df7c0e09006", (&MsgUser{FromUid: "a", ToUid: "b"}).GetHash())
	as.Equal("3cb4ed8953674f565af786ae924b559dc1f41342", (&MsgGroup{Gid: "g", Uid: "a"}).GetHash())
}

// 测试回调监听
func Test_ListenAndServe(t *testing.T) {
	testConfigRead()
	as := require.New(t)
	api := &WkTeam{}
	as.Nil(api.ListenAndServe(Settings.CallbackLocal, Settings.CallbackPublic))
}
