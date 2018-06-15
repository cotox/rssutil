// Copyright 2018 cotox. All rights reserved.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

package rssutil

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const DefaultTTL = 20 * time.Minute

var stopServe = make(chan struct{})

// Feed creates RSS implementation from binary and return.
func Feed(b []byte) (rss *RSS, err error) {
	logTrace("feed()")

	rss = new(RSS)
	decoder := xml.NewDecoder(bytes.NewBuffer(b))
	if err := decoder.Decode(rss); err != nil {
		logErr(err)
		return nil, err
	}

	// Trim elements in string type.
	const cutset = " \t\n"
	rss.Channel.Title = strings.Trim(rss.Channel.Title, cutset)
	rss.Channel.Description = strings.Trim(rss.Channel.Description, cutset)
	rss.Channel.Copyright = strings.Trim(rss.Channel.Copyright, cutset)
	for i := range rss.Channel.Items {
		item := &rss.Channel.Items[i]
		item.Title = strings.Trim(item.Title, cutset)
		item.Description = strings.Trim(item.Description, cutset)
	}

	rss.origin = b

	rss.lastUpdateAt = time.Now()

	return rss, nil
}

// FeedFromFile creates RSS implementation from specific file and return.
func FeedFromFile(filename string) (rss *RSS, err error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		logErr(err)
		return nil, err
	}

	rss, err = Feed(b)
	if err != nil {
		logErr(err)
		return nil, err
	}

	rss.source = filename

	return rss, nil
}

// FeedFromURL creates RSS implementation from specific URL and return.
func FeedFromURL(url string) (rss *RSS, err error) {
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		logErr(err)
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logErr(err)
		return nil, err
	}

	rss, err = Feed(b)
	if err != nil {
		logErr(err)
		return nil, err
	}

	rss.source = url

	return rss, nil
}

// Update updates RSS content and returns the newer RSSItem list.
func (rss *RSS) Update() (newItems []RSSItem, err error) {
	logTrace("rss.Update()")

	latestItem := rss.latestItem()

	if rss.source == "" {
		return nil, fmt.Errorf("empty rss.source")
	}

	var rss2 *RSS
	if rss.source[:4] == "http" {
		rss2, err = FeedFromURL(rss.source)
		if err != nil {
			logErr(err)
			return nil, err
		}
	} else {
		rss2, err = FeedFromFile(rss.source)
		if err != nil {
			logErr(err)
			return nil, err
		}
	}
	rss.Channel.Items = rss2.Channel.Items
	rss.lastUpdateAt = time.Now()

	if latestItem == nil {
		return nil, nil
	}

	items := rss.Channel.Items
	for i := range items {
		if items[i].PubDate.After(latestItem.PubDate) {
			newItems = append(newItems, items[i])
		}
	}

	return newItems, nil
}

// Serve updated RSS content in background automatically.
// And calls registered RSSUpdateNotifiers when new RSSItems come.
//
// The RSS content will update every ttl minutes. If ttl is 0, it tries
// to use TTL specified in RSSChannel, then DefaultTTL if RSSChannel.TTL
// is not specified.
func (rss *RSS) Serve(ttl time.Duration) error {
	if ttl == 0 {
		if rss.Channel.TTL > 0 {
			ttl = time.Duration(rss.Channel.TTL) * time.Minute
		} else {
			ttl = DefaultTTL
		}
	}

	// time.Sleep(ttl - time.Now().Sub(rss.lastUpdateAt))
	ticker := time.NewTicker(ttl)
	defer ticker.Stop()

serveLoop:
	for {
		select {
		case <-stopServe:
			break serveLoop
		case <-ticker.C:
			newItems, err := rss.Update()
			if err != nil {
				logErr(err)
				return err
			}
			if newItems != nil {
				for _, f := range rss.rssUpdateNotifiers {
					go f(newItems)
				}
			}
		}
	}

	return nil
}

// Stop to serve.
func (rss *RSS) Stop() { stopServe <- struct{}{} }

func (rss *RSS) RegisterRSSUpdateNotifier(f func([]RSSItem)) {
	rss.mu.Lock()
	rss.rssUpdateNotifiers = append(rss.rssUpdateNotifiers, f)
	rss.mu.Unlock()
}

// Serve create an RSS implementation and keep auto update in background.
//
// Argument source specifies the URL of RSS.
// The RSS content will update every ttl minutes. If ttl is 0, it tries
// to use TTL specified in RSSChannel, then DefaultTTL if RSSChannel.TTL
// is not specified.
func Serve(source string, f RSSUpdateNotifier, ttl time.Duration) error {
	var rss *RSS
	var err error
	if source[:4] == "http" {
		rss, err = FeedFromURL(source)
		if err != nil {
			logDebugln("ERROR:", err)
			return err
		}
	} else {
		rss, err = FeedFromFile(source)
		if err != nil {
			logDebugln("ERROR", err)
			return err
		}
	}

	rss.RegisterRSSUpdateNotifier(f)

	if rss.Channel.Items != nil {
		go f(rss.Channel.Items)
	}

	return rss.Serve(ttl)
}

// Stop to serve.
func Stop() { stopServe <- struct{}{} }

func (rss *RSS) latestItem() (latestItem *RSSItem) {
	items := rss.Channel.Items
	if len(items) < 1 {
		return nil
	}
	latestItem = &items[0]
	for i := 1; i < len(items); i++ {
		if items[i].PubDate.After(latestItem.PubDate) {
			latestItem = &items[i]
		}
	}
	return latestItem
}
