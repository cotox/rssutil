# rssutil

The RSS util speak in Go.

**RSSUTIL IS IN DRAFT**.

### Trivial

```golang
package main

import (
	"log"
	"time"

	"github.com/cotox/rssutil"
)

func onRSSUpdate(items []rssutil.RSSItem) {
	log.Printf("Received %d new RSS item(s).", len(items))
	for _, item := range items {
		log.Printf("  %s, %s", item.PubDate, item.Title)
	}
}

func main() {
	source := "http://localhost:8000/rss"
	interval := 4 * time.Second
	log.Fatal(rssutil.Serve(source, onRSSUpdate, interval))
}
```

### Licence

GPLv3
