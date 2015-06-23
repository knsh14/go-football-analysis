package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

type Match struct {
	homeTeam string
	awayTeam string
	homeGoal int
	awayGoal int
}

func NewMatch(homeTeam, awayTeam, homeGoal, awayGoal string) *Match {
	hoge := &Match{}
	hoge.homeTeam = homeTeam
	hoge.awayTeam = awayTeam
	hoge.homeGoal, _ = strconv.Atoi(homeGoal)
	hoge.awayGoal, _ = strconv.Atoi(awayGoal)
	return hoge
}

func (match *Match) Is(teamName string) bool {
	if match.homeTeam == teamName || match.awayTeam == teamName {
		return true
	}
	return false
}

func (match *Match) GetScore(teamName string) (int, int) {
	switch teamName {
	case match.homeTeam:
		return match.homeGoal, match.awayGoal
	case match.awayTeam:
		return match.awayGoal, match.homeGoal
	default:
		return 0, 0
	}
}

func (match *Match) IsWin(teamName string) bool {
	var s, a = 0, 0
	switch teamName {
	case match.homeTeam:
		s, a = match.homeGoal, match.awayGoal
	case match.awayTeam:
		s, a = match.awayGoal, match.homeGoal
	}
	return s > a
}

func (match *Match) IsDraw() bool {
	return match.homeGoal == match.awayGoal
}

func main() {

	var fp *os.File
	var err error
	fp, err = os.Open("result.csv")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	reader.Comma = ','
	reader.LazyQuotes = true
	var scored, allowed, match, win = 0, 0, 0, 0.0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		m := NewMatch(record[0], record[1], record[2], record[3])

		team := "Man United"
		if m.Is(team) {
			s, a := m.GetScore(team)
			scored += s
			allowed += a
			match++
			if m.IsWin(team) {
				win++
			}
			if m.IsDraw() {
				win += 0.3333
			}
			fmt.Printf("%f\t%f\n", PythagoreanExpectation(scored, allowed, match), float64(win)/float64(match)*100)
		}
	}
}

func PythagoreanExpectation(scored, allowed, match int) float64 {
	var scoredf, allowedf = float64(scored), float64(allowed)
	var n = math.Pow(((scoredf + allowedf) / float64(match)), 0.287)
	return math.Pow(scoredf, n) / (math.Pow(scoredf, n) + math.Pow(allowedf, n)) * 100
}
