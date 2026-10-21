package ali_cloud_mt

import (
	"context"
	"errors"
	"testing"

	"anto/domain/service/translator"
)

func TestConnectionRejectsUnknownRegionBeforeNetwork(t *testing.T) {
	_, err := New().TestConnection(context.Background(), "unknown-region")
	if !errors.Is(err, translator.ErrInvalidConnectionRegion) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDisplayName(t *testing.T) {
	if got := New().GetName(); got != "机器翻译-阿里云" {
		t.Fatalf("display name = %q", got)
	}
}

func TestClientUsesDocumentedPublicEndpoint(t *testing.T) {
	tests := []struct {
		region   string
		endpoint string
	}{
		{"cn-hangzhou", "mt.cn-hangzhou.aliyuncs.com"},
		{"ap-southeast-1", "mt.ap-southeast-1.aliyuncs.com"},
		{"me-central-1", "mt.me-central-1.aliyuncs.com"},
		{"me-east-1", "alimt.me-east-1.aliyuncs.com"},
	}
	for _, test := range tests {
		t.Run(test.region, func(t *testing.T) {
			engine := New()
			engine.Init(&Config{AKId: "test-id", AKSecret: "test-secret", Region: test.region})
			if engine.mtClient == nil {
				t.Fatal("Aliyun client was not created")
			}
			if got := engine.mtClient.Domain; got != test.endpoint {
				t.Fatalf("public endpoint = %q, want %q", got, test.endpoint)
			}
		})
	}
}
