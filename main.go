package main

import (
	"./chessBoard"
)

func main() {
	CreateAndPrintBoard(10)
}

func CreateAndPrintBoard(boardSize uint) {
	chessBoard.New(boardSize).PrintBoard()
}
