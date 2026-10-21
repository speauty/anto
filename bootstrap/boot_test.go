package bootstrap

import (
	"anto/domain/repository"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestBootRegistersTencentAndAliyunFromExistingConfig(t *testing.T) {
	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, "cfg.yml"), []byte("ali_cloud_mt:\n  ak_id: \"\"\n  ak_secret: \"\"\n  region: cn-hangzhou\n"), 0600); err != nil {
		t.Fatal(err)
	}
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(directory); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(previous) })
	Boot(context.Background())
	for _, id := range []string{"ali_cloud_mt", "tencent_cloud_mt", "baidu_cloud_mt", "caiyunai", "deepl_pro", "deepseek", "youdao_batch"} {
		if repository.GetTranslators().GetById(id) == nil {
			t.Fatalf("translator %s was not registered", id)
		}
	}
}
