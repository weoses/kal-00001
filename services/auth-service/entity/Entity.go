package entity

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	Id        string `gorm:"type:uuid;primaryKey"`
	AccountId string `gorm:"type:uuid;uniqueIndex;not null"`
	CreatedAt time.Time
}

func (User) TableName() string { return "users" }

type Permission struct {
	Id   int16  `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex;not null"`
}

func (Permission) TableName() string { return "permissions" }

type UserPermission struct {
	UserId       string `gorm:"type:uuid;primaryKey"`
	PermissionId int16  `gorm:"primaryKey"`
}

func (UserPermission) TableName() string { return "user_permissions" }

// IntegrationTelegram links a user to every Telegram numeric ID that should
// be treated as that same user. A Telegram ID not present in any row's
// TelegramIds is anonymous, even if it belongs to a person who is a
// registered user via a different linked ID.
type IntegrationTelegram struct {
	Id            string        `gorm:"type:uuid;primaryKey"`
	UserId        string        `gorm:"type:uuid;not null"`
	TelegramIds   pq.Int64Array `gorm:"type:bigint[];not null"`
	AnswerUnknown bool
}

func (IntegrationTelegram) TableName() string { return "integration_telegram" }

type IntegrationWebappBasic struct {
	Id           string `gorm:"type:uuid;primaryKey"`
	UserId       string `gorm:"type:uuid;not null"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
}

func (IntegrationWebappBasic) TableName() string { return "integration_webapp_basic" }
