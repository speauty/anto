package translator

import (
	"anto/lib/log"
	"anto/lib/restrictor"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func RequestSimpleGet(ctx context.Context, engine ImplTranslator, url string) ([]byte, error) {

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	return RequestSimpleHttp(ctx, engine, req)
}

func RequestSimplePost(ctx context.Context, engine ImplTranslator, httpUrl string, bodyParams interface{}) ([]byte, error) {
	reqBytes, _ := json.Marshal(bodyParams)
	req, _ := http.NewRequest(http.MethodPost, httpUrl, bytes.NewReader(reqBytes))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	return RequestSimpleHttp(ctx, engine, req)
}

func RequestSimpleHttp(ctx context.Context, engine ImplTranslator, r *http.Request) ([]byte, error) {
	if err := restrictor.Singleton().Wait(engine.GetId(), ctx); err != nil {
		return nil, fmt.Errorf("限流异常, 错误: %s", err.Error())
	}

	httpResp, err := new(http.Client).Do(r)
	defer func() {
		if httpResp != nil && httpResp.Body != nil {
			_ = httpResp.Body.Close()
		}
	}()

	if err != nil {
		log.Singleton().ErrorF("调用接口失败, 引擎: %s, 错误: %s", engine.GetName(), err)
		return nil, fmt.Errorf("网络请求异常, 错误: %s", err.Error())
	}

	if httpResp.StatusCode != 200 {
		log.Singleton().ErrorF("调用接口失败, 引擎: %s, 错误: %d(%s)", engine.GetName(), httpResp.StatusCode, httpResp.Status)
		return nil, fmt.Errorf("网络响应异常, 错误:  %d(%s)", httpResp.StatusCode, httpResp.Status)
	}

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		log.Singleton().ErrorF("读取报文异常, 引擎: %s, 错误: %s", engine.GetName(), err)
		return nil, fmt.Errorf("读取报文异常, 错误: %s", err.Error())
	}

	return respBytes, nil
}
