package catdog_data

import "sync"

var data sync.Map

func Add(key string, value interface{}) {
	data.Store(key, value)
}

func Delete(key string) {
	data.Delete(key)
}

func Update(key string, value interface{}) {
	data.Store(key, value)
}

func Get(key string) interface{} {
	val, _ := data.Load(key)
	return val
}

func List() map[string]interface{} {
	dt := make(map[string]interface{})
	data.Range(func(key, value interface{}) bool {
		dt[key.(string)] = value
		return true
	})
	return dt
}
