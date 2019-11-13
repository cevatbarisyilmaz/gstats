package files

type File struct {
	Source []byte
	Type   string
}

var Paths = map[string]*File{
	"/": {
		Source: indexHTML,
		Type:   "text/html",
	},
	"/script.js": {
		Source: scriptJS,
		Type:   "application/javascript",
	},
	"/style.css": {
		Source: styleCSS,
		Type:   "text/css",
	},
}
