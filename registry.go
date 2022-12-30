/*
 * Copyright (c) 2022. Wolfgang Popp
 *
 * This file is part of tagwatch.
 *
 * tagwatch is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * tagwatch is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with tagwatch.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

type TagsResponse struct {
	Name string
	Tags []string
}

type AuthResponse struct {
	Token string
}

type Manifest struct {
	Digest   string
	Platform struct {
		Architecture string
		OS           string
		Variant      string
	}
	Size int
}

type ManifestResponse struct {
	SchemaVersion int
	MediaType     string
	Manifests     []Manifest
	contentDigest string
}

type TagDigest struct {
	Tag    string
	Digest string
}

type RegistryClient struct {
	client        http.Client
	manifestCache map[string]*ManifestResponse
	authToken     string
	Auth          bool
	AuthUsername  string
	AuthPassword  string
	AuthURL       string
	BaseURL       string
}

func NewRegistryClientFromConf(reg *Registry) *RegistryClient {
	return &RegistryClient{
		client:        http.Client{Timeout: 60 * time.Second},
		manifestCache: make(map[string]*ManifestResponse, 20),
		Auth:          reg.Auth,
		AuthUsername:  reg.AuthUsername,
		AuthPassword:  reg.AuthPassword,
		AuthURL:       reg.AuthURL,
		BaseURL:       reg.BaseURL,
	}
}

func (c *RegistryClient) login(repo string) {
	if c.Auth {
		var authResponse AuthResponse
		c.authToken = ""
		err := c.get(c.AuthURL+"&scope=repository:"+repo+":pull", &authResponse, true)
		if err != nil {
			log.Println(err)
		}
		c.authToken = authResponse.Token
	}
}

func (c *RegistryClient) request(method, url string, basicAuth bool) (*http.Response, error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", AgentStr)
	request.Header.Set("Accept", "application/vnd.docker.distribution.manifest.list.v2+json")

	if basicAuth && c.AuthUsername != "" && c.AuthPassword != "" {
		request.SetBasicAuth(c.AuthUsername, c.AuthPassword)
	}

	if c.Auth && c.authToken != "" {
		request.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("request '%s %s'failed with status code %d", method, url, response.StatusCode)
	}
	return response, nil
}

func (c *RegistryClient) get(url string, result interface{}, basicAuth bool) error {
	response, err := c.request("GET", url, basicAuth)
	if err != nil {
		return err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(responseBody, result); err != nil {
		return err
	}
	return nil
}

func (c *RegistryClient) head(url string, basicAuth bool) (http.Header, error) {
	response, err := c.request("HEAD", url, basicAuth)
	if err != nil {
		return nil, err
	}

	return response.Header, nil
}

func matchesAny(tagPattern []string, str string) bool {
	for _, pattern := range tagPattern {
		matched, err := regexp.MatchString(pattern, str)
		if err != nil {
			log.Println(err)
		}
		if matched {
			return true
		}
	}
	return false
}

func (c *RegistryClient) fetchManifest(repo, tag string) (*ManifestResponse, error) {
	url := c.BaseURL + repo + "/manifests/" + tag

	getManifest := func() (*ManifestResponse, error) {
		var manifestResponse ManifestResponse
		if err := c.get(url, &manifestResponse, false); err != nil {
			return nil, err
		}
		return &manifestResponse, nil
	}

	manifestHeaders, err := c.head(url, false)
	if err != nil {
		// registry does seem to support HEAD /<repo>/manifests/<tag>
		return getManifest()
	}

	digest := manifestHeaders.Get("Docker-Content-Digest")
	if digest == "" {
		// Docker-Content-Digest header missing in HEAD /<repo>/manifests/<tag>
		return getManifest()
	}

	manifest := c.manifestCache[url]
	if manifest == nil || manifest.contentDigest != digest {
		// manifest cache empty or expired. Get and save response
		manifest, err = getManifest()
		if err != nil {
			return nil, err
		}
		manifest.contentDigest = digest
		c.manifestCache[url] = manifest
	}

	// cache hit. return the hit
	return manifest, nil
}

func (c *RegistryClient) ListTags(repo string) []string {
	var tagsResponse TagsResponse
	c.login(repo)
	url := c.BaseURL + repo + "/tags/list"
	if err := c.get(url, &tagsResponse, false); err != nil {
		log.Println(err)
		return nil
	}
	return tagsResponse.Tags
}

func (c *RegistryClient) ListTagDigests(repo, architecture string, allTags, tagPattern []string) []TagDigest {
	res := make([]TagDigest, 0, 100)
	c.login(repo)
	for _, tag := range allTags {
		if !matchesAny(tagPattern, tag) {
			continue
		}

		manifestResponse, err := c.fetchManifest(repo, tag)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, manifest := range manifestResponse.Manifests {
			if manifest.Platform.Architecture == architecture {
				res = append(res, TagDigest{tag, manifest.Digest})
			}
		}
	}

	return res
}
