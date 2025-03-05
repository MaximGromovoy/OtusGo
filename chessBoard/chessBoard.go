package chessBoard

import "fmt"

type chessBoard struct {
	size uint
}

func New(s uint) *chessBoard {
	return &chessBoard{size: s}
}

func (cb chessBoard) PrintBoard() {
	for i := uint(0); i < cb.size; i++ {
		for j := uint(0); j < cb.size; j++ {
			if (i+j)%2 == 0 {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
