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

	"appengine"
	"appengine/taskqueue"
)

func kickoffGenreDiscoveryHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	t := taskqueue.NewPOSTTask(
		discoverPlaylistsPath,
		map[string][]string{
			categoryTitleFormKey: []string{"Category:Music_Genres"},
		},
	)
	if _, err := taskqueue.Add(c, t, "playlistSources"); err != nil {
		c.Criticalf(err.Error())
		http.Error(
			w,
			"Error adding task to queue",
			http.StatusInternalServerError,
		)
		return
	}
	fmt.Fprintln(w, "Task added successfully")
}
