package site

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/galanhao/blog-forge/internal/content"
	"github.com/galanhao/blog-forge/internal/theme"
)

// rssXML defines the top-level RSS structure.
type rssXML struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	AtomNS  string   `xml:"xmlns:atom,attr"`
	Channel channel  `xml:"channel"`
}

// channel represents the RSS channel.
type channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	AtomLink    atomLink  `xml:"atom:link"`
	Description string    `xml:"description"`
	Generator   string    `xml:"generator"`
	LastBuild   string    `xml:"lastBuildDate"`
	Items       []rssItem `xml:"item"`
}

// atomLink represents the self-referencing atom:link element required by RSS spec.
type atomLink struct {
	XMLName xml.Name `xml:"atom:link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
}

// rssItem represents a single RSS entry.
type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
	Content     string `xml:"content:encoded"`
	PubDate     string `xml:"pubDate"`
}

// writeRSS generates the RSS feed with the latest posts.
func (b *Builder) writeRSS(distDir string, siteCtx theme.SiteCtx, posts []*content.Post) error {
	limit := 20
	if len(posts) < limit {
		limit = len(posts)
	}

	items := make([]rssItem, 0, limit)
	for i := 0; i < limit; i++ {
		p := posts[i]
		permalink := siteCtx.URL + p.Permalink
		items = append(items, rssItem{
			Title:       p.Title,
			Link:        permalink,
			GUID:        permalink,
			Description: p.Excerpt,
			Content:     p.HTML,
			PubDate:     p.Date.UTC().Format(time.RFC1123Z),
		})
	}

	rss := rssXML{
		Version: "2.0",
		AtomNS:  "http://www.w3.org/2005/Atom",
		Channel: channel{
			Title:       siteCtx.Title,
			Link:        siteCtx.URL,
			AtomLink:    atomLink{Href: siteCtx.URL + "/index.xml", Rel: "self", Type: "application/rss+xml"},
			Description: siteCtx.Subtitle,
			Generator:   "blog-forge",
			LastBuild:   time.Now().UTC().Format(time.RFC1123Z),
			Items:       items,
		},
	}

	data, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return fmt.Errorf("rss: marshal: %w", err)
	}

	out := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n%s\n", data)
	return writeFile(fmt.Sprintf("%s/index.xml", distDir), out)
}
