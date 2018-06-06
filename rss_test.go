// Copyright 2018 cotox. All rights reserved.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

package rssutil

import (
	"testing"
	"time"
)

/**
 * [RSS 0.91 Specification]: http://backend.userland.com/rss091
 * [RSS 0.91 Sample]: https://cyber.harvard.edu/rss/examples/sampleRss091.xml
 *
 * [RSS 0.92 Specification]: http://backend.userland.com/rss092
 * [RSS 0.92 Sample]: https://cyber.harvard.edu/rss/examples/sampleRss092.xml
 *
 * [RSS 2.0 Specification]: https://cyber.harvard.edu/rss/rss.html
 * [RSS 2.0 Sample]: https://cyber.harvard.edu/rss/examples/rss2sample.xml
 */

var rss091Text = ``

var rss092Text = ``

var rss20Text = `
	<?xml version="1.0" encoding="UTF-8"?>
	<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
		<channel>
			<title>最新更新 &#8211; Solidot</title>
			<link>https://www.solidot.org</link>
			<description><![CDATA[奇客的资讯，重要的东西。]]></description>
			<atom:link href="https://www.solidot.org/index.rss" rel="self" type="application/rss+xml"></atom:link>
			<language>zh-cn</language>
			<lastBuildDate>Fri, 11 May 2018 16:45:56 +0800</lastBuildDate>
			<docs>http://blogs.law.harvard.edu/tech/rss</docs>
			<generator>Weblog Editor 2.0</generator>
			<managingEditor>editor@example.com</managingEditor>
			<webMaster>webmaster@example.com</webMaster>
			<ttl>20</ttl>
			<item>
				<title><![CDATA[中国年轻一代不愿意长时间工作]]></title>
				<link><![CDATA[https://www.solidot.org/story?sid=56470]]></link>
				<description><![CDATA[中国科技行业流行的 <a href="https://www.solidot.org/story?sid=51481" target="_blank">996 工作制</a>——早九点到晚九点每周六天——正遭到年轻一代专业人士的挑战。<a href="https://www.solidot.org/story?sid=51481" target="_blank">千禧年一代不愿意长时间工作</a>，他们通常受到更好的教育，更了解自己的权利，比上一代人更有兴趣找到能实现个人抱负的东西。作为独生子女（中国的独生子女政策在 2015 年才终止），他们也更坦率和骄纵。劳工权利专家称，尤其是 90 后一代不愿意加班，他们自我中心。中国快速的经济转型创造了一个庞大的中产阶级，在 2012 年七成的城市人口收入在 9000 到 3.4 万美元之间。而在 2000 年，这一比例仅为 4%。作为独生子女，千禧年一代得到了家庭的支持，为他们提供了安全的经济保障，即使他们的事业不如意。年轻人不再愿意为了微薄的薪水而干太长时间。<p><img src="https://img.solidot.org/0/446/liiLIZF8Uh6yM.jpg" height="120" style="display:block"/></p>]]></description>
				<pubDate>Fri, 11 May 2018 16:28:39 +0800</pubDate>
				<guid>http://liftoff.msfc.nasa.gov/2003/06/03.html#item573</guid>
			</item>
		</channel>
	</rss>`

//
func TestRSS20Feed(t *testing.T) {
	rss, err := Feed([]byte(rss20Text))
	if err != nil {
		t.Error("decode failed")
	}
	if rss.Version != "2.0" {
		t.Error("rss.Version != \"2.0\"")
	}
}

func TestRSS20FeedFromFile(t *testing.T) {
	rss, err := FeedFromFile("rss2sample.xml")
	if err != nil {
		t.Error("decode failed")
	}
	if rss.Version != "2.0" {
		t.Error("rss.Version != \"2.0\"")
	}
}

func TestRSS20FeedFromURL(t *testing.T) {
	rss, err := FeedFromURL("https://cyber.harvard.edu/rss/examples/rss2sample.xml")
	if err != nil {
		t.Error("decode failed")
	}
	if rss.Version != "2.0" {
		t.Error("rss.Version != \"2.0\"")
	}
}

func TestRSS20Channel(t *testing.T) {
	rss, _ := Feed([]byte(rss20Text))
	ch := rss.Channel

	if ch.Title != "最新更新 – Solidot" {
		t.Error("rss.Channel.Title != \"最新更新 – Solidot\"")
	}

	if ch.Link != "https://www.solidot.org" {
		t.Errorf("rss.Channel.Link != \"https://www.solidot.org\", %#v`", ch.Link)
	}

	if ch.Description != "奇客的资讯，重要的东西。" {
		t.Error("rss.Channel.Description != \"奇客的资讯，重要的东西。\"")
	}

	if ch.Language != "zh-cn" {
		t.Error("rss.Channel.Language != \"zh-cn\"")
	}

	// if ch.Copyright != ""      { t.Error("ch.Copyright != \"\"") }

	if ch.ManagingEditor != "editor@example.com" {
		t.Error("ch.ManagingEditor != \"editor@example.com\"")
	}

	if ch.WebMaster != "webmaster@example.com" {
		t.Error("ch.WebMaster != \"webmaster@example.com\"")
	}

	// if ch.PubDate != ""        { t.Error("ch.PubDate != \"\"") }

	if !ch.LastBuildDate.Equal(time.Date(2018, 5, 11, 8, 45, 56, 0, time.UTC)) {
		t.Error("ch.LastBuildDate != \"Fri, 11 May 2018 16:45:56 +0800\"")
	}

	// if ch.Category != ""       { t.Error("ch.Category != \"\"") }

	if ch.Generator != "Weblog Editor 2.0" {
		t.Error("ch.Generator != \"Weblog Editor 2.0\"")
	}

	if ch.Docs != "http://blogs.law.harvard.edu/tech/rss" {
		t.Error("ch.Docs != \"http://blogs.law.harvard.edu/tech/rss\"")
	}

	// if ch.Cloud != ""          { t.Error("ch.Cloud != \"\"") }

	if ch.TTL != 20 {
		t.Error("ch.TTL != 20")
	}

	// if ch.Image != ""          { t.Error("ch.Image != \"\"") }

	// if ch.Rating != ""         { t.Error("ch.Rating != \"\"") }

	// if ch.TextInput != ""      { t.Error("ch.TextInput != \"\"") }

	if ch.SkipHours != 0 {
		t.Error("ch.SkipHours != 0")
	}

	if ch.SkipDays != 0 {
		t.Error("ch.SkipDays != 0")
	}
}

func TestRSS20Items(t *testing.T) {
	rss, _ := Feed([]byte(rss20Text))
	its := rss.Channel.Items

	if len(its) != 1 {
		t.Error("len(its) != 1")
	}

	it0 := its[0]

	if it0.Title != "中国年轻一代不愿意长时间工作" {
		t.Error("it0.Title != \"中国年轻一代不愿意长时间工作\"")
	}

	if it0.Link != "https://www.solidot.org/story?sid=56470" {
		t.Error("it0.Link != \"https://www.solidot.org/story?sid=56470\"")
	}

	if it0.Description != "中国科技行业流行的 <a href=\"https://www.solidot.org/story?sid=51481\" target=\"_blank\">996 工作制</a>——早九点到晚九点每周六天——正遭到年轻一代专业人士的挑战。<a href=\"https://www.solidot.org/story?sid=51481\" target=\"_blank\">千禧年一代不愿意长时间工作</a>，他们通常受到更好的教育，更了解自己的权利，比上一代人更有兴趣找到能实现个人抱负的东西。作为独生子女（中国的独生子女政策在 2015 年才终止），他们也更坦率和骄纵。劳工权利专家称，尤其是 90 后一代不愿意加班，他们自我中心。中国快速的经济转型创造了一个庞大的中产阶级，在 2012 年七成的城市人口收入在 9000 到 3.4 万美元之间。而在 2000 年，这一比例仅为 4%。作为独生子女，千禧年一代得到了家庭的支持，为他们提供了安全的经济保障，即使他们的事业不如意。年轻人不再愿意为了微薄的薪水而干太长时间。<p><img src=\"https://img.solidot.org/0/446/liiLIZF8Uh6yM.jpg\" height=\"120\" style=\"display:block\"/></p>" {
		t.Error("it0.Description != \"中国科技行业流行的 <a href=\"https://www.solidot.org/story?sid=51481\" target=\"_blank\">996 工作制</a>——早九点到晚九点每周六天——正遭到年轻一代专业人士的挑战。<a href=\"https://www.solidot.org/story?sid=51481\" target=\"_blank\">千禧年一代不愿意长时间工作</a>，他们通常受到更好的教育，更了解自己的权利，比上一代人更有兴趣找到能实现个人抱负的东西。作为独生子女（中国的独生子女政策在 2015 年才终止），他们也更坦率和骄纵。劳工权利专家称，尤其是 90 后一代不愿意加班，他们自我中心。中国快速的经济转型创造了一个庞大的中产阶级，在 2012 年七成的城市人口收入在 9000 到 3.4 万美元之间。而在 2000 年，这一比例仅为 4%。作为独生子女，千禧年一代得到了家庭的支持，为他们提供了安全的经济保障，即使他们的事业不如意。年轻人不再愿意为了微薄的薪水而干太长时间。<p><img src=\"https://img.solidot.org/0/446/liiLIZF8Uh6yM.jpg\" height=\"120\" style=\"display:block\"/></p>\"")
	}

	// if it0.Author != ""      { t.Error("it0.Author != \"\"") }

	// if it0.Category != ""    { t.Error("it0.Category != \"\"") }

	// if it0.Category != ""    { t.Error("it0.Category != \"\"") }

	// if it0.Comments != ""    { t.Error("it0.Comments != \"\"") }

	// if it0.Enclosure != ""   { t.Error("it0.Enclosure != \"\"") }

	g := GUID{"http://liftoff.msfc.nasa.gov/2003/06/03.html#item573", false}
	if it0.GUID != g {
		t.Error("it0.GUID != \"http://liftoff.msfc.nasa.gov/2003/06/03.html#item573\"")
	}

	if !it0.PubDate.Equal(time.Date(2018, 5, 11, 8, 28, 39, 0, time.UTC)) {
		t.Error("it0.PubDate != \"2018-05-11T08:28:39Z\"")
	}

	// if it0.Source != ""      { t.Error("it0.Source != \"\"") }
}

func TestRequiredChannelElements(t *testing.T) {
	// RSS 2.0 Specification has 3 required channel elements. They are,
	//
	// - title
	//
	//   The name of the channel. It's how people refer to your service. If
	//   you have an HTML website that contains the same information as your
	//   RSS file, the title of your channel should be the same as the title
	//   of your website.
	//
	// - link
	//
	//   The URL to the HTML website corresponding to the channel.
	//
	// - description
	//
	//   Phrase or sentence describing the channel.
}

func TestOptionalChannelElements(t *testing.T) {
	// RSS 2.0 Specification has 3 required channel elements. They are,
	//
	// 1.  language
	// 2.  copyright
	// 3.  managingEditor
	// 4.  webMaster
	// 5.  pubDate
	// 6.  lastBuildDate
	// 7.  category
	// 8.  generator
	// 9.  docs
	// 10. cloud
	// 11. ttl
	// 12. image
	// 13. textInput
	// 14. skipHours
	// 15. skipDays
}
