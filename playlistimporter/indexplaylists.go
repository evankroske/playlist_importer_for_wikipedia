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
	"fmt"
	"net/http"
	"net/url"

	"appengine"
	"appengine/urlfetch"

	"playlistimporter/unwrap"
)

const pllimit = "500"

func indexPlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	playlistTitles := r.Form["titles"]
	args := url.Values{
		"action": []string{"query"},
		"format": []string{"json"},
		"namespace": []string{"0"},
		"pllimit": []string{pllimit},
		"prop": []string{"links"},
		"titles": playlistTitles,
	}
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	userAgentString := makeUserAgentString(appengine.VersionID(c))
	jsonRsp, err := queryWikipediaAPI(client, userAgentString, args)
	if err != nil {
		http.Error(
			w,
			"I can't believe this is happening.",
			http.StatusInternalServerError,
		)
		return
	}
	pageIDToPageUntyped, err := unwrap.Unwrap(jsonRsp, ".query.pages")
	if err != nil {
		http.Error(w, "Badness", http.StatusInternalServerError)
		return
	}
	pageIDToPage, ok := pageIDToPageUntyped.(map[string]interface{})
	if !ok {
		http.Error(w, "Curses", http.StatusInternalServerError)
		return
	}
	for _, page := range pageIDToPage {
		titleUntyped, _ := unwrap.Unwrap(page, ".title")
		title, _ := titleUntyped.(string)
		linkedTitlesUntyped, err := unwrap.Unwrap(page, ".links[:].title")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		linkedTitles, ok := linkedTitlesUntyped.([]interface{})
		if !ok {
			http.Error(
				w,
				"linkedTitles type conversion failed",
				http.StatusInternalServerError,
			)
		}
		for _, linkedTitleUntyped := range(linkedTitles) {
			linkedTitle, _ := linkedTitleUntyped.(string)
			fmt.Fprintln(w, url.Values{
				"parent": []string{title},
				"linkedTitles": []string{linkedTitle},
			})
		}
		c.Infof("%v links found on page \"%v\"", len(linkedTitles), title)
	}
}
