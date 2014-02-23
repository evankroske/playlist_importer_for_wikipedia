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
	"net/url"
	"strings"

	"appengine"
	"appengine/urlfetch"

	"playlistimporter/unwrap"
)

const titlesFormKey = "childTitles"
const categoryMembersLimit = "500"
func refreshGenreListHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	r.ParseForm()
	parentTitles := r.Form[titlesFormKey]
	// "s" for "slice"
	s := func(v string) []string {
		return []string{v}
	}
	reqData := url.Values{
		"action": s("query"),
		"list": s("categorymembers"),
		"format": s("json"),
		"cmlimit": s(categoryMembersLimit),
		"cmtitle": s(strings.Join(parentTitles, "|")),
	}
	reqUrl := &url.URL{
		Scheme: "http",
		Host: "en.wikipedia.org",
		Path: "/w/api.php",
		RawQuery: reqData.Encode(),
	}
	req, err := http.NewRequest(
		"GET",
		reqUrl.String(),
		nil,
	)
	if err != nil {
		c.Criticalf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// If you want to run your own version of this app, change the user agent
	// to use your app's information.
	userAgent := fmt.Sprintf(
		"Playlist Importer for Wikipedia/%s (%s; evan@evankroske.com)",
		appengine.VersionID(c),
		 "http://playlistimporterforwikipedia.appspot.com",
	)
	req.Header.Add("User-Agent", userAgent)
	client := urlfetch.Client(c)
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
	childTitles, ok := untypedTitles.([]interface{})
    if !ok {
		c.Errorf("%v: %v", "Got wrong type from unwrap", resp.Body)
        http.Error(w, "Sorry.", http.StatusInternalServerError)
        return
    }
	w.Write([]byte(fmt.Sprintf("%#v", childTitles)))
}