package config

type WebConfig struct {
	Host                  string `mapstructure:"host"`
	Port                  int    `mapstructure:"port"`
	Name                  string `mapstructure:"name"`
	Mode                  string `mapstructure:"mode"`
	*ConsulConfig         `mapstructure:"consul" json:"consul" yaml:"consul"`
	*UserServerConfig     `mapstructure:"user_srv" json:"user_srv" yaml:"user_srv"`
	*CommentServerConfig  `mapstructure:"comment_srv" json:"comment_srv" yaml:"comment_srv"`
	*LikesServerConfig    `mapstructure:"likes_srv" json:"likes_srv" yaml:"likes_srv"`
	*FollowerServerConfig `mapstructure:"follower_srv" json:"follower_srv" yaml:"follower_srv"`
	*VideoServerConfig    `mapstructure:"video_srv" json:"video_srv" yaml:"video_srv"`
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

type LikesServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

type FollowerServerConfig struct {
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
