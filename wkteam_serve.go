package wkteam

import (
	"github.com/julienschmidt/httprouter"

	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	CallbackPrefix   = "/callback/" // 回调地址前缀
	CallbackNormal   = "normal"     // 通用回调
	CallbackMsgUser  = "messagelog" // 单聊回调
	CallbackMsgGroup = "msggroup"   // 群聊回调
)

// ListenAndServe 监听回调 addrListen:监听地址 urlPublic: 公网地址
func (api *WkTeam) ListenAndServe(urlLocal string, urlPublic string) (err error) {
	if err = api.init(); err != nil {
		return
	}
	urlLocal = strings.TrimRight(urlLocal, "/")
	urlPublic = strings.TrimRight(urlPublic, "/")
	if len(urlLocal) > 4 && strings.Contains(urlLocal, "http") == false {
		urlLocal = "http://" + urlLocal
	}
	var (
		router   = httprouter.New()
		callback = CallbackPrefix + ":category"
		uri      *url.URL
	)
	if uri, err = url.Parse(urlLocal); err != nil {
		return
	}

	// 登记回调
	if len(urlPublic) > 0 {
		if _err := api.SetCallback(urlPublic); _err != nil {
			api.Log.Warnf(`[serve-init] SetCallback err: %v`, _err)
		}
	}

	// handler
	router.HandlerFunc(http.MethodGet, uri.Path+callback, api.handleCallback)
	router.HandlerFunc(http.MethodPost, uri.Path+callback, api.handleCallback)
	router.HandlerFunc(http.MethodGet, uri.Path+"/health", api.handleHealth)

	//
	api.Log.Infof(`[serve-listen] serve on %s/ -> %s/ on %s`, urlLocal, urlPublic, uri.Path+callback)
	return http.ListenAndServe(uri.Host, router)
}

// 检查
func (api *WkTeam) handleHealth(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("ok"))
}

// 回调
func (api *WkTeam) handleCallback(rw http.ResponseWriter, req *http.Request) {
	var (
		params   = PubGetParams(req)
		category = params["category"]
		msg      = make(map[string]string)
		err      error
	)
	defer func() {
		if err != nil {
			rw.Write([]byte(err.Error()))
		} else {
			rw.Write([]byte(fmt.Sprintf("%s ok", category)))
		}
	}()
	if err = req.ParseMultipartForm(0); err != nil {
		return
	} else {
		for k, v := range req.Form {
			if len(v) > 0 {
				msg[k] = v[0]
			}
		}
	}

	// debug
	api.Log.Debugf(`[serve-callback] %s <- %s`, category, PubJSON(msg))

	//
	switch category {
	case CallbackNormal:
		// 通用回调
		break
	case CallbackMsgUser:
		// 单聊回调
		var (
			resp = &struct {
				Account     string `json:"my_account"`       // 收到消息的微信号
				Name        string `json:"my_name"`          // 收到消息的微信号
				NameAlias   string `json:"my_account_alias"` // 登录微信ID（wxid_xxxxxx开头的）
				ToUid       string `json:"to_account"`       // 好友唯一ID
				ToName      string `json:"to_name"`          // 昵称
				Type        int    `json:"type"`             // 类型：1自己发的、2好友发的
				ContentType int    `json:"content_type"`     // 消息类型 消息类型：1文字、2图片、3表情、4语音、5视频、6文件、10系统消息
				Content     string `json:"content"`          // 消息内容
				CreateTime  int64  `json:"sendtime"`         // 发送时间
			}{}
		)
		if err = json.Unmarshal([]byte(msg["data"]), resp); err != nil {
			return
		}
		data := &MsgUser{
			Account:  resp.Account,
			Category: priContentTypeToStr(resp.ContentType),
			Content:  resp.Content,
			Time:     time.Unix(resp.CreateTime, 0),
		}
		if resp.Type == 1 {
			// 自己发的
			data.IsMe = true
			data.FromUid = resp.Account
			data.FromName = resp.Name
			data.FromNameAlias = ""
			data.ToUid = resp.ToUid
			data.ToName = resp.ToName
			data.ToNameAlias = ""
		} else {
			// 别人发给我的
			data.IsMe = false
			data.FromUid = resp.ToUid
			data.FromName = resp.ToName
			data.FromNameAlias = ""
			data.ToUid = resp.Account
			data.ToName = resp.Name
			data.ToNameAlias = ""
		}
		// 回调
		call := api.HookMsgUser
		if call == nil {
			call = DefaultHookHookMsgUser
		}
		if call != nil {
			go func() {
				defer PanicRecoverError(api.Log, nil)
				if _err := call(data); _err != nil {
					api.Log.Errorf(`[serve-hook] HookMsgUser err: %v <- %s`, _err, PubJSON(data))
				}
			}()
		}
		break
	case CallbackMsgGroup:
		// 群聊回调
		var (
			resp = &struct {
				Account   string `json:"my_account"`
				Gid       string `json:"g_number"`
				Uid       string `json:"to_account"`
				Name      string `json:"to_name"`
				Content   string `json:"content"`
				Time      int64  `json:"send_time"`
				GroupName string `json:"g_name"`
			}{}
		)
		if err = json.Unmarshal([]byte(msg["data"]), resp); err != nil {
			return
		}
		data := &MsgGroup{
			Account:   resp.Account,
			Gid:       resp.Gid,
			Uid:       resp.Uid,
			Name:      resp.Name,
			NameAlias: "",
			Content:   resp.Content,
			Time:      time.Unix(resp.Time, 0),
			GroupName: resp.GroupName,
		}
		// 回调
		call := api.HookMsgGroup
		if call == nil {
			call = DefaultHookHookMsgGroup
		}
		if call != nil {
			go func() {
				defer PanicRecoverError(api.Log, nil)
				if _err := call(data); _err != nil {
					api.Log.Errorf(`[serve-hook] HookMsgGroup err: %v <- %s`, _err, PubJSON(data))
				}
			}()
		}
		break
	default:
		// 未知响应
		break
	}

	return
}
