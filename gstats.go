// Package gstats collects HTTP server statistics and and shows them inside a web page.
package gstats

import (
	"encoding/json"
	"github.com/cevatbarisyilmaz/ip2country"
	"github.com/mssola/user_agent"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const realtimeRefreshRate = time.Second * 4
const realtimeRefreshRateInt = 4

type ipKey [16]byte

type GStats struct {
	root string

	currentMu                *sync.RWMutex
	currentUniqueIPs         map[ipKey]int
	currentConnections       int
	currentOutboundBandwidth int
	currentInboundBandwidth  int

	dataMu *sync.RWMutex
	data   *data

	// Country function is used for geolocation.
	// If not overridden, it will use https://github.com/cevatbarisyilmaz/ip2country
	// which is licensed under Creative Commons Attribution Share Alike 4.0 International (CC-BY-SA-4.0)
	// as it uses GeoLite2 database created by MaxMind (which is also licensed with CC-BY-SA-4.0).
	// Therefore, with the default value, this package falls under CC-BY-SA-4.0 License.
	// However, if you decide to use another source for geolocation by overriding this attribute,
	// you can use this package under MIT License.
	Country func(net.IP) (string, error)
}

// New returns a new GStats which will use given root folder to save needed data files.
func New(root string) *GStats {
	if !strings.HasSuffix(root, "/") {
		root += "/"
	}
	g := &GStats{
		root:                     root,
		currentMu:                &sync.RWMutex{},
		currentUniqueIPs:         map[ipKey]int{},
		currentConnections:       0,
		currentOutboundBandwidth: 0,
		currentInboundBandwidth:  0,
		dataMu:                   &sync.RWMutex{},
		Country:                  ip2country.Country,
	}
	g.data = g.restoreHourlyRecord()
	go g.work()
	return g
}

// Listener wraps the given listener and returns another listener which will
// track hosts, connections and bandwidth.
func (g *GStats) Listener(sublistener net.Listener) net.Listener {
	return &listener{
		Listener: sublistener,
		g:        g,
	}
}

// Collect wraps the given handler and returns another handler which will
// track HTTP-related statistics.
func (g *GStats) Collect(subhandler http.Handler) http.Handler {
	return &collect{
		g: g,
		h: subhandler,
	}
}

// Show returns a handler which serves the web page to browse collected statistics.
// At public-facing usages, it is advised to authenticate the users first before calling this handler.
// If handler lies at a non-root path, any prefixes should be removed before invoking the handler.
// i.e. http.Handle("/statistics/", http.StripPrefix("/statistics", gs.Show()))
func (g *GStats) Show() http.Handler {
	return getShowHandler(g)
}

// PrepareToExit saves any unsaved data to disk to not lose them on exiting the program.
// Normally, GStats does disk saves at the beginning of each hour.
func (g *GStats) PrepareToExit() {
	now := time.Now().UTC()
	basePath := g.root + strconv.Itoa(now.Year()) + "/" + strconv.Itoa(int(now.Month())) + "/" + strconv.Itoa(now.Day()) + "/" + strconv.Itoa(now.Hour())
	err := os.MkdirAll(basePath, os.ModeDir)
	if err != nil {
		log.Println("gstats: mkdir failed: " + err.Error())
	} else {
		g.dataMu.Lock()
		saveData(basePath+"/data.json", g.data)
		g.dataMu.Unlock()
	}

}

func (g *GStats) notifyNewConn(raddr net.Addr) {
	key := ipToKey(getIP(raddr))
	g.currentMu.Lock()
	g.currentConnections++
	g.currentUniqueIPs[key]++
	g.currentMu.Unlock()
}

func (g *GStats) notifyConnRead(raddr net.Addr, n int) {
	g.currentMu.Lock()
	g.currentInboundBandwidth += n
	g.currentMu.Unlock()
	go func() {
		time.Sleep(realtimeRefreshRate)
		g.currentMu.Lock()
		g.currentInboundBandwidth -= n
		g.currentMu.Unlock()
	}()
	g.dataMu.Lock()
	g.data.inboundBandwidth += n
	g.dataMu.Unlock()
}

func (g *GStats) notifyConnWrite(raddr net.Addr, n int) {
	g.currentMu.Lock()
	g.currentOutboundBandwidth += n
	g.currentMu.Unlock()
	go func() {
		time.Sleep(realtimeRefreshRate)
		g.currentMu.Lock()
		g.currentOutboundBandwidth -= n
		g.currentMu.Unlock()
	}()
	g.dataMu.Lock()
	g.data.outboundBandwidth += n
	g.dataMu.Unlock()
}

func (g *GStats) notifyConnClose(raddr net.Addr) {
	key := ipToKey(getIP(raddr))
	g.currentMu.Lock()
	g.currentConnections--
	g.currentUniqueIPs[key]--
	if g.currentUniqueIPs[key] <= 0 {
		delete(g.currentUniqueIPs, key)
	}
	g.currentMu.Unlock()
}

func (g *GStats) notifyRequest(request *http.Request, statusCode int, responseSize int, responseTime time.Duration) {
	var ip net.IP
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err == nil {
		ip = net.ParseIP(host)
	}
	userAgent := user_agent.New(request.UserAgent())
	browserKey, _ := userAgent.Browser()
	var shouldRecordReferrer bool
	referrer := request.Referer()
	referrerURL, err := url.Parse(referrer)
	if err != nil {
		shouldRecordReferrer = false
	} else {
		hostName := referrerURL.Hostname()
		if hostName == "" || hostName == request.Host {
			shouldRecordReferrer = false
		} else {
			shouldRecordReferrer = true
		}
	}
	var referrerRecord *commonRecord
	request.URL.Host = request.Host
	g.dataMu.Lock()
	g.data.requests++
	reqRec := g.getRequestPathRecord(request.URL.String())
	browserRec := g.getBrowserRecord(browserKey)
	osRec := g.getOSRecord(userAgent.OSInfo().Name)
	country, _ := g.Country(ip)
	countryRec := g.getCountryRecord(country)
	if shouldRecordReferrer {
		referrerRecord = g.getReferrerRecord(referrer)
	}
	if ip != nil {
		key := ipToKey(ip)
		g.data.uniqueIPs[key] = true
		reqRec.uniqueIPs[key] = true
		browserRec.uniqueIPs[key] = true
		osRec.uniqueIPs[key] = true
		countryRec.uniqueIPs[key] = true
		if shouldRecordReferrer {
			referrerRecord.uniqueIPs[key] = true
		}
	}
	reqRec.responseSizes[responseSize]++
	reqRec.responseTimes[responseTime]++
	reqRec.statusCodes[statusCode]++

	browserRec.requests++
	osRec.requests++
	countryRec.requests++
	if shouldRecordReferrer {
		referrerRecord.requests++
	}
	g.dataMu.Unlock()
}

func (g *GStats) getRequestPathRecord(key string) (record *requestPathRecord) {
	record = g.data.requestPaths[key]
	if record == nil {
		record = newRequestPathRecord()
		g.data.requestPaths[key] = record
	}
	return
}

func (g *GStats) getCountryRecord(key string) (record *commonRecord) {
	record = g.data.countries[key]
	if record == nil {
		record = newCommonRecord()
		g.data.countries[key] = record
	}
	return
}

func (g *GStats) getOSRecord(key string) (record *commonRecord) {
	record = g.data.oss[key]
	if record == nil {
		record = newCommonRecord()
		g.data.oss[key] = record
	}
	return
}

func (g *GStats) getBrowserRecord(key string) (record *commonRecord) {
	record = g.data.browsers[key]
	if record == nil {
		record = newCommonRecord()
		g.data.browsers[key] = record
	}
	return
}

func (g *GStats) getReferrerRecord(key string) (record *commonRecord) {
	record = g.data.referrers[key]
	if record == nil {
		record = newCommonRecord()
		g.data.referrers[key] = record
	}
	return
}

func (g *GStats) restoreHourlyRecord() *data {
	now := time.Now().UTC()
	path := g.root + strconv.Itoa(now.Year()) + "/" + strconv.Itoa(int(now.Month())) + "/" + strconv.Itoa(now.Day()) + "/" + strconv.Itoa(now.Hour()) + "/data.json"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return newData()
	}
	err = os.Remove(path)
	if err != nil {
		log.Println("gstats: failed to delete restored data file: " + err.Error())
	}
	data, err := constructData(content)
	if err != nil {
		return newData()
	}
	return data.toData()
}

type serializableCurrent struct {
	UniqueIPs         int
	Connections       int
	OutboundBandwidth int
	InboundBandwidth  int
}

func (g *GStats) getSerializedCurrents() ([]byte, error) {
	g.currentMu.RLock()
	serializable := &serializableCurrent{
		UniqueIPs:         len(g.currentUniqueIPs),
		Connections:       g.currentConnections,
		OutboundBandwidth: g.currentOutboundBandwidth / realtimeRefreshRateInt,
		InboundBandwidth:  g.currentInboundBandwidth / realtimeRefreshRateInt,
	}
	g.currentMu.RUnlock()
	return json.Marshal(serializable)
}

func getIP(addr net.Addr) net.IP {
	switch addr := addr.(type) {
	case *net.TCPAddr:
		return addr.IP
	case *net.UDPAddr:
		return addr.IP
	case *net.IPAddr:
		return addr.IP
	}
	return nil
}

func ipToKey(ip net.IP) (key ipKey) {
	copy(key[:], ip.To16())
	return
}
