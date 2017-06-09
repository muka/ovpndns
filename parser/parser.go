package parser

import (
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

var updates = make(chan []Record)
var mux sync.Mutex

//Record of OpenVPN log
type Record struct {
	IP        string
	Name      string
	RemoteIP  string
	Connected time.Time
}

//GetChannel the channel used to notify updates
func GetChannel() chan []Record {
	return updates
}

//WatchFile watch for changes
func WatchFile(filename string) {
	watcher, err1 := fsnotify.NewWatcher()
	if err1 != nil {
		log.Fatal(err1)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Debugf("event: %v", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Debugf("modified file: %s", event.Name)
					go ParseFile(filename)
				}
			case err := <-watcher.Errors:
				log.Errorf("error watching %s: %s", filename, err.Error())
			}
		}
	}()

	log.Debugf("Watching %s", filename)
	err1 = watcher.Add(filename)
	if err1 != nil {
		log.Fatal(err1)
	}
	<-done
}

func readFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	return string(b), err
}

func containRecord(record Record, list []Record) bool {
	found := false
	for _, r := range list {
		if r.IP == record.IP {
			return true
		}
	}
	return found
}

//ParseFile read the content to map
func ParseFile(filename string) error {

	mux.Lock()
	content, err := readFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(content, "\n")
	records := make([]Record, 0)

	for _, line := range lines {

		// 10.0.0.6,test1,172.23.0.1:44623,Tue Jun  6 06:29:28 2017
		r, _ := regexp.Compile(`([.0-9]+)[,]([a-zA-Z0-9._-]+)[,]([\.0-9]+)[:][0-9]+[,](.*)$`)
		matches := r.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		date, err := time.Parse(time.ANSIC, matches[4])
		if err != nil {
			log.Warnf("Cannot parse date: %s", err.Error())
			date = time.Now()
		}

		record := Record{
			IP:        matches[1],
			Name:      matches[2],
			RemoteIP:  matches[3],
			Connected: date,
		}

		log.Debugf("Adding %s (%s -> %s) connected on %s", record.Name, record.IP, record.RemoteIP, record.Connected)
		records = append(records, record)
	}

	log.Debugf("Found %d records", len(records))
	updates <- records

	mux.Unlock()

	return nil
}
