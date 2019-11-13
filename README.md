# gstats

[![GoDoc](https://godoc.org/github.com/cevatbarisyilmaz/gstats?status.svg)](https://godoc.org/github.com/cevatbarisyilmaz/gstats)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/cevatbarisyilmaz/gstats?sort=semver)](https://github.com/cevatbarisyilmaz/gstats/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/cevatbarisyilmaz/gstats)](https://goreportcard.com/report/github.com/cevatbarisyilmaz/gstats)

Realtime HTTP Stats for Go.

## Example

```go
func main() {
	gs := gstats.New("gstats")
	listener, err := net.Listen("tcp", "127.0.0.1:80")
	if err != nil {
		log.Fatal(err)
	}
	glistener := gs.Listener(listener)
	http.Handle("/gstats/", gs.Collect(http.StripPrefix("/gstats", gs.Show())))
	http.Handle("/", gs.Collect(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("You are visiting " + request.URL.Path))
	})))
	err = http.Serve(glistener, nil)
	gs.PrepareToExit()
	log.Fatal(err)
}
```

Then go to http://localhost/gstats

![alt text](https://raw.githubusercontent.com/cevatbarisyilmaz/gstats/master/image.png)

## Features

* Geolocation
* Unique Visitors
* Most Visited Pages
* Average Response Sizes
* Average Response Times
* Referrers
* Visitor OS Stats
* Visitor Browser Stats
* Bandwidth

## License

Creative Commons Attribution-ShareAlike 4.0 International with default geolocation source. MIT otherwise.
See https://github.com/cevatbarisyilmaz/gstats/blob/master/LICENSE.md for details.

## Production Usage

gstats are still in early development stage. Be catious for production usages.

## Attribution

This product includes GeoLite2 data created by MaxMind, available from https://www.maxmind.com
