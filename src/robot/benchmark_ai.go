package robot

import (
	"common"
	"container/list"
	"fmt"
	"log"
	"utils"
)

// the framework of benchmark AI
type BenchMarkAI struct {
	PathMat   *common.Mat
	DistMat   *common.Mat
	SnakePred *list.List
}

func (self *BenchMarkAI) GetData() *list.List {
	return self.SnakePred
}

func (self *BenchMarkAI) GetLongestPath(status *common.Status) (string, int) {
	// reset
	// (*self.PathMat) = utils.NewMat(STATUS_X_SIZE, STATUS_Y_SIZE, 0)
	for i := 0; i < common.STATUS_Y_SIZE; i++ {
		for j := 0; j < common.STATUS_X_SIZE; j++ {
			if status.Data[i][j] == common.DEAD_LOSS {
				(*self.PathMat)[i][j] = -1
			} else {
				(*self.PathMat)[i][j] = 0
			}
			(*self.DistMat)[i][j] = -1
		}
	}
	h_pos := status.HeadPos
	h_pos.X /= common.MAGIC_NUM
	h_pos.Y /= common.MAGIC_NUM
	t_pos := status.TailPos
	t_pos.X /= common.MAGIC_NUM
	t_pos.Y /= common.MAGIC_NUM
	direct, dist := self.DFSSearch(h_pos, &t_pos)
	return direct, dist
}

func (self *BenchMarkAI) DFSSearch(pos common.Point, endpos *common.Point) (string, int) {
	// log.Printf("DFS||point=%v||endpos=%v||DistMat[13][12]=(%d,%d,%d,%d)",
	// pos, *endpos, (*self.PathMat)[17][12], (*self.PathMat)[17][13], (*self.PathMat)[12][17], (*self.PathMat)[13][17])
	if pos.X == endpos.X && pos.Y == endpos.Y {
		(*self.DistMat)[pos.Y][pos.X] = 0
		return "", 0
	}
	(*self.PathMat)[pos.Y][pos.X] = -1
	d, dist, direct := -1, -1, ""
	if pos.X+1 < common.STATUS_X_SIZE &&
		(*self.PathMat)[pos.Y][pos.X+1] >= 0 {
		if (*self.DistMat)[pos.Y][pos.X+1] >= 0 {
			d = (*self.DistMat)[pos.Y][pos.X+1]
		} else {
			_, d = self.DFSSearch(common.Point{X: pos.X + 1, Y: pos.Y}, endpos)
		}
		if d > dist {
			dist = d
			direct = "R"
		}
	}
	if pos.X-1 >= 0 &&
		(*self.PathMat)[pos.Y][pos.X-1] >= 0 {
		if (*self.DistMat)[pos.Y][pos.X-1] >= 0 {
			d = (*self.DistMat)[pos.Y][pos.X-1]
		} else {
			_, d = self.DFSSearch(common.Point{X: pos.X - 1, Y: pos.Y}, endpos)
		}
		if d > dist {
			dist = d
			direct = "L"
		}
	}
	if pos.Y+1 < common.STATUS_Y_SIZE &&
		(*self.PathMat)[pos.Y+1][pos.X] >= 0 {
		if (*self.DistMat)[pos.Y+1][pos.X] >= 0 {
			d = (*self.DistMat)[pos.Y+1][pos.X]
		} else {
			_, d = self.DFSSearch(common.Point{X: pos.X, Y: pos.Y + 1}, endpos)
		}
		if d > dist {
			dist = d
			direct = "D"
		}
	}
	if pos.Y-1 >= 0 &&
		(*self.PathMat)[pos.Y-1][pos.X] >= 0 {
		if (*self.DistMat)[pos.Y-1][pos.X] >= 0 {
			d = (*self.DistMat)[pos.Y-1][pos.X]
		} else {
			_, d = self.DFSSearch(common.Point{X: pos.X, Y: pos.Y - 1}, endpos)
		}
		if d > dist {
			dist = d
			direct = "U"
		}
	}
	(*self.DistMat)[pos.Y][pos.X] = dist + 1
	(*self.PathMat)[pos.Y][pos.X] = 0
	return direct, dist + 1
}

func (self *BenchMarkAI) Init() error {
	self.PathMat = utils.NewMat(common.STATUS_X_SIZE, common.STATUS_Y_SIZE, 0)
	if nil == self.PathMat {
		return fmt.Errorf("BenchMarkAI[Init Failed]||err_msg=Bad assignment[PathMat]")
	}
	self.DistMat = utils.NewMat(common.STATUS_X_SIZE, common.STATUS_Y_SIZE, 0)
	if nil == self.DistMat {
		return fmt.Errorf("BenchMarkAI[Init Failed]||err_msg=Bad assignment[DistMat]")
	}
	return nil
}

func (self *BenchMarkAI) InitTraining() error {
	self.PathMat = utils.NewMat(common.STATUS_X_SIZE, common.STATUS_Y_SIZE, 0)
	if nil == self.PathMat {
		return fmt.Errorf("BenchMarkAI[Init Failed]||err_msg=Bad assignment[PathMat]")
	}
	self.DistMat = utils.NewMat(common.STATUS_X_SIZE, common.STATUS_Y_SIZE, 0)
	if nil == self.DistMat {
		return fmt.Errorf("BenchMarkAI[Init Failed]||err_msg=Bad assignment[DistMat]")
	}
	return nil
}

func (self *BenchMarkAI) InitTesting() error {
	self.PathMat = utils.NewMat(common.STATUS_X_SIZE, common.STATUS_Y_SIZE, 0)
	if nil == self.PathMat {
		return fmt.Errorf("BenchMarkAI[Init Failed]||err_msg=Bad assignment[PathMat]")
	}
	self.DistMat = utils.NewMat(common.STATUS_X_SIZE, common.STATUS_Y_SIZE, 0)
	if nil == self.DistMat {
		return fmt.Errorf("BenchMarkAI[Init Failed]||err_msg=Bad assignment[DistMat]")
	}
	return nil
}

func (self *BenchMarkAI) Process(status *common.Status) *common.Action {
	ac := new(common.Action)
	fd_pos := status.FoodPos
	fd_pos.Y /= common.MAGIC_NUM
	fd_pos.X /= common.MAGIC_NUM
	dist := int(common.MAX_DIST)
	if status.Data[fd_pos.Y][fd_pos.X] == common.DEAD_LOSS {
		ac.Value, dist = self.GetLongestPath(status)
		log.Printf("LongestPathStg[Out]||fd_pos=%v||dist=%d||direct=%s", fd_pos, dist, ac.Value)
	} else {
		ac.Value, dist = self.GetShortestPath(status)
		if ac.Value == "" {
			ac.Value, dist = self.GetLongestPath(status)
			log.Printf("LongestPathStg[In]||fd_pos=%v||dist=%d||direct=%s", fd_pos, dist, ac.Value)
		} else {
			log.Printf("ShortestPathStg||fd_pos=%v||dist=%d||direct=%s", fd_pos, dist, ac.Value)
		}
	}
	if !self.ConnectEnsure(status, ac.Value, false) {
		for _, v := range common.ACTIONS {
			if v == ac.Value {
				continue
			}
			if self.ConnectEnsure(status, v, false) {
				log.Printf("ConnectEnsureStg||fd_pos=%v||dist=%d||direct_old=%s||direct_new=%s",
					fd_pos, dist, ac.Value, v)
				ac.Value = v
				break
			}
		}
	}
	return ac
}

func (self *BenchMarkAI) IsReach(headPos common.Point, statMat *common.Mat) bool {
	headPos.X /= common.MAGIC_NUM
	headPos.Y /= common.MAGIC_NUM
	if headPos.X+1 < common.STATUS_X_SIZE &&
		common.NOTHING_REWARD == (*statMat)[headPos.Y][headPos.X+1] {
		return false
	}
	if headPos.X-1 >= 0 &&
		common.NOTHING_REWARD == (*statMat)[headPos.Y][headPos.X-1] {
		return false

	}
	if headPos.Y+1 < common.STATUS_Y_SIZE &&
		common.NOTHING_REWARD == (*statMat)[headPos.Y+1][headPos.X] {
		return false
	}

	if headPos.Y-1 >= 0 &&
		common.NOTHING_REWARD == (*statMat)[headPos.Y-1][headPos.X] {
		return false
	}

	return true
}

func (self *BenchMarkAI) GetDirection(headPos common.Point, traceMat *common.Mat) string {
	headPos.X /= common.MAGIC_NUM
	headPos.Y /= common.MAGIC_NUM
	direction := ""
	if (*traceMat)[headPos.Y][headPos.X] >= 0 {
		direction = common.ACTIONS[(*traceMat)[headPos.Y][headPos.X]]
	}
	return direction
}

func (self *BenchMarkAI) GetShortestPath(status *common.Status) (string, int) {
	// get the shortest path with broad search
	// =============================Find Path==================================
	// status
	statMat := utils.NewMat2(&status.Data)
	HeadPos, FoodPos := status.HeadPos, status.FoodPos
	// head pos
	HeadPos.Y /= common.MAGIC_NUM
	HeadPos.X /= common.MAGIC_NUM
	// food pos
	FoodPos.Y /= common.MAGIC_NUM
	FoodPos.X /= common.MAGIC_NUM
	// trace mat
	(*statMat)[FoodPos.Y][FoodPos.X] = 0
	(*statMat)[HeadPos.Y][HeadPos.X] = common.NOTHING_REWARD
	traceMat := self.BFSShortestPath(FoodPos, HeadPos, statMat)
	// =========================Connection Check================================
	self.SnakePred = self.Predict(FoodPos, HeadPos, traceMat, status.Snake)
	connected := self.ConnectCheck(self.SnakePred)
	// ===============================Return====================================
	direction := ""
	dist := common.MAX_DIST
	if connected {
		direction = self.GetDirection(status.HeadPos, traceMat)
		dist = (*statMat)[HeadPos.Y][HeadPos.X]
	}
	return direction, dist
}

func (self *BenchMarkAI) BFSShortestPath(startPos, endPos common.Point, statMat *common.Mat) (traceMat *common.Mat) {

	traceMat = utils.NewMat(common.STATUS_X_SIZE, common.STATUS_X_SIZE, -1)
	// shortest path search
	queue := new(list.List)

	queue.PushBack(startPos)
	(*statMat)[startPos.Y][startPos.X] = 0
	/*is_reach := false*/
	for queue.Len() > 0 {
		p := queue.Front()
		headPos := p.Value.(common.Point)
		if headPos.Y+1 < common.STATUS_Y_SIZE &&
			common.NOTHING_REWARD == (*statMat)[headPos.Y+1][headPos.X] {
			(*statMat)[headPos.Y+1][headPos.X] = (*statMat)[headPos.Y][headPos.X] + 1
			(*traceMat)[headPos.Y+1][headPos.X] = 0 // Up
			queue.PushBack(common.Point{Y: headPos.Y + 1, X: headPos.X})
		}
		if headPos.Y-1 >= 0 &&
			common.NOTHING_REWARD == (*statMat)[headPos.Y-1][headPos.X] {
			(*statMat)[headPos.Y-1][headPos.X] = (*statMat)[headPos.Y][headPos.X] + 1
			(*traceMat)[headPos.Y-1][headPos.X] = 1 // Down
			queue.PushBack(common.Point{Y: headPos.Y - 1, X: headPos.X})
		}
		if endPos.X >= 0 && endPos.X < common.STATUS_X_SIZE &&
			endPos.Y >= 0 && endPos.Y < common.STATUS_Y_SIZE &&
			(*statMat)[endPos.Y][endPos.X] != common.NOTHING_REWARD {
			break
		}
		if headPos.X+1 < common.STATUS_X_SIZE &&
			common.NOTHING_REWARD == (*statMat)[headPos.Y][headPos.X+1] {
			(*statMat)[headPos.Y][headPos.X+1] = (*statMat)[headPos.Y][headPos.X] + 1
			(*traceMat)[headPos.Y][headPos.X+1] = 2 // Left
			queue.PushBack(common.Point{Y: headPos.Y, X: headPos.X + 1})
		}
		if headPos.X-1 >= 0 &&
			common.NOTHING_REWARD == (*statMat)[headPos.Y][headPos.X-1] {
			(*statMat)[headPos.Y][headPos.X-1] = (*statMat)[headPos.Y][headPos.X] + 1
			(*traceMat)[headPos.Y][headPos.X-1] = 3 // Right
			queue.PushBack(common.Point{Y: headPos.Y, X: headPos.X - 1})
		}
		queue.Remove(p)
	}
	return
}

func (self *BenchMarkAI) Predict(FoodPos, HeadPos common.Point, traceMat *common.Mat, snake *list.List) *list.List {
	// track the path
	path := new(list.List)
	point := HeadPos
	path.PushBack(point)
	for point.X != FoodPos.X || point.Y != FoodPos.Y {
		if 0 == (*traceMat)[point.Y][point.X] {
			point = common.Point{X: point.X, Y: point.Y - 1}
		} else if 1 == (*traceMat)[point.Y][point.X] {
			point = common.Point{X: point.X, Y: point.Y + 1}
		} else if 2 == (*traceMat)[point.Y][point.X] {
			point = common.Point{X: point.X - 1, Y: point.Y}
		} else if 3 == (*traceMat)[point.Y][point.X] {
			point = common.Point{X: point.X + 1, Y: point.Y}
		} else {
			break
		}
		path.PushBack(point)
	}
	// reverse
	rpath := new(list.List)
	for p := path.Front(); nil != p && rpath.Len() < snake.Len(); p = p.Next() {
		pt := p.Value.(common.Point)
		if rpath.Len() == 0 {
			rpath.PushBack(pt)
		} else {
			rpath.InsertBefore(pt, rpath.Front())
		}
	}
	//log.Printf("TraceMat[FoodPos=%v]=%v,path=%v,l1=%d,l2=%d,headPos=%v", FoodPos, *traceMat, common.GetListInfo(rpath, 1), path.Len(), snake.Len(), HeadPos)
	//os.Exit(-1)
	for p := snake.Front(); p != nil && path.Len() < snake.Len(); p = p.Next() {
		point = p.Value.(common.Point)
		point.X /= common.MAGIC_NUM
		point.Y /= common.MAGIC_NUM
		path.PushBack(point)
	}
	return path
}

func (self *BenchMarkAI) ConnectCheck(snake_pred *list.List) bool {
	statMat := utils.NewMat(common.STATUS_X_SIZE, common.STATUS_Y_SIZE, common.NOTHING_REWARD)
	for p := snake_pred.Front(); nil != p; p = p.Next() {
		point := p.Value.(common.Point)
		(*statMat)[point.Y][point.X] = common.DEAD_LOSS
	}
	return (self.ConnectCount(statMat) <= 1)
}

func (self *BenchMarkAI) ConnectCount(statMat *common.Mat) int {
	count := 0
	for i := 0; i < common.STATUS_Y_SIZE; i++ {
		for j := 0; j < common.STATUS_X_SIZE; j++ {
			if common.NOTHING_REWARD == (*statMat)[i][j] {
				self.BFSShortestPath(common.Point{X: j, Y: i}, common.Point{X: -1, Y: -1}, statMat)
				count += 1
			}
		}
	}
	return count
}

func (self *BenchMarkAI) ConnectEnsure(status *common.Status, direct string, print_out bool) bool {
	statMat0 := utils.NewMat(common.STATUS_X_SIZE, common.STATUS_Y_SIZE, common.NOTHING_REWARD)
	statMat := utils.NewMat(common.STATUS_X_SIZE, common.STATUS_Y_SIZE, common.NOTHING_REWARD)
	for i := 0; i < common.STATUS_Y_SIZE; i++ {
		for j := 0; j < common.STATUS_X_SIZE; j++ {
			if status.Data[i][j] == common.DEAD_LOSS {
				(*statMat)[i][j] = common.DEAD_LOSS
				(*statMat0)[i][j] = common.DEAD_LOSS
			}
		}
	}
	count0 := self.ConnectCount(statMat0)
	HeadPos, TailPos, FoodPos := status.HeadPos, status.TailPos, status.FoodPos
	HeadPos.X /= common.MAGIC_NUM
	HeadPos.Y /= common.MAGIC_NUM
	TailPos.X /= common.MAGIC_NUM
	TailPos.Y /= common.MAGIC_NUM
	FoodPos.X /= common.MAGIC_NUM
	FoodPos.Y /= common.MAGIC_NUM
	ever_in := false
	NewHead := HeadPos
	if "U" == direct && HeadPos.Y-1 >= 0 &&
		(*statMat)[HeadPos.Y-1][HeadPos.X] != common.DEAD_LOSS {
		(*statMat)[HeadPos.Y-1][HeadPos.X] = common.DEAD_LOSS
		NewHead.Y -= 1
		if HeadPos.Y-1 != FoodPos.Y || HeadPos.X != FoodPos.X {
			(*statMat)[TailPos.Y][TailPos.X] = common.NOTHING_REWARD
		}
		ever_in = true
	}
	if "D" == direct && HeadPos.Y+1 < common.STATUS_Y_SIZE &&
		(*statMat)[HeadPos.Y+1][HeadPos.X] != common.DEAD_LOSS {
		(*statMat)[HeadPos.Y+1][HeadPos.X] = common.DEAD_LOSS
		NewHead.Y += 1
		if HeadPos.Y+1 != FoodPos.Y || HeadPos.X != FoodPos.X {
			(*statMat)[TailPos.Y][TailPos.X] = common.NOTHING_REWARD
		}
		ever_in = true
	}
	if "L" == direct && HeadPos.X-1 >= 0 &&
		(*statMat)[HeadPos.Y][HeadPos.X-1] != common.DEAD_LOSS {
		NewHead.X -= 1
		(*statMat)[HeadPos.Y][HeadPos.X-1] = common.DEAD_LOSS
		if HeadPos.Y != FoodPos.Y || HeadPos.X-1 != FoodPos.X {
			(*statMat)[TailPos.Y][TailPos.X] = common.NOTHING_REWARD
		}
		ever_in = true
	}
	if "R" == direct && HeadPos.X+1 < common.STATUS_Y_SIZE &&
		(*statMat)[HeadPos.Y][HeadPos.X+1] != common.DEAD_LOSS {
		NewHead.X += 1
		(*statMat)[HeadPos.Y][HeadPos.X+1] = common.DEAD_LOSS
		if HeadPos.Y != FoodPos.Y || HeadPos.X+1 != FoodPos.X {
			(*statMat)[TailPos.Y][TailPos.X] = common.NOTHING_REWARD
		}
		ever_in = true
	}
	can_alive := false
	if NewHead.X+1 < common.STATUS_X_SIZE && !can_alive {
		if (*statMat)[NewHead.Y][NewHead.X+1] == common.NOTHING_REWARD {
			can_alive = true
		}
	}
	if NewHead.X-1 >= 0 && !can_alive {
		if (*statMat)[NewHead.Y][NewHead.X-1] == common.NOTHING_REWARD {
			can_alive = true
		}
	}
	if NewHead.Y+1 < common.STATUS_Y_SIZE && !can_alive {
		if (*statMat)[NewHead.Y+1][NewHead.X] == common.NOTHING_REWARD {
			can_alive = true
		}
	}
	if NewHead.Y-1 >= 0 && !can_alive {
		if (*statMat)[NewHead.Y-1][NewHead.X] == common.NOTHING_REWARD {
			can_alive = true
		}
	}
	if print_out {
		log.Printf("ConnectEnsure[HPos=%v,TPos=%v,FPos=%v]||statMat[direct=%v]=%v", HeadPos, TailPos, FoodPos, direct, *statMat)
	}
	ret := true
	if ever_in && can_alive {
		ret = (self.ConnectCount(statMat) <= count0)
	} else {
		ret = false
	}

	return ret
}
