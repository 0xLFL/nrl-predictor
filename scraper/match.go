package main

type Play struct {
	time int
	play string
	team string
	notes string
}

type Player struct {
	nameFirst string
	nameLast string
	position string
	number int
	playerStats PlayerStats
}

type MatchOffical struct {
	nameFirst string
	nameLast string
	role string
}

type PosAndComp struct {
	homePosPer int
	awayPosPer int

	homePosTime int
	awayPosTime int

	homeSets int
	homeSetCompleated int

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

	homeAvgSetDistance float64
	awayAvgSetDistance float64

	homeKickReturnMeters int
	awayKickReturnMeters int

	homeAvgPlayTheBallSpeed int
	awayAvgPlayTheBallSpeed int
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
	homeEffecTackle int
	awayEffecTackle int

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

type MatchStats struct {
	posAndComp PosAndComp
	attack Attack
	passing Passing
	kicking Kicking
	defence Defence
	negPlays NegPlays
}

type Match struct {
	homeTeam string
	homeHalfScore int
	homeScore int
	homeTeamList []Player

	awayTeam string
	awayHaftScore int
	awayScore int
	awayTeamList []Player

	matchOfficals []MatchOffical

	location string
	kickoffTime string
	datePlayed string
	weather	string

	playByPlay []Play

	stats MatchStats
}