package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
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

type Board struct {
	Cells []*Cell
	DimX  int
	DimY  int
}

func main() {

	dimMaxX := 3
	dimMaxY := 4

	fmt.Println("Start")
	fmt.Printf("Board (%d x %d)\n", dimMaxX, dimMaxY)
	board, err := InitBoard(0, 0, dimMaxX, dimMaxY)
	if err != nil {
		panic(err)
	}

	sol := 0

	var buf bytes.Buffer

	for err == nil {
		err = SearchSolution(board)
		if err == nil {
			sol++
			board.Print(&buf, sol)
		} else if sol == 0 && errors.Is(err, ErrNoMoreSolutions) {
			fmt.Println("No hay solución")
		}
	}
	fmt.Println("Finish")
}

func InitBoard(x, y, dimMaxX, dimMaxY int) (*Board, error) {
	if x >= dimMaxX || y >= dimMaxY {
		return nil, ErrInitValuesOutOfBounds
	}
	cells := make([]*Cell, dimMaxX*dimMaxY)

	board := &Board{
		Cells: cells,
		DimX:  dimMaxX,
		DimY:  dimMaxY,
	}

	firstCell := &Cell{}

	firstCell.Position = 1
	firstCell.X = x
	firstCell.Y = y

	board.Cells[0] = firstCell

	return board, nil
}

func SearchSolution(board *Board) error {
	var err error

	Dim := board.DimX * board.DimY

	i := 0
	if board.Cells[Dim-1] != nil { // prepare to next solution
		i = Dim - 2
		board.Cells[Dim-1] = nil
	}

	for i >= 0 && i < Dim-1 {
		cell := board.Cells[i]
		cell.Movement = cell.Movement + 1
		var nextCell *Cell
		nextCell, err = GetNextCell(board, cell.X, cell.Y, i+1, cell.Movement)
		if err == nil {
			i++
			board.Cells[i] = nextCell
		} else {
			if err == ErrNoMoreSolutions {
				return err
			}
			// break
			if err == ErrOutOfBounds || err == ErrCellNotFree {
				// cell.Movement = cell.Movement + 1
			} else if err == ErrNoMoreMovements {
				board.Cells[i] = nil
				i--
			} /* else {
				break
			} */
		}
	}

	return nil
}

func GetNextCell(board *Board, x, y, position, movement int) (*Cell, error) {
	cell := &Cell{}

	if movement > 8 {
		if position == 1 {
			return nil, ErrNoMoreSolutions
		}
		return nil, ErrNoMoreMovements
	}

	if x+movements[movement][0] < 0 || x+movements[movement][0] > board.DimX-1 || y+movements[movement][1] < 0 || y+movements[movement][1] > board.DimY-1 {
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

func (cell *Cell) IsFree(board *Board) bool {
	for i := 0; i < board.DimX*board.DimY-1; i++ {
		if board.Cells[i] == nil {
			break
		}
		if board.Cells[i].X == cell.X && board.Cells[i].Y == cell.Y {
			return false
		}
	}
	return true
}

func (board *Board) Print(w io.Writer, sol int) {
	fmt.Fprintln(w)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Solución %d", sol)
	fmt.Fprintln(w)
	fmt.Fprintln(w)

	solution := make([][]int, board.DimX)
	for i := range solution {
		solution[i] = make([]int, board.DimY)
	}

	for _, cell := range board.Cells {
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
			fmt.Fprintln(w, rowToSeparateRows)
		}
		fmt.Fprintln(w, rowToPrint)
		fmt.Fprintln(w, rowToSeparateRows)
	}

	fmt.Fprintln(w)
}
