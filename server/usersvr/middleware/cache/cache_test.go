package cache

import (
	"testing"
	"usersvr/config"
	"usersvr/log"
)

func TestCachePing(t *testing.T) {
	if err := config.Init(); err != nil {
		log.Fatalf("init config failed, err:%v\n", err)
	}
	log.InitLog()
	log.Info("log init success...")

	client := GetRedisCli()
	t.Log("client======", client)
	if client == nil {
		t.Errorf("client is nil")
	}
}
