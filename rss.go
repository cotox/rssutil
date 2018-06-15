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

var DefaultTTL = 20 * time.Minute

// Feed creates RSS implementation from binary and return.
func Feed(b []byte) (rss *RSS, err error) {
	logTrace("feed()")

	rss = new(RSS)
	decoder := xml.NewDecoder(bytes.NewBuffer(b))
	if err := decoder.Decode(rss); err != nil {
		logErr(err)
		return nil, err
	}

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

func (rss *RSS) Update() (newItems []RSSItem, err error) {
	logTrace("rss.Update()")

	latestItem := rss.latestItem()

	if rss.source == "" {
		return nil, fmt.Errorf("empty rss.source")
	}

	if rss.source[:4] == "http" {
		rss, err = FeedFromURL(rss.source)
	} else {
		rss, err = FeedFromFile(rss.source)
	}
	if err != nil {
		logErr(err)
		return nil, err
	}

	items := rss.Channel.Items
	for i := range items {
		if items[i].PubDate.After(latestItem.PubDate) {
			newItems = append(newItems, items[i])
		}
	}

	return newItems, nil
}

func (rss RSS) Serve(forceTTL int) {
	var ttl time.Duration
	if rss.Channel.TTL != 0 {
		ttl = time.Duration(rss.Channel.TTL) * time.Minute
	} else {
		ttl = DefaultTTL
	}
	if forceTTL != 0 {
		ttl = time.Duration(forceTTL) * time.Second
	}

	time.Sleep(ttl - time.Now().Sub(rss.lastUpdateAt))

	for {
		newItems, err := rss.Update()
		if err != nil {
			logErr(err)
		}

		if newItems != nil && rss.OnRSSUpdate != nil {
			rss.OnRSSUpdate(newItems)
		}

		time.Sleep(ttl)
	}
}

func (rss RSS) latestItem() (latestItem *RSSItem) {
	items := rss.Channel.Items
	latestItem = &items[0]
	for i := 1; i < len(items); i++ {
		if items[i].PubDate.After(latestItem.PubDate) {
			latestItem = &items[i]
		}
	}
	return latestItem
}
