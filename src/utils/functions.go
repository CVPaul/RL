package utils

import (
	"common"
	"container/list"
	"fmt"
	"math"
	"os"
)

func NewMat(cols, rows, default_val int) (data *common.Mat) {
	var data_t common.Mat = make([][]int, rows)
	data = &data_t
	for i := 0; i < rows; i++ {
		(*data)[i] = make([]int, cols)
	}
	if default_val == 0 {
		return
	}
	for i := 0; i < rows; i += 1 {
		for j := 0; j < cols; j += 1 {
			(*data)[i][j] = default_val
		}
	}
	return
}

func NewMat2(ori *common.Mat) *common.Mat {
	rows, cols := GetMatSize(ori)
	cur := NewMat(cols, rows, 0)
	for i := 0; i < rows; i += 1 {
		for j := 0; j < cols; j += 1 {
			(*cur)[i][j] = (*ori)[i][j]
		}
	}
	return cur
}

func GetMatSize(mat *common.Mat) (rows, cols int) {
	if nil == mat {
		return 0, 0
	}
	rows = len(*mat)
	if rows < 1 {
		cols = 0
	} else {
		cols = len((*mat)[0])
	}
	return
}

func MultiMats(mat1 *common.Mat, mat2 *common.Mat) (err error, resMat *common.Mat) {
	if nil == mat1 || nil == mat2 {
		err = fmt.Errorf("One of the input is nil:mat1=%v,mat2=%v", mat1, mat2)
		return err, nil
	}
	r1, c1 := GetMatSize(mat1)
	r2, c2 := GetMatSize(mat2)
	if c1 != r2 {
		err := fmt.Errorf("Size of two matrices must be match:(%d,%d)X(%d,%d)", r1, c1, r2, c2)
		return err, nil
	}
	resMat = NewMat(r1, c2, 0)
	for i := 0; i < r1; i++ {
		for j := 0; j < c2; j++ {
			for k := 0; k < r2; k++ {
				(*resMat)[i][j] += ((*mat1)[i][k]) * ((*mat2)[k][j])
			}
		}
	}
	return nil, resMat
}

func Pos2Idx(pos common.Point) int {
	return int(pos.Y)*int(common.STATUS_X_SIZE) + int(pos.X)
}
func Idx2Pos(idx int) (pos common.Point) {
	pos.X = int(idx % common.STATUS_X_SIZE)
	pos.Y = int(idx / common.STATUS_X_SIZE)
	return
}
func GetListInfo(link *list.List, Factor int) string {
	snake_info := ""
	for p := link.Front(); p != nil; p = p.Next() {
		pt := p.Value.(common.Point)
		pt.X /= Factor
		pt.Y /= Factor
		if "" == snake_info {
			snake_info = fmt.Sprintf("%v", pt)
		} else {
			snake_info = fmt.Sprintf("%s,%v", snake_info, pt)
		}
	}
	return snake_info
}
func GetNewPos(cur common.Point, direct string) *common.Point {
	if "U" == direct {
		cur.Y -= 1
	}
	if "D" == direct {
		cur.Y += 1
	}
	if "L" == direct {
		cur.X -= 1
	}
	if "R" == direct {
		cur.X += 1
	}
	return &cur
}
func IsOutOfRange(pt common.Point) bool {
	if pt.X < 0 || pt.X >= common.STATUS_X_SIZE ||
		pt.Y < 0 || pt.Y >= common.STATUS_Y_SIZE {
		return true
	}
	return false
}
func GetPhase(pt0, pt1 common.Point) int {
	h := float64(pt1.X - pt0.X)
	v := float64(pt1.Y - pt0.Y)
	if 0.0 == h && v > 0 {
		return 0
	}
	if 0.0 == h && v <= 0 {
		return 1
	}
	if 0.0 == v && h > 0 {
		return 2
	}
	if 0.0 == v && h <= 0 {
		return 3
	}
	degree := int(180*math.Atan2(v, h)/math.Pi) + 180

	return 4 + degree/common.DEGREE_SPLIT
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
