package setting

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	Port      int    `mapstructure:"port"`
	MachineId int64  `mapstructure:"machine_id"`

	*PublishConfig `mapstructure:"publish"`
	*MysqlConfig   `mapstructure:"mysql"`
	*RedisConfig   `mapstructure:"redis"`
}

type PublishConfig struct {
	Mode                     bool   `mapstructure:"mode"`
	QiNiuCloudPlayUrlPrefix  string `mapstructure:"qi_niu_cloud_play_url_prefix"`
	QiNiuCloudCoverUrlPrefix string `mapstructure:"qi_niu_cloud_cover_url_prefix"`
	LocalIP                  string `mapstructure:"local_ip_address"`
	VideoPathPrefix          string `mapstructure:"video_path_prefix"`
	CoverPathPrefix          string `mapstructure:"cover_path_prefix"`
	AccessKey                string `mapstructure:"access_key"`
	SecretKey                string `mapstructure:"secret_key"`
	BucketName               string `mapstructure:"bucket_name"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
}

func Init(filePath string) (err error) {
	viper.SetConfigFile(filePath)
	//读取配置信息
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed, err:%v \n", err)
		return
	}
	//读取到的配置信息反系列化到Conf变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}
	//热加载
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})
	return
}
