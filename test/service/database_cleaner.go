package service

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

const hookName = "cleanupHook"

type entity struct {
	table   string
	keyname string
	key     interface{}
}

type databaseClener struct {
	entries []entity
}

func (dc *databaseClener) InitDatabaseCleaner() {
	if ts.database.Callback().Create().After("gorm:create").Get(hookName) != nil {
		return
	}

	ts.database.Callback().Create().After("gorm:create").Register(hookName, func(scope *gorm.Scope) {
		dc.entries = append(dc.entries, entity{
			table:   scope.TableName(),
			keyname: scope.PrimaryKey(),
			key:     scope.PrimaryKeyValue(),
		})
	})
}

func (dc *databaseClener) ClearDatabase() {
	_, inTransaction := ts.database.CommonDB().(*sql.Tx)

	tx := ts.database

	if !inTransaction {
		tx = ts.database.Begin()
	}

	for i := range dc.entries {
		entry := dc.entries[i]

		tx.Table(entry.table).Where(entry.keyname+" = ?", entry.key).Delete("")
	}

	if !inTransaction {
		tx.Commit()
	}
}
