package timewheel

import (
	"container/list"
	"time"
	//"fmt"
)


type Job func(interface{})

type TimeWheelRelative struct {
	since_from_start time.Duration
	interval  		 time.Duration
	ticker  		 *time.Ticker
	currentPos       int
	slotNum          int
	slots   		 []*list.List
	addTask 		 chan Task
	job 			 Job
}

type Data interface{}

type Task struct {
	entry  		time.Duration 
	circle 		int
	data   	    interface{} 
}

func New(interval time.Duration, slot_num int, job Job) *TimeWheelRelative {
	return &TimeWheelRelative {
		since_from_start : time.Duration(0),
		interval 		 : interval,
		currentPos		 : 0,
		slotNum			 : slot_num,
		slots 			 : make([]*list.List, slot_num),
		addTask 		 : make(chan Task),
		job				 : job,
	}	
}

func (tw *TimeWheelRelative) Start() {

	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}

	tw.run()
}

func (tw *TimeWheelRelative) run() {

	tw.ticker = time.NewTicker(tw.interval)

	for {
		select {
		case <- tw.ticker.C :
			//fmt.Printf("since : %v\n", tw.since_from_start)
			tw.tickHandle()
			tw.since_from_start += tw.interval
		case task := <- tw.addTask :
			tw.addTaskHandle(&task)
		}
	}
}

func (tw *TimeWheelRelative) tickHandle() {
	l := tw.slots[tw.currentPos]

	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.circle > 0 {
			task.circle--
			e = e.Next()
			continue
		}
	
		go tw.job(task.data)

		n := e.Next()
		l.Remove(e)
		e = n
	}

	tw.currentPos = (tw.currentPos + 1) % tw.slotNum
}

func (tw *TimeWheelRelative) AddTask(entry time.Duration, data interface{}) {
	tw.addTask <- Task{entry : entry, data : data}
}

func (tw *TimeWheelRelative) addTaskHandle(task *Task) {
	//fmt.Printf("------------------%v:%v\n", task.entry, tw.since_from_start)
	if task.entry < tw.since_from_start { //expiration task direct run
		go tw.job(task.data)
		return
	}

	delay  		:= task.entry - tw.since_from_start
	task.circle  = ( int( delay / tw.interval ) / tw.slotNum)
	pos    		:= (tw.currentPos + int( delay / tw.interval ) ) % tw.slotNum

	//fmt.Printf("%v:%v:%v\n", delay, task.circle, pos)

	tw.slots[pos].PushBack(task)

}