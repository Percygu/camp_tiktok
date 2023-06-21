package initialize

import (
	"fmt"
	"gatewaysvr/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig() (err error) {
	// 自动推导项目根目录
	configFile := global.RootDir + "/config/config.yaml"
	viper.SetConfigFile(configFile)
	// viper.SetConfigFile("./config.yaml") //指定配置文件（带后缀，可写绝对路径和相对路径两种）
	// viper.SetConfigName("config") //指定配置文件的名字（不带后缀）
	// 基本上是配合远程配置中心使用的，告诉viper当前的数据使用什么格式去解析
	viper.SetConfigType("yaml") // 远程配置文件传输 确定配置文件的格式
	viper.AddConfigPath(".")    // 指定配置文件的一个寻找路径
	err = viper.ReadInConfig()  // 读取配置信息

	if err != nil {
		// 读取配置信息错误
		fmt.Printf("viper.ReadInConfig() failed: %v\n", err)
		return
	}

	// 把读取到的信息反序列化到 Conf 变量中
	if err = viper.Unmarshal(global.C); err != nil {
		fmt.Printf("viper.Unmarshal failed: %v\n", err)
		return
	}

	// // 创建文件夹
	// err = util.Mkdir(global.Conf.PathConfig.VideoFile)
	// if err != nil {
	// 	panic("mkdir videofile error")
	// }
	// err = util.Mkdir(global.Conf.PathConfig.PicFile)
	// if err != nil {
	// 	panic("mkdir picfile error")
	// }

	viper.WatchConfig() // 实时监控配置文件（热加载）
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改...")
		// 当配置文件信息发生变化 就修改 Conf 变量
		if err := viper.Unmarshal(global.Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed: %v\n", err)
		}
	})

	return
}
