package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
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
}

type TagDigest struct {
	Tag    string
	Digest string
}

type RegistryClient struct {
	client       http.Client
	Auth         bool
	AuthUsername string
	AuthPassword string
	AuthToken    string
	AuthURL      string
	BaseURL      string
}

func NewRegistryClientFromConf(reg *Registry) *RegistryClient {
	return &RegistryClient{
		client:       http.Client{Timeout: 60 * time.Second},
		Auth:         reg.Auth,
		AuthUsername: reg.AuthUsername,
		AuthPassword: reg.AuthPassword,
		AuthURL:      reg.AuthURL,
		BaseURL:      reg.BaseURL,
	}
}

func (c *RegistryClient) get(url string, result interface{}, basicAuth bool) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	request.Header.Set("User-Agent", AgentStr)
	request.Header.Set("Accept", "application/vnd.docker.distribution.manifest.list.v2+json")

	if basicAuth && c.AuthUsername != "" && c.AuthPassword != "" {
		request.SetBasicAuth(c.AuthUsername, c.AuthPassword)
	}

	if c.Auth && c.AuthToken != "" {
		request.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("request failed with status code %d", response.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(responseBody, result); err != nil {
		return err
	}
	return nil
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

func (c *RegistryClient) Login(repo string) {
	if c.Auth {
		var authResponse AuthResponse
		err := c.get(c.AuthURL+"&scope=repository:"+repo+":pull", &authResponse, true)
		if err != nil {
			log.Println(err)
		}
		c.AuthToken = authResponse.Token
	}
}

func (c *RegistryClient) ListTags(repo string) []string {
	var tagsResponse TagsResponse
	url := c.BaseURL + repo + "/tags/list"
	if err := c.get(url, &tagsResponse, false); err != nil {
		log.Println(err)
		return nil
	}
	return tagsResponse.Tags
}

func (c *RegistryClient) ListTagDigests(repo, architecture string, allTags, tagPattern []string) []TagDigest {
	wg := sync.WaitGroup{}
	results := make(chan struct {
		tag      string
		response ManifestResponse
	})
	for _, tag := range allTags {
		wg.Add(1)
		go func(tag string) {
			defer wg.Done()
			if !matchesAny(tagPattern, tag) {
				return
			}

			var manifestResponse ManifestResponse
			url := c.BaseURL + repo + "/manifests/" + tag
			if err := c.get(url, &manifestResponse, false); err != nil {
				log.Println(err)
				return
			}
			results <- struct {
				tag      string
				response ManifestResponse
			}{tag: tag, response: manifestResponse}
		}(tag)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	res := make([]TagDigest, 0, 100)
	for r := range results {
		for _, manifest := range r.response.Manifests {
			if manifest.Platform.Architecture == architecture {
				res = append(res, TagDigest{r.tag, manifest.Digest})
			}
		}
	}

	return res
}
