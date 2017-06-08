package main

import (
	"io/ioutil"
)

//WatchFile watch for changes
func WatchFile(filename string) {

}

//ParseFile read the content to map
func ParseFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	return string(b), err
}
