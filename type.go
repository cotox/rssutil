// Copyright 2018 cotox. All rights reserved.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

package rssutil

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// RSS is a Web content syndication format.
//
// Its name is an acronym for Really Simple Syndication.
//
// RSS is a dialect of XML. All RSS files must conform to the XML 1.0
// specification, as published on the World Wide Web Consortium (W3C)
// website.
//
// A summary of RSS version history.
//
// At the top level, a RSS document is a <rss> element, with a mandatory
// attribute called version, that specifies the version of RSS that the
// document conforms to. If it conforms to this specification, the
// version attribute must be 2.0.
//
// Subordinate to the <rss> element is a single <channel> element, which
// contains information about the channel (metadata) and its contents.
type RSS struct {
	Version string     `xml:"version,attr" json:"version"`
	Channel RSSChannel `xml:"channel"      json:"channel"`

	origin       []byte
	source       string
	lastUpdateAt time.Time

	OnRSSUpdate func(newItems []RSSItem)
}

func (rss RSS) String() string {
	return "Version: \"" + rss.Version + "\", Channel: {" + rss.Channel.String() + "}"
}

func (rss RSS) ToJSON() string {
	data := struct {
		Source  string     `json:"source"`
		Version string     `json:"version"`
		Channel RSSChannel `json:"channel"`
	}{rss.source, rss.Version, rss.Channel}
	b, err := json.MarshalIndent(data, "", "  ")
	// b, err := json.Marshal(data)
	if err != nil {
		logErr(err)
		return err.Error()
	}
	return string(b)
}

type RSSChannel struct {

	/*************************** Required elements ***************************/

	// The name of the channel. It's how people refer to your service. If
	// you have an HTML website that contains the same information as
	// your RSS file, the title of your channel should be the same as the
	// title of your website.
	//
	// Sample:
	//   GoUpstate.com News Headlines
	Title string `xml:"title" json:"title"`

	// The URL to the HTML website corresponding to the channel.
	//
	// Sample:
	//   http://www.goupstate.com/
	Link string `xml:"link" json:"link"`

	// Phrase or sentence describing the channel.
	//
	// Sample:
	//   The latest news from GoUpstate.com, a Spartanburg Herald-Journal Web site.
	Description string `xml:"description" json:"description"`

	/*************************** Optional elements ***************************/

	// The language the channel is written in. This allows aggregators to
	// group all Italian language sites, for example, on a single page.
	// A list of allowable values for this element, as provided by
	// Netscape, is [here](https://cyber.harvard.edu/rss/languages.html).
	// You may also use [values defined](https://www.w3.org/TR/REC-html40/struct/dirlang.html#langcodes)
	// by the W3C.
	//
	// Sample:
	//   en-us
	Language string `xml:"language,omitempty" json:"language,omitempty"`

	// Copyright notice for content in the channel.
	//
	// Sample:
	//   Copyright 2002, Spartanburg Herald-Journal
	Copyright string `xml:"copyright,omitempty" json:"copyright,omitempty"`

	// Email address for person responsible for editorial content.
	//
	// Sample:
	//   geo@herald.com (George Matesky)
	ManagingEditor string `xml:"managingEditor,omitempty" json:"managingEditor,omitempty"`

	// Email address for person responsible for technical issues relating to channel.
	//
	// Sample:
	//   betty@herald.com (Betty Guernsey)
	WebMaster string `xml:"webMaster,omitempty" json:"webMaster,omitempty"`

	// The publication date for the content in the channel. For example,
	// the New York Times publishes on a daily basis, the publication
	// date flips once every 24 hours. That's when the pubDate of the
	// channel changes. All date-times in RSS conform to the Date and
	// Time Specification of [RFC 822](http://asg.web.cmu.edu/rfc/rfc822.html),
	// with the exception that the year may be expressed with two
	// characters or four characters (four preferred).
	//
	// Sample:
	//   Sat, 07 Sep 2002 00:00:01 GMT
	PubDate *RFC822 `xml:"pubDate,omitempty" json:"pubDate,omitempty"`

	// The last time the content of the channel changed.
	//
	// Sample:
	//   Sat, 07 Sep 2002 09:42:31 GMT
	LastBuildDate *RFC822 `xml:"lastBuildDate,omitempty" json:"lastBuildDate,omitempty"`

	// Specify one or more categories that the channel belongs to.
	// Follows the same rules as the <item>-level
	// [category](https://cyber.harvard.edu/rss/rss.html#ltcategorygtSubelementOfLtitemgt)
	// element. More [info](https://cyber.harvard.edu/rss/rss.html#syndic8).
	//
	// Sample:
	//   <category>Newspapers</category>
	Categories []RSSCategory `xml:"category,omitempty" json:"category,omitempty"`

	// A string indicating the program used to generate the channel.
	//
	// Sample:
	//   MightyInHouse Content System v2.3
	Generator string `xml:"generator,omitempty" json:"generator,omitempty"`

	// A URL that points to the documentation for the format used in the
	// RSS file. It's probably a pointer to this page. It's for people
	// who might stumble across an RSS file on a Web server 25 years from
	// now and wonder what it is.
	//
	// Sample:
	//   http://blogs.law.harvard.edu/tech/rss
	Docs string `xml:"docs,omitempty" json:"docs,omitempty"`

	// Allows processes to register with a cloud to be notified of
	// updates to the channel, implementing a lightweight
	// publish-subscribe protocol for RSS feeds. More info
	// [here](https://cyber.harvard.edu/rss/rss.html#ltcloudgtSubelementOfLtchannelgt).
	//
	// Sample:
	//   <cloud domain="rpc.sys.com" port="80" path="/RPC2" registerProcedure="pingMe" protocol="soap"/>
	Cloud *RSSCloud `xml:"cloud,omitempty" json:"cloud,omitempty"`

	// TTL stands for time to live. It's a number of minutes that
	// indicates how long a channel can be cached before refreshing from
	// the source. More info [here](https://cyber.harvard.edu/rss/rss.html#ltttlgtSubelementOfLtchannelgt).
	//
	// Sample:
	//   <ttl>60</ttl>
	TTL int `xml:"ttl,omitempty" json:"ttl,omitempty"`

	// Specifies a GIF, JPEG or PNG image that can be displayed with the
	// channel.
	// More info [here](https://cyber.harvard.edu/rss/rss.html#ltimagegtSubelementOfLtchannelgt).
	Image *RSSImage `xml:"image,omitempty" json:"image,omitempty"`

	// The [PICS](https://www.w3.org/PICS/) rating for the channel.
	Rating string `xml:"rating,omitempty" json:"rating,omitempty"`

	// Specifies a text input box that can be displayed with the channel.
	// More info [here](https://cyber.harvard.edu/rss/rss.html#lttextinputgtSubelementOfLtchannelgt).
	TextInput *RSSTextInput `xml:"textInput,omitempty" json:"textInput,omitempty"`

	// A hint for aggregators telling them which hours they can skip.
	// More info [here](https://cyber.harvard.edu/rss/skipHoursDays.html#skiphours).
	SkipHours []int `xml:"skipHours>hour,omitempty" json:"skipHours,omitempty"`

	// A hint for aggregators telling them which days they can skip.
	// More info [here](https://cyber.harvard.edu/rss/skipHoursDays.html#skipdays).
	SkipDays []time.Weekday `xml:"skipDays>day,omitempty" json:"skipDays,omitempty"`

	Items []RSSItem `xml:"item,omitempty" json:"item,omitempty"`
}

func (c RSSChannel) String() string {
	var a []string

	// Required elements
	a = append(a, "Title: \""+c.Title+"\"")
	a = append(a, "Link: \""+c.Link+"\"")
	a = append(a, "Description: \""+c.Description+"\"")

	// Optional elements
	if c.Language != "" {
		a = append(a, "Language: \""+c.Language+"\"")
	}
	if c.Copyright != "" {
		a = append(a, "Copyright: \""+c.Copyright+"\"")
	}
	if c.ManagingEditor != "" {
		a = append(a, "ManagingEditor: \""+c.ManagingEditor+"\"")
	}
	if c.WebMaster != "" {
		a = append(a, "WebMaster: \""+c.WebMaster+"\"")
	}
	if !c.PubDate.IsZero() {
		a = append(a, "PubDate: "+c.PubDate.String())
	}
	if !c.LastBuildDate.IsZero() {
		a = append(a, "LastBuildDate: "+c.LastBuildDate.String())
	}
	if c.Categories != nil {
		var b []string
		for _, ca := range c.Categories {
			b = append(b, ca.String())
		}
		a = append(a, "Category: ["+strings.Join(b, ", ")+"]")
	}
	if c.Generator != "" {
		a = append(a, "Generator: \""+c.Generator+"\"")
	}
	if c.Docs != "" {
		a = append(a, "Docs: \""+c.Docs+"\"")
	}
	if c.Cloud != nil {
		a = append(a, "Cloud: {"+c.Cloud.String()+"}")
	}
	if c.TTL != 0 {
		a = append(a, "TTL: "+string(c.TTL))
	}
	if c.Image != nil {
		a = append(a, "Image: {"+c.Image.String()+"}")
	}
	if c.Rating != "" {
		a = append(a, "Rating: \""+c.Rating+"\"")
	}
	if c.TextInput != nil {
		a = append(a, "TextInput: {"+c.TextInput.String()+"}")
	}
	if c.SkipHours != nil {
		var b []string
		for _, v := range c.SkipHours {
			b = append(b, string(v))
		}
		a = append(a, "SkipHours: ["+strings.Join(b, ", ")+"]")
	}
	if c.SkipDays != nil {
		var b []string
		for _, v := range c.SkipDays {
			b = append(b, string(v))
		}
		a = append(a, "SkipDays: ["+strings.Join(b, ", ")+"]")
	}
	if c.Items != nil {
		var b []string
		for i := range c.Items {
			b = append(b, c.Items[i].String())
		}
		a = append(a, "Items: [{"+strings.Join(b, "}, {")+"}]")
	}

	return strings.Join(a, ", ")
}

// RSSCategory is an optional sub-element of RSSChannel/RSSItem.
//
// It has one optional attribute, domain, a string that identifies a
// categorization taxonomy.
//
// The value of the element is a forward-slash-separated string that
// identifies a hierarchic location in the indicated taxonomy. Processors
// may establish conventions for the interpretation of categories. Two
// examples are provided below:
//
// <category>Grateful Dead</category>
//
// <category domain="http://www.fool.com/cusips">MSFT</category>
//
// You may include as many category elements as you need to, for
// different domains, and to have an item cross-referenced in different
// parts of the same domain.
type RSSCategory struct {

	/*************************** Required elements ***************************/

	Value string `xml:",chardata"             json:"value"`

	/*************************** Optional elements ***************************/

	Domain string `xml:"domain,attr,omitempty" json:"domain,omitempty"`
}

func (c RSSCategory) String() string {
	if c.Domain == "" {
		return "\"" + c.Value + "\""
	}
	return fmt.Sprintf("\"%s\", domain=\"%s\"", c.Value, c.Domain)
}

// RSSCloud is an optional sub-element of RSSChannel. It specifies a web
// service that supports the RSSCloud interface which can be implemented
// in HTTP-POST, XML-RPC or SOAP 1.1.
//
// Its purpose is to allow processes to register with a cloud to be
// notified of updates to the channel, implementing a lightweight
// publish-subscribe protocol for RSS feeds.
//
// <cloud domain="rpc.sys.com" port="80" path="/RPC2" registerProcedure="myCloud.rssPleaseNotify" protocol="xml-rpc" />
//
// In this example, to request notification on the channel it appears in,
// you would send an XML-RPC message to rpc.sys.com on port 80, with a
// path of /RPC2. The procedure to call is myCloud.rssPleaseNotify.
//
// A full explanation of this element and the RSSCloud interface is
// [here](https://cyber.harvard.edu/rss/soapMeetsRss.html#rsscloudInterface).
type RSSCloud struct {

	/*************************** Required elements ***************************/

	Domain            string `xml:"domain,attr"            json:"domain"`
	Port              int    `xml:"port,attr"              json:"port"`
	Path              string `xml:"path,attr"              json:"path"`
	RegisterProcedure string `xml:"registerProcedure,attr" json:"registerProcedure"`
	Protocol          string `xml:"protocol,attr"          json:"protocol"`

	/*************************** Optional elements ***************************/

	// No optional element.
}

func (c RSSCloud) String() string {
	// All attributes are required.
	return fmt.Sprintf(
		"Domain: \"%s\", Port: %d, Path: \"%s\", RegisterProcedure: \"%s\", Protocol: \"%s\"",
		c.Domain, c.Port, c.Path, c.RegisterProcedure, c.Protocol)
}

// RSSImage is an optional sub-element of RSSChannel, which contains
// three required and three optional sub-elements.
type RSSImage struct {

	/*************************** Required elements ***************************/

	// URL is the URL of a GIF, JPEG or PNG image that represents the
	// Channel.
	URL string `xml:"url" json:"url"`

	// Title describes the image, it's used in the ALT attribute of the
	// HTML <img> tag when the channel is rendered in HTML.
	Title string `xml:"title" json:"title"`

	// Link is the URL of the site, when the channel is rendered, the
	// image is a link to the site. (Note, in practice the image Title
	// and Link should have the same value as the Channel's Title and Link.
	Link string `xml:"link" json:"link"`

	/*************************** Optional elements ***************************/

	// Width is an optional elements, in numbers, indicating the width of
	// the image in pixels.
	//
	// Maximum value for width is 144, default value is 88.
	Width int `xml:"width,omitempty" json:"width,omitempty"`

	// Height is an optional elements, in numbers, indicating the height
	// of the image in pixels.
	//
	// Maximum value for height is 400, default value is 31.
	Height int `xml:"height,omitempty" json:"height,omitempty"`

	// Description is an optional elements, which contains text that is
	// included in the TITLE attribute of the link formed around the
	// image in the HTML rendering.
	Description string `xml:"description,omitempty" json:"description,omitempty"`
}

func (img RSSImage) String() string {
	var a []string

	// Required elements
	a = append(a, "URL: \""+img.URL+"\"")
	a = append(a, "Title: \""+img.Title+"\"")
	a = append(a, "Link: \""+img.Link+"\"")

	// Optional elements
	if img.Width != 0 {
		a = append(a, fmt.Sprintf("Width: %d", img.Width))
	}
	if img.Height != 0 {
		a = append(a, fmt.Sprintf("Height: %d", img.Height))
	}
	if img.Description != "" {
		a = append(a, "Description: \""+img.Description+"\"")
	}

	return strings.Join(a, ", ")
}

// RSSTextInput is an optional sub-element of RSSChannel, which contains
// four required sub-elements.
//
// The purpose of the TextInput element is something of a mystery. You
// can use it to specify a search engine box. Or to allow a reader to
// provide feedback. Most aggregators ignore it.
type RSSTextInput struct {

	/*************************** Required elements ***************************/

	// The label of the Submit button in the text input area.
	Title string `xml:"title" json:"title"`

	// Explains the text input area.
	Description string `xml:"decsription" json:"decsription"`

	// The name of the text object in the text input area.
	Name string `xml:"name" json:"name"`

	// The URL of the CGI script that processes text input requests.
	Link string `xml:"link" json:"link"`

	/*************************** Optional elements ***************************/

	// No optional element.
}

func (ti RSSTextInput) String() string {
	// All attributes are required.
	return fmt.Sprintf(
		"Title: \"%s\", Description: \"%s\", Name: \"%s\", Link: \"%s\"",
		ti.Title, ti.Description, ti.Name, ti.Link)
}

// A channel may contain any number of <item>s. An item may represent a
// "story" -- much like a story in a newspaper or magazine; if so its
// description is a synopsis of the story, and the link points to the
// full story. An item may also be complete in itself, if so, the
// description contains the text (entity-encoded HTML is allowed; see
// [examples](https://cyber.harvard.edu/rss/encodingDescriptions.html)),
// and the link and title may be omitted. All elements of an item are
// optional, however at least one of title or description must be present.
type RSSItem struct {
	// The title of the item.
	//
	// Sample:
	//   Venice Film Festival Tries to Quit Sinking
	Title string `xml:"title,omitempty" json:"title,omitempty"`

	// The URL of the item.
	//
	// Sample:
	//   http://nytimes.com/2004/12/07FEST.html
	Link string `xml:"link,omitempty" json:"link,omitempty"`

	// The item synopsis.
	//
	// Sample:
	//   Some of the most heated chatter at the Venice Film Festival this
	//   week was about the way that the arrival of the stars at the
	//   Palazzo del Cinema was being staged.
	Description string `xml:"description,omitempty" json:"description,omitempty"`

	// Email address of the author of the item.
	// [More](https://cyber.harvard.edu/rss/rss.html#ltauthorgtSubelementOfLtitemgt).
	//
	// Sample:
	//   oprah@oxygen.net
	Author string `xml:"author,omitempty" json:"author,omitempty"`

	// Includes the item in one or more categories.
	// [More](https://cyber.harvard.edu/rss/rss.html#ltcategorygtSubelementOfLtitemgt).
	Categories []RSSCategory `xml:"category,omitempty" json:"category,omitempty"`

	// URL of a page for comments relating to the item.
	// [More](https://cyber.harvard.edu/rss/rss.html#ltcommentsgtSubelementOfLtitemgt).
	//
	// Sample:
	//   http://www.myblog.org/cgi-local/mt/mt-comments.cgi?entry_id=290
	Comments string `xml:"comments,omitempty" json:"comments,omitempty"`

	// Describes a media object that is attached to the item.
	// [More](https://cyber.harvard.edu/rss/rss.html#ltenclosuregtSubelementOfLtitemgt).
	Enclosure *RSSEnclosure `xml:"enclosure,omitempty" json:"enclosure,omitempty"`

	// A string that uniquely identifies the item.
	// [More](https://cyber.harvard.edu/rss/rss.html#ltguidgtSubelementOfLtitemgt).
	//
	// Sample:
	//   http://inessential.com/2002/09/01.php#a2
	GUID string `xml:"guid,omitempty" json:"guid,omitempty"`

	// Indicates when the item was published.
	// [More](https://cyber.harvard.edu/rss/rss.html#ltpubdategtSubelementOfLtitemgt).
	//
	// Sample:
	//   Sun, 19 May 2002 15:21:36 GMT
	PubDate *RFC822 `xml:"pubDate,omitempty" json:"pubDate,omitempty"`

	// The RSS channel that the item came from.
	// [More](https://cyber.harvard.edu/rss/rss.html#ltsourcegtSubelementOfLtitemgt).
	//
	// Sample:
	//   <source url="http://www.tomalak.org/links2.xml">Tomalak's Realm</source>
	Source *RSSSource `xml:"source,omitempty" json:"source,omitempty"`
}

func (it RSSItem) String() string {
	// All elements of an item are optional, however at least one of title or description must be present.
	var a []string
	if it.Title != "" {
		a = append(a, "Title: \""+it.Title+"\"")
	}
	if it.Description != "" {
		desc := strings.Replace(it.Description, "\n", "\\n", -1)
		a = append(a, "Description: \""+desc+"\"")
	}
	if it.Link != "" {
		a = append(a, "Link: \""+it.Link+"\"")
	}
	if it.Author != "" {
		a = append(a, "Author: \""+it.Author+"\"")
	}
	if it.Categories != nil {
		var b []string
		for _, ca := range it.Categories {
			b = append(b, ca.String())
		}
		a = append(a, "Category: ["+strings.Join(b, ", ")+"]")
	}
	if it.Comments != "" {
		a = append(a, "Comments: \""+it.Comments+"\"")
	}
	if it.Enclosure != nil {
		a = append(a, "Enclosure: {"+it.Enclosure.String()+"}")
	}
	if it.GUID != "" {
		a = append(a, "GUID: \""+it.GUID+"\"")
	}
	if !it.PubDate.IsZero() {
		a = append(a, "PubDate: "+it.PubDate.String())
	}
	if it.Source != nil {
		a = append(a, "Source: {"+it.Source.String()+"}")
	}

	return strings.Join(a, ", ")
}

// RSSEnclosure is an optional sub-element of RSSItem.
//
// It has three required attributes. url says where the enclosure is
// located, length says how big it is in bytes, and type says what its
// type is, a standard MIME type.
//
// The url must be an http url.
//
// <enclosure url="http://www.scripting.com/mp3s/weatherReportSuite.mp3" length="12216320" type="audio/mpeg" />
//
// A use-case narrative for this element is [here](http://www.thetwowayweb.com/payloadsforrss).
type RSSEnclosure struct {

	/*************************** Required elements ***************************/

	URL    string `xml:"url,attr"    json:"url"`
	Length int    `xml:"length,attr" json:"length"`
	Type   string `xml:"type,attr"   json:"type"`

	/*************************** Optional elements ***************************/

	// No optional element.
}

func (ec RSSEnclosure) String() string {
	// All attributes are required.
	return fmt.Sprintf(
		"URL: \"%s\", Length: %d, Type: \"%s\"", ec.URL, ec.Length, ec.Type)
}

// RSSSource is an optional sub-element of RSSItem.
//
// Its value is the name of the RSSChannel that the item came from,
// derived from its <title>. It has one required attribute, url, which
// links to the XMLization of the source.
//
// <source url="http://www.tomalak.org/links2.xml">Tomalak's Realm</source>
//
// The purpose of this element is to propagate credit for links, to
// publicize the sources of news items. It can be used in the Post
// command of an aggregator. It should be generated automatically when
// forwarding an item from an aggregator to a weblog authoring tool.
type RSSSource struct {

	/*************************** Required elements ***************************/

	Value string `xml:",chardata" json:"value"`
	URL   string `xml:"url,attr"  json:"url"`

	/*************************** Optional elements ***************************/

	// No optional element.
}

func (s RSSSource) String() string {
	// All attributes are required.
	return fmt.Sprintf("\"%s\", URL: \"%s\"", s.Value, s.URL)
}

type RFC822 time.Time

var rfc822layout = [2]string{
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05 -0700",
}

// UnmarshalXML implements the xml.Unmarshal interface.
func (r *RFC822) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v, layout string
	var t time.Time
	var err error
	d.DecodeElement(&v, &start)
	for _, layout = range rfc822layout {
		t, err = time.Parse(layout, v)
		if err == nil {
			*r = RFC822(t)
			return nil
		}
	}
	return err
}

// MarshalJSON implements the json.Marshal interface.
func (r *RFC822) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// IsZero reports whether r represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (r RFC822) IsZero() bool { return time.Time(r).IsZero() }

func (r RFC822) String() string { return time.Time(r).Format(time.RFC3339) }

// After reports whether the RFC822 instant r is after t.
func (r RFC822) After(t *RFC822) bool { return time.Time(r).After(time.Time(*t)) }
