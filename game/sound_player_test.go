package game

import "testing"

func TestMovingBlockPlaysSoundsOnUpdate(t *testing.T) {
	spy := &spyGameSounds{}
	player := NewSoundPlayer(spy)

	player.BlockMovedHorizontally(0)
	player.BlockMovedDown(0)
	player.BlockRotated(0)
	if spy.log != "" {
		t.Errorf("sounds played before signal, log was: '%s'", spy.log)
	}

	player.PlaySounds()
	if spy.log != "horizontal down rotate " {
		t.Errorf("sound log was: '%s'", spy.log)
	}
}

func TestSoundsAlreadyPlayed_AreNotRepeated(t *testing.T) {
	spy := &spyGameSounds{}
	player := NewSoundPlayer(spy)
	player.BlockMovedHorizontally(0)
	player.BlockMovedDown(0)
	player.BlockRotated(0)

	player.PlaySounds()
	player.PlaySounds()
	if spy.log != "horizontal down rotate " {
		t.Errorf("sounds played more than once, log: '%s'", spy.log)
	}
}

func TestOnlyOnMoveSoundIsPlayedAtATime(t *testing.T) {
	spy := &spyGameSounds{}
	player := NewSoundPlayer(spy)
	player.BlockMovedHorizontally(0)
	player.BlockMovedHorizontally(0)
	player.BlockMovedDown(0)
	player.BlockMovedDown(0)
	player.BlockRotated(0)
	player.BlockRotated(0)

	player.PlaySounds()
	if spy.log != "horizontal down rotate " {
		t.Errorf("sounds played more than once, log: '%s'", spy.log)
	}
}

func TestGroundHitAndCollisionSoundIsPlayed(t *testing.T) {
	spy := &spyGameSounds{}
	player := NewSoundPlayer(spy)
	player.BlockHitLeftOrRight(1)
	player.BlockHitGround(2)
	player.PlaySounds()
	if spy.log != "collision grounded " {
		t.Errorf("sound log was: '%s'", spy.log)
	}

	player.PlaySounds()
	if spy.log != "collision grounded " {
		t.Errorf("sounds played again, log: '%s'", spy.log)
	}
}

func TestAllKindsOfCollisionPlayTheSameSound(t *testing.T) {
	spy := &spyGameSounds{}
	player := NewSoundPlayer(spy)
	player.BlockHitOtherBlock(1)
	player.PlaySounds()
	player.BlockCouldNotRotate(1)
	player.PlaySounds()
	if spy.log != "collision collision " {
		t.Errorf("sound log was: '%s'", spy.log)
	}
}

func TestLineRemovalDraggingDownBlockDoesNotPlaySound(t *testing.T) {
	spy := &spyGameSounds{}
	player := NewSoundPlayer(spy)
	player.BlockDraggedDownByLineRemoval(1)
	player.PlaySounds()
	if spy.log != "" {
		t.Errorf("sound log was: '%s'", spy.log)
	}
}

type spyGameSounds struct {
	log string
}

func (spy *spyGameSounds) PlayDown()       { spy.log += "down " }
func (spy *spyGameSounds) PlayHorizontal() { spy.log += "horizontal " }
func (spy *spyGameSounds) PlayRotate()     { spy.log += "rotate " }
func (spy *spyGameSounds) PlayCollision()  { spy.log += "collision " }
func (spy *spyGameSounds) PlayGroundHit()  { spy.log += "grounded " }
