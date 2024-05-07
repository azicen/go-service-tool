package database

import (
	"time"

	"github.com/azicen/go-service-tool/uuid"
	"gorm.io/gorm"
)

type IDModel struct {
	ID uint64 `json:"id" gorm:"column:id; primaryKey"`
}

type UniqueUUIDModel struct {
	UniqueUUID uuid.UUID `json:"unique_uuid" gorm:"column:unique_uuid; uniqueIndex; not null"`
}

type MetadataTimeModel struct {
	Ctime time.Time `json:"ctime" gorm:"column:ctime; not null"` // 状态最后一次更改
	Mtime time.Time `json:"mtime" gorm:"column:mtime; not null"` // 数据最后一次修改
}

type DeleteModel struct {
	Delete bool `json:"deleted" gorm:"column:deleted; not null"` // 数据软删除
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
		// 创建时间
		createTimeField := db.Statement.Schema.LookUpField("ctime")
		if createTimeField != nil {
			err := createTimeField.Set(db.Statement.Context, db.Statement.ReflectValue, nowTime)
			if err != nil {
				db.Logger.Error(db.Statement.Context, "Failed to initialization creation timel: %v", err)
			}
		}

		// 更改时间
		modifyTimeField := db.Statement.Schema.LookUpField("mtime")
		if modifyTimeField != nil {
			err := modifyTimeField.Set(db.Statement.Context, db.Statement.ReflectValue, nowTime)
			if err != nil {
				db.Logger.Error(db.Statement.Context, "Failed to initialization modification time: %v", err)
			}
		}

		// 删除状态
		modifyDeleteField := db.Statement.Schema.LookUpField("deleted")
		if modifyDeleteField != nil {
			err := modifyDeleteField.Set(db.Statement.Context, db.Statement.ReflectValue, false)
			if err != nil {
				db.Logger.Error(db.Statement.Context, "Failed to initialization deleted: %v", err)
			}
		}

		// 唯一uuid
		uniqueUUIDField := db.Statement.Schema.LookUpField("unique_uuid")
		if uniqueUUIDField != nil {
			ou, err := uuid.NewOrderedUUID()
			if err != nil {
				db.Logger.Error(db.Statement.Context, "Failed to generate uuid: %v", err)
			}
			err = uniqueUUIDField.Set(db.Statement.Context, db.Statement.ReflectValue, ou)
			if err != nil {
				db.Logger.Error(db.Statement.Context, "Failed to initialization unique uuid: %v", err)
			}
		}
	}
}

func updateCallback(db *gorm.DB) {
	if db.Statement.Schema != nil {
		nowTime := time.Now()
		modifyTimeField := db.Statement.Schema.LookUpField("mtime")
		if modifyTimeField != nil {
			err := modifyTimeField.Set(db.Statement.Context, db.Statement.ReflectValue, nowTime)
			if err != nil {
				db.Logger.Error(db.Statement.Context, "Failed to update modification time: %v", err)
			}
		}
	}
}

func deleteCallback(db *gorm.DB) {
	if db.Statement.Schema != nil {
		modifyDeleteField := db.Statement.Schema.LookUpField("deleted")
		if modifyDeleteField != nil {
			err := modifyDeleteField.Set(db.Statement.Context, db.Statement.ReflectValue, true)
			if err != nil {
				db.Logger.Error(db.Statement.Context, "Failed to set soft delete: %v", err)
			}
		}

		nowTime := time.Now()
		modifyTimeField := db.Statement.Schema.LookUpField("mtime")
		if modifyTimeField != nil {
			err := modifyTimeField.Set(db.Statement.Context, db.Statement.ReflectValue, nowTime)
			if err != nil {
				db.Logger.Error(db.Statement.Context, "Failed to set modification time to soft delete operation time: %v", err)
			}
		}
	}
}
