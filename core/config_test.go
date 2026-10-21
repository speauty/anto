package core

import "testing"

func TestTencentConfigLimitAndCredentialRedaction(t *testing.T) {
	if err := validateProviderValues("tencent_cloud_mt", map[string]interface{}{"max_single_text_length": 2001}); err == nil {
		t.Fatal("Tencent text limit was not validated")
	}
	if err := validateProviderValues("tencent_cloud_mt", map[string]interface{}{"max_single_text_length": 2000}); err != nil {
		t.Fatal(err)
	}
	redacted := redactValues(map[string]interface{}{"secret_id": "id", "secret_key": "key", "region": "ap-guangzhou"})
	if redacted["secret_id"] != redactedSecret || redacted["secret_key"] != redactedSecret || redacted["region"] != "ap-guangzhou" {
		t.Fatalf("Tencent credentials were exposed: %+v", redacted)
	}
}

func TestBaiduTierValidationAndSecretRedaction(t *testing.T) {
	for _, tier := range []string{"standard", "advanced", "premium"} {
		if err := validateProviderValues("baidu_cloud_mt", map[string]interface{}{"service_tier": tier}); err != nil {
			t.Fatalf("valid tier %s rejected: %v", tier, err)
		}
	}
	for _, values := range []map[string]interface{}{
		{"service_tier": "invalid"},
		{"qps": 10},
		{"max_single_text_length": 6000},
	} {
		if err := validateProviderValues("baidu_cloud_mt", values); err == nil {
			t.Fatalf("invalid Baidu parameters accepted: %v", values)
		}
	}
	redacted := redactValues(map[string]interface{}{"app_id": "id", "app_secret": "secret", "service_tier": "standard"})
	if redacted["app_id"] != "id" || redacted["app_secret"] != redactedSecret || redacted["service_tier"] != "standard" {
		t.Fatalf("unexpected Baidu settings response: %v", redacted)
	}
}

func TestDeepSeekConfigValidationAndSecretRedaction(t *testing.T) {
	if err := validateProviderValues("deepseek", map[string]interface{}{"base_url": "http://127.0.0.1:8000/v1", "model": "local-model", "prompt": "简洁", "max_single_text_length": 6000}); err != nil {
		t.Fatal(err)
	}
	for _, values := range []map[string]interface{}{
		{"base_url": "http://example.com"},
		{"model": ""},
		{"prompt": string(make([]byte, 8001))},
		{"max_single_text_length": 6001},
	} {
		if err := validateProviderValues("deepseek", values); err == nil {
			t.Fatalf("invalid DeepSeek parameters accepted: %v", values)
		}
	}
	redacted := redactValues(map[string]interface{}{"api_key": "secret", "model": "deepseek-flash", "prompt": "自然"})
	if redacted["api_key"] != redactedSecret || redacted["model"] != "deepseek-flash" || redacted["prompt"] != "自然" {
		t.Fatalf("unexpected DeepSeek settings response: %v", redacted)
	}
}
