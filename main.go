package main

import (
	"errors"
	"fmt"
	"strings"
)

const (
	DimMaxX = 8
	DimMaxY = 8
)

var (
	ErrNoMoreMovements = errors.New("No more movements.")
	ErrNoMoreSolutions = errors.New("No more solutions.")
	ErrOutOfBounds     = errors.New("Out of bounds.")
	ErrCellNotMatching = errors.New("Cell not matching.")
	ErrCellNotFree     = errors.New("Cell not free.")
)

var board = [DimMaxX * DimMaxY]*Cell{}

var movements = [][]int{{0, 0}, {1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}}

type Cell struct {
	Position int `json:"position"`
	Movement int `json:"movement"`
	X        int `json:"x"`
	Y        int `json:"y"`
}

func main() {

	fmt.Println("Start")
	firstCell := &Cell{}
	i := 0
	j := 0
	firstCell.Position = i*DimMaxX + j + 1
	firstCell.X = 0
	firstCell.Y = 0

	board[0] = firstCell

	sol := 0
	for index, cell := range board {
		if cell == nil {
			sol = index - 1
			break
		}
	}
	// fmt.Println(sol)
	var err error
	for err == nil {
		err = searchSolution()
		if err == nil {
			sol++
			print(sol)
		}
	}
	fmt.Println("Finish")
}

func searchSolution() error {
	var err error

	Dim := DimMaxX * DimMaxY

	i := 0
	if board[Dim-1] != nil {
		i = Dim - 2
		board[Dim-1] = nil
	}
	/*
		for index, cell := range board {
			i = index - 1
			if cell == nil {
				break
			}
		}
		fmt.Println(i)
	*/

	for i >= 0 && i < Dim-1 {
		cell := board[i]
		cell.Movement = cell.Movement + 1
		var nextCell *Cell
		nextCell, err = getNextCell(cell.X, cell.Y, i+1, cell.Movement)
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
			} else {
				break
			}
		}
	}

	return nil
}

func getNextCell(x, y, position, movement int) (*Cell, error) {
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

	if !cell.isFree() {
		return nil, ErrCellNotFree
	}

	return cell, nil
}

func print(sol int) {
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
				// fmt.Print("| ")
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
}

func (cell *Cell) isFree() bool {
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
