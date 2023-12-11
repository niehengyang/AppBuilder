package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	UserApi    userApi
	Log        log
	JwtAuth    jwtAuth
	Mysql      mysql
	RabbitMQ   rabbitMQ
	CacheRedis cacheRedis
	System     system
	FileServer fileServer
}

type userApi rest.RestConf

type log struct {
	Level    string
	Path     string
	KeepDays int8
}

// JWT 认证需要的密钥和过期时间配置
type jwtAuth struct {
	AccessSecret string
	AccessExpire int64 //单位s
}

type mysql struct {
	DataSource   string
	MaxIdleConns int // 最大空闲连接数, <=MaxOpenConns
	MaxOpenConns int // 最大连接数
}

type rabbitMQ struct {
	Host     string
	Port     int
	User     string
	Password string
}

type cacheRedis struct {
	Host     string
	Port     int
	Password string
}

type system struct {
	AdminDefaultPwd string
	UserDefaultPwd  string
}

type fileServer struct {
	PrefixPath        string
	RelativePath      string
	AudioSamplePath   string //声音采样资源
	VenueFileSavePath string //场馆采样资源
	AudioFileSavePath string //音频文件存储路径
	ImageFileSavePath string //图片文件存储路径
	PptFileSavePath   string //ppt文件存储路径
	NaviResourcePath  string //数字人资源路径
	PresentationVideo string //构建生成的宣讲视频
}
