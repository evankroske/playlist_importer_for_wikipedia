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
	"strings"

	"appengine"
	"appengine/urlfetch"

	"playlistimporter/unwrap"
)

const categoryMembersLimit = "500"

func refreshGenreListHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	r.ParseForm()
	parentTitles := r.Form[categoryTitleFormKey]
	if len(parentTitles) != 1 {
		c.Criticalf("len(parentTitles) == %v", len(parentTitles))
		http.Error(w, "Whoops.", http.StatusInternalServerError)
		return
	}
	reqData := url.Values{
		"action": []string{"query"},
		"list": []string{"categorymembers"},
		"format": []string{"json"},
		"cmlimit": []string{categoryMembersLimit},
		"cmtitle": parentTitles,
	}
	client := urlfetch.Client(c)
	userAgentStr := makeUserAgentString(appengine.VersionID(c))
	jsonRsp, err := queryWikipediaAPI(client, userAgentStr, reqData)
	untypedTitles, err := unwrap.Unwrap(jsonRsp, ".query.categorymembers[:].title")
    if err != nil {
		c.Errorf("%v", err.Error())
        http.Error(w, "My bad.", http.StatusInternalServerError)
        return
    }
	childTitles, ok := untypedTitles.([]interface{})
    if !ok {
		c.Errorf("%v", "Got wrong type from unwrap")
        http.Error(w, "Sorry.", http.StatusInternalServerError)
        return
    }
	var playlistTitles, subCategoryTitles []string
	for _, untypedTitle := range childTitles {
		title, ok := untypedTitle.(string)
		if !ok {
			c.Criticalf(`Expected string found %t: %v`, untypedTitle)
			http.Error(w, "I messed up.", http.StatusInternalServerError)
			return
		}
		if strings.HasPrefix(title, "Category:") {
			subCategoryTitles = append(subCategoryTitles, title)
		} else {
			playlistTitles = append(playlistTitles, title)
		}
	}
	c.Infof(
		"Found %d subcategories and %d playlists.",
		len(subCategoryTitles),
		len(playlistTitles),
	)
	fmt.Println(w, "Success!")
}
