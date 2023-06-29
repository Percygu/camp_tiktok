package utils

import (
	"context"
	"gatewaysvr/config"
	"gatewaysvr/log"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"testing"
)

func TestVideoConnection(t *testing.T) {
	config.Init()
	log.InitLog()
	client := NewVideoSvrClient("videosvr")
	t.Log(client)
	if client == nil {

		t.Errorf("NewVideoSvrClient err")
	}

	resp, err := client.PublishVideo(context.Background(), &pb.PublishVideoRequest{})
	t.Log(resp, err)
}
