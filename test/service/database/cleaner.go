package database

import (
	"database/sql"
	"errors"
	"log"
	"runtime"
	"sync"

	"github.com/jinzhu/gorm"
)

const numberOfstackFramesToAscend = 2

var dc *databaseCleaner
var once sync.Once

type databaseCleaner struct {
	enabledHooks []*cleanupHook
	database     *gorm.DB
}

func Cleaner(db *gorm.DB) *databaseCleaner {
	once.Do(func() {
		dc = &databaseCleaner{database: db}
	})

	testFunctionName, err := retrieveTestFunctionName()

	if err != nil {
		log.Fatal(err)
	}

	if !dc.checkByNameWhetherHookEnabled(testFunctionName) {
		log.Fatal("Hook " + testFunctionName + " already enabled")
	}

	hook := newCleanUpHook(testFunctionName)

	db.Callback().Create().After("gorm:create").Register(hook.name(), func(scope *gorm.Scope) {
		hook.pushEntity(
			newEntity(scope.TableName(), scope.PrimaryKey(), scope.PrimaryKeyValue()),
		)
	})

	dc.enabledHooks = append(dc.enabledHooks, hook)

	return dc
}

func (dc *databaseCleaner) CleanUp() {
	testFuctionName, err := retrieveTestFunctionName()

	if err != nil {
		log.Fatal(err)
	}

	hook, err := dc.retrieveHookByName(testFuctionName)

	if err != nil {
		log.Fatal(err)
	}

	defer dc.database.Callback().Create().Remove(hook.name())

	_, inTransaction := dc.database.CommonDB().(*sql.Tx)

	tx := dc.database

	if !inTransaction {
		tx = dc.database.Begin()
	}

	entities := hook.entities()

	for i := range entities {
		tx.Table(entities[i].table()).Where(entities[i].keyName()+" = ?", entities[i].keyValue()).Delete("")
	}

	if !inTransaction {
		tx.Commit()
	}
}

func (dc *databaseCleaner) checkByNameWhetherHookEnabled(hookName string) bool {
	for i := range dc.enabledHooks {
		if dc.enabledHooks[i].name() == hookName {
			return false
		}
	}

	return true
}

func (dc *databaseCleaner) retrieveHookByName(hookName string) (*cleanupHook, error) {
	var ch *cleanupHook

	for i := range dc.enabledHooks {
		if hookName == dc.enabledHooks[i].name() {
			return dc.enabledHooks[i], nil
		}
	}

	return ch, errors.New("Hook hasn't enabled " + hookName)
}

func retrieveTestFunctionName() (string, error) {
	pc, _, _, ok := runtime.Caller(numberOfstackFramesToAscend)

	if !ok {
		return "", errors.New("An error has occurred during retrieving a test function name")
	}

	return runtime.FuncForPC(pc).Name(), nil
}
