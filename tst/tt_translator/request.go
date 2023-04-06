package tt_translator

import (
	"fmt"
	"io"
	"net/http"
	"translator/tst/tt_log"
)

func RequestSimpleGet(engine ITranslator, url string) ([]byte, error) {
	httpResp, err := http.DefaultClient.Get(url)
	defer func() {
		if httpResp != nil && httpResp.Body != nil {
			_ = httpResp.Body.Close()
		}
	}()
	if err != nil {
		tt_log.GetInstance().Error(fmt.Sprintf("调用接口失败, 引擎: %s, 错误: %s", engine.GetName(), err))
		return nil, fmt.Errorf("网络请求异常, 错误: %s", err.Error())
	}

	if httpResp.StatusCode != 200 {
		tt_log.GetInstance().Error(fmt.Sprintf("调用接口失败, 引擎: %s, 错误: %d(%s)", engine.GetName(), httpResp.StatusCode, httpResp.Status))
		return nil, fmt.Errorf("网络响应异常, 错误:  %d(%s)", httpResp.StatusCode, httpResp.Status)
	}
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		tt_log.GetInstance().Error(fmt.Sprintf("读取报文异常, 引擎: %s, 错误: %s", engine.GetName(), err))
		return nil, fmt.Errorf("读取报文出现异常, 错误: %s", err.Error())
	}
	return respBytes, nil
}
