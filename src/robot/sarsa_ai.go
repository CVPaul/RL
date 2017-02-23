package robot

import (
	"bufio"
	"bytes"
	"common"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"utils"
)

// the framework of prime ai
type SarsaAI struct {
	IterCount     int64
	Alpha         float64
	Gama          float64
	Epsilon       float64
	Epsilon0      float64
	EpsilonFactor float64
	PowerParam    float64

	TheQ   *map[string]float64
	suffix string

	CurStatus  string
	LastStatus string
	LastReward float64

	IsTestPhase bool
}

func (self *SarsaAI) Init() error {
	return nil
}

func (self *SarsaAI) GetData() *list.List {
	return nil
}

func (self *SarsaAI) InitTraining() error {
	// defalut param
	self.IterCount = 0
	self.Alpha = 0.1
	self.Gama = 0.9
	self.PowerParam = 0.8 // sqrt
	self.EpsilonFactor = common.EPSILON_FAC
	self.Epsilon0 = 1.0
	self.suffix = "sarsa"
	self.IsTestPhase = false
	err, last_param := self.LoadTheQ()
	// last save param
	if v, ok := (*last_param)["IterCount"]; ok {
		self.IterCount = int64(v)
	}
	if v, ok := (*last_param)["Alpha"]; ok {
		self.Alpha = v
	}
	if v, ok := (*last_param)["Gama"]; ok {
		self.Gama = v
	}
	if v, ok := (*last_param)["Epsilon"]; ok {
		self.Epsilon0 = v
	}
	self.Epsilon = self.Epsilon0
	// last staus
	self.LastStatus = "ZERO#SRL"
	self.LastReward = common.DEFAULT_REWARD
	(*self.TheQ)[self.LastStatus] = common.DEFAULT_VALUE // default
	if nil != err {
		log.Print("Load %s failed, please check it!", common.PARAM_POLICY_PATH+self.suffix)
		os.Exit(-1)
	}
	return nil
}

func (self *SarsaAI) InitTesting() error {
	// defalut param
	self.PowerParam = 0.5 // sqrt
	self.EpsilonFactor = common.EPSILON_FAC
	self.suffix = "sarsa"
	self.IsTestPhase = true
	err, last_param := self.LoadTheQ()
	// last save param
	if v, ok := (*last_param)["IterCount"]; ok {
		self.IterCount = int64(v)
	} else {
		return fmt.Errorf("Miss IterCount")
	}
	if v, ok := (*last_param)["Alpha"]; ok {
		self.Alpha = v
	} else {
		return fmt.Errorf("Miss Alpha")
	}
	if v, ok := (*last_param)["Gama"]; ok {
		self.Gama = v
	} else {
		return fmt.Errorf("Miss Gama")
	}
	if v, ok := (*last_param)["Epsilon"]; ok {
		self.Epsilon = v
	} else {
		return fmt.Errorf("Miss Epsilon")
	}
	self.Epsilon0 = 0.0
	self.Epsilon = self.Epsilon0
	// last staus
	self.LastStatus = "ZERO#SRL"
	self.LastReward = common.DEFAULT_REWARD
	(*self.TheQ)[self.LastStatus] = common.DEFAULT_VALUE // default
	if nil != err {
		log.Print("Load %s failed, please check it!", common.PARAM_POLICY_PATH+self.suffix)
		os.Exit(-1)
	}
	/*/ =============For check===============
	fmt.Printf("param=%v\n", *last_param)
	fmt.Printf("theQ=%v\n", *self.TheQ)
	var c byte
	fmt.Scanf("%c", &c)
	// =============End for check===============*/
	return nil
}

func (self *SarsaAI) Process(status *common.Status) *common.Action {
	// Convert the status first
	self.CurStatus = self.StatusConvert(status)
	// Strategy
	ac := new(common.Action)
	action, st := self.EpsilonGreedy(status)
	reward := self.CalReward(status, action)
	TheQValue := (*self.TheQ)[self.LastStatus] +
		self.Alpha*(self.LastReward+self.Gama*(*self.TheQ)[st]-(*self.TheQ)[self.LastStatus])
	if common.DEBUG {
		/*log.Printf("fuck[cur=%v,lastStatus=%s]:[0][0]=%.4f,[0][1]=%.4f,[0][2]=%.4f,[0][3]=%.4f",
			status.HeadPos, self.LastStatus,
			(*self.TheQ)["Lrs000000000400#S"], (*self.TheQ)["Lrs010001000400#S"],
			(*self.TheQ)["Lrs020002000400#S"], (*self.TheQ)["Lrs030003000400#S"])
		log.Printf("fuck:[04][L]=%.4f,[04][R]=%.4f,[04][S]=%.4f",
			(*self.TheQ)["Lrs040003001715#L"], (*self.TheQ)["Lrs040003001715#R"], (*self.TheQ)["Lrs040003001715#S"])*/
		/*if self.LastStatus == "Lrs000000000400#S" || st == "Lrs000000000400#S" {
			log.Printf("fuck[last=%s,cur=%s,lrwd=%.0f]:lastQ=%.4f,curQ=%.4f", self.LastStatus, st, self.LastReward,
				(*self.TheQ)[self.LastStatus], (*self.TheQ)[st])
		}*/
	}

	// ======================================
	(*self.TheQ)[self.LastStatus] = TheQValue
	if common.DEBUG {
		var input_c byte
		fmt.Scanf("please input %c", &input_c)
	}
	// ======================================
	self.LastReward, self.LastStatus = float64(reward), st
	ac.Value = action
	// after process
	self.IterCount += 1
	if self.IterCount >= common.SAVE_EVERY_ITER &&
		0 == self.IterCount%common.SAVE_EVERY_ITER {
		if common.VERBOSE == false {
			log.Printf("save at iteration=%d", self.IterCount)
		}
		err := self.SaveTheQ()
		if nil != err {
			log.Printf("Iteration=%d||err_msg=%s", self.IterCount, err.Error())
		}
		//log.Printf("TheQ=%v", *self.TheQ)
		//os.Exit(-1)
	}
	if !self.IsTestPhase && common.VERBOSE {
		log.Printf("Info:{epsilon:%.4f,IterCount:%d,status=%s,action=%s}\n",
			self.Epsilon, self.IterCount, st, ac.Value)
	}
	return ac
}

func (self *SarsaAI) UpdateEpsilon() {
	//log.Printf("before||Epsilon=%.4f||fac=%.4f",
	//	self.Epsilon, common.MINIMIZE_EPSILON)
	if self.IsTestPhase {
		self.Epsilon = 0.0
		return
	}
	if self.Epsilon <= common.MINIMIZE_EPSILON {
		self.Epsilon = common.MINIMIZE_EPSILON
	} else {
		factor := self.EpsilonFactor / (self.EpsilonFactor + math.Pow(float64(self.IterCount), self.PowerParam))
		self.Epsilon = factor * self.Epsilon0
	}
	//log.Printf("after||Epsilon=%.4f||fac=%.4f",
	//	self.Epsilon, common.MINIMIZE_EPSILON)
}

func (self *SarsaAI) LoadTheQ() (err error, param *map[string]float64) {
	filename := common.PARAM_POLICY_PATH + self.suffix
	fin, err := os.Open(filename)
	if nil != err {
		fin, err = os.Create(filename)
	}
	Q_t := make(map[string]float64)
	param_t := make(map[string]float64)
	defer func() {
		fin.Close()
		self.TheQ = &Q_t
		param = &param_t
	}()
	reader := bufio.NewReader(fin)
	if nil == reader {
		err = fmt.Errorf("Open file=%s failed", filename)
		return
	}
	for {
		line, read_err := reader.ReadString('\n')
		if read_err == io.EOF {
			break
		} else if read_err != nil {
			err = read_err
			return
		}
		line = strings.TrimSpace(line)
		line = strings.Replace(line, "\r", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		items := strings.Split(line, "||")
		if len(items) == 2 && items[0] == "Param" {
			json.Unmarshal([]byte(items[1]), &param_t)
		}
		for _, item := range items {
			sa_v := strings.Split(item, "=")
			if len(sa_v) < 2 || len(sa_v[0]) < 1 || len(sa_v[1]) < 1 {
				continue
			}
			sa_v[0] = strings.TrimSpace(sa_v[0])
			sa_v[1] = strings.TrimSpace(sa_v[1])
			value, err_t := strconv.ParseFloat(sa_v[1], 64)
			if nil != err_t {
				log.Printf("Error[Parsing]||item=%s||err_msg=%s", items, err_t.Error())
				continue
			}
			Q_t[sa_v[0]] = value
		}
	}
	return nil, &param_t
}

func (self *SarsaAI) SaveTheQ() (err error) {
	filename := common.PARAM_POLICY_PATH + self.suffix
	fout, err := os.Create(filename)
	if nil != err {
		err = fmt.Errorf("Open file error||file_name=%s||err_msg=%s", filename, err.Error())
		return
	}
	outputWriter := bufio.NewWriter(fout)
	defer func() {
		outputWriter.Flush()
		fout.Close()
	}()
	// save the basic informatio for human
	info := fmt.Sprintf("Param||{\"Epsilon\":%.4f,\"Alpha\":%.4f,\"Gama\":%.4f,\"IterCount\":%d.0000}\n",
		self.Epsilon, self.Alpha, self.Gama, self.IterCount)
	outputWriter.WriteString(info)
	var buffer bytes.Buffer //Buffer是一个实现了读写方法的可变大小的字节缓冲
	var n int = 0
	for k, v := range *self.TheQ {
		n, err = buffer.WriteString(fmt.Sprintf("%s=%f\n", k, v))
		if n > 10000 { //准备写入
			_, err = outputWriter.WriteString(buffer.String())
			buffer.Reset()
		}
		/*
		   func (b *Buffer) WriteString(s string) (n int, err error)
		   Write将s的内容写入缓冲中，如必要会增加缓冲容量。返回值n为len(p)，err总是nil。如果缓冲变得太大，Write会采用错误值ErrTooLarge引发panic。
		*/
	}
	_, err = outputWriter.WriteString(buffer.String())
	return
}

func (self *SarsaAI) StatusConvert(status *common.Status) string {
	w_s, w_l, w_r := "s", "l", "r"
	// straight
	d := common.GoStraight[status.Direction]
	HeadPos := status.HeadPos
	HeadPos.X /= common.MAGIC_NUM
	HeadPos.Y /= common.MAGIC_NUM
	NewPos := utils.GetNewPos(HeadPos, d)
	if nil == NewPos ||
		utils.IsOutOfRange(*NewPos) ||
		common.DEAD_LOSS == status.Data[NewPos.Y][NewPos.X] {
		w_s = "S"
	}
	// left
	d = common.TurnLeft[status.Direction]
	NewPos = utils.GetNewPos(HeadPos, d)
	if nil == NewPos ||
		utils.IsOutOfRange(*NewPos) ||
		common.DEAD_LOSS == status.Data[NewPos.Y][NewPos.X] {
		w_l = "L"
	}
	// right
	d = common.TurnRight[status.Direction]
	NewPos = utils.GetNewPos(HeadPos, d)
	if nil == NewPos ||
		utils.IsOutOfRange(*NewPos) ||
		common.DEAD_LOSS == status.Data[NewPos.Y][NewPos.X] {
		w_r = "R"
	}
	qf := utils.GetPhase(status.HeadPos, status.FoodPos)
	qt := utils.GetPhase(status.TailPos, status.FoodPos)
	return fmt.Sprintf("%s%s%s%d%d%d", w_l, w_r, w_s, status.Snake.Len(), qf, qt)
}

func (self *SarsaAI) EpsilonGreedy(st *common.Status) (action string, status string) {
	self.UpdateEpsilon()
	if rand.Float64() <= self.Epsilon { // random
		idx := int(rand.Float64() * 3)
		switch idx {
		case 0:
			action, status = common.GoStraight[st.Direction], self.CurStatus+"#S"
		case 1:
			action, status = common.TurnLeft[st.Direction], self.CurStatus+"#L"
		case 2:
			action, status = common.TurnRight[st.Direction], self.CurStatus+"#R"
		}
	} else {
		// straight
		s := self.CurStatus + "#S"
		sv := float64(common.DEFAULT_VALUE) // default value
		if _, ok := (*self.TheQ)[s]; ok {
			sv = (*self.TheQ)[s]
		}
		// greedy init
		value := sv
		action, status = common.GoStraight[st.Direction], s
		// left
		s = self.CurStatus + "#L"
		sv = float64(common.DEFAULT_VALUE) // default value
		if _, ok := (*self.TheQ)[s]; ok {
			sv = (*self.TheQ)[s]
		}
		if sv > value {
			action, value, status = common.TurnLeft[st.Direction], sv, s
		}
		// right
		s = self.CurStatus + "#R"
		sv = float64(common.DEFAULT_VALUE) // default value
		if _, ok := (*self.TheQ)[s]; ok {
			sv = (*self.TheQ)[s]
		}
		if sv > value {
			action, value, status = common.TurnRight[st.Direction], sv, s
		}
		// log.Printf("Debug||Value=%.4f", value)
	}
	// log.Printf("Debug||Status=%v,Action=%v", status, action)
	return
}

func (self *SarsaAI) CalReward(st *common.Status, direct string) int {
	HeadPos := st.HeadPos
	HeadPos.X /= common.MAGIC_NUM
	HeadPos.Y /= common.MAGIC_NUM
	NewPos := utils.GetNewPos(HeadPos, direct)
	if utils.IsOutOfRange(*NewPos) {
		//log.Printf("Will HitWall:HeadPos=%v,Direct=%v,NewPos=%v", HeadPos, direct, *NewPos)
		return common.HIT_WALL
	} else if st.Data[NewPos.Y][NewPos.X] == common.DEAD_LOSS {
		//log.Printf("Will HitWall:HeadPos=%v,Direct=%v,NewPos=%v", HeadPos, direct, *NewPos)
		return common.HIT_WALL
	} else if st.Data[NewPos.Y][NewPos.X] == common.FOOD_REWARD {
		//log.Printf("Will EatFood:HeadPos=%v,Direct=%v,NewPos=%v", HeadPos, direct, *NewPos)
		return common.EAT_FOOD
	} else {
		//log.Printf("Will Defalut:HeadPos=%v,Direct=%v,NewPos=%v", HeadPos, direct, *NewPos)
		return common.DEFAULT_REWARD
	}
}
