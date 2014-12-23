package game

type GameSounds interface {
	PlayDown()
	PlayHorizontal()
	PlayRotate()
	PlayCollision()
	PlayGroundHit()
}

type SoundPlayer struct {
	sounds     GameSounds
	down       bool
	horizontal bool
	rotate     bool
	collision  bool
	ground     bool
}

func NewSoundPlayer(sounds GameSounds) *SoundPlayer {
	return &SoundPlayer{sounds: sounds}
}

func (p *SoundPlayer) BlockMovedHorizontally(int) {
	p.horizontal = true
}

func (p *SoundPlayer) BlockMovedDown(int) {
	p.down = true
}

func (p *SoundPlayer) BlockRotated(int) {
	p.rotate = true
}

func (p *SoundPlayer) BlockHitLeftOrRight(int) {
	p.collision = true
}

func (p *SoundPlayer) BlockHitGround(int) {
	p.ground = true
}

func (p *SoundPlayer) BlockHitOtherBlock(int) {
	p.collision = true
}

func (p *SoundPlayer) BlockCouldNotRotate(int) {
	p.collision = true
}

func (ü *SoundPlayer) BlockDraggedDownByLineRemoval(int) {}

func (p *SoundPlayer) PlaySounds() {
	if p.horizontal {
		p.sounds.PlayHorizontal()
	}
	if p.down {
		p.sounds.PlayDown()
	}
	if p.rotate {
		p.sounds.PlayRotate()
	}
	if p.collision {
		p.sounds.PlayCollision()
	}
	if p.ground {
		p.sounds.PlayGroundHit()
	}
	p.resetFlags()
}

func (p *SoundPlayer) resetFlags() {
	p.down = false
	p.horizontal = false
	p.rotate = false
	p.ground = false
	p.collision = false
}
