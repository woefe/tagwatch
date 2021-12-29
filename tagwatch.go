package tagwatch

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
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

func MakeFeed(conf *Conf) *[]byte {
	feed := NewFeed("Docker registry tags", "https://github.com/woefe/tagwatch", makeDescription(conf.Tagwatch))
	for _, watchConf := range conf.Tagwatch {
		repo := watchConf.Repo
		arch := watchConf.Arch
		reg := watchConf.Registry
		client := NewRegistryClient(reg.Auth, reg.AuthURL, reg.BaseURL)
		for _, taggedDigest := range client.ListTags(repo, arch, watchConf.Tags) {
			feed.AppendItems(&Item{
				Title:       fmt.Sprintf("%s in %s (%s)", taggedDigest.Tag, repo, arch),
				Link:        "",
				Description: taggedDigest.Tag + ": " + taggedDigest.Digest,
				Guid:        makeGuid(taggedDigest, reg, repo, arch),
			})
		}
	}
	xml:= feed.ToXML()
	return &xml
}
