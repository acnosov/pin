package handler

import (
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

type Market struct {
	BetType      string
	PeriodNumber int
	Team         string
	Side         string
	Handicap     *float64
}

const (
	SPREAD            = "SPREAD"
	MONEYLINE         = "MONEYLINE"
	TOTAL_POINTS      = "TOTAL_POINTS"
	TEAM_TOTAL_POINTS = "TEAM_TOTAL_POINTS"
	OVER              = "OVER"
	UNDER             = "UNDER"
	Team1             = "Team1"
	Team2             = "Team2"
	Draw              = "Draw"
)

var teamSide = map[string]string{"1": Team1, "2": Team2, "Х": Draw, "X": Draw}
var totalSide = map[string]string{"Б": OVER, "М": UNDER}

var moneyRe = regexp.MustCompile(`П(\d)$`)
var threeRe = regexp.MustCompile(`(^|\s)([12ХX]$)`)
var handicapRe = regexp.MustCompile(`Ф(\d)\((-?\d+,?\d{0,3})\)`)
var totalRe = regexp.MustCompile(`Т([МБ])\((-?\d+,?\d{0,3})\)`)
var teamTotalRe = regexp.MustCompile(`ИТ(\d)([МБ])\((-?\d+,?\d{0,3})\)`)

//BetType: Available values : SPREAD, MONEYLINE, TOTAL_POINTS, TEAM_TOTAL_POINTS
//PeriodNumber: For example, for soccer we have 0 (Game), 1 (1st Half) & 2 (2nd Half)
//Handicap: This is needed for SPREAD, TOTAL_POINTS and TEAM_TOTAL_POINTS bet types
//Side: This is needed only for TOTAL_POINTS and TEAM_TOTAL_POINTS  Available values : OVER, UNDER
//Team: This is needed only for SPREAD, MONEYLINE and TEAM_TOTAL_POINTS bet types  Available values : Team1, Team2, Draw
//var abc = map[string]Market{
//	//"1X": {Name: "Double Chance", Side: "1x"},
//	//"X2": {Name: "Double Chance", Side: "x2"},
//	//"12": {Name: "Double Chance", Side: "12"},
//}

func Convert(name string) (Market, error) {
	if name == "" {
		return Market{}, errors.Errorf("cannot convert empty str")
	}
	var found []string

	found = moneyRe.FindStringSubmatch(name)
	if len(found) == 2 {
		return processMoney(name, found)
	}

	found = threeRe.FindStringSubmatch(name)
	if len(found) == 3 {
		return processThreeWay(name, found)
	}

	found = handicapRe.FindStringSubmatch(name)
	if len(found) == 3 {
		return processHandicap(name, found)
	}

	found = totalRe.FindStringSubmatch(name)
	if len(found) == 3 {
		return processTotal(name, found)
	}
	found = teamTotalRe.FindStringSubmatch(name)
	if len(found) == 4 {
		return processTeamTotal(name, found)
	}
	return Market{}, nil
}
func processPoint(point string) (float64, error) {
	pointsStr := strings.Replace(point, ",", ".", -1)
	return strconv.ParseFloat(pointsStr, 64)
}

func processPeriodNumber(market string) int {
	var periodNumber = 0
	switch {
	case strings.Index(market, "1/2") != -1:
		periodNumber = 1
	case strings.Index(market, "1/3") != -1:
		periodNumber = 1
	case strings.Index(market, "2/3") != -1:
		periodNumber = 2
	case strings.Index(market, "3/3") != -1:
		periodNumber = 3
	case strings.Index(market, "2/2") != -1:
		periodNumber = 2
	case strings.Index(market, "1/4") != -1:
		periodNumber = 3
	case strings.Index(market, "3/2") != -1:
		periodNumber = 3
	case strings.Index(market, "2/4") != -1:
		periodNumber = 4
	case strings.Index(market, "4/2") != -1:
		periodNumber = 4
	case strings.Index(market, "3/4") != -1:
		periodNumber = 5
	case strings.Index(market, "5/2") != -1:
		periodNumber = 5
	case strings.Index(market, "4/4") != -1:
		periodNumber = 6
	case strings.Index(market, "6/2") != -1:
		periodNumber = 6
	case strings.Index(market, "7/2") != -1:
		periodNumber = 7
	case strings.Index(market, "8/2") != -1:
		periodNumber = 8
	}
	return periodNumber
}
func processMoney(market string, found []string) (Market, error) {
	periodNumber := processPeriodNumber(market)
	m := Market{
		BetType:      MONEYLINE,
		PeriodNumber: periodNumber,
		Team:         teamSide[found[1]],
	}
	return m, nil
}
func processThreeWay(market string, found []string) (Market, error) {
	periodNumber := processPeriodNumber(market)
	m := Market{
		BetType:      MONEYLINE,
		PeriodNumber: periodNumber,
		Team:         teamSide[found[2]],
	}
	return m, nil
}
func processHandicap(market string, found []string) (Market, error) {
	point, err := processPoint(found[2])
	if err != nil {
		return Market{}, err
	}
	periodNumber := processPeriodNumber(market)
	m := Market{
		BetType:      SPREAD,
		PeriodNumber: periodNumber,
		Team:         teamSide[found[1]],
		Handicap:     &point,
	}
	return m, nil
}

func processTotal(market string, found []string) (Market, error) {
	point, err := processPoint(found[2])
	if err != nil {
		return Market{}, err
	}
	periodNumber := processPeriodNumber(market)
	m := Market{
		BetType:      TOTAL_POINTS,
		PeriodNumber: periodNumber,
		Side:         totalSide[found[1]],
		Handicap:     &point,
	}
	return m, nil
}
func processTeamTotal(market string, found []string) (Market, error) {
	point, err := processPoint(found[3])
	if err != nil {
		return Market{}, err
	}
	periodNumber := processPeriodNumber(market)
	m := Market{
		BetType:      TEAM_TOTAL_POINTS,
		PeriodNumber: periodNumber,
		Team:         teamSide[found[1]],
		Side:         totalSide[found[2]],
		Handicap:     &point,
	}
	return m, nil
}

//"П1":     {BetType: MONEYLINE, PeriodNumber: 0, Side: Team1},
//"П2":     {BetType: MONEYLINE, PeriodNumber: 0, Side: Team2},
//"1/2 П1": {BetType: MONEYLINE, PeriodNumber: 1, Side: Team1},
//"1/2 П2": {BetType: MONEYLINE, PeriodNumber: 1, Side: Team2},
//"2/2 П1": {BetType: MONEYLINE, PeriodNumber: 2, Side: Team1},
//"2/2 П2": {BetType: MONEYLINE, PeriodNumber: 2, Side: Team2},
//"1/4 П1": {BetType: MONEYLINE, PeriodNumber: 3, Side: Team1},
//"1/4 П2": {BetType: MONEYLINE, PeriodNumber: 3, Side: Team2},
//"2/4 П1": {BetType: MONEYLINE, PeriodNumber: 4, Side: Team1},
//"2/4 П2": {BetType: MONEYLINE, PeriodNumber: 4, Side: Team2},
//"3/4 П1": {BetType: MONEYLINE, PeriodNumber: 5, Side: Team1},
//"3/4 П2": {BetType: MONEYLINE, PeriodNumber: 5, Side: Team2},
//"4/4 П1": {BetType: MONEYLINE, PeriodNumber: 6, Side: Team1},
//"4/4 П2": {BetType: MONEYLINE, PeriodNumber: 6, Side: Team2},
//"Чёт":       {Name: "Odd/Even", Side: "a"},
//"Нечёт":     {Name: "Odd/Even", Side: "h"},
//"1/2 Чёт":   {Name: "First Half O/E", Side: "a"},
//"1/2 Нечёт": {Name: "First Half O/E", Side: "h"},
//"1/4 Чёт":   {Name: "1st Quarter - Odd/Even", Side: "a"},
//"1/4 Нечёт": {Name: "1st Quarter - Odd/Even", Side: "h"},
//"2/4 Чёт":   {Name: "2nd Quarter - Odd/Even", Side: "a"},
//"2/4 Нечёт": {Name: "2nd Quarter - Odd/Even", Side: "h"},
//"3/4 Чёт":   {Name: "3rd Quarter - Odd/Even", Side: "a"},
//"3/4 Нечёт": {Name: "3rd Quarter - Odd/Even", Side: "h"},
//"4/4 Чёт":   {Name: "4th Quarter - Odd/Even", Side: "a"},
//"4/4 Нечёт": {Name: "4th Quarter - Odd/Even", Side: "h"},

//"1":     {BetType: MONEYLINE, PeriodNumber: 0, Side: Team1},
//"2":     {BetType: MONEYLINE, PeriodNumber: 0, Side: Team2},
//"Х":     {BetType: MONEYLINE, PeriodNumber: 0, Side: Draw},
//"1/2 1": {BetType: MONEYLINE, PeriodNumber: 1, Side: Team1},
//"1/2 2": {BetType: MONEYLINE, PeriodNumber: 1, Side: Team2},
//"1/2 Х": {BetType: MONEYLINE, PeriodNumber: 1, Side: Draw},
//"2/2 1": {BetType: MONEYLINE, PeriodNumber: 2, Side: Team1},
//"2/2 2": {BetType: MONEYLINE, PeriodNumber: 2, Side: Team2},
//"2/2 Х": {BetType: MONEYLINE, PeriodNumber: 2, Side: Draw},
//
//"УГЛ 1":     {BetType: MONEYLINE, PeriodNumber: 0, Side: Team1},
//"УГЛ 2":     {BetType: MONEYLINE, PeriodNumber: 0, Side: Team2},
//"УГЛ Х":     {BetType: MONEYLINE, PeriodNumber: 0, Side: Draw},
//"УГЛ 1/2 1": {BetType: MONEYLINE, PeriodNumber: 1, Side: Team1},
//"УГЛ 1/2 2": {BetType: MONEYLINE, PeriodNumber: 1, Side: Team2},
//"УГЛ 1/2 Х": {BetType: MONEYLINE, PeriodNumber: 1, Side: Draw},
//"УГЛ 2/2 1": {BetType: MONEYLINE, PeriodNumber: 2, Side: Team1},
//"УГЛ 2/2 2": {BetType: MONEYLINE, PeriodNumber: 2, Side: Team2},
//"УГЛ 2/2 Х": {BetType: MONEYLINE, PeriodNumber: 2, Side: Draw},
//"ЖК 1":     {BetType: MONEYLINE, PeriodNumber: 0, Side: Team1},
//"ЖК 2":     {BetType: MONEYLINE, PeriodNumber: 0, Side: Team2},
//"ЖК Х":     {BetType: MONEYLINE, PeriodNumber: 0, Side: Draw},
//"ЖК 1/2 1": {BetType: MONEYLINE, PeriodNumber: 1, Side: Team1},
//"ЖК 1/2 2": {BetType: MONEYLINE, PeriodNumber: 1, Side: Team2},
//"ЖК 1/2 Х": {BetType: MONEYLINE, PeriodNumber: 1, Side: Draw},
//"ЖК 2/2 1": {BetType: MONEYLINE, PeriodNumber: 2, Side: Team1},
//"ЖК 2/2 2": {BetType: MONEYLINE, PeriodNumber: 2, Side: Team2},
//"ЖК 2/2 Х": {BetType: MONEYLINE, PeriodNumber: 2, Side: Draw},
