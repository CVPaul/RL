// this file is a sample implement of  pseudo random
package utils

type Rander struct {
	Next uint32
}

// using the persudo random function in C
func (self *Rander) Rand() int {
	self.Next = self.Next*1103515245 + 12345
	return int((self.Next / 65536) % 32768)
}
func (self *Rander) Srand(seed uint32) {
	self.Next = seed
}
