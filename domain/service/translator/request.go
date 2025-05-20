package translator

import (
	"anto/lib/log"
	"anto/lib/restrictor"
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"io"
	"strings"
)

func RequestSimpleHttp(ctx context.Context, engine ImplTranslator, url string, isPost bool, body interface{}, headers map[string]string) ([]byte, error) {
	if err := restrictor.Singleton().Wait(engine.GetId(), ctx); err != nil {
		return nil, fmt.Errorf("限流异常, 错误: %s", err.Error())
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["content-type"] = "application/json"
	headers["accept"] = "application/json"

	client := req.C().SetCommonHeaders(headers).SetCommonRetryCount(3)
	if strings.Contains(url, "api.openai.com") {
		client.SetProxyURL("http://127.0.0.1:7890")
	}
	request := client.R()
	if isPost && body != nil {
		request.SetBody(body)
	}
	var httpResp *req.Response
	var err error
	if isPost {
		httpResp, err = request.Post(url)
	} else {
		httpResp, err = request.Get(url)
	}
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
