package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	dimX = 3
	dimY = 4
	tens = 10
)

var (
	ErrInitValuesOutOfBounds = errors.New("init values out of board")
	ErrNoMoreMovements       = errors.New("no more movements")
	ErrNoSolution            = errors.New("no solution")
	ErrNoMoreSolutions       = errors.New("no more solutions")
	ErrOutOfBounds           = errors.New("out of bound")
	ErrCellNotFree           = errors.New("cell not free")
	// ErrCellNotMatching       = errors.New("cell not matching")
)

var movements = [][]int{
	{0, 0},
	{1, 2}, {2, 1}, {2, -1}, {1, -2},
	{-1, -2}, {-2, -1}, {-2, 1}, {-1, 2},
}

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
	if err := Run(dimX, dimY); err != nil {
		panic(err)
	}
}

func Run(dimensionX, dimensionY int) error {
	// revive:disable:unhandled-error

	logrus.Info("Start")
	logrus.Infof("Board (%d x %d)", dimensionX, dimensionY)
	// revive:enable:unhandled-error
	board, err := InitBoard(0, 0, dimensionX, dimensionY)
	if err != nil {
		logrus.Error(err)
		return err
	}

	sol := 0

	for err == nil {
		err = SearchSolution(board)
		if err == nil {
			sol++
			board.Print(sol)
		} else if sol == 0 && errors.Is(err, ErrNoMoreSolutions) {
			// revive:disable:unhandled-error
			logrus.Info("No hay solución")
			return ErrNoSolution
			// revive:enable:unhandled-error
		}
	}
	logrus.Info("Finish") // revive:disable-line:unhandled-error

	return nil
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

	firstCell := &Cell{
		Position: 1,
		X:        x,
		Y:        y,
		Movement: 0,
	}

	board.Cells[0] = firstCell

	return board, nil
}

func prepareToNextSolution(board *Board) int {
	dim := board.DimX * board.DimY
	if board.Cells[dim-1] != nil { // prepare to next solution
		board.Cells[dim-1] = nil
		return dim - 2
	}

	return 0
}

func SearchSolution(board *Board) error {
	dim := board.DimX * board.DimY

	i := prepareToNextSolution(board)

	var err error
	for i >= 0 && i < dim-1 {
		err = nextStepSearchSolution(board, &i)
		if err != nil {
			return err
		}
	}

	return err
}

func nextStepSearchSolution(board *Board, i *int) error {
	cell := board.Cells[*i]
	cell.Movement = cell.Movement + 1
	nextCell, err := GetNextCell(board, cell, *i+1)
	if err == nil {
		*i++
		board.Cells[*i] = nextCell
	} else if errors.Is(err, ErrNoMoreSolutions) {
		return err
	} else if errors.Is(err, ErrNoMoreMovements) {
		board.Cells[*i] = nil
		*i--
	}

	return nil
}

// func GetNextCell(board *Board, x, y, position, movement int) (*Cell, error) {
func GetNextCell(board *Board, currentCell *Cell, position int) (*Cell, error) {
	// nolint:exhaustruct
	cell := &Cell{}

	if currentCell.Movement >= len(movements) {
		if position == 1 {
			return nil, ErrNoMoreSolutions
		}
		return nil, ErrNoMoreMovements
	}

	if currentCell.X+movements[currentCell.Movement][0] < 0 ||
		currentCell.X+movements[currentCell.Movement][0] > board.DimX-1 ||
		currentCell.Y+movements[currentCell.Movement][1] < 0 ||
		currentCell.Y+movements[currentCell.Movement][1] > board.DimY-1 {
		return nil, ErrOutOfBounds
	}
	cell.X = currentCell.X + movements[currentCell.Movement][0]
	cell.Y = currentCell.Y + movements[currentCell.Movement][1]
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

func prepareSolutionToPrint(board *Board) [][]string {
	solution := make([][]string, board.DimX)
	for i := range solution {
		solution[i] = make([]string, board.DimY)
	}

	for _, cell := range board.Cells {
		initCell := ""
		if cell.Y == 0 {
			initCell = "|"
		}
		extraTens := ""
		if cell.Position < tens {
			extraTens = "0"
		}
		solution[cell.X][cell.Y] =
			initCell + " " + extraTens + fmt.Sprint(cell.Position) + " |"
	}
	return solution
}

func (board *Board) Print(sol int) {
	logrus.Infof("Solución %d", sol)

	solution := prepareSolutionToPrint(board)

	rowToSeparateRows := ""
	for indexi, row := range solution {
		rowToPrint := ""
		for _, pos := range row {
			rowToPrint += pos
		}
		if indexi == 0 {
			length := len(strings.Trim(rowToPrint, " "))
			for i := 0; i < length; i++ {
				rowToSeparateRows += "-"
			}
			logrus.Info(rowToSeparateRows)
		}
		logrus.Info(rowToPrint)
		logrus.Info(rowToSeparateRows)
	}
}
