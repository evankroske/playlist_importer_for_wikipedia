package plcreator

import (
	"encoding/json"
    "fmt"
    "net/http"

	"appengine"
	"appengine/urlfetch"
)

func init() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/refreshgenrelist", refreshGenreListHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
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
	w.Write([]byte("Writing parser"))
}
