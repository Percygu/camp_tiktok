package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var globalConfig GlobalConfig

type GlobalConfig struct {
	Host                  string `mapstructure:"host"`
	Port                  int    `mapstructure:"port"`
	Name                  string `mapstructure:"name"`
	Mode                  string `mapstructure:"mode"`
	VideoPath             string `mapstructure:"video_path"`
	*LogConfig            `mapstructure:"log" json:"log" yaml:"log"`
	*ConsulConfig         `mapstructure:"consul" json:"consul" yaml:"consul"`
	*UserServerConfig     `mapstructure:"user_srv" json:"user_srv" yaml:"user_srv"`
	*CommentServerConfig  `mapstructure:"comment_srv" json:"comment_srv" yaml:"comment_srv"`
	*RelationServerConfig `mapstructure:"relation_srv" json:"relation_srv" yaml:"relation_srv"`
	*FavoriteServerConfig `mapstructure:"favorite_srv" json:"favorite_srv" yaml:"favorite_srv"`
	*VideoServerConfig    `mapstructure:"video_srv" json:"video_srv" yaml:"video_srv"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"file_name"`
	LogPath    string `mapstructure:"log_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type UserServerConfig struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
	Name string `mapstructure:"name" json:"name" yaml:"name"`
}

type CommentServerConfig struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
	Name string `mapstructure:"name" json:"name" yaml:"name"`
}

type RelationServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

type FavoriteServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

type VideoServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

type ConsulConfig struct {
	Host string   `mapstructure:"host" json:"host" yaml:"host"`
	Port int      `mapstructure:"port" json:"port" yaml:"port"`
	Tags []string `mapstructure:"tags" json:"tags" yaml:"tags"`
}

func Init() (err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml") // 远程配置文件传输 确定配置文件的格式
	viper.AddConfigPath(".")    // 指定配置文件的一个寻找路径
	viper.AddConfigPath("../")  // 指定配置文件的一个寻找路径
	err = viper.ReadInConfig()  // 读取配置信息

	if err != nil {
		// 读取配置信息错误
		fmt.Printf("viper.ReadInConfig() failed: %v\n", err)
		return fmt.Errorf("viper.ReadInConfig() failed: %v\n", err)
	}

	// 把读取到的信息反序列化到 Conf 变量中
	if err = viper.Unmarshal(&globalConfig); err != nil {
		fmt.Printf("viper.Unmarshal failed: %v\n", err)
		return fmt.Errorf("viper.Unmarshal failed: %v\n", err)
	}

	viper.WatchConfig() // 实时监控配置文件（热加载）
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改...")
		// 当配置文件信息发生变化 就修改 Conf 变量
		if err := viper.Unmarshal(&globalConfig); err != nil {
			fmt.Printf("viper.Unmarshal failed: %v\n", err)
		}
	})
	return nil
}

func GetGlobalConfig() *GlobalConfig {
	return &globalConfig
}
