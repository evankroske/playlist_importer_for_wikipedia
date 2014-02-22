package plcreator

import (
	"encoding/json"
    "fmt"
	"io"
    "net/http"

	"appengine"
	"appengine/urlfetch"

	"plcreator/unwrap"
)

func init() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/refreshgenrelist", refreshGenreListHandler)
	http.HandleFunc("/kickoffgenrediscovery", kickoffGenreDiscoveryHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprint(w, r.URL.Path)
	} else {
		http.NotFound(w, r)
	}
}

func refreshGenreListHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	resp, err := client.Get("http://en.wikipedia.org/w/api.php?action=query&list=categorymembers&format=json&cmtitle=Category%3AMusic_genres&cmlimit=100")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	dec := json.NewDecoder(resp.Body)
	var v interface{}
	dec.Decode(&v)
	untypedTitles, err := unwrap.Unwrap(v, ".query.categorymembers[:].title")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	titles, ok := untypedTitles.([]interface{})
    if !ok {
        http.Error(w, "Guessed the wrong type", http.StatusInternalServerError)
        return
    }
	w.Write([]byte(fmt.Sprintf("%v", titles)))
}

func kickoffGenreDiscoveryHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Get ready!")
}
