package setting

import (
	"gopkg.in/ini.v1"
	"log"
	"time"
)

type App struct {
	JwtSecret       string
	PageSize        int
	RuntimeRootPath string

	ImagePrefixUrl string
	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DataBase struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var (
	AppSetting      = &App{}
	ServerSetting   = &Server{}
	DataBaseSetting = &DataBase{}
	RedisSetting    = &Redis{}
)

/*
	编写(App、Server、DataBase)与 app.ini 一直的结构体
	使用 MapTo 将配置项映射到结构上面定义的结构体上
	对一些特殊设置的配置项进行在赋值
*/
var cfg *ini.File

func Setup() {
	var e error
	cfg, e = ini.Load("conf/app.ini")
	if e != nil {
		log.Fatalf("Fail to parse 'conf/app.in' : %v", e)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DataBaseSetting)
	mapTo("redis", RedisSetting)

	//将  MB  转换为 B
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024

	//超时时间单位设置为 秒
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second

}

func mapTo(section string, v interface{}) {
	e := cfg.Section(section).MapTo(v)
	if e != nil {
		log.Fatalf("cfg.MapTo %sSetting err: %v", section, e)
	}
}
