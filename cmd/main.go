package main

import (
	"fmt"
	"os"
	"tagwatch"
)

func makeGuid(taggedDigest tagwatch.TagDigest) string {
	return taggedDigest.Digest + taggedDigest.Tag
}

func main() {
	feed := tagwatch.NewFeed("Docker registry tags", "https://hub.docker.com", "new tags for patterns")
	repo := "library/ubuntu"
	arch := "amd64"
	for _, taggedDigest := range tagwatch.ListTags(repo, arch, []string{"20\\.04"}) {
		fmt.Println(taggedDigest.Digest, taggedDigest.Tag)
		feed.AppendItems(&tagwatch.Item{
			Title:       fmt.Sprintf("%s in %s (%s)", taggedDigest.Tag, repo, arch),
			Link:        "",
			Description: taggedDigest.Tag + ": " + taggedDigest.Digest,
			Guid:        makeGuid(taggedDigest),
		})
	}
	os.Stdout.Write(feed.ToXML())
}
