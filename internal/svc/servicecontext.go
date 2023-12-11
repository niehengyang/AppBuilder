package svc

import (
	"appBuilder/ebyte/logger"
	"appBuilder/internal/config"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"path"
)

type ServiceContext struct {
	DB     *gorm.DB
	Config config.Config
	Redis  *redis.Client
}

func NewServiceContext(c config.Config, moduleName string) (*ServiceContext, error) {
	accessLogFile := path.Join(c.Log.Path, moduleName+"-access.log")
	errorLogFile := path.Join(c.Log.Path, moduleName+"-error.log")

	//初始化Zap日志
	var zapConfig = zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:      []string{"stdout", accessLogFile}, // 控制台和日志文件路径
		ErrorOutputPaths: []string{"stderr", errorLogFile},  // 控制台和错误日志文件路径
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	}
	err := logger.InitLogger(zapConfig)
	if err != nil {
		return nil, fmt.Errorf("初始化日志包失败：%s\n", err.Error())
	}
	//初始化数据库
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 数据迁移时不生成外键
	})
	if err != nil {
		return nil, fmt.Errorf("初始化数据库失败：%s\n", err.Error())
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.CacheRedis.Host, 6379), // Redis 服务器地址
		Password: c.CacheRedis.Password,                         // Redis 密码
		DB:       0,                                             // 选择的 Redis 数据库
	})

	return &ServiceContext{
		DB:     db,
		Config: c,
		Redis:  redisClient,
	}, nil
}
