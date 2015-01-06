package game

import "testing"

func TestPlayerCanBeAssignedToTeam(t *testing.T) {
	s := NewTeamScorer()
	s.AssignPlayerToTeam(0, 2)
}

func TestInitiallyAllScoresAreZero(t *testing.T) {
	s := NewTeamScorer()
	for i := 0; i < 4; i++ {
		if score := s.ScoreForTeam(i); score != 0 {
			t.Error("score for player", i, "not 0 but was", score)
		}
	}
}

func TestEachTeamIsScoredByItsRemovedLines(t *testing.T) {
	s := NewTeamScorer()
	s.AssignPlayerToTeam(1, 3)

	s.LinesRemoved([][]int{
		{},
		{2, 3},
		{},
		{},
	})

	if score := s.ScoreForTeam(3); score != lineScores[2] {
		t.Error("team 3 scored", score, "but expected", lineScores[2])
	}
}

func TestTeamScoresIncreaseWithRemovedLines(t *testing.T) {
	s := NewTeamScorer()
	s.AssignPlayerToTeam(0, 0)
	s.LinesRemoved([][]int{{1}})
	s.LinesRemoved([][]int{{1, 2, 3}})
	expected := lineScores[1] + lineScores[3]
	if score := s.ScoreForTeam(0); score != expected {
		t.Errorf("expected %v but score was %v", expected, score)
	}
}
