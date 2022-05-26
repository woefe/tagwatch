package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html"
	"html/template"
	"log"
	"sort"
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
	tagsItemTemplate        = `Available tags of <em>{{ .Repo }}</em> have changed. Current tags are:<br/><ul>{{ range .Tags }}<li><code>{{ . }}</code></li>{{ end }}</ul>`
	feedDescriptionTemplate = "New tags for{{ range . }}\n  - {{ .Repo }}({{ .Arch }})): {{ join .Tags \", \" }}{{ end }}"
)

func makeGuid(first string, strings ...string) string {
	hasher := sha256.New()
	hasher.Write([]byte(first))
	for _, str := range strings {
		hasher.Write([]byte(str))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func makeDescription(tpl *template.Template, watches []*WatchConf) string {
	var b strings.Builder
	err := tpl.Execute(&b, watches)

	if err != nil {
		log.Println(err)
		return ""
	}

	return b.String()
}

func makeDigestLink(baseURL string, repo string, taggedDigest TagDigest) string {
	if baseURL != "https://registry.hub.docker.com/v2/" {
		return baseURL + repo + "/manifests/" + taggedDigest.Tag
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

func makeTagsLink(baseURL, repo string) string {
	if baseURL != "https://registry.hub.docker.com/v2/" {
		return baseURL + repo + "/tags/list"
	}
	if strings.HasPrefix(repo, "library/") {
		return "https://hub.docker.com/_/" + strings.TrimPrefix(repo, "library/") + "?tab=tags"
	}
	return "https://hub.docker.com/r/" + repo + "/tags"
}

func makeTagsDescription(tpl *template.Template, repo string, allTags []string) string {
	var b strings.Builder
	err := tpl.Execute(&b, struct {
		Repo string
		Tags []string
	}{
		Repo: repo,
		Tags: allTags,
	})

	if err != nil {
		log.Println(err)
		return ""
	}

	return b.String()
}

func MakeFeed(conf *Conf) *[]byte {
	descriptionTpl, err := template.New("description").Funcs(template.FuncMap{"join": strings.Join}).Parse(feedDescriptionTemplate)
	tagsTpl, err := template.New("tags").Parse(tagsItemTemplate)

	if err != nil {
		log.Println(err)
		return nil
	}

	feed := NewFeed("Docker registry tags", "https://github.com/woefe/tagwatch", makeDescription(descriptionTpl, conf.Tagwatch))
	for _, watchConf := range conf.Tagwatch {
		repo := watchConf.Repo
		arch := watchConf.Arch
		reg := watchConf.Registry
		client := NewRegistryClientFromConf(reg)
		client.Login(repo)
		allTags := client.ListTags(repo)
		sort.Sort(sort.Reverse(ByVersion(allTags)))
		if watchConf.WatchNew {
			feed.AppendItems(NewItem(
				"Available tags of "+html.EscapeString(repo)+" have changed",
				makeTagsLink(reg.BaseURL, repo),
				makeTagsDescription(tagsTpl, repo, allTags),
				makeGuid(repo, allTags...),
			))
		}

		for _, taggedDigest := range client.ListTagDigests(repo, arch, allTags, watchConf.Tags) {
			title := html.EscapeString(fmt.Sprintf("%s:%s (%s)", repo, taggedDigest.Tag, arch))
			digest := html.EscapeString(taggedDigest.Digest)
			feed.AppendItems(NewItem(
				title,
				makeDigestLink(reg.BaseURL, repo, taggedDigest),
				"Digest of <em>"+title+"</em> has changed. New digest is:<br/><code>"+digest+"</code>",
				makeGuid(reg.BaseURL, repo, taggedDigest.Tag, arch, taggedDigest.Digest),
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
