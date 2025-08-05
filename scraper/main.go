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

type Season struct {
	year string
	rounds []*Round
}

type Competition struct {
	seasons map[string]*Season
}

type Fetcher interface {
	Fetch(url string, instructions chromedp.Tasks, require bool) (body string, err error)
	FetchSeasons(wg *sync.WaitGroup) ([]string)
	ParseList(html string, selector string) ([]string, error)
	IsCached(url string) (cachedAt int)
}

type PageFetcher struct {
	allocCtx context.Context
	browserCtx context.Context
	semaphore chan struct{}
}

func createListStr[T fmt.Stringer](list []T) string {
	var sb strings.Builder
	sb.WriteString("[")

	for i, v := range list {
		sb.WriteString(v.String())
		if i < len(list)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString("]")
	return sb.String()
}

func (c Competition) String() string {
	result := "{"
	i := 0
	for _, round := range c.seasons {
		result += round.String()
		if i < len(c.seasons)-1 {
			result += ","
		}
		i++
	}

	return result + "}"
}

func (s *Season) String() string {
	return fmt.Sprintf("\"%s\": %s", s.year, createListStr(s.rounds))
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
	require bool,
) (string, error) {
	p.semaphore <- struct{}{}
	defer func() { <-p.semaphore }() 

	ctx, browserCancel := chromedp.NewContext(p.browserCtx)
	defer browserCancel()

	ctx, timeoutCancel:= context.WithTimeout(ctx, 60*time.Second)
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

	err := chromedp.Run(ctx, tasks)
	for i := 0; i < 4; i++ {
		if err == nil {
			return html, nil
		} else if !require {
			return "", err
		}

		err = chromedp.Run(ctx, tasks);
	}

	return "", err
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
			true,
	))

	years, err := pf.ParseList(content, "#season-dropdown li button div")
	if err != nil {
		panic(err)
	}

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

	writeToFile(Scrape(fetcher, &wg).String(), "out.json")
}

func Scrape(f Fetcher, wg *sync.WaitGroup) (comp Competition) {
	seasons := f.FetchSeasons(wg)
	comp = Competition{
		seasons: map[string]*Season{},
	}

	stats := &StatsTracker{}

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Printf("Started: %d, Finished: %d, Active: %d\n",
					stats.Started(), stats.Finished(), stats.Active())
				// writeToFile(comp.String(), "out.json") 
			case <-done:
				return
			}
		}
	}()

	
	for _, s := range seasons {
		if s == "2016" {
			wg.Add(1)
			season := &Season{
				year:   s,
				rounds: []*Round{},
			}
		
			comp.seasons[s] = season
		
			go ScrapeSeason(season, f, wg, stats)
		}
	}

	wg.Wait()
	close(done)
	fmt.Println("All jobs complete.")
	return
}

func ScrapeSeason(season *Season, f Fetcher, wg *sync.WaitGroup, stats *StatsTracker) {
	defer wg.Done()
	stats.Start()
	defer stats.Finish()

	content := fmt.Sprintf(
		f.Fetch(
			fmt.Sprintf("https://www.nrl.com/draw/?competition=111&round=1&season=%s", season.year),
			chromedp.Tasks{
				chromedp.WaitVisible(`[aria-controls="round-dropdown"]`, chromedp.ByQuery),
				chromedp.Click(`[aria-controls="round-dropdown"]`, chromedp.ByQuery),
				chromedp.Sleep(2 * time.Second),
			},
			true,
	))

	rounds, err := f.ParseList(content, "#round-dropdown li button div")
	if err != nil {
		panic(err)
	}

	Reverse(rounds)
	wg.Add(1)
	scrapeRounds(rounds, season, f, wg, stats)
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