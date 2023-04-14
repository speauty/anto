package translator

import (
	"anto/lib/log"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func RequestSimpleGet(engine InterfaceTranslator, url string) ([]byte, error) {
	httpResp, err := http.DefaultClient.Get(url)
	defer func() {
		if httpResp != nil && httpResp.Body != nil {
			_ = httpResp.Body.Close()
		}
	}()
	if err != nil {
		log.Singleton().Error(fmt.Sprintf("调用接口失败, 引擎: %s, 错误: %s", engine.GetName(), err))
		return nil, fmt.Errorf("网络请求异常, 错误: %s", err.Error())
	}

	if httpResp.StatusCode != 200 {
		log.Singleton().Error(fmt.Sprintf("调用接口失败, 引擎: %s, 错误: %d(%s)", engine.GetName(), httpResp.StatusCode, httpResp.Status))
		return nil, fmt.Errorf("网络响应异常, 错误:  %d(%s)", httpResp.StatusCode, httpResp.Status)
	}
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		log.Singleton().Error(fmt.Sprintf("读取报文异常, 引擎: %s, 错误: %s", engine.GetName(), err))
		return nil, fmt.Errorf("读取报文出现异常, 错误: %s", err.Error())
	}
	return respBytes, nil
}

func RequestSimplePost(engine InterfaceTranslator, httpUrl string, bodyParams interface{}) ([]byte, error) {
	reqBytes, _ := json.Marshal(bodyParams)
	req, _ := http.NewRequest(http.MethodPost, httpUrl, bytes.NewReader(reqBytes))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	return RequestSimpleHttp(engine, req)
}

func RequestSimpleHttp(engine InterfaceTranslator, r *http.Request) ([]byte, error) {
	httpResp, err := new(http.Client).Do(r)
	defer func() {
		if httpResp != nil && httpResp.Body != nil {
			_ = httpResp.Body.Close()
		}
	}()
	if err != nil {
		log.Singleton().Error(fmt.Sprintf("调用接口失败, 引擎: %s, 错误: %s", engine.GetName(), err))
		return nil, fmt.Errorf("网络请求出现异常, 错误: %s", err.Error())
	}
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		log.Singleton().Error(fmt.Sprintf("读取报文异常, 引擎: %s, 错误: %s", engine.GetName(), err))
		return nil, fmt.Errorf("读取报文出现异常, 错误: %s", err.Error())
	}
	return respBytes, nil
}
