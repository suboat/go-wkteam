package wkteam

import (
	"github.com/stretchr/testify/require"

	"testing"
)

// 测试回调监听
func Test_ListenAndServe(t *testing.T) {
	testConfigRead()
	as := require.New(t)
	api := &WkTeam{}
	as.Nil(api.ListenAndServe(Settings.CallbackLocal, Settings.CallbackPublic))
}
