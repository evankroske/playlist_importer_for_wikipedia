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

import "fmt"

func makeUserAgentString(version string) string {
	// If you want to run your own version of this app, change the user agent
	// to use your app's information.
	return fmt.Sprintf(
		"Playlist Importer for Wikipedia/%s (%s; evan@evankroske.com)",
		version,
		"http://playlistimporterforwikipedia.appspot.com",
	)
}
