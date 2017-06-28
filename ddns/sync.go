package ddns

import (
	"errors"
	"sync"

	"github.com/muka/ddns/client"
	"github.com/muka/ddns/client/d_dns_service"
	"github.com/muka/ddns/models"
	"github.com/muka/ovpndns/parser"
	log "github.com/sirupsen/logrus"
)

var state = make([]*parser.Record, 0)
var mux sync.Mutex

var dnsclient *client.APIService

//CreateClient create a DDNS api client
func CreateClient(host string) {
	if dnsclient == nil {

		// create the API client, with the transport
		cfg := client.TransportConfig{
			BasePath: "",
			Host:     host,
			Schemes:  []string{"http"},
		}
		dnsclient = client.NewHTTPClientWithConfig(nil, &cfg)
	}
}

func getClient() *client.APIService {
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

	record := models.APIRecord{}

	record.Domain = domain
	record.IP = ip
	record.Type = "A"
	record.PTR = true

	params := d_dns_service.NewSaveRecordParams()
	params.SetBody(&record)

	_, err := getClient().DDNSService.SaveRecord(params)

	return err
}

//DeleteRecord remove a record
func DeleteRecord(domain string) error {

	params := d_dns_service.NewDeleteRecordParams()
	params.SetDomain(domain)
	params.SetType("A")

	_, err := getClient().DDNSService.DeleteRecord(params)

	return err
}
