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
)

const (
	discoverPlaylistsPath = "/admin/discoverplaylists"
	indexPlaylistsPath = "/admin/indexplaylists"
)

func init() {
    http.HandleFunc("/", handler)
    http.HandleFunc(discoverPlaylistsPath, refreshGenreListHandler)
	http.HandleFunc(
		"/admin/startplaylistdiscovery",
		kickoffGenreDiscoveryHandler,
	)
	http.HandleFunc("/admin/login", loginHandler)
	http.HandleFunc("/admin/logout", logoutHandler)
	http.HandleFunc(
		indexPlaylistsPath,
		muxByMethod(serveIndexPlaylistsForm, indexPlaylistsHandler),
	)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprint(w, r.URL.Path)
	} else {
		http.NotFound(w, r)
	}
}

func muxByMethod(
	getHandler http.HandlerFunc,
	postHandler http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getHandler(w, r)
		case "POST":
			postHandler(w, r)
		default:
			http.Error(w, "Badness", http.StatusInternalServerError)
		}
	}
}
