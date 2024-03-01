package database

import (
	"time"

	"github.com/cypunsource/cypunsource-tool/uuid"
	"gorm.io/gorm"
)

type IDModel struct {
	ID uint64 `json:"id" gorm:"column:id; primaryKey"`
}

type UniqueUUIDModel struct {
	UniqueUUID *uuid.UUID `json:"unique_uuid" gorm:"column:unique_uuid; uniqueIndex; not null"`
}

func (model *UniqueUUIDModel) BeforeCreate(tx *gorm.DB) (err error) {
	if model.UniqueUUID == nil {
		var ou uuid.UUID
		ou, err = uuid.NewOrderedUUID()
		model.UniqueUUID = &ou
	}
	return
}

type MetadataTimeModel struct {
	Ctime time.Time `json:"ctime" gorm:"column:ctime"` // 状态最后一次更改
	Mtime time.Time `json:"mtime" gorm:"column:mtime"` // 数据最后一次修改
}

type DeleteModel struct {
	Delete bool `json:"deleted" gorm:"column:deleted"` // 数据软删除
}

type BaseModel struct {
	IDModel
	MetadataTimeModel
	DeleteModel
}

func RegisterCallbacks(db *gorm.DB) {
	err := db.Callback().Create().Before("gorm:create").Replace("gorm:create_time_stamp", createCallback)
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().Before("gorm:update").Replace("gorm:update_time_stamp", updateCallback)
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().Before("gorm:delete").Replace("gorm:delete_time_stamp", deleteCallback)
	if err != nil {
		panic(err)
	}
}

func createCallback(db *gorm.DB) {
	if db.Statement.Schema != nil {
		nowTime := time.Now()
		createTimeField := db.Statement.Schema.LookUpField("ctime")
		if createTimeField != nil {
			_ = createTimeField.Set(db.Statement.Context, db.Statement.ReflectValue, nowTime)
		}
		modifyTimeField := db.Statement.Schema.LookUpField("mtime")
		if modifyTimeField != nil {
			_ = modifyTimeField.Set(db.Statement.Context, db.Statement.ReflectValue, nowTime)
		}
	}
}

func updateCallback(db *gorm.DB) {
	if db.Statement.Schema != nil {
		nowTime := time.Now()
		modifyTimeField := db.Statement.Schema.LookUpField("mtime")
		if modifyTimeField != nil {
			_ = modifyTimeField.Set(db.Statement.Context, db.Statement.ReflectValue, nowTime)
		}
	}
}

func deleteCallback(db *gorm.DB) {
	if db.Statement.Schema != nil {
		modifyDeleteField := db.Statement.Schema.LookUpField("deleted")
		if modifyDeleteField != nil {
			_ = modifyDeleteField.Set(db.Statement.Context, db.Statement.ReflectValue, true)
		}

		nowTime := time.Now()
		modifyTimeField := db.Statement.Schema.LookUpField("mtime")
		if modifyTimeField != nil {
			_ = modifyTimeField.Set(db.Statement.Context, db.Statement.ReflectValue, nowTime)
		}
	}
}
