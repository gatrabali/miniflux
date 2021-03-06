// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package rewrite // import "miniflux.app/reader/rewrite"

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	youtubeRegex  = regexp.MustCompile(`youtube\.com/watch\?v=(.*)`)
	imgRegex      = regexp.MustCompile(`<img [^>]+>`)
	textLinkRegex = regexp.MustCompile(`(?mi)(\bhttps?:\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])`)
)

func addImageTitle(entryURL, entryContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(entryContent))
	if err != nil {
		return entryContent
	}

	matches := doc.Find("img[src][title]")

	if matches.Length() > 0 {
		matches.Each(func(i int, img *goquery.Selection) {
			altAttr := img.AttrOr("alt", "")
			srcAttr, _ := img.Attr("src")
			titleAttr, _ := img.Attr("title")

			img.ReplaceWithHtml(`<figure><img src="` + srcAttr + `" alt="` + altAttr + `"/><figcaption><p>` + titleAttr + `</p></figcaption></figure>`)
		})

		output, _ := doc.Find("body").First().Html()
		return output
	}

	return entryContent
}

func addDynamicImage(entryURL, entryContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(entryContent))
	if err != nil {
		return entryContent
	}

	// Ordered most preferred to least preferred.
	candidateAttrs := []string{
		"data-src",
		"data-original",
		"data-orig",
		"data-url",
		"data-orig-file",
		"data-large-file",
		"data-medium-file",
		"data-2000src",
		"data-1000src",
		"data-800src",
		"data-655src",
		"data-500src",
		"data-380src",
	}

	changed := false

	doc.Find("img,div").Each(func(i int, img *goquery.Selection) {
		for _, candidateAttr := range candidateAttrs {
			if srcAttr, found := img.Attr(candidateAttr); found {
				changed = true

				if img.Is("img") {
					img.SetAttr("src", srcAttr)
				} else {
					altAttr := img.AttrOr("alt", "")
					img.ReplaceWithHtml(`<img src="` + srcAttr + `" alt="` + altAttr + `"/>`)
				}

				break
			}
		}
	})

	if !changed {
		doc.Find("noscript").Each(func(i int, noscript *goquery.Selection) {
			matches := imgRegex.FindAllString(noscript.Text(), 2)

			if len(matches) == 1 {
				changed = true

				noscript.ReplaceWithHtml(matches[0])
			}
		})
	}

	if changed {
		output, _ := doc.Find("body").First().Html()
		return output
	}

	return entryContent
}

func addYoutubeVideo(entryURL, entryContent string) string {
	matches := youtubeRegex.FindStringSubmatch(entryURL)

	if len(matches) == 2 {
		video := `<iframe width="650" height="350" frameborder="0" src="https://www.youtube-nocookie.com/embed/` + matches[1] + `" allowfullscreen></iframe>`
		return video + "<p>" + replaceLineFeeds(replaceTextLinks(entryContent)) + "</p>"
	}
	return entryContent
}

func addPDFLink(entryURL, entryContent string) string {
	if strings.HasSuffix(entryURL, ".pdf") {
		return fmt.Sprintf(`<a href="%s">PDF</a><br>%s`, entryURL, entryContent)
	}
	return entryContent
}

func replaceTextLinks(input string) string {
	return textLinkRegex.ReplaceAllString(input, `<a href="${1}">${1}</a>`)
}

func replaceLineFeeds(input string) string {
	return strings.Replace(input, "\n", "<br>", -1)
}

// -- Gatra Bali specific rewriter functions -- //

// hideFirstImage replaces the first image found on body with span tag '<span data-minifux-enclosure=""/>'
// Before the content displayed, we can use the 'data-minifux-enclosure' value as an enclosure object
func hideFirstImage(entryURL, entryContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(entryContent))
	if err != nil {
		return entryContent
	}

	matches := doc.Find("img")

	if matches.Length() > 0 {
		// we only need to hide the first image
		img := matches.First()
		srcAttr, _ := img.Attr("src")
		img.ReplaceWithHtml(`<span data-miniflux-enclosure="` + srcAttr + `"/>`)

		output, _ := doc.Find("body").First().Html() // the whole output
		return output
	}

	return entryContent
}

func cleanupBacaJuga(s *goquery.Selection) bool {
	// if element has class 'IRRP_kangoo'
	if s.HasClass("IRRP_kangoo") {
		s.Remove()
		return true
	}

	// If text contains 'baca juga'
	text := strings.ToLower(s.Text())
	if strings.Contains(text, "baca juga") {
		s.Parent().Remove()
		return true
	}
	return false
}

func cleanupBaliPost(entryURL, entryContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(entryContent))
	if err != nil {
		return entryContent
	}

	changed := false

	// Remove 'Baca Juga' Links
	bacaJuga := doc.Find("span")
	bacaJuga.Each(func(i int, bj *goquery.Selection) {
		removed := cleanupBacaJuga(bj)
		if removed {
			changed = true
		}
	})

	if changed {
		output, _ := doc.Find("body").First().Html()
		return output
	}
	return entryContent
}

func cleanupMetroBali(entryURL, entryContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(entryContent))
	if err != nil {
		return entryContent
	}

	changed := false

	// Remove 'Baca Juga' Links
	bacaJuga := doc.Find("a")
	bacaJuga.Each(func(i int, bj *goquery.Selection) {
		removed := cleanupBacaJuga(bj)
		if removed {
			changed = true
		}
	})

	// Remove Related Posts
	relatedPostSectionHeader := doc.Find("h3")
	relatedPostSectionHeader.Each(func(i int, h3 *goquery.Selection) {
		if h3.Text() == "Related Posts" {

			// remove all elements after '<h3>Related Posts</h3>'
			nexts := h3.NextAll()
			nexts.Each(func(i int, next *goquery.Selection) {
				next.Remove()
			})

			// remove the h3 itself
			h3.Remove()
			changed = true
		}
	})

	// Remove Ad Links
	ad := doc.Find(".advertising_content_single")
	ad.Each(func(i int, ad *goquery.Selection) {
		ad.Remove()
		changed = true
	})

	if changed {
		output, _ := doc.Find("body").First().Html()
		return output
	}
	return entryContent
}

func cleanupBaliPuspaNews(entryURL, entryContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(entryContent))
	if err != nil {
		return entryContent
	}

	changed := false

	// Remove Ads
	ads := doc.Find(".td-all-devices")
	ads.Each(func(i int, ad *goquery.Selection) {
		ad.Remove()
		changed = true
	})

	if changed {
		output, _ := doc.Find("body").First().Html()
		return output
	}
	return entryContent
}
