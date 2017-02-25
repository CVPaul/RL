package games

import (
	"common"
	"container/list"
	"fmt"
	"image"
	"log"
	"os"
	"robot"
	"time"
	"utils"

	"azul3d.org/engine/gfx"
	"azul3d.org/engine/gfx/camera"
	"azul3d.org/engine/gfx/window"
	"azul3d.org/engine/keyboard"
	//math "azul3d.org/engine/lmath"
)

var g_rand *utils.Rander

type Snake struct {
	Body   *list.List
	Direct string
}

type Game struct {
	Snake_  *Snake
	FoodPos common.Point
	Score   int64
	Steps   int64
	// enviroment
	cam    *camera.Camera
	shader *gfx.Shader
	// ai
	agent common.AI
}

func (self *Snake) Destroy() {
	self.Body = nil
}

func (self *Game) Destroy() {
	self.Snake_ = nil
	self.Score = 0
	self.Steps = 0
}

func (self *Game) ReStart() error {
	// stop
	self.Snake_.Destroy()
	self.Destroy()
	// start
	// re-init rander
	g_rand.Srand(common.SEED)
	// init game
	self.Snake_ = new(Snake)
	err := self.Snake_.Init() // initialize the snake
	if nil != err {
		return fmt.Errorf("Game re-init failed[snake]||err_msg=%s", err.Error())
	}
	// set food position
	self.FoodPos.X = (g_rand.Rand() % common.STATUS_X_SIZE) * common.MAGIC_NUM
	self.FoodPos.Y = (g_rand.Rand() % common.STATUS_Y_SIZE) * common.MAGIC_NUM
	// init Score
	self.Score = 0
	// init Steps
	self.Steps = 0
	return nil
}

func (self *Snake) Init() error {
	self.Body = list.New()
	if nil == self.Body {
		return fmt.Errorf("Snake init failed")
	}
	//self.Body.PushBack(common.Point{STATUS_X_SIZE / 2, STATUS_Y_SIZE / 2})
	self.Body.PushBack(common.Point{0, 0})
	self.Direct = "R"
	return nil
}

func (self *Game) Init() error {
	// init rander
	g_rand = &utils.Rander{}
	g_rand.Srand(common.SEED)
	// init game
	self.Snake_ = new(Snake)
	err := self.Snake_.Init() // initialize the snake
	if nil != err {
		return fmt.Errorf("Game init failed[snake]||err_msg=%s", err.Error())
	}
	// set food position
	self.FoodPos.X = (g_rand.Rand() % common.STATUS_X_SIZE) * common.MAGIC_NUM
	self.FoodPos.Y = (g_rand.Rand() % common.STATUS_Y_SIZE) * common.MAGIC_NUM
	// init Score
	self.Score = 0
	// init Steps
	self.Steps = 0
	return nil
}

func (self *Game) InitTraining(ai string) error {
	// init rander
	g_rand = &utils.Rander{}
	g_rand.Srand(common.SEED)
	// init ai
	self.agent = AIFactory(ai)
	if nil == self.agent {
		return fmt.Errorf("Do not have agent named:\"%s\", please look at the end of \"greedy_snake.go\" for ai name that support!", ai)
	}
	err := self.agent.InitTraining()
	if nil != err {
		return fmt.Errorf("Agent init failed||err_msg=%s", err.Error())
	}
	// init game
	self.Snake_ = new(Snake)
	err = self.Snake_.Init() // initialize the snake
	if nil != err {
		return fmt.Errorf("Game init failed||err_msg=%s", err.Error())
	}
	// set food position
	self.FoodPos.X = (g_rand.Rand() % common.STATUS_X_SIZE) * common.MAGIC_NUM
	self.FoodPos.Y = (g_rand.Rand() % common.STATUS_Y_SIZE) * common.MAGIC_NUM
	// init Score
	self.Score = 0
	// init Steps
	self.Steps = 0
	return nil
}

func (self *Game) InitTesting(ai string) error {
	// init rander
	g_rand = &utils.Rander{}
	g_rand.Srand(common.SEED)
	// init ai
	self.agent = AIFactory(ai)
	if nil == self.agent {
		return fmt.Errorf("Do not have agent named:\"%s\", please look at the end of \"greedy_snake.go\" for ai name that support!", ai)
	}
	err := self.agent.InitTesting()
	if nil != err {
		return fmt.Errorf("Agent init failed||err_msg=%s", err.Error())
	}
	// init game
	self.Snake_ = new(Snake)
	err = self.Snake_.Init() // initialize the snake
	if nil != err {
		return fmt.Errorf("Game init failed||err_msg=%s", err.Error())
	}
	// set food position
	self.FoodPos.X = (g_rand.Rand() % common.STATUS_X_SIZE) * common.MAGIC_NUM
	self.FoodPos.Y = (g_rand.Rand() % common.STATUS_Y_SIZE) * common.MAGIC_NUM
	// init Score
	self.Score = 0
	// init Steps
	self.Steps = 0
	return nil
}

func (self *Game) GoAhead() {
	direct := self.Snake_.Direct
	head := self.Snake_.Body.Front().Value.(common.Point)
	if "U" == direct {
		head.Y -= common.MAGIC_NUM
	} else if "D" == direct {
		head.Y += common.MAGIC_NUM
	} else if "R" == direct {
		head.X += common.MAGIC_NUM
	} else if "L" == direct {
		head.X -= common.MAGIC_NUM
	} else {
		return // Do nothing
	}
	if head.X != self.FoodPos.X ||
		head.Y != self.FoodPos.Y {
		tail := self.Snake_.Body.Back()
		self.Snake_.Body.Remove(tail)
	} else {
		self.Score += 1
		// reset food position
		self.FoodPos.X = (g_rand.Rand() % common.STATUS_X_SIZE) * common.MAGIC_NUM
		self.FoodPos.Y = (g_rand.Rand() % common.STATUS_Y_SIZE) * common.MAGIC_NUM
		// log.Printf("FoodPos=%v", self.FoodPos)
	}
	self.Snake_.Body.PushFront(head)
	self.Steps += 1
	// log.Printf("GoAhead||headPos=%v||Direction=%s", head, direct)
}

func (self *Snake) IsDead() bool {
	headPos := self.Body.Front().Value.(common.Point)
	if headPos.X < 0 || (headPos.X+common.MAGIC_NUM) > common.MAGIC_NUM*common.STATUS_X_SIZE ||
		headPos.Y < 0 || (headPos.Y+common.MAGIC_NUM) > common.MAGIC_NUM*common.STATUS_Y_SIZE {
		//log.Printf("Dead[Bound]||headPos=%v", headPos)
		return true
	}
	for p := self.Body.Front(); p != nil; p = p.Next() {
		point := p.Value.(common.Point)
		if p != self.Body.Front() && point.X == headPos.X && point.Y == headPos.Y {
			//log.Printf("Dead[Body]||point=%v||head=%v", point, headPos)
			return true
		}
	}
	return false
}

func (self *Game) Point2ColorBox(p common.Point) *image.Rectangle {
	rect := image.Rect(p.X*common.BODY_SIZE, p.Y*common.BODY_SIZE,
		(p.X+common.MAGIC_NUM)*common.BODY_SIZE, (p.Y+common.MAGIC_NUM)*common.BODY_SIZE)
	return &rect
}

// gfxLoop is responsible for drawing things to the window.
func (self *Game) GfxLoop(w window.Window, d gfx.Device) {

	// Create a channel of events.
	events := make(chan window.Event, 256)

	go func() {
		// Have the window notify our channel whenever events occur.
		w.Notify(events, window.KeyboardEvents)
	}()
	for {
		// Handle each pending event.
		window.Poll(events, func(e window.Event) {
			//log.Printf("event[%T]=%v\n", e, e)
		})
		// Depending on keyboard state, transform the triangle.
		kb := w.Keyboard()
		if kb.Down(keyboard.ArrowLeft) && self.Snake_.Direct != "R" { // it is not allowed to go to the opposite side
			self.Snake_.Direct = "L"
		}
		if kb.Down(keyboard.ArrowRight) && self.Snake_.Direct != "L" { // it is not allowed to go to the opposite side
			self.Snake_.Direct = "R"
		}
		if kb.Down(keyboard.ArrowDown) && self.Snake_.Direct != "U" { // it is not allowed to go to the opposite side
			self.Snake_.Direct = "D"
		}
		if kb.Down(keyboard.ArrowUp) && self.Snake_.Direct != "D" { // it is not allowed to go to the opposite side
			self.Snake_.Direct = "U"
		}
		// Clear color and depth buffers.
		d.Clear(d.Bounds(), gfx.Color{1, 1, 1, 1})
		self.GoAhead()
		for p := self.Snake_.Body.Front(); nil != p; p = p.Next() {
			rect := self.Point2ColorBox(p.Value.(common.Point))
			d.Clear(*rect, gfx.Color{0, 1, 1, 1})
		}
		// draw the food to screen
		food := self.Point2ColorBox(self.FoodPos)
		d.Clear(*food, gfx.Color{1, 1, 0, 1})
		// sleep
		start_time := time.Now().UnixNano() / common.NANOS_TO_MILLISECOND
		current := time.Now().UnixNano() / common.NANOS_TO_MILLISECOND
		for {
			if current >= start_time+common.TIME_FOR_SPEED_CTR {
				break
			}
			// also check the keyboard event
			if kb.Down(keyboard.ArrowLeft) && self.Snake_.Direct != "R" { // it is not allowed to go to the opposite side
				self.Snake_.Direct = "L"
			}
			if kb.Down(keyboard.ArrowRight) && self.Snake_.Direct != "L" { // it is not allowed to go to the opposite side
				self.Snake_.Direct = "R"
			}
			if kb.Down(keyboard.ArrowDown) && self.Snake_.Direct != "U" { // it is not allowed to go to the opposite side
				self.Snake_.Direct = "D"
			}
			if kb.Down(keyboard.ArrowUp) && self.Snake_.Direct != "D" { // it is not allowed to go to the opposite side
				self.Snake_.Direct = "U"
			}
			current = time.Now().UnixNano() / common.NANOS_TO_MILLISECOND
			// sleep
			time.Sleep(common.SLEEP_TIME * time.Millisecond)
		}
		// Render the whole frame.
		d.Render()
		if self.Score == 98 {
			var key_type byte
			fmt.Scanf("Press Any Key To Continue:%c", &key_type)
		}
		if self.Snake_.IsDead() {
			log.Printf("Game Over[Score=%d||Steps=%d]\n", self.Score, self.Steps)
			os.Exit(-1)
		}
	}
}

// gfxLoopTraining is responsible for drawing things to the window.
func (self *Game) GfxLoopTraining(w window.Window, d gfx.Device) {
	for {
		// call the ai
		self.PlayerApi()
		// go ahead
		self.GoAhead()
		if common.SHOW {
			//Clear color and depth buffers.
			d.Clear(d.Bounds(), gfx.Color{1, 1, 1, 1})
			//update && draw to screen
			//Draw the snake to the screen.
			for p := self.Snake_.Body.Front(); nil != p; p = p.Next() {
				rect := self.Point2ColorBox(p.Value.(common.Point))
				d.Clear(*rect, gfx.Color{0, 1, 1, 1})
			}
			// draw the food to screen
			food := self.Point2ColorBox(self.FoodPos)
			d.Clear(*food, gfx.Color{1, 1, 0, 1})
			// Render the whole frame.
			d.Render()
		}
		// check if dead
		if self.Snake_.IsDead() {
			log.Printf("Game Over[Score=%d||Steps=%d]\n", self.Score, self.Steps)
			self.ReStart()
		}
	}
}

// gfxLoopTesting is responsible for drawing things to the window.
func (self *Game) GfxLoopTesting(w window.Window, d gfx.Device) {
	for {
		// call the ai
		self.PlayerApi()
		// go ahead
		self.GoAhead()
		if true {
			//Clear color and depth buffers.
			d.Clear(d.Bounds(), gfx.Color{1, 1, 1, 1})
			//update && draw to screen
			//Draw the snake to the screen.
			for p := self.Snake_.Body.Front(); nil != p; p = p.Next() {
				rect := self.Point2ColorBox(p.Value.(common.Point))
				d.Clear(*rect, gfx.Color{0, 1, 1, 1})
			}
			// draw the food to screen
			food := self.Point2ColorBox(self.FoodPos)
			d.Clear(*food, gfx.Color{1, 1, 0, 1})
			// Render the whole frame.
			d.Render()
		}
		// check if dead
		if self.Snake_.IsDead() {
			log.Printf("Game Over[Score=%d||Steps=%d]\n", self.Score, self.Steps)
			self.ReStart()
		}
	}
}

func AIFactory(ai string) common.AI {
	if "greedy" == ai {
		return &robot.GreedyAI{}
	} else if "sarsa" == ai {
		return &robot.SarsaAI{}
	} else if "benchmark" == ai {
		return &robot.BenchMarkAI{} // benchmark ai
	} else {
		return nil
	}
}
