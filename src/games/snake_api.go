// About matrix: http://stackoverflow.com/questions/16536029/go-matrix-library
package games

import (
	"common"
)

func (self *Game) PlayerApi() {
	st := self.GetStatus()
	ac := self.agent.Process(st)
	//log.Printf("Check FoodPos=%v,Action=%v", st.FoodPos, ac.Value)
	if len(ac.Value) > 0 {
		if ac.Value == "L" && self.Snake_.Direct != "R" { // it is not allowed to go to the opposite side
			self.Snake_.Direct = "L"
		}
		if ac.Value == "R" && self.Snake_.Direct != "L" { // it is not allowed to go to the opposite side
			self.Snake_.Direct = "R"
		}
		if ac.Value == "D" && self.Snake_.Direct != "U" { // it is not allowed to go to the opposite side
			self.Snake_.Direct = "D"
		}
		if ac.Value == "U" && self.Snake_.Direct != "D" { // it is not allowed to go to the opposite side
			self.Snake_.Direct = "U"
		}
		self.Snake_.Direct = ac.Value
	}
}

func (self *Game) GetStatus() (status *common.Status) {
	status = new(common.Status)
	status.Init(common.STATUS_X_SIZE, common.STATUS_Y_SIZE)

	fx := self.FoodPos.X / common.MAGIC_NUM
	fy := self.FoodPos.Y / common.MAGIC_NUM
	status.Data[fy][fx] = common.FOOD_REWARD

	for p := self.Snake_.Body.Front(); p != nil; p = p.Next() {
		point := p.Value.(common.Point)
		x := point.X / common.MAGIC_NUM
		y := point.Y / common.MAGIC_NUM
		status.Data[y][x] = common.DEAD_LOSS
	}
	// get head position
	p := self.Snake_.Body.Front()
	status.HeadPos = p.Value.(common.Point)
	// get tail position
	p = self.Snake_.Body.Back()
	status.TailPos = p.Value.(common.Point)
	// get food position
	status.FoodPos = self.FoodPos
	// snake info
	status.Direction = self.Snake_.Direct
	status.Snake = self.Snake_.Body
	return status
}
