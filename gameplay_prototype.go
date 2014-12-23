package main

import (
	"fmt"
	"github.com/gonutz/multiblocks/game"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var window *sdl.Window
var renderer *sdl.Renderer

var playerCount int

const blockSize = 20

func RunGame() {
	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()
	playerCount = 2
	if len(os.Args) == 2 {
		count, err := strconv.ParseInt(os.Args[1], 10, 32)
		if err == nil {
			playerCount = int(count)
		}
	}
	window, renderer = sdl.CreateWindowAndRenderer(800, 600, 0)
	defer window.Destroy()
	defer renderer.Destroy()
	ttf.Init()
	defer ttf.Quit()
	window.SetTitle("Multiblocks")
	initControllers()
	initAssets()
	defer closeAssets()
	g := game.NewLogic(randomBlock)
	g.SetBoardSizeForPlayerCount(1, game.BoardSize{10, 18})
	g.SetBlockStartPositions(1, []game.Point{{5, 16}})
	g.SetBoardSizeForPlayerCount(2, game.BoardSize{10, 18})
	g.SetBlockStartPositions(2, []game.Point{{7, 16}, {2, 16}})
	g.SetBoardSizeForPlayerCount(3, game.BoardSize{13, 18})
	g.SetBlockStartPositions(3, []game.Point{{6, 16}, {2, 16}, {10, 16}})
	g.SetBoardSizeForPlayerCount(4, game.BoardSize{16, 18})
	g.SetBlockStartPositions(4, []game.Point{{10, 16}, {2, 16}, {14, 16}, {6, 16}})
	g.SetDropTimer(&dropTimer{})
	g.SetLineAnimation(animation)
	g.SetInitialLeftRightKeyDelay(9)
	g.SetFastLeftRightKeyDelay(2)
	g.SetInitialDownKeyDelay(2)
	g.SetFastDownKeyDelay(1)
	g.SetScorer(scorer)
	g.SetSoundPlayer(game.NewSoundPlayer(sounds))
	g.StartNewGame(playerCount)
	animation.board = g.Board()

	running := true
	var inputs []game.InputEvent
	for running {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch event := e.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyDownEvent:
				if event.Keysym.Sym == sdl.K_ESCAPE {
					running = false
				}
				if event.Repeat == 0 {
					input, ok := keyDownMap[event.Keysym.Sym]
					if ok {
						inputs = append(inputs, input)
					}
				}
			case *sdl.KeyUpEvent:
				switch event.Keysym.Sym {
				case sdl.K_ESCAPE:
					running = false
				case sdl.K_n:
					g.StartNewGame(playerCount)
					scorer.reset()
				}

				input, ok := keyUpMap[event.Keysym.Sym]
				if ok {
					inputs = append(inputs, input)
				}
			case *sdl.ControllerButtonEvent:
				player := int(event.Which)
				switch event.Button {
				case sdl.CONTROLLER_BUTTON_A:
					if event.State == sdl.PRESSED {
						inputs = append(inputs, game.InputEvent{player, game.RotateLeft})
					}
				case sdl.CONTROLLER_BUTTON_B:
					if event.State == sdl.PRESSED {
						inputs = append(inputs, game.InputEvent{player, game.RotateRight})
					}
				case sdl.CONTROLLER_BUTTON_DPAD_LEFT, sdl.CONTROLLER_BUTTON_LEFTSHOULDER:
					state := game.LeftPressed
					if event.State == sdl.RELEASED {
						state = game.LeftReleased
					}
					inputs = append(inputs, game.InputEvent{player, state})
				case sdl.CONTROLLER_BUTTON_DPAD_RIGHT, sdl.CONTROLLER_BUTTON_RIGHTSHOULDER:
					state := game.RightPressed
					if event.State == sdl.RELEASED {
						state = game.RightReleased
					}
					inputs = append(inputs, game.InputEvent{player, state})
				case sdl.CONTROLLER_BUTTON_DPAD_DOWN:
					state := game.DownPressed
					if event.State == sdl.RELEASED {
						state = game.DownReleased
					}
					inputs = append(inputs, game.InputEvent{player, state})
				}
			}
		}
		update(g, &inputs)
		sdl.Delay(30)
		draw(g)
	}
}

func update(g *game.Logic, inputs *[]game.InputEvent) {
	g.Update((*inputs)...)
	*inputs = make([]game.InputEvent, 0)
}

func initControllers() {
	joystickCount := sdl.NumJoysticks()
	fmt.Println(joystickCount, "joysticks detected")
	for i := 0; i < joystickCount; i++ {
		joystick := sdl.JoystickOpen(i)
		fmt.Println("joystick", joystick)
		guid := sdl.JoystickGetGUIDString(joystick.GetGUID())
		mapping := "a:b2,b:b1,x:b3,y:b0,start:b9,back:b8," +
			"dpup:h0.1,dpleft:h0.8,dpdown:h0.4,dpright:h0.2"
		fmt.Println("add mapping", sdl.GameControllerAddMapping(guid+",USB GamePad,"+mapping))
		controller := sdl.GameControllerOpen(i)
		fmt.Println("controller", controller)
		fmt.Println("SDL error", sdl.GetError())
	}
}

func initAssets() {
	initKeys()
	initColors()
	initFactory()
	initScorer()
	initSounds()
	animation = &lineAnimation{}
	rand.Seed(time.Now().UnixNano())
}

var keyDownMap map[sdl.Keycode]game.InputEvent
var keyUpMap map[sdl.Keycode]game.InputEvent

func initKeys() {
	keyDownMap = map[sdl.Keycode]game.InputEvent{
		sdl.K_LEFT:  game.InputEvent{0, game.LeftPressed},
		sdl.K_RIGHT: game.InputEvent{0, game.RightPressed},
		sdl.K_DOWN:  game.InputEvent{0, game.DownPressed},
		sdl.K_UP:    game.InputEvent{0, game.RotateRight},
		sdl.K_a:     game.InputEvent{1, game.LeftPressed},
		sdl.K_d:     game.InputEvent{1, game.RightPressed},
		sdl.K_s:     game.InputEvent{1, game.DownPressed},
		sdl.K_w:     game.InputEvent{1, game.RotateRight},
		sdl.K_KP_4:  game.InputEvent{2, game.LeftPressed},
		sdl.K_KP_6:  game.InputEvent{2, game.RightPressed},
		sdl.K_KP_5:  game.InputEvent{2, game.DownPressed},
		sdl.K_KP_8:  game.InputEvent{2, game.RotateRight},
		sdl.K_g:     game.InputEvent{3, game.LeftPressed},
		sdl.K_j:     game.InputEvent{3, game.RightPressed},
		sdl.K_h:     game.InputEvent{3, game.DownPressed},
		sdl.K_z:     game.InputEvent{3, game.RotateRight},
	}

	keyUpMap = map[sdl.Keycode]game.InputEvent{
		sdl.K_LEFT:  game.InputEvent{0, game.LeftReleased},
		sdl.K_RIGHT: game.InputEvent{0, game.RightReleased},
		sdl.K_DOWN:  game.InputEvent{0, game.DownReleased},
		sdl.K_a:     game.InputEvent{1, game.LeftReleased},
		sdl.K_d:     game.InputEvent{1, game.RightReleased},
		sdl.K_s:     game.InputEvent{1, game.DownReleased},
		sdl.K_KP_4:  game.InputEvent{2, game.LeftReleased},
		sdl.K_KP_6:  game.InputEvent{2, game.RightReleased},
		sdl.K_KP_5:  game.InputEvent{2, game.DownReleased},
		sdl.K_g:     game.InputEvent{3, game.LeftReleased},
		sdl.K_j:     game.InputEvent{3, game.RightReleased},
		sdl.K_h:     game.InputEvent{3, game.DownReleased},
	}
}

func initColors() {
	colors = [][]color{
		[]color{{255, 32, 32}, {210, 22, 22}},
		[]color{{0, 192, 0}, {0, 167, 0}},
		[]color{{64, 64, 255}, {54, 54, 235}},
		[]color{{235, 0, 235}, {205, 0, 205}},
	}
}

type color [3]uint8

var colors [][]color

func initFactory() {
	factory := game.NewBlockFactory()
	blockNewers = []func() game.Block{
		factory.CreateO,
		factory.CreateT,
		factory.CreateI,
		factory.CreateL,
		factory.CreateJ,
		factory.CreateS,
		factory.CreateZ,
	}
}

var blockNewers []func() game.Block

func randomBlock() game.Block {
	return blockNewers[rand.Int()%len(blockNewers)]()
}

func initScorer() {
	linesToScoreMap := []int{0, 2, 5, 15, 60}
	scorer = &perPlayerScorer{
		make([]int, playerCount),
		make([]int, playerCount),
		linesToScoreMap,
	}
}

var scorer *perPlayerScorer

type perPlayerScorer struct {
	lines           []int
	scores          []int
	linesToScoreMap []int
}

func (s *perPlayerScorer) LinesRemoved(lines game.PlayersToLines) {
	for p := 0; p < playerCount; p++ {
		lineCount := len(lines.LinesForPlayer(p))
		s.lines[p] += lineCount
		s.scores[p] += s.linesToScoreMap[lineCount]
	}
}

func (s *perPlayerScorer) draw() {
	_, boardH := animation.board.Size()
	for p, s := range s.scores {
		y := (p + 1 + boardH) * blockSize
		c := dark(p)
		renderer.SetDrawColor(c[0], c[1], c[2], 255)
		for x := 0; x < s; x++ {
			drawX := x + 10 + x*2
			renderer.DrawLine(drawX, y, drawX, y+blockSize*2/3)
		}
	}
}

func (s *perPlayerScorer) reset() {
	for i := range s.scores {
		s.scores[i] = 0
	}
}

func draw(g *game.Logic) {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	board := g.Board()
	blocks := g.Blocks()
	previews := g.PreviewBlocks()
	w, h := board.Size()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := light(board.At(x, y))
			drawPiece(int32(x), int32(h-1-y), c, c)
		}
	}

	for player, b := range blocks {
		for _, p := range b.Points {
			drawPiece(int32(p.X), int32(h-1-p.Y), light(player), dark(player))
		}
	}

	for player, b := range previews {
		for _, p := range b.Points {
			drawPiece(int32(w+2+p.X), int32(h-1-player*5-p.Y), light(player), dark(player))
		}
	}

	animation.draw()
	scorer.draw()

	renderer.Present()
}

var backGroundColor color = color{64, 64, 64}

func light(player int) color {
	if player == game.NoPlayer {
		return backGroundColor
	}
	return colors[player][0]
}

func dark(player int) color {
	if player == game.NoPlayer {
		return backGroundColor
	}
	return colors[player][1]
}

func drawPiece(x, y int32, light, dark color) {
	renderer.SetDrawColor(dark[0], dark[1], dark[2], 255)
	r := &sdl.Rect{x * blockSize, y * blockSize, blockSize, blockSize}
	renderer.FillRect(r)
	renderer.SetDrawColor(light[0], light[1], light[2], 255)
	r = &sdl.Rect{x*blockSize + 1, y*blockSize + 1, blockSize - 2, blockSize - 2}
	renderer.FillRect(r)
}

type dropTimer struct {
	timer int
}

func (t *dropTimer) IsTimeToDrop() bool {
	return t.timer == 0
}

func (t *dropTimer) Reset() {
	t.timer = 27
}

func (t *dropTimer) Update() {
	t.timer--
	if t.timer < 0 {
		t.Reset()
	}
}

type lineAnimation struct {
	board    game.Board
	lines    []int
	blinking bool
	timer    int
}

var animation *lineAnimation

func (a *lineAnimation) IsRunning() bool {
	return a.timer > 0
}

func (a *lineAnimation) Start(lines []int) {
	a.lines = lines
	a.timer = 15
	a.blinking = true
	sounds.PlayRemove()
}

func (a *lineAnimation) Update() {
	if a.timer%3 == 2 {
		a.blinking = !a.blinking
	}
	a.timer--
}

func (a *lineAnimation) draw() {
	if a.IsRunning() {
		w, h := a.board.Size()
		if a.blinking {
			for _, line := range a.lines {
				for x := 0; x < w; x++ {
					white := color{255, 255, 255}
					drawPiece(int32(x), int32(h-1-line), white, white)
				}
			}
		}
	}
}

func initSounds() {
	sounds = &gameSounds{}
	sounds.init()
}

var sounds *gameSounds

type gameSounds struct {
	down       *mix.Chunk
	horizontal *mix.Chunk
	rotate     *mix.Chunk
	collision  *mix.Chunk
	ground     *mix.Chunk
	remove     *mix.Chunk
	music      *mix.Music
}

func (s *gameSounds) init() {
	if !mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 1024) {
		panic("could not open audio")
	}
	folder := "sfx/"
	s.down = mix.LoadWAV(folder + "down.wav")
	s.horizontal = mix.LoadWAV(folder + "horizontal.wav")
	s.rotate = mix.LoadWAV(folder + "rotate.wav")
	s.collision = mix.LoadWAV(folder + "collision.wav")
	s.ground = mix.LoadWAV(folder + "ground.wav")
	s.remove = mix.LoadWAV(folder + "remove.wav")
	if s.down == nil {
		panic("unable to load sounds")
	}
	s.music = mix.LoadMUS(folder + "track_a.ogg")
	if s.music != nil {
		s.music.Play(-1)
	}
}

func (s *gameSounds) close() {
	mix.CloseAudio()
	s.down.Free()
	s.horizontal.Free()
	s.rotate.Free()
	s.collision.Free()
	s.ground.Free()
	s.remove.Free()
	if s.music != nil {
		s.music.Free()
	}
}

func (s *gameSounds) PlayDown() {
	//s.down.PlayChannel(-1, 0)
}

func (s *gameSounds) PlayHorizontal() {
	s.horizontal.PlayChannel(-1, 0)
}

func (s *gameSounds) PlayRotate() {
	s.rotate.PlayChannel(-1, 0)
}

func (s *gameSounds) PlayCollision() {
	//s.collision.PlayChannel(-1, 0)
}

func (s *gameSounds) PlayGroundHit() {
	s.ground.PlayChannel(-1, 0)
}

func (s *gameSounds) PlayRemove() {
	s.remove.PlayChannel(-1, 0)
}

func closeAssets() {
	sounds.close()
}
