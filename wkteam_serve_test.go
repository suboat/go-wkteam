package wkteam

import (
	"github.com/stretchr/testify/require"

	"testing"
)

const (
	// 测试web服务地址
	testListenAndServe = "http://127.0.0.1:8080/v1/wkteam/"
	// 测试web发布地址
	testListenAndPublic = "https://yourhost/api/wkteam2/"
)

// 测试回调监听
func Test_ListenAndServe(t *testing.T) {
	testConfigRead()
	as := require.New(t)
	as.Nil(nil)
	api := &WkTeam{}
	as.Nil(api.ListenAndServe(testListenAndServe, testListenAndPublic))
}
