package robot

import (
	"common"
	"container/list"
	"log"
	"os"
)

// the framework of greedy ai
type GreedyAI struct {
}

func (self *GreedyAI) Init() error {
	// DO NOTHING
	return nil
}
func (self *GreedyAI) GetData() *list.List {
	return nil
}

func (self *GreedyAI) InitTraining() error {
	log.Print("GreedyAI do not need training")
	os.Exit(-1)
	return nil
}

func (self *GreedyAI) InitTesting() error {
	return nil
}

func (self *GreedyAI) Process(status *common.Status) *common.Action {
	HeadPos := status.HeadPos
	HeadPos.X /= common.MAGIC_NUM
	HeadPos.Y /= common.MAGIC_NUM
	ac := new(common.Action)
	if 0 == HeadPos.Y%2 {
		if 0 == HeadPos.X {
			if 0 == HeadPos.Y {
				ac.Value = "R"
			} else {
				ac.Value = "U"
			}
		} else {
			if HeadPos.X < common.STATUS_X_SIZE-1 {
				ac.Value = "R"
			} else {
				ac.Value = "D"
			}
		}
	} else {
		if common.STATUS_Y_SIZE-1 > HeadPos.Y {
			if HeadPos.X == 1 {
				ac.Value = "D"
			} else {
				if HeadPos.X > 0 {
					ac.Value = "L"
				} else {
					ac.Value = "U"
				}
			}
		} else {
			if HeadPos.X > 0 {
				ac.Value = "L"
			} else {
				ac.Value = "U"
			}
		}
	}
	return ac
}
