package main

import (
	"fmt"
	"sync"
	"github.com/chromedp/chromedp"
	"context"
	"time"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type Fetcher interface {
	Fetch(url string, instructions chromedp.Tasks) (body string, err error)
	FetchSeasons(f Fetcher, wg *sync.WaitGroup) ([]string)
	ParseList(html string, selector string) ([]string, error)
	IsCached(url string) (cachedAt int)
}

type PageFetcher struct {}

func (PageFetcher) Fetch(
	url string,
	instructions chromedp.Tasks,
) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var html string

	// Base tasks: navigate and wait for page
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
	}

	// Add user instructions if provided
	if instructions != nil {
		tasks = append(tasks, instructions...)
	}

	// Always get the HTML at the end
	tasks = append(tasks, chromedp.OuterHTML("html", &html))

	if err := chromedp.Run(ctx, tasks); err != nil {
		return "", err
	}

	return html, nil
}

func (PageFetcher) FetchSeasons(f Fetcher, wg *sync.WaitGroup) (years []string) {
	content := fmt.Sprintf(
		f.Fetch(
			"https://www.nrl.com/draw/?competition=111&round=22&season=2025",
			chromedp.Tasks{
					chromedp.WaitVisible(`[aria-controls="season-dropdown"]`, chromedp.ByQuery),
					chromedp.Click(`[aria-controls="season-dropdown"]`, chromedp.ByQuery),
					chromedp.Sleep(2 * time.Second),
			},
		))

		years, err := f.ParseList(content, "#season-dropdown li button div")
		if err != nil {
			panic(err)
		}

		fmt.Println("Years found:", years)
		return
}

func (PageFetcher) ParseList(html string, selector string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var results []string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			results = append(results, text)
		}
	})

	return results, nil
}

func (PageFetcher) IsCached(url string) (int) {
	return 0
}

func main() {
	var wg sync.WaitGroup
	fetcher := PageFetcher{}
	wg.Add(1)
	go Scrape(fetcher, &wg)
	wg.Wait()
}

func Scrape(f Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()
	seasons := f.FetchSeasons(f, wg)
	
	for _, v := range seasons {
		fmt.Println(v)
	}
}