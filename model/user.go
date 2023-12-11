package model

// 普通用户
type User struct {
	BaseModel
	Phonenum  string `gorm:"type:varchar(20);unique;not null;comment:'手机号'" json:"phonenum"`
	Password  string `gorm:"type:varchar(200);not null;comment:'登录密码'" json:"-"`
	Email     string `gorm:"type:varchar(50);default null;comment:'邮箱'" json:"email"`
	Avatar    string `gorm:"type:varchar(200);default null;comment:'头像路径'" json:"avatar"`
	Status    string `gorm:"type:varchar(1);not null;default:1;comment:'账号状态：1：正常，0：禁用'" json:"status"`
	Token     string `gorm:"type:varchar(500);default null;comment:'身份令牌'" json:"token"`
	LastLogin string `gorm:"type:varchar(50);default null;comment:'上次登录时间'" json:"last_login"`
}

func (User) TableName() string {
	return "gb_user"
}

// 账号状态：1：正常，0：禁用
const (
	AccountStatus_Avaliable = "1"
	AccountStatus_Disabled  = "0"
)
