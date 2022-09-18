package tootls

import (
	"errors"
	"log"
	"reflect"
	"sync"
	"time"
)

var cornCount = 0

var countMutex sync.Mutex

var task map[int]chan bool

var data map[int]interface{}

var taskMutex sync.Mutex

// AddCornFunc 添加定时任务
// interval 时间间隔
// count 重复次数
// fn 任务函数
func AddCornFunc(interval time.Duration, count int, fn interface{}, ad interface{}) (id int) {
	taskMutex.Lock()
	if task == nil {
		task = make(map[int]chan bool)
		data = make(map[int]interface{})
	}
	taskMutex.Unlock()

	typeOfFn := reflect.TypeOf(fn)
	if typeOfFn.Kind() != reflect.Func {
		panic("fn is not function!")
	}

	countMutex.Lock()
	cornCount++
	id = cornCount
	countMutex.Unlock()

	data[id] = ad

	go func(id int) {
		t := time.NewTicker(interval)
		flag := make(chan bool)
		task[id] = flag
		c := 0
		for {
			select {
			case <-flag:
				t.Stop()
				delete(task, id)
				delete(data, id)
				return
			case <-t.C:
				c++
				if c > count {
					t.Stop()
					delete(task, id)
					delete(data, id)
					return
				}
				log.Printf("task_id: %d c: %d\n", id, c)
				fn.(func())()
			}
		}
	}(id)
	return
}

func StopCorn(id int) error {
	if task[id] == nil {
		return errors.New("id is not exist")
	}
	task[id] <- true
	return nil
}

func GetTaskDataList() map[int]interface{} {
	return data
}
