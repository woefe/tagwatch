package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
)

const (
	Version   = "1.0"
	AgentStr  = "tagwatch/" + Version
	EmptyFeed = `
		<rss version="2.0">
			<channel>
				<title>Docker registry tags</title>
				<link>https://github.com/woefe/tagwatch</link>
				<description>No tags available</description>
				<generator>` + AgentStr + `</generator>
			</channel>
		</rss>
	`
)

func makeGuid(taggedDigest TagDigest, reg *Registry, repo string, arch string) string {
	hasher := sha256.New()
	hasher.Write([]byte(reg.BaseURL))
	hasher.Write([]byte(repo))
	hasher.Write([]byte(taggedDigest.Tag))
	hasher.Write([]byte(arch))
	hasher.Write([]byte(taggedDigest.Digest))
	return hex.EncodeToString(hasher.Sum(nil))
}

func makeDescription(watches []*WatchConf) string {
	var b strings.Builder
	b.WriteString("New tags for\n")
	for _, watchConf := range watches {
		tags := strings.Join(watchConf.Tags, ", ")
		b.WriteString("  - ")
		b.WriteString(watchConf.Repo)
		b.WriteString(" (")
		b.WriteString(watchConf.Arch)
		b.WriteString("): ")
		b.WriteString(tags)
		b.WriteRune('\n')
	}
	return b.String()
}

func makeLink(baseUrl string, repo string, taggedDigest TagDigest) string {
	if baseUrl != "https://registry.hub.docker.com/v2/" {
		return baseUrl + repo + "/manifests/" + taggedDigest.Tag
	}
	if strings.HasPrefix(repo, "library/") {
		repo = strings.TrimPrefix(repo, "library/") + "/" + repo
	}
	return fmt.Sprintf(
		"https://hub.docker.com/layers/%s/%s/images/%s?context=explore",
		repo,
		taggedDigest.Tag,
		strings.Replace(taggedDigest.Digest, ":", "-", 1),
	)
}

func MakeFeed(conf *Conf) *[]byte {
	feed := NewFeed("Docker registry tags", "https://github.com/woefe/tagwatch", makeDescription(conf.Tagwatch))
	for _, watchConf := range conf.Tagwatch {
		repo := watchConf.Repo
		arch := watchConf.Arch
		reg := watchConf.Registry
		client := NewRegistryClient(reg.Auth, reg.AuthURL, reg.BaseURL)
		for _, taggedDigest := range client.ListTags(repo, arch, watchConf.Tags) {
			title := fmt.Sprintf("%s:%s (%s)", repo, taggedDigest.Tag, arch)
			feed.AppendItems(NewItem(
				title,
				makeLink(reg.BaseURL, repo, taggedDigest),
				"<p>Digest of "+title+" changed. Digest now is:</p><pre>"+taggedDigest.Digest+"</pre>",
				makeGuid(taggedDigest, reg, repo, arch),
			))
		}
	}
	xml, err := feed.ToXML()
	if err != nil {
		log.Println(err)
		xml = []byte(EmptyFeed)
	}
	return &xml
}
