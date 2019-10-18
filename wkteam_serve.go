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
	api.Log.Infof(`[serve-callback] %s <- %s`, category, PubJSON(msg))

	//
	switch category {
	case CallbackNormal:
		// 通用回调
		break
	case CallbackMsgGroup:
		// 群聊回调: 解析为群消息
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
