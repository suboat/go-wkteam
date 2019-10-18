package wkteam

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/suboat/go-contrib"
	"github.com/suboat/go-contrib/log"
	"runtime"
	"strings"

	"encoding/json"
	"net/http"
)

//
func PubJSON(v_ interface{}) (r string) {
	r = "{}"
	if v_ == nil {
		return
	}
	switch v := v_.(type) {
	case string:
		if len(v) > 2 && v[0] == '{' && v[len(v)-1] == '}' {
			r = v
		}
		break
	case *string:
		_v := *v
		if len(_v) > 2 && _v[0] == '{' && _v[len(_v)-1] == '}' {
			r = _v
		}
		break
	default:
		// null: v_不为nil但其实是空
		if r1, _err := json.Marshal(v_); _err == nil && string(r1) != "null" {
			r = string(r1)
		}
		break
	}
	return
}

// PubGetParams 取url参数
func PubGetParams(req *http.Request) (ret map[string]string) {
	ret = make(map[string]string)
	for _, _ps := range httprouter.ParamsFromContext(req.Context()) {
		ret[_ps.Key] = _ps.Value
	}
	return
}

// PanicRecoverError 统一处理panic, 并更新error
func PanicRecoverError(logger contrib.Logger, err *error) {
	r := recover()
	if r == nil {
		return
	}
	if logger == nil {
		logger = log.Log
	}
	logger.Errorf(`[panic-recover] "%s" %v`, panicIdentify(), r)
	if err != nil {
		*err = contrib.ErrPanicRecover.SetVars(r)
	}
	return
}

// 定位panic位置 参考自: https://gist.github.com/swdunlop/9629168
func panicIdentify() string {
	var (
		pc [16]uintptr
		n  = runtime.Callers(3, pc[:])
	)

	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		_fnName := fn.Name()
		if strings.HasPrefix(_fnName, "runtime.") {
			continue
		}
		file, line := fn.FileLine(pc)

		//
		var (
			_fnNameDir = strings.Split(_fnName, "/")
			_fnNameLis = strings.Split(_fnName, ".")
			_fnNameSrc string
		)
		if len(_fnNameDir) > 1 {
			_fnNameSrc = _fnNameDir[0] + "/" + _fnNameDir[1] + "/"
		} else {
			_fnNameSrc = _fnNameDir[0]
		}
		fnName := _fnNameLis[len(_fnNameLis)-1]

		// file
		_pcLis := strings.Split(file, _fnNameSrc)
		filePath := strings.Join(_pcLis[1:], "")

		return fmt.Sprintf("%s:%d|%s", filePath, line, fnName)
	}

	return "unknown"
}
