# go-wkteam

微控管家 SDK in Golang.

参考文档：https://wkteam.gitbook.io/doc/

## 快速开始

#### 下载与安装

    go get github.com/suboat/go-wkteam

#### 快速启动 `hello.go`
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

## License

The [MIT License](LICENSE)
