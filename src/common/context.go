package common

import (
	"container/list"
)

type Mat [][]int

type AI interface {
	Process(status *Status) *Action
	Init() error
	InitTraining() error
	InitTesting() error
	GetData() *list.List
}

type Context struct {
	Policy []Mat
	Return Mat
}

type Status struct {
	HeadPos    Point
	FoodPos    Point
	TailPos    Point
	Cols, Rows int
	Data       Mat
	Direction  string
	Snake      *list.List
}

type Action struct {
	Value string
}

type Point struct {
	X, Y int
}

func (self *Status) Init(cols, rows int) {
	self.Cols = cols
	self.Rows = rows
	self.Data = make([][]int, self.Rows)
	for i := 0; i < self.Rows; i++ {
		self.Data[i] = make([]int, self.Cols)
	}
	for i := 0; i < self.Rows; i += 1 {
		for j := 0; j < self.Cols; j += 1 {
			self.Data[i][j] = NOTHING_REWARD
		}
	}
}
