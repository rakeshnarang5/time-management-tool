package main

import (
	//"./res"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ctcpip/notifize"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	work       = time.Minute * 25
	shortBreak = time.Minute * 10
	longBreak  = time.Minute * 30
)

type Task struct {
	Name  string
	Count int
}

type TaskMap struct {
	Tasks   map[string]Task
	Count   int
	Session int
}

func (t *TaskMap) addTask() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Name your task")
	name, _ := reader.ReadString('\n')
	t.Count++
	key := strconv.Itoa(t.Count)
	t.Tasks[key] = Task{name[:len(name)-1], 1}
}

func timeconv(pomo int) time.Duration {
	var retVal time.Duration
	for pomo > 0 {
		retVal += work
		pomo--
	}
	return retVal
}

func (t *TaskMap) status(flag bool) {
	fmt.Printf("Total tasks: %d, Sessions: %d\n", t.Count, t.Session)
	var str string
	for _, tsk := range t.Tasks {
		totalTime := timeconv(tsk.Count)
		//fmt.Printf("Key: %s\nName: %s\nSessions: %d\nTime Spent: %02d:%02d\n\n",key, tsk.Name,tsk.Count, hh, mm)
		str += fmt.Sprintf("%s ~ %v ~ %d", tsk.Name, totalTime, tsk.Count)
	}
	fmt.Println(str)
	if flag {
		Notify(str)
	}
}

func ticker(d time.Duration) {
	for d >= time.Second*0 {
		print("\033[H\033[2J")
		if d == time.Second*0 {
			fmt.Println("timer finished!")
		} else {
			fmt.Println(d)
		}
		d -= time.Second * 1
		time.Sleep(time.Second * 1)
	}
}

func (t *TaskMap) update() {
	reader := bufio.NewReader(os.Stdin)
	if len(t.Tasks) == 0 {
		fmt.Println("Press 2 to enter a new task")
	} else {
		fmt.Println("Press 1 if you work on an existing task\nPress 2 if you worked on a new task")
	}

	flag := false
	var input string

	for !flag {
		input, _ = reader.ReadString('\n')

		if input == "1\n" || input == "2\n" {
			flag = true
		} else {
			fmt.Println("wrong input")
		}
	}

	if input == "1\n" {
		t.status(false)
		fmt.Println("Which task did you work on?")
		chc, _ := reader.ReadString('\n')
		tsk := t.Tasks[chc[:len(chc)-1]]
		tsk.Count++
		t.Tasks[chc[:len(chc)-1]] = tsk
		fmt.Println("updated!")
		t.status(true)
	}
	if input == "2\n" {
		t.addTask()
	}
	t.saveState()
	t.timer("break")
}

func Notify(msg string) {
	if msg == "" {
		msg = "No tasks to show!"
	}
	notifize.Display("Status", msg, false, "/home/nagarro/workspace/src/timeManager/img/time.jpg")
	//res.SendMessage(msg)
}

func (t *TaskMap) timer(tsk string) {
	var d time.Duration
	if tsk == "work" {
		d = work
	} else if tsk == "break" {
		if t.Session%4 == 0 {
			tsk = "longBreak"
			d = longBreak
		} else {
			tsk = "shortBreak"
			d = shortBreak
		}
	}

	go ticker(d)

	timer := time.NewTimer(d)
	<-timer.C
	if tsk == "work" {
		Notify("Timer finished!\nWhat did you work on?")
		t.Session++
		t.update()
	} else if tsk == "shortBreak" {
		Notify("short break finished")
	} else if tsk == "longBreak" {
		Notify("long break finished")
	}
}

func (t *TaskMap) saveState() {
	m, _ := json.Marshal(t)
	ioutil.WriteFile("./task.json", m, 0755)
}

func main() {
	taskmap := initializeTaskMap()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Would you like to continue with saved tasks (Y/n)")
	input, _ := reader.ReadString('\n')
	//fmt.Println([]byte(input))

	if input == "Y\n" {
		data, _ := ioutil.ReadFile("./task.json")
		json.Unmarshal(data, &taskmap)
		taskmap.Session = 0
	}

	if input == "n\n" {
		taskmap.saveState()
	}

	//fmt.Println(taskmap)

	var wg sync.WaitGroup

	wg.Add(1)

	go menu(taskmap)

	//fmt.Println(taskmap)

	// task1 := Task{"first task", 1}
	// task2 := Task{"second task", 2}
	// task3 := Task{"third task", 3}
	// taskmap.Tasks["1"] = task1
	// taskmap.Tasks["2"] = task2
	// taskmap.Tasks["3"] = task3
	// fmt.Println(taskmap.Tasks)
	// m, _ := json.Marshal(taskmap.Tasks)
	// fmt.Println(string(m))
	// x, _ := json.Marshal(taskmap)
	// fmt.Println(string(x))
	// ioutil.WriteFile("./task.json", x, 0755)
	//unmarshal
	//var y TaskMap
	// data, _ := ioutil.ReadFile("./task.json")
	// json.Unmarshal(data, &y)
	// fmt.Println(y)
	wg.Wait()
}

func menu(t *TaskMap) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Choose an Action:-\nPress 1 to see current tasks\nPress 2 to start working\nPress 3 to end the day")
		input, _ := reader.ReadString('\n')
		if input == "1\n" {
			t.status(true)
		}
		if input == "2\n" {
			t.timer("work")
		}
		if input == "3\n" {
			t.saveState()
			os.Exit(0)
		}
	}

}

func initializeTaskMap() *TaskMap {
	var taskmap TaskMap
	taskmap.Tasks = make(map[string]Task)
	taskmap.Count = 0
	taskmap.Session = 0
	return &taskmap
}
