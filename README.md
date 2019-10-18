# go-wkteam

微控管家 SDK in Golang.

参考文档：https://wkteam.gitbook.io/doc/

## 快速开始

### 下载与安装

    go get github.com/suboat/go-wkteam

### 快速启动 `hello.go`
```go
package main

import "github.com/suboat/go-wkteam"

func main() {
	api := &wkteam.WkTeam{
		Phone:    "13700000000",                                                      // 微控平台手机号
		Password: "wd2dQ6enmYcj2V5K",                                                 // 微控平台密码
		Secret:   "SQjDw77MEdyWRnTu2ZejfzAtNQBt6arxjGUAhA2QApQAuUCnVq5NCpXDv8nMDbmK", // 微控平台密钥 64字节
		Account:  "wx_2dQ6enmYcj2V5K",                                                // 托管在微控上的微信账号
	}
	if info, err := api.GetAgent(); err != nil {
		panic(err)
	} else {
		println(wkteam.PubJSON(info))
	}
}
```

### 监听消息 `serve.go`
```go
package main

import "github.com/suboat/go-wkteam"

func main() {
	api := &wkteam.WkTeam{
		Phone:    "13700000000",                                                      // 微控平台手机号
		Password: "wd2dQ6enmYcj2V5K",                                                 // 微控平台密码
		Secret:   "SQjDw77MEdyWRnTu2ZejfzAtNQBt6arxjGUAhA2QApQAuUCnVq5NCpXDv8nMDbmK", // 微控平台密钥 64字节
		Account:  "wx_2dQ6enmYcj2V5K",                                                // 托管在微控上的微信账号
	}

	// hook新单聊消息
	api.HookMsgUser = func(msg *wkteam.MsgUser) error {
		api.Log.Infof("%s 发了一条消息给 %s，内容是：%s", msg.GetFromName(), msg.GetToName(), msg.Content)
		return nil
	}

	// hook群聊消息
	api.HookMsgGroup = func(msg *wkteam.MsgGroup) error {
		api.Log.Infof("%s 在群【%s】发了一条消息:%s", msg.GetName(), msg.GroupName, msg.Content)
		return nil
	}

	// 监听地址
	if err := api.ListenAndServe(
		"http://0.0.0.0:8080/api/wkteam/",  // 本机监听地址，还需要nginx等工具暴露给公网
		"https://yourhost.com/api/wkteam/", // 公网地址，微控收到新消息会往这个地址POST数据
	); err != nil {
		panic(err)
	}
}
```

### 拉取消息 `sync.go`
```go
package main

import (
	"github.com/suboat/go-wkteam"
	"time"
)

func main() {
	api := &wkteam.WkTeam{
		Phone:    "13700000000",                                                      // 微控平台手机号
		Password: "wd2dQ6enmYcj2V5K",                                                 // 微控平台密码
		Secret:   "SQjDw77MEdyWRnTu2ZejfzAtNQBt6arxjGUAhA2QApQAuUCnVq5NCpXDv8nMDbmK", // 微控平台密钥 64字节
		Account:  "wx_2dQ6enmYcj2V5K",                                                // 托管在微控上的微信账号
	}

	// 拉取12小时内的与某人的聊天记录
	if data, err := api.GetMsgUserSince("某好友的微信ID", time.Now().Add(-time.Hour*12)); err != nil {
		panic(err)
	} else {
		for i, d := range data {
			api.Log.Infof("#%d %s %s -> %s : %s", i+1, d.Time, d.GetFromName(), d.GetToName(), d.Content)
		}
	}

	// 拉取6小时内的某群的聊天记录
	if data, err := api.GetMsgGroupSince("某群的微信ID", time.Now().Add(-time.Hour*6)); err != nil {
		panic(err)
	} else {
		for i, d := range data {
			api.Log.Infof("#%d %s %s(%s) : %s", i+1, d.Time, d.Uid, d.GetName(), d.Content)
		}
	}
}
```

## License

The [MIT License](LICENSE)
