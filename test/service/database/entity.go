package database

type entity struct {
	tableName string
	kName     string
	kValue    interface{}
}

func newEntity(table, keyName string, keyValue interface{}) entity {
	return entity{
		tableName: table,
		kName:     keyName,
		kValue:    keyValue,
	}
}

func (e entity) table() string {
	return e.tableName
}

func (e entity) keyName() string {
	return e.kName
}

func (e entity) keyValue() interface{} {
	return e.kValue
}
