package game

type TeamScorer struct {
	playerToTeam [4]int
	teamScores   [4]int
}

var lineScores = [...]int{
	0,
	1, 3, 6, 10,
	15, 21, 28, 36,
	45, 55, 66, 78,
	91, 105, 120, 136,
}

func NewTeamScorer() *TeamScorer {
	return &TeamScorer{}
}

func (s *TeamScorer) AssignPlayerToTeam(player, team int) {
	s.playerToTeam[player] = team
}

func (s *TeamScorer) ScoreForTeam(team int) int {
	return s.teamScores[team]
}

func (s *TeamScorer) LinesRemoved(linesForPlayer [][]int) {
	teamLines := s.assembleLinesForAllTeamsOfAllPlayers(linesForPlayer)
	for team, lines := range teamLines {
		lineCount := countDistinct(lines)
		s.teamScores[team] += lineScores[lineCount]
	}
}

func (s *TeamScorer) assembleLinesForAllTeamsOfAllPlayers(linesForPlayer [][]int) [4][]int {
	var teamLines [4][]int
	for player, lines := range linesForPlayer {
		team := s.playerToTeam[player]
		teamLines[team] = append(teamLines[team], lines...)
	}
	return teamLines
}

func countDistinct(lines []int) int {
	count := 0
	for i, line := range lines {
		if !contains(lines[:i], line) {
			count++
		}
	}
	return count
}

func contains(lines []int, line int) bool {
	for _, l := range lines {
		if l == line {
			return true
		}
	}
	return false
}
