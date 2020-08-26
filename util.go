package hlsm

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func readParseHTMLBody(reader io.Reader) (selector *goquery.Selection, title string, err error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return
	}

	mainContainer := doc.Find("#mainContent")
	title = mainContainer.Find(".pagetitle > h2").Text()
	selector = mainContainer.Find(".post")

	if strings.ToLower(strings.TrimSpace(title)) == "kesalahan" {
		err = fmt.Errorf("Terjadi kesalahan: %s", selector.Find("p").Text())
	}

	return
}
