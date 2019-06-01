package ddns

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/muka/ddns/api"
	"github.com/muka/ovpndns/parser"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var state = make([]*parser.Record, 0)
var mux sync.Mutex

var dnsclient api.DDNSServiceClient

//CreateClient create a DDNS api client
func CreateClient(host string) {
	if dnsclient == nil {

		conn, err := grpc.Dial(host, grpc.WithInsecure())
		if err != nil {
			fmt.Printf("Failed to connect to ddns at %s: %s", host, err)
			os.Exit(1)
		}

		dnsclient = api.NewDDNSServiceClient(conn)
	}
}

func getClient() api.DDNSServiceClient {
	if dnsclient == nil {
		panic(errors.New("Client not initialized. Call CreateClient first"))
	}

	return dnsclient
}

func has(key string, store []*parser.Record) bool {
	for i := 0; i < len(store); i++ {
		if store[i].IP == key {
			return true
		}
	}
	return false
}

// Compare a map and sync with ddns
func Compare(records []*parser.Record, domain string) error {

	mux.Lock()

	var werr error

	// find new
	for _, record := range records {

		domainName := record.Name + "." + domain

		if !has(record.IP, state) {
			log.Debugf("Saving DNS record %s", domainName)

			err := SaveRecord(domainName, record.IP)
			if err != nil {
				log.Errorf("Error saving %s: %s", domainName, err.Error())
				werr = err
				continue
			}

			state = append(state, record)
		} else {
			log.Debugf("Skip %s", domainName)
		}
	}

	// find deleted
	for i, record := range state {

		domainName := record.Name + "." + domain

		if !has(record.IP, records) {

			log.Debugf("Removing DNS record %s", domainName)

			err := DeleteRecord(domainName)
			if err != nil {
				log.Errorf("Error removing %s: %s", domainName, err.Error())
				werr = err
				continue
			}

			if i < len(state) {
				// unreference for GC
				state[i] = nil
				// delete element
				state = state[:i+copy(state[i:], state[i+1:])]
			}
		}
	}

	log.Debugf("State has %d, records has %d", len(state), len(records))

	mux.Unlock()

	return werr
}

//SaveRecord store a record
func SaveRecord(domain string, ip string) error {

	record := new(api.Record)
	record.Domain = domain
	record.Ip = ip
	record.Type = "A"
	record.PTR = true

	c := getClient()

	ctx1 := context.Background()
	ctx, cancel := context.WithTimeout(ctx1, time.Millisecond*500)
	defer cancel()
	_, err := c.SaveRecord(ctx, record)

	return err
}

//DeleteRecord remove a record
func DeleteRecord(domain string) error {

	record := new(api.Record)
	record.Domain = domain
	record.Type = "A"

	c := getClient()
	ctx1 := context.Background()
	ctx, cancel := context.WithTimeout(ctx1, time.Millisecond*500)
	defer cancel()
	_, err := c.DeleteRecord(ctx, record)
	return err
}
