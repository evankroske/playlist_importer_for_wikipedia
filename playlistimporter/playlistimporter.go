/*
Copyright 2014 Google Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package playlistimporter

import (
	"encoding/json"
    "fmt"
    "net/http"

	"appengine"
	"appengine/urlfetch"

	"playlistimporter/unwrap"
)

func init() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/admin/discoverplaylists", refreshGenreListHandler)
	http.HandleFunc(
		"/admin/startplaylistdiscovery",
		kickoffGenreDiscoveryHandler,
	)
	http.HandleFunc("/admin/login", loginHandler)
	http.HandleFunc("/admin/logout", logoutHandler)
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
	req, err := http.NewRequest(
		"GET",
		"http://en.wikipedia.org/w/api.php?action=query&list=categorymembers&format=json&cmtitle=Category%3AMusic_genres&cmlimit=100",
		nil,
	)
	if err != nil {
		c.Criticalf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	const homePage = "http://playlistimporterforwikipedia.appspot.com"
	userAgent := fmt.Sprintf(
		"Playlist Importer for Wikipedia/%s (%s; evan@evankroske.com)",
		appengine.VersionID(c),
		homePage,
	)
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dec := json.NewDecoder(resp.Body)
	var v interface{}
	dec.Decode(&v)
	untypedTitles, err := unwrap.Unwrap(v, ".query.categorymembers[:].title")
    if err != nil {
		c.Errorf("%v: %v", err.Error(), resp.Body)
        http.Error(w, "My bad.", http.StatusInternalServerError)
        return
    }
	titles, ok := untypedTitles.([]interface{})
    if !ok {
		c.Errorf("%v: %v", "Got wrong type from unwrap", resp.Body)
        http.Error(w, "Sorry.", http.StatusInternalServerError)
        return
    }
	w.Write([]byte(fmt.Sprintf("%#v", titles)))
}

