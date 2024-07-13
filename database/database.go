package database

import (
	"time"

	"github.com/azicen/go-service-tool/uuid"
)

type IDModel struct {
	ID uint64 `json:"id" db:"id" gorm:"column:id; primaryKey"`
}

type UniqueUUIDModel struct {
	UUID uuid.UUID `json:"uuid" db:"uuid" gorm:"column:uuid; uniqueIndex; not null"`
}

func (m UniqueUUIDModel) Init() {
	ou, err := uuid.NewOrderedUUID()
	if err != nil {
		// 降级操作，使用UUIDv4
		ou = uuid.New()
	}
	m.UUID = ou
}

type MetadataTimeModel struct {
	Ctime time.Time `json:"ctime" db:"ctime" gorm:"column:ctime; not null"` // 状态最后一次更改
	Mtime time.Time `json:"mtime" db:"mtime" gorm:"column:mtime; not null"` // 数据最后一次修改
}

func (m MetadataTimeModel) Init() {
	m.Ctime = time.Now()
	m.Mtime = m.Ctime
}

type DeleteModel struct {
	Delete bool `json:"deleted" db:"deleted" gorm:"column:deleted; not null"` // 数据软删除
}

func (m DeleteModel) Init() {
	m.Delete = false
}

type BaseModel struct {
	IDModel
	MetadataTimeModel
	DeleteModel
}

func (m BaseModel) Init() {
	m.MetadataTimeModel.Init()
	m.DeleteModel.Init()
}

type BaseUUIDModel struct {
	BaseModel
	UniqueUUIDModel
}

func (m BaseUUIDModel) Init() {
	m.BaseModel.Init()
	m.UniqueUUIDModel.Init()
}
