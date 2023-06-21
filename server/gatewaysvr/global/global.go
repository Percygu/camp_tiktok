package global

import (
	"gatewaysvr/config"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"os"
	"path/filepath"
)

var (
	Conf             = new(config.WebConfig) // Conf 全局配置变量
	UserSrvClient    pb.UserServiceClient    // 用户服务客户端
	CommentSrvClient pb.CommentServiceClient

	UserSrvClient pb.VideoServiceClient
	UserSrvClient pb.RelationServiceClient

	// UserSrvClient proto.UserClient
	// UserSrvClient proto.UserClient

)

// 项目主目录
var RootDir string

func init() {
	inferRootDir()
	// 初始化配置
}

// 推断 Root目录（copy就行）
func inferRootDir() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var infer func(string) string
	infer = func(dir string) string {
		if exists(dir + "/main.go") {
			return dir
		}

		// 查看dir的父目录
		parent := filepath.Dir(dir)
		return infer(parent)
	}

	RootDir = infer(pwd)
}

func exists(dir string) bool {
	// 查找主机是不是存在 dir
	_, err := os.Stat(dir)
	return err == nil || os.IsExist(err)
}
