package iowrap

import (
	"bufio"
	"bytes"
	"common"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"utils"
)

func LoadMat(file string) (err error, data common.Mat) {
	fin, err := os.Open(file)
	data = common.Mat{}
	if nil != err {
		err = fmt.Errorf("Open file error||file_name=%s||err_msg=%s", file, err.Error())
		return
	}
	defer fin.Close()
	inputReader := bufio.NewReader(fin)
	cols := -1
	for {
		line, read_err := inputReader.ReadString('\n')
		if read_err == io.EOF {
			break
		} else if read_err != nil {
			err = read_err
			return
		}
		line = strings.TrimSpace(line)
		line = strings.Replace(line, "\r", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		items := strings.Split(line, " ")
		if cols < 0 {
			cols = len(items)
		}
		data_l := make([]int, cols)
		for idx, v := range items {
			value, _ := strconv.Atoi(v)
			data_l[idx] = value
		}
		data = append(data, data_l)
	}
	return
}

func StringCombine(data []int) string {
	if nil == data {
		return ""
	}
	var buffer bytes.Buffer //Buffer是一个实现了读写方法的可变大小的字节缓冲
	for _, v := range data {
		buffer.WriteString(fmt.Sprintf("%d ", v))
		/*
		   func (b *Buffer) WriteString(s string) (n int, err error)
		   Write将s的内容写入缓冲中，如必要会增加缓冲容量。返回值n为len(p)，err总是nil。如果缓冲变得太大，Write会采用错误值ErrTooLarge引发panic。
		*/
	}
	return buffer.String()
}

func SaveMat(file string, data common.Mat) (err error) {
	fout, err := os.Open(file)
	if nil != err {
		err = fmt.Errorf("Open file error||file_name=%s||err_msg=%s", file, err.Error())
		return
	}
	outputWriter := bufio.NewWriter(fout)
	defer func() {
		outputWriter.Flush()
		fout.Close()
	}()
	rows, _ := utils.GetMatSize(&data)
	for i := 0; i < rows; i++ {
		line := StringCombine(data[i])
		_, err = outputWriter.WriteString(line)
		if nil != err {
			return
		}
	}
	return
}
