package gstats

import (
	"encoding/json"
	"github.com/cevatbarisyilmaz/gstats/internal/files"
	"net/http"
	"time"
)

type responseData struct {
	Today     *unitData
	Yesterday *unitData
	ThisMonth *unitData
	LastMonth *unitData
	ThisYear  *unitData
	LastYear  *unitData
}

type unitData struct {
	Highlights    *highlightedData
	SubHighlights []*highlightedData
}

func getAPIHandler(g *GStats) http.Handler {
	apiHandler := http.NewServeMux()
	apiHandler.HandleFunc("/current", func(writer http.ResponseWriter, request *http.Request) {
		buffer, err := g.getSerializedCurrents()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			_, _ = writer.Write(buffer)
		}
	})
	apiHandler.HandleFunc("/data", func(writer http.ResponseWriter, request *http.Request) {
		response := &responseData{
			Today: &unitData{
				Highlights:    g.getTodaysHighlights(),
				SubHighlights: g.getTodaysHourlyHighlights(),
			},
			Yesterday: &unitData{
				Highlights:    g.getYesterdaysHighlights(),
				SubHighlights: g.getYesterdaysHourlyHighlights(),
			},
			ThisMonth: &unitData{
				Highlights:    g.getThisMonthsHighlights(),
				SubHighlights: g.getThisMonthsDailyHighlights(),
			},
			LastMonth: &unitData{
				Highlights:    g.getLastMonthsHighlights(),
				SubHighlights: g.getLastMonthsDailyHighlights(),
			},
			ThisYear: &unitData{
				Highlights:    g.getThisYearsHighlights(),
				SubHighlights: g.getThisYearsMonthlyHighlights(),
			},
			LastYear: &unitData{
				Highlights:    g.getLastYearsHighlights(),
				SubHighlights: g.getLastYearsMonthlyHighlights(),
			},
		}
		buffer, err := json.Marshal(response)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			_, _ = writer.Write(buffer)
		}
	})
	return apiHandler
}

func getShowHandler(g *GStats) http.Handler {
	show := http.NewServeMux()
	show.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		file := files.Paths[request.URL.Path]
		if file == nil {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			writer.Header().Add("content-type", file.Type)
			_, _ = writer.Write(file.Source)
		}
	})
	show.Handle("/api/", http.StripPrefix("/api", getAPIHandler(g)))
	return show
}

func (g *GStats) getTodaysHighlights() *highlightedData {
	record := g.collectDailyRecord(time.Now().UTC())
	g.dataMu.RLock()
	record.merge(g.data)
	g.dataMu.RUnlock()
	return record.toHighlightedData()
}

func (g *GStats) getTodaysHourlyHighlights() []*highlightedData {
	now := time.Now().UTC()
	var hourlyHighlights []*highlightedData
	for i := 0; i < now.Hour(); i++ {
		highlights, err := g.lookupHourlyRecord(time.Date(now.Year(), now.Month(), now.Day(), i, 0, 0, 0, time.UTC))
		if err != nil {
			hourlyHighlights = append(hourlyHighlights, newData().toHighlightedData())
		} else {
			hourlyHighlights = append(hourlyHighlights, highlights)
		}
	}
	g.dataMu.RLock()
	hourlyHighlights = append(hourlyHighlights, g.data.toHighlightedData())
	g.dataMu.RUnlock()
	return hourlyHighlights
}

func (g *GStats) getYesterdaysHighlights() *highlightedData {
	now := time.Now().UTC()
	yesterday := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, time.UTC).Add(-24 * time.Hour)
	d, err := g.lookupDailyRecord(yesterday)
	if err != nil {
		d = newData().toHighlightedData()
	}
	return d
}

func (g *GStats) getYesterdaysHourlyHighlights() []*highlightedData {
	now := time.Now().UTC()
	yesterday := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, time.UTC).Add(-24 * time.Hour)
	var hourlyHighlights []*highlightedData
	for i := 0; i < 24; i++ {
		highlights, err := g.lookupHourlyRecord(time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), i, 0, 0, 0, time.UTC))
		if err != nil {
			hourlyHighlights = append(hourlyHighlights, newData().toHighlightedData())
		} else {
			hourlyHighlights = append(hourlyHighlights, highlights)
		}
	}
	return hourlyHighlights
}

func (g *GStats) getThisMonthsHighlights() *highlightedData {
	record := g.collectMonthlyRecord(time.Now().UTC())
	g.dataMu.RLock()
	record.merge(g.data)
	g.dataMu.RUnlock()
	return record.toHighlightedData()
}

func (g *GStats) getThisMonthsDailyHighlights() []*highlightedData {
	now := time.Now().UTC()
	var dailyHighlights []*highlightedData
	for i := 1; i < now.Day(); i++ {
		highlights, err := g.lookupDailyRecord(time.Date(now.Year(), now.Month(), i, 0, 0, 0, 0, time.UTC))
		if err != nil {
			dailyHighlights = append(dailyHighlights, newData().toHighlightedData())
		} else {
			dailyHighlights = append(dailyHighlights, highlights)
		}
	}
	today := g.collectDailyRecord(now)
	g.dataMu.RLock()
	today.merge(g.data)
	g.dataMu.RUnlock()
	dailyHighlights = append(dailyHighlights, today.toHighlightedData())
	return dailyHighlights
}

func (g *GStats) getLastMonthsHighlights() *highlightedData {
	now := time.Now().UTC()
	lastMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Add(-24 * time.Hour)
	d, err := g.lookupMonthlyRecord(lastMonth)
	if err != nil {
		d = newData().toHighlightedData()
	}
	return d
}

func (g *GStats) getLastMonthsDailyHighlights() []*highlightedData {
	now := time.Now().UTC()
	lastMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Add(-24 * time.Hour)
	lastMonth = time.Date(lastMonth.Year(), lastMonth.Month(), 1, 1, 0, 0, 0, time.UTC)
	month := lastMonth.Month()
	var dailyHighlights []*highlightedData
	for lastMonth.Month() == month {
		highlights, err := g.lookupDailyRecord(lastMonth)
		if err != nil {
			dailyHighlights = append(dailyHighlights, newData().toHighlightedData())
		} else {
			dailyHighlights = append(dailyHighlights, highlights)
		}
		lastMonth = lastMonth.Add(time.Hour * 24)
	}
	return dailyHighlights

}

func (g *GStats) getThisYearsHighlights() *highlightedData {
	record := g.collectYearlyRecord(time.Now().UTC())
	g.dataMu.RLock()
	record.merge(g.data)
	g.dataMu.RUnlock()
	return record.toHighlightedData()
}

func (g *GStats) getThisYearsMonthlyHighlights() []*highlightedData {
	now := time.Now().UTC()
	var monthlyHighlights []*highlightedData
	for i := time.Month(1); i < now.Month(); i++ {
		highlights, err := g.lookupMonthlyRecord(time.Date(now.Year(), i, 1, 0, 0, 0, 0, time.UTC))
		if err != nil {
			monthlyHighlights = append(monthlyHighlights, newData().toHighlightedData())
		} else {
			monthlyHighlights = append(monthlyHighlights, highlights)
		}
	}
	thisMonth := g.collectMonthlyRecord(now)
	g.dataMu.RLock()
	thisMonth.merge(g.data)
	g.dataMu.RUnlock()
	monthlyHighlights = append(monthlyHighlights, thisMonth.toHighlightedData())
	return monthlyHighlights
}

func (g *GStats) getLastYearsHighlights() *highlightedData {
	now := time.Now().UTC()
	lastYear := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC)
	d, err := g.lookupYearlyRecord(lastYear)
	if err != nil {
		d = newData().toHighlightedData()
	}
	return d
}

func (g *GStats) getLastYearsMonthlyHighlights() []*highlightedData {
	now := time.Now().UTC()
	lastYear := time.Date(now.Year()-1, 1, 1, 1, 0, 0, 0, time.UTC)
	year := lastYear.Year()
	var monthlyHighlights []*highlightedData
	for lastYear.Year() == year {
		highlights, err := g.lookupMonthlyRecord(lastYear)
		if err != nil {
			monthlyHighlights = append(monthlyHighlights, newData().toHighlightedData())
		} else {
			monthlyHighlights = append(monthlyHighlights, highlights)
		}
		lastYear = lastYear.Add(31 * 24 * time.Hour)
	}
	return monthlyHighlights
}
