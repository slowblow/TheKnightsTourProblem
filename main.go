package main

import (
	"errors"
	"fmt"
	"strings"
)

const (
	DimMaxX = 3
	DimMaxY = 4
)

var (
	ErrInitValuesOutOfBounds = errors.New("init values out of board")
	ErrNoMoreMovements       = errors.New("no more movements")
	ErrNoMoreSolutions       = errors.New("no more solutions")
	ErrOutOfBounds           = errors.New("out of bound.")
	ErrCellNotFree           = errors.New("cell not free")
	// ErrCellNotMatching       = errors.New("cell not matching")
)

var movements = [][]int{{0, 0}, {1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}}

type Cell struct {
	Position int `json:"position"`
	Movement int `json:"movement"`
	X        int `json:"x"`
	Y        int `json:"y"`
}

func main() {
	fmt.Println("Start")
	fmt.Printf("Board (%d x %d)\n", DimMaxX, DimMaxY)
	board, err := InitBoard(0, 0)
	if err != nil {
		panic(err)
	}

	sol := 0

	for err == nil {
		err = SearchSolution(board)
		if err == nil {
			sol++
			print(sol, board)
		} else if sol == 0 && errors.Is(err, ErrNoMoreSolutions) {
			fmt.Println("No hay solución")
		}
	}
	fmt.Println("Finish")
}

func InitBoard(x, y int) ([]*Cell, error) {
	if x >= DimMaxX || y >= DimMaxY {
		return nil, ErrInitValuesOutOfBounds
	}
	board := make([]*Cell, DimMaxX*DimMaxY)

	firstCell := &Cell{}

	firstCell.Position = 1
	firstCell.X = x
	firstCell.Y = y

	board[0] = firstCell

	return board, nil
}

func SearchSolution(board []*Cell) error {
	var err error

	Dim := DimMaxX * DimMaxY

	i := 0
	if board[Dim-1] != nil { // prepare to next solution
		i = Dim - 2
		board[Dim-1] = nil
	}

	for i >= 0 && i < Dim-1 {
		cell := board[i]
		cell.Movement = cell.Movement + 1
		var nextCell *Cell
		nextCell, err = GetNextCell(board, cell.X, cell.Y, i+1, cell.Movement)
		if err == nil {
			i++
			board[i] = nextCell
		} else {
			if err == ErrNoMoreSolutions {
				return err
			}
			// break
			if err == ErrOutOfBounds || err == ErrCellNotFree {
				// cell.Movement = cell.Movement + 1
			} else if err == ErrNoMoreMovements {
				board[i] = nil
				i--
			} /* else {
				break
			} */
		}
	}

	return nil
}

func GetNextCell(board []*Cell, x, y, position, movement int) (*Cell, error) {
	cell := &Cell{}

	if movement > 8 {
		if position == 1 {
			return nil, ErrNoMoreSolutions
		}
		return nil, ErrNoMoreMovements
	}

	if x+movements[movement][0] < 0 || x+movements[movement][0] > DimMaxX-1 || y+movements[movement][1] < 0 || y+movements[movement][1] > DimMaxY-1 {
		return nil, ErrOutOfBounds
	}
	cell.X = x + movements[movement][0]
	cell.Y = y + movements[movement][1]
	cell.Position = position + 1
	cell.Movement = 0

	if !cell.IsFree(board) {
		return nil, ErrCellNotFree
	}

	return cell, nil
}

func Print(sol int, board []*Cell) {
	fmt.Println()
	fmt.Println()
	fmt.Println("Solución", sol)
	fmt.Println()
	solution := [DimMaxX][DimMaxY]int{}

	for _, cell := range board {
		solution[cell.X][cell.Y] = cell.Position
	}

	rowToSeparateRows := ""
	for indexi, row := range solution {
		rowToPrint := ""
		for indexj, pos := range row {
			if indexj == 0 {
				rowToPrint = "| "
			}
			extra := ""
			if pos < 10 {
				extra = "0"
			}
			rowToPrint += (extra + fmt.Sprint(pos) + " | ")
		}
		if indexi == 0 {
			length := len(strings.Trim(rowToPrint, " "))
			for i := 0; i < length; i++ {
				rowToSeparateRows += "-"
			}
			fmt.Println(rowToSeparateRows)
		}
		fmt.Println(rowToPrint)
		fmt.Println(rowToSeparateRows)
	}
	for _, cell := range board {
		fmt.Printf("\nX: %d, Y: %d, position: %d, movement: %d", cell.X, cell.Y, cell.Position, cell.Movement)
	}
}

func (cell *Cell) IsFree(board []*Cell) bool {
	for i := 0; i < DimMaxX*DimMaxY-1; i++ {
		if board[i] == nil {
			break
		}
		if board[i].X == cell.X && board[i].Y == cell.Y {
			return false
		}
	}
	return true
}
