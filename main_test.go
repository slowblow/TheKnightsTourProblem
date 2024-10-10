package main_test

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"
	"the-knights-tour-problem/internal/mocks"

	main "the-knights-tour-problem"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testCmd         = "cmd"
	testDimIntOne   = 1
	testDimIntThree = 3
	testDimIntFour  = 4
	testDimOne      = "1"
	testDimThree    = "3"
	testDimFour     = "4"
)

var testExpectedOutput3x4 = `level=info msg="Solución 1"
level=info msg=---------------------
level=info msg="| 01 | 04 | 07 | 10 |"
level=info msg=---------------------
level=info msg="| 12 | 09 | 02 | 05 |"
level=info msg=---------------------
level=info msg="| 03 | 06 | 11 | 08 |"
level=info msg=---------------------
`

var testSolution3x4 = []*main.Cell{
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

func TestIsFree(t *testing.T) {
	board, err := main.InitBoard(0, 0, 3, 4) // revive:disable-line:add-constant
	require.NoError(t, err)

	// Cleaning the board
	for i := range board.Cells {
		board.Cells[i] = nil
	}

	// Case 1: the new cell is not present on the board then is free
	cell := &main.Cell{
		X:        1,
		Y:        1,
		Position: 0,
		Movement: 0,
	}
	assert.True(t, cell.IsFree(board))

	// Caso 2: the new main.Cell is present on the board then is not free
	board.Cells[0] = &main.Cell{
		X:        1,
		Y:        1,
		Position: 0,
		Movement: 0,
	}

	assert.False(t, cell.IsFree(board))

	// Case 3: a new cell is not present on the board
	otherCell := &main.Cell{
		X:        2,
		Y:        2,
		Position: 0,
		Movement: 0,
	}
	assert.True(t, otherCell.IsFree(board))
}

// Test para la función getNextCell
func TestGetNextCell(t *testing.T) {
	board, err := main.InitBoard(0, 0, 3, 4) // revive:disable-line:add-constant
	require.NoError(t, err)

	// revive:disable:add-constant
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
	// revive:enable-line:add-constant

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/*
				cell, err := main.GetNextCell(
					board, tt.x, tt.y, tt.position, tt.movement,
				)
			*/
			cell := &main.Cell{
				X:        tt.x,
				Y:        tt.y,
				Position: tt.position,
				Movement: tt.movement,
			}
			cell, err := main.GetNextCell(
				board, cell, cell.Position,
			)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				// Verificamos que la celda generada coincida con la esperada
				if tt.expectedCell != nil {
					assert.NotNil(t, cell)
					assert.Equal(t, tt.expectedCell, cell)
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

	board, err := main.InitBoard(initCell.X, initCell.Y, 3, 4)
	require.NoError(t, err)
	assert.NotNil(t, board)
	assert.Equal(t, initCell, board.Cells[0])

	// Case 2: Valores fuera de los límites
	board, err = main.InitBoard(9, 9, 3, 4) // revive:disable-line:add-constant
	assert.Nil(t, board)
	assert.Equal(t, main.ErrInitValuesOutOfBounds, err)

	// Caso 3: end of board
	// revive:disable:add-constant
	initCell = &main.Cell{
		X:        2,
		Y:        3,
		Position: 1,
		Movement: 0,
	}
	board, err = main.InitBoard(initCell.X, initCell.Y, 3, 4)
	// revive:enable:add-constant
	require.NoError(t, err)
	assert.NotNil(t, board)
	assert.Equal(t, initCell, board.Cells[0])
}

func TestSearchSolution(t *testing.T) {
	// Init board
	board, err := main.InitBoard(0, 0, 3, 4) // revive:disable-line:add-constant
	require.NoError(t, err)

	// Case 1: successful solution
	lastCell := &main.Cell{
		X:        1, // revive:disable-line:add-constant
		Y:        2, // revive:disable-line:add-constant
		Position: 2, // revive:disable-line:add-constant
		Movement: 4, // revive:disable-line:add-constant
	}
	err = main.SearchSolution(board)
	require.NoError(t, err)
	assert.NotNil(t, board.Cells[1])
	assert.Equal(t, lastCell, board.Cells[1])

	// Case 2: Error no more movements
	for index := range board.Cells {
		board.Cells[index] = nil
	}
	// revive:disable:add-constant
	board.Cells[0] = &main.Cell{X: 0, Y: 0, Position: 1, Movement: 8}
	// revive:enable:add-constant

	err = main.SearchSolution(board)
	require.ErrorIs(t, err, main.ErrNoMoreSolutions)

	// Init board
	board, err = main.InitBoard(0, 0, 3, 4) // revive:disable-line:add-constant
	require.NoError(t, err)

	cellsCopy := make([]*main.Cell, len(testSolution3x4))
	copy(cellsCopy, testSolution3x4)

	board.Cells = cellsCopy
	err = main.SearchSolution(board)
	require.NoError(t, err)
	assert.NotNil(t, board.Cells[len(board.Cells)-1])
}

func TestPrint(t *testing.T) {
	// Config a buffer to bring the logs
	var buf bytes.Buffer

	// Replace the global logrus logger
	logrus.SetOutput(&buf)
	//nolint:exhaustruct
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		FullTimestamp:          false,
	})
	defer logrus.SetOutput(os.Stdout) // Restore the output

	board, err := main.InitBoard(0, 0, 3, 4) // revive:disable-line:add-constant
	require.NoError(t, err)

	cellsCopy := make([]*main.Cell, len(testSolution3x4))
	copy(cellsCopy, testSolution3x4)

	board.Cells = cellsCopy

	// Call the Print method with the buffer
	board.Print(1)

	// Check if the output matches
	output := buf.String()
	assert.Equal(t, testExpectedOutput3x4, output)
}

func TestRun(t *testing.T) {
	initValues := new(main.InitValues)

	runner := &main.Runner{
		InitValues: initValues,
	}

	err := runner.Run()
	require.ErrorIs(t, err, main.ErrInitValuesOutOfBounds)

	// revive:disable:add-constant
	initValues.LoadValues(0, 0, 4, 3)
	runner.InitValues = initValues
	err = runner.Run()
	require.NoError(t, err)

	initValues.LoadValues(0, 0, 2, 2)
	err = runner.Run()
	require.ErrorIs(t, err, main.ErrNoSolution)
	// revive:enable:add-constant
}

func TestGetInitValues(t *testing.T) {
	/*
		ValidArgs
	*/

	// test values
	initX := testDimIntOne
	initY := testDimIntOne
	dimX := testDimIntThree
	dimY := testDimIntFour

	// simulate arguments in command line
	os.Args = []string{
		"cmd",
		fmt.Sprint(initX),
		fmt.Sprint(initY),
		fmt.Sprint(dimX),
		fmt.Sprint(dimY),
	}

	//nolint:exhaustruct
	runner := &main.Runner{}

	err := runner.GetInitValues()
	require.NoError(t, err)

	// Using reflection in order to access to private fields
	values := reflect.ValueOf(runner.InitValues).Elem()

	tests := []struct {
		field    string
		expected int
	}{
		{"initX", initX},
		{"initY", initY},
		{"dimX", dimX},
		{"dimY", dimY},
	}

	for _, test := range tests {
		fieldValue := values.FieldByName(test.field)

		assert.True(t, fieldValue.IsValid())
		assert.Equal(t, int64(test.expected), fieldValue.Int())
	}

	/*
		InsufficientArgs
	*/
	os.Args = []string{"cmd", testDimOne, testDimOne}

	err = runner.GetInitValues()
	require.Error(t, err)
	require.ErrorIs(t, err, main.ErrInitArgumentsRequired)

	/*
		InvalidArgs
	*/
	invalidArgTests := []struct {
		args   []string
		errMsg string
	}{
		{
			args: []string{
				testCmd, "one", testDimOne, testDimThree, testDimFour,
			},
			errMsg: "invalid syntax: initX not valid",
		},
		{
			args: []string{
				testCmd, testDimOne, "one", testDimThree, testDimFour,
			},
			errMsg: "invalid syntax: initY not valid",
		},
		{
			args: []string{
				testCmd, testDimOne, testDimOne, "three", testDimFour,
			},
			errMsg: "invalid syntax: dimX not valid",
		},
		{
			args: []string{
				testCmd, testDimOne, testDimOne, testDimThree, "four",
			},
			errMsg: "invalid syntax: dimY not valid",
		},
	}

	for _, test := range invalidArgTests {
		os.Args = test.args
		err := runner.GetInitValues()
		require.Error(t, err)
		require.ErrorContains(t, err, test.errMsg)
	}
}

func TestRealMain(t *testing.T) {
	// Successful
	mock := mocks.NewMainRunnable(t)

	mock.EXPECT().GetInitValues().Return(nil)
	mock.EXPECT().Run().Return(nil)

	retorno := main.RealMain(mock)
	assert.Equal(t, main.StateSuccess, retorno)

	// Error on init values
	mock = mocks.NewMainRunnable(t)

	mock.EXPECT().GetInitValues().Return(main.ErrInitArgumentsRequired)
	retorno = main.RealMain(mock)
	assert.Equal(t, main.StateInitValuesFail, retorno)

	// Error on run task
	mock = mocks.NewMainRunnable(t)

	mock.EXPECT().GetInitValues().Return(nil)
	mock.EXPECT().Run().Return(main.ErrNoSolution)

	retorno = main.RealMain(mock)
	assert.Equal(t, main.StateRunFail, retorno)
}
