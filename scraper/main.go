package main

import (
	"fmt"
	"os"
	"sync"
	"github.com/chromedp/chromedp"
	"context"
	"time"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type Seasons struct {
	year string
	rounds []Round
}

type Competition struct {
	seasons []Seasons
}

type Fetcher interface {
	Fetch(url string, instructions chromedp.Tasks) (body string, err error)
	FetchSeasons(wg *sync.WaitGroup) ([]string)
	ParseList(html string, selector string) ([]string, error)
	IsCached(url string) (cachedAt int)
}

type PageFetcher struct {
	allocCtx context.Context
	browserCtx context.Context
	semaphore chan struct{}
}

func NewPageFetcher(maxConcurrent int) (*PageFetcher, error) {
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),                          // run in headless mode
			chromedp.Flag("disable-gpu", true),                       // disable GPU
			chromedp.Flag("blink-settings", "imagesEnabled=false"),   // disables images
			chromedp.Flag("mute-audio", true),                        // mutes audio (optional)
		)...,
	)

	browserCtx, cancel := chromedp.NewContext(allocCtx)
	// Optional: run something small to ensure it boots
	if err := chromedp.Run(browserCtx); err != nil {
		allocCancel()
		cancel()
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	return &PageFetcher{
		allocCtx: allocCtx,
		browserCtx: browserCtx,
		semaphore: make(chan struct{}, maxConcurrent),
	}, nil
}

func (p *PageFetcher) Fetch(
	url string,
	instructions chromedp.Tasks,
) (string, error) {
	p.semaphore <- struct{}{}
	defer func() { <-p.semaphore }() 

	ctx, browserCancel := chromedp.NewContext(p.browserCtx)
	defer browserCancel()

	ctx, timeoutCancel:= context.WithTimeout(ctx, 30*time.Second)
	defer timeoutCancel()

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

func (pf PageFetcher) FetchSeasons(wg *sync.WaitGroup) (years []string) {
	content := fmt.Sprintf(
		pf.Fetch(
			"https://www.nrl.com/draw/?competition=111&round=1&season=2025",
			chromedp.Tasks{
					chromedp.WaitVisible(`[aria-controls="season-dropdown"]`, chromedp.ByQuery),
					chromedp.Click(`[aria-controls="season-dropdown"]`, chromedp.ByQuery),
					chromedp.Sleep(2 * time.Second),
			},
	))

	years, err := pf.ParseList(content, "#season-dropdown li button div")
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
	fetcher, err := NewPageFetcher(10)
	
	if err != nil {
		fmt.Println("unable to creater page fatcher", err)
		return
	}

	wg.Add(1)
	go Scrape(fetcher, &wg)
	wg.Wait()
}

func Scrape(f Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()
	seasons := f.FetchSeasons(wg)
	
	for _, s := range seasons {
		wg.Add(1)
		go ScrapeSeason(s, f, wg)
	}
}

func ScrapeSeason(season string, f Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()

	content := fmt.Sprintf(
		f.Fetch(
			fmt.Sprintf("https://www.nrl.com/draw/?competition=111&round=1&season=%s", season),
			chromedp.Tasks{
				chromedp.WaitVisible(`[aria-controls="round-dropdown"]`, chromedp.ByQuery),
				chromedp.Click(`[aria-controls="round-dropdown"]`, chromedp.ByQuery),
				chromedp.Sleep(2 * time.Second),
			},
	))

	rounds, err := f.ParseList(content, "#round-dropdown li button div")
	if err != nil {
		panic(err)
	}

	Reverse(rounds)
	wg.Add(1)
	scrapeRounds(rounds, season, f, wg)
	return
}

func writeToFile (content string, fileName string) {
	// Create or truncate the file
	file, err := os.Create(fileName)
	if err != nil {
			fmt.Println("Error creating file:", err)
			return
	}
	defer file.Close()

	// Write to the file
	_, err = file.WriteString(content)
	if err != nil {
			fmt.Println("Error writing to file:", err)
			return
	}

	fmt.Println("File written successfully.")
}

func Reverse[T any](arr []T) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}