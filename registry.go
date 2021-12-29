package tagwatch

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

func get(url string, result interface{}, headers map[string]string) error {
	client := http.Client{Timeout: 60 * time.Second}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}
	request.Header.Set("User-Agent", "registry-update-check/1.0")

	response, err := client.Do(request)
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
			log.Fatalln(err)
		}
		if matched {
			return true
		}
	}
	return false
}

func ListTags(repo, architecture string, tagPattern []string) []TagDigest {
	var authResponse AuthResponse
	url := "https://auth.docker.io/token?service=registry.docker.io&scope=repository:" + repo + ":pull"
	if err := get(url, &authResponse, nil); err != nil {
		log.Fatalln(err)
		return nil
	}

	headers := map[string]string{
		"Authorization": "Bearer " + authResponse.Token,
		"Accept":        "application/vnd.docker.distribution.manifest.list.v2+json",
	}

	var tagsResponse TagsResponse
	url = "https://registry.hub.docker.com/v2/" + repo + "/tags/list"
	if err := get(url, &tagsResponse, headers); err != nil {
		log.Fatalln(err)
		return nil
	}

	wg := sync.WaitGroup{}
	results := make(chan struct {
		tag      string
		response ManifestResponse
	})
	for _, tag := range tagsResponse.Tags {
		wg.Add(1)
		go func(tag string) {
			defer wg.Done()
			if !matchesAny(tagPattern, tag) {
				return
			}

			var manifestResponse ManifestResponse
			url := "https://registry.hub.docker.com/v2/" + repo + "/manifests/latest"
			if err := get(url, &manifestResponse, headers); err != nil {
				log.Fatalln(err)
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
