package main

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/plotutil"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
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
	teams := [...]string{"Arsenal", "Aston Villa", "Burnley", "Chelsea", "Crystal Palace", "Everton", "Hull", "Leicester", "Liverpool", "Man City", "Man United", "Newcastle", "QPR", "Southampton", "Stoke", "Sunderland", "Swansea", "Tottenham", "West Brom", "West Ham"}
	for _, t := range teams {
		analysis(t)
	}
}

func analysis(team string) {
	resp, err := http.Get("http://www.football-data.co.uk/mmz4281/1415/E0.csv")
	if err != nil {
		// handle error
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	reader.LazyQuotes = true
	var scored, allowed, match, win = 0, 0, 0, 0.0
	reals := make([]float64, 38)
	pythagorean := make([]float64, 38)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		m := NewMatch(record[2], record[3], record[4], record[5])

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
			reals[match-1] = float64(win) / float64(match) * 100
			pythagorean[match-1] = PythagoreanExpectation(scored, allowed, match)
		}
	}
	GenerateGraph(team, pythagorean, reals)
}

func GeneratePlot(parcentage []float64) plotter.XYs {
	pts := make(plotter.XYs, len(parcentage))
	for i := range parcentage {
		pts[i].X = float64(i)
		pts[i].Y = parcentage[i]
	}
	return pts
}

func GenerateGraph(teamName string, pythagorean, reals []float64) {
	p, _ := plot.New()
	p.Title.Text = teamName
	p.X.Label.Text = "Day"
	p.Y.Label.Text = "Parcentage"
	p.X.Min = 0.0
	p.X.Max = 38.0
	p.Y.Min = 0.0
	p.Y.Max = 100.0
	plotutil.AddLinePoints(p, "", GeneratePlot(pythagorean), GeneratePlot(reals))
	width := 4.0
	height := 4.0
	p.Save(width, height, fmt.Sprintf("%s.png", teamName))
}

func PythagoreanExpectation(scored, allowed, match int) float64 {
	var scoredf, allowedf = float64(scored), float64(allowed)
	var n = math.Pow(((scoredf + allowedf) / float64(match)), 0.287)
	return math.Pow(scoredf, n) / (math.Pow(scoredf, n) + math.Pow(allowedf, n)) * 100
}
