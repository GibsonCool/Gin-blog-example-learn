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

var (
	AppSetting      = &App{}
	ServerSetting   = &Server{}
	DataBaseSetting = &DataBase{}
)

/*
	编写(App、Server、DataBase)与 app.ini 一直的结构体
	使用 MapTo 将配置项映射到结构上面定义的结构体上
	对一些特殊设置的配置项进行在赋值
*/
func Setup() {
	log.Printf("读取 app.ini 配置项....")

	Cfg, e := ini.Load("conf/app.ini")
	if e != nil {
		log.Fatalf("Fail to parse 'conf/app.in' : %v", e)
	}

	//不再使用直接读取key的方式，而是用 MapTo 映射到结构体中
	e = Cfg.Section("app").MapTo(AppSetting)
	if e != nil {
		log.Fatalf("Cfg.MapTo AppSetting err: %v", e)
	}
	//将  MB  转换为 B
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024

	e = Cfg.Section("server").MapTo(ServerSetting)
	if e != nil {
		log.Fatalf("Cfg.MapTo ServerSetting err: %v", e)
	}

	//超时时间单位设置为 秒
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	e = Cfg.Section("database").MapTo(DataBaseSetting)
	if e != nil {
		log.Fatalf("Cfg.MapTo DataBaseSetting err: %v", e)
	}

	if ServerSetting.RunMode == "debug" {
		log.Printf("app : %v", AppSetting)
		log.Printf("server : %v", ServerSetting)
		log.Printf("database : %v", DataBaseSetting)
	}
}
