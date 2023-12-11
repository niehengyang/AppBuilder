package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// XTime
// @Description: 1. 创建 time.Time 类型的副本 XTime；
type XTime struct {
	time.Time
}

// TimeFormat 时间常量
const TimeFormat = "2006-01-02 15:04:05"

// UnmarshalJSON
//
//	@Description: 解析JSON数据
//	@receiver t 时间
//	@param data 字节数据
//	@return error 返回错误消息
func (t *XTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = XTime{now}
	return err
}

// MarshalJSON
//
//	@Description: 为 Xtime 重写 MarshaJSON 方法，在此方法中实现自定义格式的转换；
//	@receiver t 时间
//	@return []byte 字节数据
//	@return error 返回错误消息
func (t XTime) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(output), nil
}

// Value
//
//	@Description: 为 Xtime 实现 Value 方法，写入数据库时会调用该方法将自定义时间类型转换并写入数据库；
//	@receiver t 时间
//	@return driver.Value 操作值
//	@return error 返回错误消息
func (t XTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan
//
//	@Description: 为 Xtime 实现 Scan 方法，读取数据库时会调用该方法将时间数据转换成自定义时间类型；
//	@receiver t 时间
//	@param v 对象参数
//	@return error 返回错误消息
func (t *XTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = XTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// GetFormattedTime 获取当前时间并按"2006-01-02 15:04:05"格式转换为字符串
func GetFormattedTime() string {
	currentTime := time.Now()
	return currentTime.Format(TimeFormat)
}

type BaseModel struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	CreatedAt XTime          `json:"created_at"`
	UpdatedAt XTime          `json:"-"`
	DeletedAt gorm.DeletedAt `sql:"index" json:"-"`
}

// LoginAccount
//
//	@Description: 根据账号ID过滤
//	@param ID 登录用户ID
//	@return func(db *gorm.DB) *gorm.DB 返回Mysql数据表操作句柄
func LoginAccount(ID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("userid = ?", ID)
	}
}

// VisibleDataFilter
//
//	@Description: 过滤账号可见数据，根据各数据表的创建人ID字段进行过滤；
//	@param account 账号
//	@return func(db *gorm.DB) *gorm.DB MySQL操作句柄
func VisibleDataFilter(account *User) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_by = ?", account.ID)
	}
}
