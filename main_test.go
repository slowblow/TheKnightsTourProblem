package main_test

import (
	"testing"

	main "the-knights-tour-problem"

	"github.com/stretchr/testify/assert"
)

func Test_IsFree(t *testing.T) {

	board, err := main.InitBoard(0, 0)
	assert.NoError(t, err)

	// Cleaning the board
	for i := range board {
		board[i] = nil
	}

	// Case 1: the new cell is not present on the board then is free
	cell := &main.Cell{X: 1, Y: 1}
	assert.Equal(t, cell.IsFree(board), true)

	// Caso 2: the new main.Cell is present on the board then is not free
	board[0] = &main.Cell{X: 1, Y: 1}

	assert.Equal(t, cell.IsFree(board), false)

	// Case 3: a new cell is not present on the board
	otherCell := &main.Cell{X: 2, Y: 2}
	assert.Equal(t, otherCell.IsFree(board), true)
}

// Test para la función getNextCell
func TestGetNextCell(t *testing.T) {
	board, err := main.InitBoard(0, 0)
	assert.NoError(t, err)

	tests := []struct {
		name                     string
		x, y, position, movement int
		expectedErr              error
		expectedCell             *main.Cell
	}{
		{
			name:        "Valid movement in the board",
			x:           0,
			y:           0,
			position:    1,
			movement:    1, // movimiento {2, 1}
			expectedErr: nil,
			expectedCell: &main.Cell{
				X:        1,
				Y:        2,
				Position: 2,
				Movement: 0,
			},
		},
		{
			name:         "Movement oit of bounds",
			x:            7,
			y:            7,
			position:     1,
			movement:     0, // movimiento {2, 1}, fuera del tablero
			expectedErr:  main.ErrOutOfBounds,
			expectedCell: nil,
		},
		{
			name:         "Cell not free",
			x:            0,
			y:            0,
			position:     1,
			movement:     0, // movimiento {2, 1}, pero esa celda está ocupada
			expectedErr:  main.ErrCellNotFree,
			expectedCell: nil,
		},
		{
			name:         "No more solutions",
			x:            0,
			y:            0,
			position:     1,
			movement:     9, // Movimiento mayor que 8
			expectedErr:  main.ErrNoMoreSolutions,
			expectedCell: nil,
		},
		{
			name:         "No more movements",
			x:            1,
			y:            2,
			position:     2,
			movement:     9, // Movimiento mayor que 8
			expectedErr:  main.ErrNoMoreMovements,
			expectedCell: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell, err := main.GetNextCell(board, tt.x, tt.y, tt.position, tt.movement)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				// Verificamos que la celda generada coincida con la esperada
				if tt.expectedCell != nil {
					assert.NotNil(t, cell)
					assert.Equal(t, cell, tt.expectedCell)
					/*
						if cell.X != tt.expectedCell.X || cell.Y != tt.expectedCell.Y || cell.Position != tt.expectedCell.Position {
							t.Errorf("expected cell %+v, got %+v", tt.expectedCell, cell)
						}
					*/
				}
			}

		})
	}
}

func TestInitBoard(t *testing.T) {
	// Case 1: Values into the board
	initCell := &main.Cell{
		X:        0,
		Y:        0,
		Position: 1,
		Movement: 0,
	}

	board, err := main.InitBoard(initCell.X, initCell.Y)
	assert.Nil(t, err)
	assert.NotNil(t, board)
	assert.Equal(t, initCell, board[0])

	// Case 2: Valores fuera de los límites
	board, err = main.InitBoard(9, 9) // DimMaxX y DimMaxY are 8
	assert.Nil(t, board)
	assert.Equal(t, main.ErrInitValuesOutOfBounds, err)

	// Caso 3: end of board
	initCell = &main.Cell{
		X:        main.DimMaxX - 1,
		Y:        main.DimMaxY - 1,
		Position: 1,
		Movement: 0,
	}
	board, err = main.InitBoard(initCell.X, initCell.Y)
	assert.Nil(t, err)
	assert.NotNil(t, board)
	assert.Equal(t, initCell, board[0])
}

func TestSearchSolution(t *testing.T) {
	// Init board
	board, err := main.InitBoard(0, 0)
	assert.NoError(t, err)

	// Case 1: successful solution
	lastCell := &main.Cell{
		X:        1,
		Y:        2,
		Position: 2,
		Movement: 4,
	}
	err = main.SearchSolution(board)
	assert.Nil(t, err)
	assert.NotNil(t, board[1])
	assert.Equal(t, lastCell, board[1])

	// Case 2: Error no more movements
	for index := range board {
		board[index] = nil
	}
	board[0] = &main.Cell{X: 0, Y: 0, Position: 1, Movement: 8}

	err = main.SearchSolution(board)
	assert.ErrorIs(t, err, main.ErrNoMoreSolutions)

	/*
		// Case 3: Out of bounds error
		for index := range board {
			board[index] = nil
		}

		board[0] = &Cell{X: DimMaxX - 1, Y: DimMaxY - 1, Position: 2, Movement: 8}
		err = searchSolution(board)
		assert.ErrorIs(t, err, ErrNoMoreMovements)
	*/

	/*

		---------------------
		| 01 | 04 | 07 | 10 |
		---------------------
		| 12 | 09 | 02 | 05 |
		---------------------
		| 03 | 06 | 11 | 08 |
		---------------------

		X: 0, Y: 0, position: 1, movement: 1
		X: 1, Y: 2, position: 2, movement: 4
		X: 2, Y: 0, position: 3, movement: 7
		X: 0, Y: 1, position: 4, movement: 1
		X: 1, Y: 3, position: 5, movement: 4
		X: 2, Y: 1, position: 6, movement: 7
		X: 0, Y: 2, position: 7, movement: 2
		X: 2, Y: 3, position: 8, movement: 5
		X: 1, Y: 1, position: 9, movement: 8
		X: 0, Y: 3, position: 10, movement: 3
		X: 2, Y: 2, position: 11, movement: 5
		X: 1, Y: 0, position: 12, movement: 0
	*/

	// Init board
	board, err = main.InitBoard(0, 0)
	assert.NoError(t, err)
	board = []*main.Cell{
		{
			X:        0,
			Y:        0,
			Position: 1,
			Movement: 1,
		}, {
			X:        1,
			Y:        2,
			Position: 2,
			Movement: 4,
		}, {
			X:        2,
			Y:        0,
			Position: 3,
			Movement: 7,
		}, {
			X:        0,
			Y:        1,
			Position: 4,
			Movement: 1,
		}, {
			X:        1,
			Y:        3,
			Position: 5,
			Movement: 4,
		}, {
			X:        2,
			Y:        1,
			Position: 6,
			Movement: 7,
		}, {
			X:        0,
			Y:        2,
			Position: 7,
			Movement: 2,
		}, {
			X:        2,
			Y:        3,
			Position: 8,
			Movement: 5,
		}, {
			X:        1,
			Y:        1,
			Position: 9,
			Movement: 8,
		}, {
			X:        0,
			Y:        3,
			Position: 10,
			Movement: 3,
		}, {
			X:        2,
			Y:        2,
			Position: 11,
			Movement: 5,
		}, {
			X:        1,
			Y:        0,
			Position: 12,
			Movement: 0,
		},
	}
	err = main.SearchSolution(board)
	assert.NoError(t, err)
	print(2, board)
	assert.NotNil(t, board[len(board)-1])

}
