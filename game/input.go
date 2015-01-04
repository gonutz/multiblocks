package game

// InputEvent abstracts the button presses of players.
type InputEvent struct {
	Player  int
	Command Command
}

type Command int

const (
	DownPressed Command = iota
	LeftPressed
	RightPressed
	DownReleased
	LeftReleased
	RightReleased
	RotateLeft
	RotateRight
	Pause
)
