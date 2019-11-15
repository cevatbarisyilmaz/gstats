package gstats

import (
	"encoding/json"
	"math"
	"net"
	"net/http"
	"sort"
	"time"
)

type data struct {
	requests          int
	uniqueIPs         map[ipKey]bool
	inboundBandwidth  int
	outboundBandwidth int
	requestPaths      map[string]*requestPathRecord
	countries         map[string]*commonRecord
	browsers          map[string]*commonRecord
	oss               map[string]*commonRecord
	referrers         map[string]*commonRecord
}

func newData() *data {
	return &data{
		requests:          0,
		uniqueIPs:         map[ipKey]bool{},
		inboundBandwidth:  0,
		outboundBandwidth: 0,
		requestPaths:      map[string]*requestPathRecord{},
		countries:         map[string]*commonRecord{},
		browsers:          map[string]*commonRecord{},
		oss:               map[string]*commonRecord{},
		referrers:         map[string]*commonRecord{},
	}
}

func (d *data) merge(o *data) {
	d.requests += o.requests
	for ip := range o.uniqueIPs {
		d.uniqueIPs[ip] = true
	}
	d.inboundBandwidth += o.inboundBandwidth
	d.outboundBandwidth += o.outboundBandwidth
	for path, record := range o.requestPaths {
		currentRecord := d.requestPaths[path]
		if currentRecord == nil {
			currentRecord = newRequestPathRecord()
			d.requestPaths[path] = currentRecord
		}
		currentRecord.merge(record)
	}
	for identifier, record := range o.countries {
		currentRecord := d.countries[identifier]
		if currentRecord == nil {
			currentRecord = newCommonRecord()
			d.countries[identifier] = currentRecord
		}
		currentRecord.merge(record)
	}
	for identifier, record := range o.browsers {
		currentRecord := d.browsers[identifier]
		if currentRecord == nil {
			currentRecord = newCommonRecord()
			d.browsers[identifier] = currentRecord
		}
		currentRecord.merge(record)
	}
	for identifier, record := range o.oss {
		currentRecord := d.oss[identifier]
		if currentRecord == nil {
			currentRecord = newCommonRecord()
			d.oss[identifier] = currentRecord
		}
		currentRecord.merge(record)
	}
	for identifier, record := range o.referrers {
		currentRecord := d.referrers[identifier]
		if currentRecord == nil {
			currentRecord = newCommonRecord()
			d.referrers[identifier] = currentRecord
		}
		currentRecord.merge(record)
	}
}

func (d *data) toSerializableData() *serializableData {
	var uniqueIPs []string
	for ip := range d.uniqueIPs {
		uniqueIPs = append(uniqueIPs, net.IP(ip[:]).String())
	}
	var requestPaths []*serializableRequestPathRecord
	for path, record := range d.requestPaths {
		serializedRecord := record.toSerializableRequestPathRecord()
		serializedRecord.Path = path
		requestPaths = append(requestPaths, serializedRecord)
	}
	var countries []*serializableCommonRecord
	for country, record := range d.countries {
		countryRecord := record.toSerializableCommonRecord()
		countryRecord.Identifier = country
		countries = append(countries, countryRecord)
	}
	var browsers []*serializableCommonRecord
	for browser, record := range d.browsers {
		browserRecord := record.toSerializableCommonRecord()
		browserRecord.Identifier = browser
		browsers = append(browsers, browserRecord)
	}
	var oss []*serializableCommonRecord
	for os, record := range d.oss {
		osRecord := record.toSerializableCommonRecord()
		osRecord.Identifier = os
		oss = append(oss, osRecord)
	}
	var referrers []*serializableCommonRecord
	for referrer, record := range d.referrers {
		referrerRecord := record.toSerializableCommonRecord()
		referrerRecord.Identifier = referrer
		referrers = append(referrers, referrerRecord)
	}
	return &serializableData{
		Requests:          d.requests,
		UniqueIPs:         uniqueIPs,
		InboundBandwidth:  d.inboundBandwidth,
		OutboundBandwidth: d.outboundBandwidth,
		RequestPaths:      requestPaths,
		Countries:         countries,
		Browsers:          browsers,
		OSs:               oss,
		Referrers:         referrers,
	}
}

func (d *data) toHighlightedData() *highlightedData {
	var requestPaths []*highlightedRequestPathRecord
	for path, record := range d.requestPaths {
		requestPath := record.toHighlightedRequestPathRecord()
		requestPath.Path = path
		requestPaths = append(requestPaths, requestPath)
	}
	sort.Slice(requestPaths, func(i, j int) bool {
		return math.Sqrt(float64(requestPaths[i].Requests))+float64(requestPaths[i].UniqueIPs) > math.Sqrt(float64(requestPaths[j].Requests))+float64(requestPaths[j].UniqueIPs)
	})
	if len(requestPaths) > 64 {
		requestPaths = requestPaths[:64]
	}
	var countries []*highlightedCommonRecord
	for identifier, record := range d.countries {
		commonRecord := record.toHighlightedCommonRecord()
		commonRecord.Identifier = identifier
		countries = append(countries, commonRecord)
	}
	sort.Slice(countries, func(i, j int) bool {
		return countries[i].UniqueIPs > countries[j].UniqueIPs
	})
	if len(countries) > 32 {
		countries = countries[:32]
	}
	var browsers []*highlightedCommonRecord
	for identifier, record := range d.browsers {
		commonRecord := record.toHighlightedCommonRecord()
		commonRecord.Identifier = identifier
		browsers = append(browsers, commonRecord)
	}
	sort.Slice(browsers, func(i, j int) bool {
		return browsers[i].UniqueIPs > browsers[j].UniqueIPs
	})
	if len(browsers) > 16 {
		browsers = browsers[:16]
	}
	var oss []*highlightedCommonRecord
	for identifier, record := range d.oss {
		commonRecord := record.toHighlightedCommonRecord()
		commonRecord.Identifier = identifier
		oss = append(oss, commonRecord)
	}
	sort.Slice(oss, func(i, j int) bool {
		return oss[i].UniqueIPs > oss[j].UniqueIPs
	})
	if len(oss) > 16 {
		oss = oss[:16]
	}
	var referrers []*highlightedCommonRecord
	for identifier, record := range d.referrers {
		commonRecord := record.toHighlightedCommonRecord()
		commonRecord.Identifier = identifier
		referrers = append(referrers, commonRecord)
	}
	sort.Slice(referrers, func(i, j int) bool {
		return referrers[i].UniqueIPs > referrers[j].UniqueIPs
	})
	if len(referrers) > 64 {
		referrers = referrers[:64]
	}
	return &highlightedData{
		Requests:          d.requests,
		UniqueIPs:         len(d.uniqueIPs),
		InboundBandwidth:  d.inboundBandwidth,
		OutboundBandwidth: d.outboundBandwidth,
		RequestPaths:      requestPaths,
		Countries:         countries,
		Browsers:          browsers,
		OSs:               oss,
		Referrers:         referrers,
	}
}

type serializableData struct {
	Requests          int
	UniqueIPs         []string
	InboundBandwidth  int
	OutboundBandwidth int
	RequestPaths      []*serializableRequestPathRecord
	Countries         []*serializableCommonRecord
	Browsers          []*serializableCommonRecord
	OSs               []*serializableCommonRecord
	Referrers         []*serializableCommonRecord
}

func constructData(data []byte) (*serializableData, error) {
	d := &serializableData{}
	err := json.Unmarshal(data, d)
	return d, err
}

func constructHighlights(data []byte) (*highlightedData, error) {
	d := &highlightedData{}
	err := json.Unmarshal(data, d)
	return d, err
}

func (d *serializableData) toData() *data {
	uniqueIPs := map[ipKey]bool{}
	for _, ip := range d.UniqueIPs {
		uniqueIPs[ipToKey(net.ParseIP(ip))] = true
	}
	requestPaths := map[string]*requestPathRecord{}
	for _, record := range d.RequestPaths {
		requestPaths[record.Path] = record.toRequestPathRecord()
	}
	countries := map[string]*commonRecord{}
	for _, record := range d.Countries {
		countries[record.Identifier] = record.toCommonRecord()
	}
	browsers := map[string]*commonRecord{}
	for _, record := range d.Browsers {
		browsers[record.Identifier] = record.toCommonRecord()
	}
	oss := map[string]*commonRecord{}
	for _, record := range d.OSs {
		oss[record.Identifier] = record.toCommonRecord()
	}
	referrers := map[string]*commonRecord{}
	for _, record := range d.Referrers {
		referrers[record.Identifier] = record.toCommonRecord()
	}
	return &data{
		requests:          d.Requests,
		uniqueIPs:         uniqueIPs,
		inboundBandwidth:  d.InboundBandwidth,
		outboundBandwidth: d.OutboundBandwidth,
		requestPaths:      requestPaths,
		countries:         countries,
		browsers:          browsers,
		oss:               oss,
		referrers:         referrers,
	}
}

func (d *serializableData) serialize() ([]byte, error) {
	return json.Marshal(d)
}

type highlightedData struct {
	Requests          int
	UniqueIPs         int
	InboundBandwidth  int
	OutboundBandwidth int
	RequestPaths      []*highlightedRequestPathRecord
	Countries         []*highlightedCommonRecord
	Browsers          []*highlightedCommonRecord
	OSs               []*highlightedCommonRecord
	Referrers         []*highlightedCommonRecord
}

func (d *highlightedData) serialize() ([]byte, error) {
	return json.Marshal(d)
}

type requestPathRecord struct {
	statusCodes   map[int]int
	responseTimes map[time.Duration]int
	responseSizes map[int]int
	uniqueIPs     map[ipKey]bool
}

func newRequestPathRecord() *requestPathRecord {
	return &requestPathRecord{
		statusCodes:   map[int]int{},
		responseTimes: map[time.Duration]int{},
		responseSizes: map[int]int{},
		uniqueIPs:     map[ipKey]bool{},
	}
}

func (r *requestPathRecord) merge(o *requestPathRecord) {
	for statusCode, amount := range o.statusCodes {
		r.statusCodes[statusCode] += amount
	}
	for responseTime, amount := range o.responseTimes {
		r.responseTimes[responseTime] += amount
	}
	for responseSize, amount := range o.responseSizes {
		r.responseSizes[responseSize] += amount
	}
	for ip := range o.uniqueIPs {
		r.uniqueIPs[ip] = true
	}
}

func (r *requestPathRecord) toSerializableRequestPathRecord() *serializableRequestPathRecord {
	var uniqueIPs []string
	for ip := range r.uniqueIPs {
		uniqueIPs = append(uniqueIPs, net.IP(ip[:]).String())
	}
	return &serializableRequestPathRecord{
		StatusCodes:   r.statusCodes,
		ResponseTimes: r.responseTimes,
		ResponseSizes: r.responseSizes,
		UniqueIPs:     uniqueIPs,
	}
}

func (r *requestPathRecord) toHighlightedRequestPathRecord() *highlightedRequestPathRecord {
	var requests int
	for _, count := range r.statusCodes {
		requests += count
	}
	var successStatusCodeAmount int
	var totalAmount int
	var highestNonSuccessStatusCode int
	var highestNonSuccessStatusCodeAmount int
	for statusCode, amount := range r.statusCodes {
		totalAmount += amount
		if (statusCode >= 200 && statusCode < 300) || statusCode == http.StatusNotModified {
			successStatusCodeAmount += amount
		} else if amount > highestNonSuccessStatusCodeAmount {
			highestNonSuccessStatusCode = statusCode
			highestNonSuccessStatusCodeAmount = amount
		}
	}
	successfulStatusCodeRate := float64(successStatusCodeAmount) / float64(totalAmount)
	topNonSuccessfulStatusCodeRate := float64(highestNonSuccessStatusCodeAmount) / float64(totalAmount)
	var averageResponseTime float64
	t := float64(1)
	for duration, amount := range r.responseTimes {
		for i := 0; i < amount; i++ {
			averageResponseTime += (float64(duration) - averageResponseTime) / t
		}
	}
	var averageResponseSize float64
	t = 1
	for size, amount := range r.responseSizes {
		for i := 0; i < amount; i++ {
			averageResponseSize += (float64(size) - averageResponseSize) / t
		}
	}
	return &highlightedRequestPathRecord{
		Requests:                       requests,
		SuccessfulStatusCodeRate:       successfulStatusCodeRate,
		TopNonSuccessfulStatusCode:     highestNonSuccessStatusCode,
		TopNonSuccessfulStatusCodeRate: topNonSuccessfulStatusCodeRate,
		AverageResponseTime:            time.Duration(averageResponseTime),
		AverageResponseSize:            int(averageResponseSize),
		UniqueIPs:                      len(r.uniqueIPs),
	}
}

type serializableRequestPathRecord struct {
	Path          string
	StatusCodes   map[int]int
	ResponseTimes map[time.Duration]int
	ResponseSizes map[int]int
	UniqueIPs     []string
}

func (r *serializableRequestPathRecord) toRequestPathRecord() *requestPathRecord {
	uniqueIPs := map[ipKey]bool{}
	for _, ip := range r.UniqueIPs {
		uniqueIPs[ipToKey(net.ParseIP(ip))] = true
	}
	return &requestPathRecord{
		statusCodes:   r.StatusCodes,
		responseTimes: r.ResponseTimes,
		responseSizes: r.ResponseSizes,
		uniqueIPs:     uniqueIPs,
	}
}

type highlightedRequestPathRecord struct {
	Path                           string
	Requests                       int
	SuccessfulStatusCodeRate       float64
	TopNonSuccessfulStatusCode     int
	TopNonSuccessfulStatusCodeRate float64
	AverageResponseTime            time.Duration
	AverageResponseSize            int
	UniqueIPs                      int
}

type commonRecord struct {
	requests  int
	uniqueIPs map[ipKey]bool
}

func newCommonRecord() *commonRecord {
	return &commonRecord{
		requests:  0,
		uniqueIPs: map[ipKey]bool{},
	}
}

func (r *commonRecord) merge(o *commonRecord) {
	r.requests += o.requests
	for ip := range o.uniqueIPs {
		r.uniqueIPs[ip] = true
	}
}

func (r *commonRecord) toSerializableCommonRecord() *serializableCommonRecord {
	var uniqueIPs []string
	for ip := range r.uniqueIPs {
		uniqueIPs = append(uniqueIPs, net.IP(ip[:]).String())
	}
	return &serializableCommonRecord{
		Requests:  r.requests,
		UniqueIPs: uniqueIPs,
	}
}

func (r *commonRecord) toHighlightedCommonRecord() *highlightedCommonRecord {
	return &highlightedCommonRecord{
		Requests:  r.requests,
		UniqueIPs: len(r.uniqueIPs),
	}
}

type serializableCommonRecord struct {
	Identifier string
	Requests   int
	UniqueIPs  []string
}

func (r *serializableCommonRecord) toCommonRecord() *commonRecord {
	uniqueIPs := map[ipKey]bool{}
	for _, ip := range r.UniqueIPs {
		uniqueIPs[ipToKey(net.ParseIP(ip))] = true
	}
	return &commonRecord{
		requests:  r.Requests,
		uniqueIPs: uniqueIPs,
	}
}

type highlightedCommonRecord struct {
	Identifier string
	Requests   int
	UniqueIPs  int
}
