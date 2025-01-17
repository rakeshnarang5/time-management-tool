package main

import (
	"./res"
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
	// work       = time.Minute * 25
	// shortBreak = time.Minute * 5
	// longBreak  = time.Minute * 20

	work       = time.Second * 1
	shortBreak = time.Second * 1
	longBreak  = time.Second * 1
)

type Task struct {
	Name  string
	Count []string
}

type TaskMap struct {
	Tasks   map[string]*Task
	Count   int
	Session int
}

func (t *TaskMap) addTask() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Name your task")
	name, _ := reader.ReadString('\n')
	fmt.Println("(comment) What did you do?")
	comment, _ := reader.ReadString('\n')
	t.Count++
	tsk := initializeTask(name[:len(name)-1], comment[:len(comment)-1])
	key := strconv.Itoa(t.Count)
	t.Tasks[key] = tsk
}

func initializeTask(name string, comment string) *Task {
	var retVal Task
	retVal.Count = make([]string, 0, 0)
	retVal.Name = name
	retVal.Count = append(retVal.Count, comment)
	return &retVal
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
	for key, tsk := range t.Tasks {
		totalTime := timeconv(len(tsk.Count))
		str += fmt.Sprintf("[K: %s] [N: %s] [T: %v] [C: %d]\n", key, tsk.Name, totalTime, len(tsk.Count))
	}
	fmt.Println(str)
	if flag {
		Notify(str, false)
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
		fmt.Println("(comment) What did you do?")
		comment, _ := reader.ReadString('\n')
		tsk.addEntry(comment[:len(comment)-1])
		//tsk.Count++
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

func (t *Task) addEntry(comment string) {
	t.Count = append(t.Count, comment)
}

func Notify(msg string, sendEmail bool) {
	if msg == "" {
		msg = "No tasks to show!"
	}
	notifize.Display("Status", msg, false, "/home/nagarro/workspace/src/timeManager/img/time.jpg")
	//res.SendMessage(msg)
	if sendEmail {
		res.SendEmail(msg, "rakeshnarang5@gmail.com")
	}
}

func (t *TaskMap) prepareEmail() string {
	var retVal string
	retVal += "<h1>Current Complete Report:-</h1>"
	retVal += "<table border=\"1\"><thead><th>Key</th><th>Name</th><th>Sessions</th><th>Entries</th></thead><tbody><tbody>"

	for key, val := range t.Tasks {
		retVal += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%d</td><td>%s</td></tr>", key, val.Name, len(val.Count), val.Count[0])
		if len(val.Count) > 1 {
			for i := 1; i < len(val.Count); i++ {
				retVal += fmt.Sprintf("<tr><td></td><td></td><td></td><td>%s</td></tr>", val.Count[i])
			}
		}
	}

	retVal += "</tbody></table>"
	return retVal
}

// func (t *TaskMap) returnJSON() string {
// 	m, _ := json.MarshalIndent(t, "", "    ")
// 	return string(m)
// }

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
		Notify("Timer finished!\nWhat did you work on?", false)
		t.Session++
		t.update()
	} else if tsk == "shortBreak" {
		Notify("short break finished", false)
	} else if tsk == "longBreak" {
		Notify("long break finished", false)
	}
}

func (t *TaskMap) saveState() {
	m, _ := json.MarshalIndent(t, "", "    ")
	ioutil.WriteFile("./task.json", m, 0755)
}

func main() {
	taskmap := initializeTaskMap()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Would you like to continue with saved tasks (Y/n)")
	input, _ := reader.ReadString('\n')

	if input == "Y\n" {
		data, _ := ioutil.ReadFile("./task.json")
		json.Unmarshal(data, &taskmap)
		taskmap.Session = 0
	}

	if input == "n\n" {
		taskmap.saveState()
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go menu(taskmap)

	wg.Wait()
}

func menu(t *TaskMap) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Choose an Action:-\nPress 1 to see current tasks\nPress 2 to start working\nPress 3 to end the day\nPress 4 to email current report")
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
		if input == "4\n" {
			msg := t.prepareEmail()
			Notify(msg, true)
		}
	}

}

func initializeTaskMap() *TaskMap {
	var taskmap TaskMap
	taskmap.Tasks = make(map[string]*Task)
	taskmap.Count = 0
	taskmap.Session = 0
	return &taskmap
}
