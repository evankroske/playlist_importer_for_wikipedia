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
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func queryWikipediaAPI(
	client *http.Client,
	userAgent string,
	reqData url.Values,
) (jsonRsp interface{}, err error) {
	reqURL := makeWikipediaEndpoint()
	reqURL.RawQuery = reqData.Encode()
	req, err := http.NewRequest(
		"GET",
		reqURL.String(),
		nil,
	)
	if err != nil {
		return
	}
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	if !strings.HasPrefix(resp.Status, "2") {
		err = errors.New(fmt.Sprintf("HTTP error: %v", resp.Status))
		return
	}
	dec := json.NewDecoder(resp.Body)
	// Decode takes a pointer.
	err = dec.Decode(&jsonRsp)
	if err != nil {
		return
	}
	return
}
