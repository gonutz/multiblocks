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
	for player, lines := range linesForPlayer {
		s.teamScores[s.playerToTeam[player]] += lineScores[len(lines)]
	}
}
