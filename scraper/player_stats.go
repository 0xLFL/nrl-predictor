package main

type PlayerStats struct {
	points int
	tries int
	conversions int
	conversionAttempted int
	penGoals int
	conversionRate int
	feildGoal1 int
	feildGoal2 int

	runs int
	runMeters int
	kickReturnMeters int

	postContactMeters int
	lineBreaks int
	lineBreakAssists int
	tryAssists int
	lineEngagedRuns int
	tackleBreaks int
	hitUps int
	AvgPlayBallSpeed int
	dummyHalfRuns int
	dummyHalfRunMeters int
	oneOnOneSteal int

	offloads int
	dummyPasses int
	passes int
	receipts int
	passesToRunRatio float64

	tackleEff int
	tacklesMade int
	tacklesMissed int
	ineffectiveTackles int
	intercepts int
	kicksDefused int

	kicks int
	kickMeters int
	forcedDropOuts int
	bombKicks int
	grubbers int
	fourtyTwenty int
	twentyFourty int
	crossFieldKicks int
	kicksDead int
	
	errors int
	handlingErr int
	OneOnOneLost int
	pen int
	ruckInf int
	inside10 int
	onReport int
	sinBins int
	sendOffs int
}
