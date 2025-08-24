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
  "github.com/google/uuid"
)

type Season struct {
	id uuid.UUID
	year string
	rounds []*Round
}

type Competition struct {
	id int
	name string
	seasons []*Season
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
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(os.Getenv("CHOMEDP_CHROME_PATH")),
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("blink-settings", "imagesEnabled=false"),
			chromedp.Flag("mute-audio", true),
		)...,
	)

	browserCtx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithErrorf(func(string, ...any) {}),
    chromedp.WithDebugf(func(string, ...any) {}),
    chromedp.WithLogf(func(string, ...any) {}),
	)
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

	chromedp.ListenTarget(ctx, func(ev interface{}) {})

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
			fmt.Sprintf("https://www.nrl.com/draw/?competition=111&round=1&season=2025"),
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
	/*var wg sync.WaitGroup
	fetcher, err := NewPageFetcher(10)
	
	if err != nil {
		fmt.Println("unable to creater page fatcher", err)
		return
	}

	Scrape(111, fetcher, &wg)*/

	db, err := NewDB()
	if err != nil {
		panic(err)
	}

	defer db.Conn.Close()
	comp, _ := db.GetCompetition(111)
	writeToFile(fmt.Sprint(comp), "/app/output/results.json")
}

func Scrape(compID int, f Fetcher, wg *sync.WaitGroup) (comp Competition) {
	db, err := NewDB()
	if err != nil {
		panic(err)
	}
	defer db.Conn.Close() 

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.CreateCompIfNotExist(ctx, "Mens NRL Premiership", compID)
	
	seasons := f.FetchSeasons(wg)

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
			case <-done:
				return
			}
		}
	}()

	for _, s := range seasons {
		wg.Add(1)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		seasonID, err := db.CreateSeasonIfNotExist(ctx, compID, s)
		if err != nil {
			wg.Done() 
			continue
		}

		go ScrapeSeason(compID, s, seasonID, f, wg, stats)
	}

	wg.Wait()
	close(done)
	fmt.Println("All jobs complete.")
	return
}

func ScrapeSeason(compID int, season string, seasonID uuid.UUID, f Fetcher, wg *sync.WaitGroup, stats *StatsTracker) {
	defer wg.Done()
	stats.Start()
	defer stats.Finish()

	if season != "2024" {
		return
	}

	content := fmt.Sprintf(
		f.Fetch(
			fmt.Sprintf("https://www.nrl.com/draw/?competition=%d&round=1&season=%s", compID, season),
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
	go scrapeRounds(rounds, season, seasonID, compID, f, wg, stats)
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