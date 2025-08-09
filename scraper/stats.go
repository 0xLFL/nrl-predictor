package main

import (
	"sync"
	"strings"
	"strconv"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type MatchStats struct {
	posAndComp *PosAndComp
	attack *Attack
	passing *Passing
	kicking *Kicking
	defence *Defence
	negPlays *NegPlays
}

type PosAndComp struct {
	homePosPer int
	awayPosPer int

	homePosTime string
	awayPosTime string

	homeSets int
	homeSetsCompleated int

	awaySets int
	awaySetsCompleated int
}

type Attack struct {
	homeRuns int
	awayRuns int

	homeRunMeters int
	awayRunMeters int

	homePostContactMeters int
	awayPostContactMeters int

	homeLineBreaks int
	awayLineBreaks int

	homeTackleBreaks int
	awayTackleBreaks int

	homeAvgSetDistance float64
	awayAvgSetDistance float64

	homeKickReturnMeters int
	awayKickReturnMeters int

	homeAvgPlayTheBallSpeed float64
	awayAvgPlayTheBallSpeed float64
}

type Passing struct {
	homeOffloads int
	awayOffloads int

	homeReceipts int
	awayReceipts int

	homeTotalPasses int
	awayTotalPasses int

	homeDummyPasses int
	awayDummyPasses int
}

type Kicking struct {
	homeKicks int
	awayKicks int

	homeKickingMeters int
	awayKickingMeters int

	homeForcedDropOuts int
	awayForcedDropOuts int

	homeKickDefusal int
	awayKickDefusal int

	homeBombs int
	awayBombs int

	homeGrubbers int
	awayGrubbers int
}

type Defence struct {
	homeEffecTackle float64
	awayEffecTackle float64

	homeTacklesMade int
	awayTacklesMade int

	homeMissedTackles int
	awayMissedTackles int

	homeIntercepts int
	awayIntercepts int

	homeIneffecTackles int
	awayIneffecTackles int
}

type NegPlays struct {
	homeErrors int
	awayErrors int

	homePenCon int
	awayPenCon int

	homeRuckInf int
	awayRuckInf int

	homeInside10 int
	awayInside10 int

	homeOnReport int
	awayOnReport int
}

func (ms *MatchStats) String() string {
	return fmt.Sprintf(`{
			"posAndComp": %s,
			"attack": %s,
			"passing": %s,
			"kicking": %s,
			"defence": %s,
			"negPlays": %s
		}`,
		ms.posAndComp,
		ms.attack,
		ms.passing,
		ms.kicking,
		ms.defence,
		ms.negPlays,
	)
}

func (pc *PosAndComp) String() string {
	return fmt.Sprintf(`{
			"homePosPer": %d,
			"awayPosPer": %d,
			"homePosTime": "%s",
			"awayPosTime": "%s",
			"homeSets": %d,
			"homeSetsCompleated": %d,
			"awaySets": %d,
			"awaySetsCompleated": %d
		}`,
		pc.homePosPer,
		pc.awayPosPer,
		pc.homePosTime,
		pc.awayPosTime,
		pc.homeSets,
		pc.homeSetsCompleated,
		pc.awaySets,
		pc.awaySetsCompleated,
	)
}

func (a *Attack) String() string {
	return fmt.Sprintf(`{
			"homeRuns": %d,
			"awayRuns": %d,
			"homeRunMeters": %d,
			"awayRunMeters": %d,
			"homePostContactMeters": %d,
			"awayPostContactMeters": %d,
			"homeLineBreaks": %d,
			"awayLineBreaks": %d,
			"homeTackleBreaks": %d,
			"awayTackleBreaks": %d,
			"homeAvgSetDistance": %f,
			"awayAvgSetDistance": %f,
			"homeKickReturnMeters": %d,
			"awayKickReturnMeters": %d,
			"homeAvgPlayTheBallSpeed": %f,
			"awayAvgPlayTheBallSpeed": %f
		}`,
		a.homeRuns,
		a.awayRuns,
		a.homeRunMeters,
		a.awayRunMeters,
		a.homePostContactMeters,
		a.awayPostContactMeters,
		a.homeLineBreaks,
		a.awayLineBreaks,
		a.homeTackleBreaks,
		a.awayTackleBreaks,
		a.homeAvgSetDistance,
		a.awayAvgSetDistance,
		a.homeKickReturnMeters,
		a.awayKickReturnMeters,
		a.homeAvgPlayTheBallSpeed,
		a.awayAvgPlayTheBallSpeed,
	)
}

func (p *Passing) String() string {
	return fmt.Sprintf(`{
			"homeOffloads": %d,
			"awayOffloads": %d,
			"homeReceipts": %d,
			"awayReceipts": %d,
			"homeTotalPasses": %d,
			"awayTotalPasses": %d,
			"homeDummyPasses": %d,
			"awayDummyPasses": %d
		}`,
		p.homeOffloads,
		p.awayOffloads,
		p.homeReceipts,
		p.awayReceipts,
		p.homeTotalPasses,
		p.awayTotalPasses,
		p.homeDummyPasses,
		p.awayDummyPasses,
	)
}

func (k *Kicking) String() string {
	return fmt.Sprintf(`{
			"homeKicks": %d,
			"awayKicks": %d,
			"homeKickingMeters": %d,
			"awayKickingMeters": %d,
			"homeForcedDropOuts": %d,
			"awayForcedDropOuts": %d,
			"homeKickDefusal": %d,
			"awayKickDefusal": %d,
			"homeBombs": %d,
			"awayBombs": %d,
			"homeGrubbers": %d,
			"awayGrubbers": %d
		}`,
		k.homeKicks,
		k.awayKicks,
		k.homeKickingMeters,
		k.awayKickingMeters,
		k.homeForcedDropOuts,
		k.awayForcedDropOuts,
		k.homeKickDefusal,
		k.awayKickDefusal,
		k.homeBombs,
		k.awayBombs,
		k.homeGrubbers,
		k.awayGrubbers,
	)
}

func (d *Defence) String() string {
	return fmt.Sprintf(`{
			"homeEffecTackle": %f,
			"awayEffecTackle": %f,
			"homeTacklesMade": %d,
			"awayTacklesMade": %d,
			"homeMissedTackles": %d,
			"awayMissedTackles": %d,
			"homeIntercepts": %d,
			"awayIntercepts": %d,
			"homeIneffecTackles": %d,
			"awayIneffecTackles": %d
		}`,
		d.homeEffecTackle,
		d.awayEffecTackle,
		d.homeTacklesMade,
		d.awayTacklesMade,
		d.homeMissedTackles,
		d.awayMissedTackles,
		d.homeIntercepts,
		d.awayIntercepts,
		d.homeIneffecTackles,
		d.awayIneffecTackles,
	)
}

func (np *NegPlays) String() string {
	return fmt.Sprintf(`{
			"homeErrors": %d,
			"awayErrors": %d,
			"homePenCon": %d,
			"awayPenCon": %d,
			"homeRuckInf": %d,
			"awayRuckInf": %d,
			"homeInside10": %d,
			"awayInside10": %d,
			"homeOnReport": %d,
			"awayOnReport": %d
		}`,
		np.homeErrors,
		np.awayErrors,
		np.homePenCon,
		np.awayPenCon,
		np.homeRuckInf,
		np.awayRuckInf,
		np.homeInside10,
		np.awayInside10,
		np.homeOnReport,
		np.awayOnReport,
	)
}

func parseMatchStats(stats *MatchStats, content string, wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(1)
	go parsePosAndCompStats(stats.posAndComp, content, wg)

	barChartHandlers := parseAttackStats(stats.attack, content, wg)
	MergeInto(barChartHandlers, parsePassingStats(stats.passing, content, wg))
	MergeInto(barChartHandlers, parseKickingStats(stats.kicking, content, wg))
	MergeInto(barChartHandlers, parseDefenceStats(stats.defence, content, wg))
	MergeInto(barChartHandlers, parseNegPlayStats(stats.negPlays, content, wg))

	wg.Add(1)
	go parseBarChart(barChartHandlers, content, wg)
}

func parsePosAndCompStats(stats *PosAndComp, content string, wg *sync.WaitGroup) {
	defer wg.Done()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return;
	}

	doc.Find(".match-centre-card-donut__value--home").Each(func(_ int, s *goquery.Selection) {
		posStr := strings.TrimSpace(s.Text())
		stats.homePosPer, _ = strconv.Atoi(strings.TrimSuffix(posStr, "%"))
	})

	doc.Find(".match-centre-card-donut__value--away").Each(func(_ int, s *goquery.Selection) {
		posStr := strings.TrimSpace(s.Text())
		stats.awayPosPer, _ = strconv.Atoi(strings.TrimSuffix(posStr, "%"))
	})

	doc.Find("figure.stats-bar-chart").EachWithBreak(func(i int, s *goquery.Selection) bool {
    title := s.Find("figcaption.stats-bar-chart__title").Text()
    title = strings.TrimSpace(title)

    if title == "Time In Possession" {
			stats.homePosTime = strings.TrimSpace(s.Find("dd.stats-bar-chart__label--home").Text())
			stats.awayPosTime = strings.TrimSpace(s.Find("dd.stats-bar-chart__label--away").Text())

			return false
    }

		return true
	})

	doc.Find(".u-spacing-pb-24.u-spacing-pt-16.u-width-100").EachWithBreak(func(i int, s *goquery.Selection) bool {
    title := s.Find("h3.stats-bar-chart__title").Text()
    title = strings.TrimSpace(title)

    if title == "Completion Rate" {
			compRate := s.Find(".match-centre-card-donut__value.match-centre-card-donut__value--footer")
			homeCompRate_ := strings.TrimSpace(compRate.Eq(0).Text())
			awayCompRate_ := strings.TrimSpace(compRate.Eq(1).Text())

			homeCompRate := strings.Split(homeCompRate_, "/")
			awayCompRate := strings.Split(awayCompRate_, "/")

			stats.homeSets, _ = strconv.Atoi(homeCompRate[1])
			stats.homeSetsCompleated, _ = strconv.Atoi(homeCompRate[0])

			stats.awaySets, _ = strconv.Atoi(awayCompRate[1])
			stats.awaySetsCompleated, _ = strconv.Atoi(awayCompRate[0])

			return false
    }

		return true
	})
}

func parseAttackStats(a *Attack, content string, wg *sync.WaitGroup) (handlers map[string]func(string, string)) {
	handlers = map[string]func(string, string) {
		"All Runs": func(homeStr string, awayStr string) {
			a.homeRuns, _ = strconv.Atoi(homeStr)
			a.awayRuns, _ = strconv.Atoi(awayStr)
		},
		"All Run Metres": func(homeStr string, awayStr string) {
			a.homeRunMeters, _ = strconv.Atoi(homeStr)
			a.awayRunMeters, _ = strconv.Atoi(awayStr)
		},
		"Post Contact Metres": func(homeStr string, awayStr string) {
			a.homePostContactMeters, _ = strconv.Atoi(homeStr)
			a.awayPostContactMeters, _ = strconv.Atoi(awayStr)
		},
		"Line Breaks": func(homeStr string, awayStr string) {
			a.homeLineBreaks, _ = strconv.Atoi(homeStr)
			a.awayLineBreaks, _ = strconv.Atoi(awayStr)
		},
		"Tackle Breaks": func(homeStr string, awayStr string) {
			a.homeTackleBreaks, _ = strconv.Atoi(homeStr)
			a.awayTackleBreaks, _ = strconv.Atoi(awayStr)
		},
		"Average Set Distance": func(homeStr string, awayStr string) {
			a.homeAvgSetDistance, _ = strconv.ParseFloat(homeStr, 64)
			a.awayAvgSetDistance, _ = strconv.ParseFloat(awayStr, 64)
		},
		"Kick Return Metres": func(homeStr string, awayStr string) {
			a.homeKickReturnMeters, _ = strconv.Atoi(homeStr)
			a.awayKickReturnMeters, _ = strconv.Atoi(awayStr)
		},
	}

	wg.Add(1)
	go func () {
		defer wg.Done()

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			return;
		}
	
		doc.Find(".u-spacing-pb-24.u-spacing-pt-16.u-width-100").EachWithBreak(func(i int, s *goquery.Selection) bool {
			title := s.Find("h3.stats-bar-chart__title").Text()
			title = strings.TrimSpace(title)

			if title == "Average Play The Ball Speed" {
				compRate := s.Find(".donut-chart-stat__value > span > span:not(.donut-chart__unit)")
				homeStr := strings.TrimSpace(compRate.Eq(0).Text())
				awayStr := strings.TrimSpace(compRate.Eq(1).Text())

				a.homeAvgPlayTheBallSpeed, _ = strconv.ParseFloat(homeStr, 64)
				a.awayAvgPlayTheBallSpeed, _ = strconv.ParseFloat(awayStr, 64)

				return false
			}

			return true
		})
	}()

	return
}

func parsePassingStats(p *Passing, content string, wg *sync.WaitGroup) (handlers map[string]func(string, string)) {
	handlers = map[string]func(string, string){
		"Offloads": func(homeStr string, awayStr string) {
			p.homeOffloads, _ = strconv.Atoi(homeStr)
			p.awayOffloads, _ = strconv.Atoi(awayStr)
		},
		"Receipts": func(homeStr string, awayStr string) {
			p.homeReceipts, _ = strconv.Atoi(homeStr)
			p.awayReceipts, _ = strconv.Atoi(awayStr)
		},
		"Total Passes": func(homeStr string, awayStr string) {
			p.homeTotalPasses, _ = strconv.Atoi(homeStr)
			p.awayTotalPasses, _ = strconv.Atoi(awayStr)
		},
		"Dummy Passes": func(homeStr string, awayStr string) {
			p.homeDummyPasses, _ = strconv.Atoi(homeStr)
			p.awayDummyPasses, _ = strconv.Atoi(awayStr)
		},
	}

	return
}

func parseKickingStats(k *Kicking, content string, wg *sync.WaitGroup) (handlers map[string]func(string, string)) {
	handlers = map[string]func(string, string){
		"Kicks": func(homeStr string, awayStr string) {
			k.homeKicks, _ = strconv.Atoi(homeStr)
			k.awayKicks, _ = strconv.Atoi(awayStr)
		},
		"Kicking Metres": func(homeStr string, awayStr string) {
			k.homeKickingMeters, _ = strconv.Atoi(homeStr)
			k.awayKickingMeters, _ = strconv.Atoi(awayStr)
		},
		"Forced Drop Outs": func(homeStr string, awayStr string) {
			k.homeForcedDropOuts, _ = strconv.Atoi(homeStr)
			k.awayForcedDropOuts, _ = strconv.Atoi(awayStr)
		},
		"Bombs": func(homeStr string, awayStr string) {
			k.homeBombs, _ = strconv.Atoi(homeStr)
			k.awayBombs, _ = strconv.Atoi(awayStr)
		},
		"Grubbers": func(homeStr string, awayStr string) {
			k.homeGrubbers, _ = strconv.Atoi(homeStr)
			k.awayGrubbers, _ = strconv.Atoi(awayStr)
		},
	}

	wg.Add(1)
	go func () {
		defer wg.Done()

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			return;
		}
	
		doc.Find(".u-spacing-pb-24.u-spacing-pt-16.u-width-100").EachWithBreak(func(i int, s *goquery.Selection) bool {
			title := s.Find("h3.stats-bar-chart__title").Text()
			title = strings.TrimSpace(title)

			if title == "Kick Defusal %" {
				compRate := s.Find(".donut-chart-stat__value > span > span:not(.donut-chart-stat__value--sup)")
				homeStr := strings.TrimSpace(compRate.Eq(0).Text())
				awayStr := strings.TrimSpace(compRate.Eq(1).Text())

				k.homeKickDefusal, _ = strconv.Atoi(homeStr)
				k.awayKickDefusal, _ = strconv.Atoi(awayStr)

				return false
			}

			return true
		})
	}()

	return
}

func parseDefenceStats(d *Defence, content string, wg *sync.WaitGroup) (handlers map[string]func(string, string)) {
	handlers = map[string]func(string, string){
		"Tackles Made": func(homeStr string, awayStr string) {
			d.homeTacklesMade, _ = strconv.Atoi(homeStr)
			d.awayTacklesMade, _ = strconv.Atoi(awayStr)
		},
		"Missed Tackles": func(homeStr string, awayStr string) {
			d.homeMissedTackles, _ = strconv.Atoi(homeStr)
			d.awayMissedTackles, _ = strconv.Atoi(awayStr)
		},
		"Ineffective Tackles": func(homeStr string, awayStr string) {
			d.homeIneffecTackles, _ = strconv.Atoi(homeStr)
			d.awayIneffecTackles, _ = strconv.Atoi(awayStr)
		},
		"Intercepts": func(homeStr string, awayStr string) {
			d.homeIntercepts, _ = strconv.Atoi(homeStr)
			d.awayIntercepts, _ = strconv.Atoi(awayStr)
		},
	}

	wg.Add(1)
	go func () {
		defer wg.Done()

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			return;
		}
	
		doc.Find(".u-spacing-pb-24.u-spacing-pt-16.u-width-100").EachWithBreak(func(i int, s *goquery.Selection) bool {
			title := s.Find("h3.stats-bar-chart__title").Text()
			title = strings.TrimSpace(title)

			if title == "Effective Tackle %" {
				compRate := s.Find(".donut-chart-stat__value > span > span:not(.donut-chart-stat__value--sup)")
				homeStr := strings.TrimSpace(compRate.Eq(0).Text())
				awayStr := strings.TrimSpace(compRate.Eq(1).Text())

				d.homeEffecTackle, _ = strconv.ParseFloat(homeStr, 64)
				d.awayEffecTackle, _ = strconv.ParseFloat(awayStr, 64)

				return false
			}

			return true
		})
	}()

	return
}

func parseNegPlayStats(ng *NegPlays, content string, wg *sync.WaitGroup) (handlers map[string]func(string, string)) {
	handlers = map[string]func(string, string){
		"Errors": func(homeStr string, awayStr string) {
			ng.homeErrors, _ = strconv.Atoi(homeStr)
			ng.awayErrors, _ = strconv.Atoi(awayStr)
		},
		"Penalties Conceded": func(homeStr string, awayStr string) {
			ng.homePenCon, _ = strconv.Atoi(homeStr)
			ng.awayPenCon, _ = strconv.Atoi(awayStr)
		},
		"Ruck Infringements": func(homeStr string, awayStr string) {
			ng.homeRuckInf, _ = strconv.Atoi(homeStr)
			ng.awayRuckInf, _ = strconv.Atoi(awayStr)
		},
		"Inside 10 Metres": func(homeStr string, awayStr string) {
			ng.homeInside10, _ = strconv.Atoi(homeStr)
			ng.awayInside10, _ = strconv.Atoi(awayStr)
		},
		"On Reports": func(homeStr string, awayStr string) {
			ng.homeOnReport, _ = strconv.Atoi(homeStr)
			ng.awayOnReport, _ = strconv.Atoi(awayStr)
		},
	}

	return
}

func parseBarChart(handlers map[string]func(string, string), content string, wg *sync.WaitGroup) {
	defer wg.Done()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return;
	}
	
	doc.Find("figure.stats-bar-chart").EachWithBreak(func(i int, s *goquery.Selection) bool {
    title := strings.TrimSpace(s.Find("figcaption.stats-bar-chart__title").Text())
		if handler, ok := handlers[title]; ok {
			homeVals := s.Find(".stats-bar-chart__label--home")
			awayVals := s.Find(".stats-bar-chart__label--away")

			homeStr := strings.TrimSpace(homeVals.Eq(0).Text())
			awayStr := strings.TrimSpace(awayVals.Eq(0).Text())
	
			handler(homeStr, awayStr)
			delete(handlers, title)
		}

		return len(handlers) > 0
	})
}

func MergeInto[K comparable, V any](m1, m2 map[K]V) {
	for k, v := range m2 {
		m1[k] = v
	}
}
