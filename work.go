package gstats

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

func (g *GStats) work() {
	for {
		utcNow := time.Now().UTC()
		nextHour := time.Date(utcNow.Year(), utcNow.Month(), utcNow.Day(), utcNow.Hour(), 0, 0, 0, time.UTC).Add(time.Hour)
		time.Sleep(nextHour.Sub(utcNow))
		g.record(utcNow, nextHour)
	}
}

func (g *GStats) record(o, n time.Time) {
	g.recordHour(o)
	if o.Day() != n.Day() {
		g.recordDay(o)
		if o.Month() != n.Month() {
			g.recordMonth(o)
			if o.Year() != n.Year() {
				g.recordYear(o)
			}
		}
	}
}

func (g *GStats) recordHour(t time.Time) {
	g.dataMu.Lock()
	data := g.data
	g.data = newData()
	g.dataMu.Unlock()
	basePath := g.root + strconv.Itoa(t.Year()) + "/" + strconv.Itoa(int(t.Month())) + "/" + strconv.Itoa(t.Day()) + "/" + strconv.Itoa(t.Hour()) + "/"
	err := os.MkdirAll(basePath, os.ModeDir)
	if err != nil {
		log.Println("gstats: error while creating the dir: " + err.Error())
	} else {
		saveData(basePath+"data.json", data)
		saveHighlights(basePath+"highlights.json", data)
	}
}

func (g *GStats) collectDailyRecord(t time.Time) *data {
	data := newData()
	basePath := g.root + strconv.Itoa(t.Year()) + "/" + strconv.Itoa(int(t.Month())) + "/" + strconv.Itoa(t.Day()) + "/"
	for i := 0; i < 24; i++ {
		path := basePath + strconv.Itoa(i) + "/data.json"
		dataContent, err := ioutil.ReadFile(path)
		if err == nil {
			serializedData, err := constructData(dataContent)
			if err == nil {
				data.merge(serializedData.toData())
			}
		}
	}
	return data
}

func (g *GStats) recordDay(t time.Time) {
	saveHighlights(g.root+strconv.Itoa(t.Year())+"/"+strconv.Itoa(int(t.Month()))+"/"+strconv.Itoa(t.Day())+"/highlights.json", g.collectDailyRecord(t))
}

func (g *GStats) collectMonthlyRecord(t time.Time) *data {
	data := newData()
	for i := 1; i < 32; i++ {
		data.merge(g.collectDailyRecord(time.Date(t.Year(), t.Month(), i, 0, 0, 0, 0, time.UTC)))
	}
	return data
}

func (g *GStats) recordMonth(t time.Time) {
	saveHighlights(g.root+strconv.Itoa(t.Year())+"/"+strconv.Itoa(int(t.Month()))+"/highlights.json", g.collectMonthlyRecord(t))
}

func (g *GStats) collectYearlyRecord(t time.Time) *data {
	data := newData()
	for i := 1; i < 13; i++ {
		data.merge(g.collectMonthlyRecord(time.Date(t.Year(), time.Month(i), 1, 0, 0, 0, 0, time.UTC)))
	}
	return data
}

func (g *GStats) recordYear(t time.Time) {
	saveHighlights(g.root+strconv.Itoa(t.Year())+"/highlights.json", g.collectYearlyRecord(t))
}

func (g *GStats) lookupHourlyRecord(t time.Time) (*highlightedData, error) {
	path := g.root + strconv.Itoa(t.Year()) + "/" + strconv.Itoa(int(t.Month())) + "/" + strconv.Itoa(t.Day()) + "/" + strconv.Itoa(t.Hour()) + "/highlights.json"
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return constructHighlights(contents)
}

func (g *GStats) lookupSerializedHourlyRecord(t time.Time) ([]byte, error) {
	path := g.root + strconv.Itoa(t.Year()) + "/" + strconv.Itoa(int(t.Month())) + "/" + strconv.Itoa(t.Day()) + "/" + strconv.Itoa(t.Hour()) + "/highlights.json"
	return ioutil.ReadFile(path)
}

func (g *GStats) lookupDailyRecord(t time.Time) (*highlightedData, error) {
	basePath := g.root + strconv.Itoa(t.Year()) + "/" + strconv.Itoa(int(t.Month())) + "/" + strconv.Itoa(t.Day())
	highlightsPath := basePath + "/highlights.json"
	contents, err := ioutil.ReadFile(highlightsPath)
	if err == nil {
		return constructHighlights(contents)
	}
	data := g.collectDailyRecord(t)
	err = os.MkdirAll(basePath, os.ModeDir)
	if err == nil {
		saveHighlights(highlightsPath, data)
	} else {
		log.Println("gstats: mkdir failed: " + err.Error())
	}
	return data.toHighlightedData(), nil
}

func (g *GStats) lookupSerializedDailyRecord(t time.Time) ([]byte, error) {
	basePath := g.root + strconv.Itoa(t.Year()) + "/" + strconv.Itoa(int(t.Month())) + "/" + strconv.Itoa(t.Day())
	highlightsPath := basePath + "/highlights.json"
	contents, err := ioutil.ReadFile(highlightsPath)
	if err == nil {
		return contents, nil
	}
	data := g.collectDailyRecord(t)
	contents, err = data.toHighlightedData().serialize()
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(basePath, os.ModeDir)
	if err == nil {
		saveHighlights(highlightsPath, data)
	} else {
		log.Println("gstats: mkdir failed: " + err.Error())
	}
	return contents, nil
}

func (g *GStats) lookupMonthlyRecord(t time.Time) (*highlightedData, error) {
	basePath := g.root + strconv.Itoa(t.Year()) + "/" + strconv.Itoa(int(t.Month()))
	highlightsPath := basePath + "/highlights.json"
	contents, err := ioutil.ReadFile(highlightsPath)
	if err == nil {
		return constructHighlights(contents)
	}
	data := g.collectMonthlyRecord(t)
	err = os.MkdirAll(basePath, os.ModeDir)
	if err == nil {
		saveHighlights(highlightsPath, data)
	} else {
		log.Println("gstats: mkdir failed: " + err.Error())
	}
	return data.toHighlightedData(), nil
}

func (g *GStats) lookupSerializedMonthlyRecord(t time.Time) ([]byte, error) {
	basePath := g.root + strconv.Itoa(t.Year()) + "/" + strconv.Itoa(int(t.Month()))
	highlightsPath := basePath + "/highlights.json"
	contents, err := ioutil.ReadFile(highlightsPath)
	if err == nil {
		return contents, nil
	}
	data := g.collectMonthlyRecord(t)
	contents, err = data.toHighlightedData().serialize()
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(basePath, os.ModeDir)
	if err == nil {
		saveHighlights(highlightsPath, data)
	} else {
		log.Println("gstats: mkdir failed: " + err.Error())
	}
	return contents, nil
}

func (g *GStats) lookupYearlyRecord(t time.Time) (*highlightedData, error) {
	basePath := g.root + strconv.Itoa(t.Year())
	highlightsPath := basePath + "/highlights.json"
	contents, err := ioutil.ReadFile(highlightsPath)
	if err == nil {
		return constructHighlights(contents)
	}
	data := g.collectYearlyRecord(t)
	err = os.MkdirAll(basePath, os.ModeDir)
	if err == nil {
		saveHighlights(highlightsPath, data)
	} else {
		log.Println("gstats: mkdir failed: " + err.Error())
	}
	return data.toHighlightedData(), nil
}

func (g *GStats) lookupSerializedYearlyRecord(t time.Time) ([]byte, error) {
	basePath := g.root + strconv.Itoa(t.Year())
	highlightsPath := basePath + "/highlights.json"
	contents, err := ioutil.ReadFile(highlightsPath)
	if err == nil {
		return contents, nil
	}
	data := g.collectYearlyRecord(t)
	contents, err = data.toHighlightedData().serialize()
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(basePath, os.ModeDir)
	if err == nil {
		saveHighlights(highlightsPath, data)
	} else {
		log.Println("gstats: mkdir failed: " + err.Error())
	}
	return contents, nil
}

func saveHighlights(path string, d *data) {
	serializedData, err := d.toHighlightedData().serialize()
	if err != nil {
		log.Println("gstats: error while serializing data: " + err.Error())
	} else {
		saveFile(path, serializedData)
	}
}

func saveData(path string, d *data) {
	serializedData, err := d.toSerializableData().serialize()
	if err != nil {
		log.Println("gstats: error while serializing data: " + err.Error())
	} else {
		saveFile(path, serializedData)
	}
}

func saveFile(path string, data []byte) {
	targetFile, err := os.Create(path)
	if err != nil {
		log.Println("gstats: error while creating data file: " + err.Error())
	} else {
		_, err = targetFile.Write(data)
		if err != nil {
			log.Println("gstats: error while writing data file: " + err.Error())
			err = targetFile.Close()
			if err != nil {
				log.Println("gstats: error while closing data file:" + err.Error())
			} else {
				err = os.Remove(path)
				if err != nil {
					log.Println("gstats: error while deleting corrupted data file:" + err.Error())
				}
			}
		}
	}
}
