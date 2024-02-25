package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	quit int = iota
	up
	down
	left
	right

	quitKey  byte = 'q'
	upKey    byte = 'w'
	downKey  byte = 's'
	leftKey  byte = 'a'
	rightKey byte = 'd'
)

type game struct {
	board            [][]int
	goal, goalDigits int
}

func ceilDiv2(n int) int {
	if float64(n)/2 > float64(n) {
		return n/2 + 1
	}

	return n / 2
}

func pow2(n int) int {
	return int(math.Pow(2, float64(n)))
}

func digits(n int) int {
	var x, count int

	for x, count = 10, 1; x <= n; x, count = x*10, count+1 {
	}

	return count
}

func rawTerm() func() {
	var (
		stdin    int
		oldState *term.State
		err      error
	)

	stdin = int(os.Stdin.Fd())

	oldState, err = term.MakeRaw(stdin)
	if err != nil {
		panic(err)
	}

	return func() {
		term.Restore(stdin, oldState)
	}
}

func getDirection() int {
	var (
		reader *bufio.Reader
		char   byte
		err    error
	)

	reader = bufio.NewReader(os.Stdin)

	for {
		char, err = reader.ReadByte()
		if err != nil {
			panic(err)
		}

		switch char {
		case quitKey:
			return quit
		case upKey:
			return up
		case downKey:
			return down
		case leftKey:
			return left
		case rightKey:
			return right
		}
	}
}

func transpose(board [][]int) {
	var size, x, y int

	size = len(board)

	for y = 0; y < size-1; y++ {
		for x = y + 1; x < size; x++ {
			board[y][x], board[x][y] = board[x][y], board[y][x]
		}
	}
}

func rotate90(board [][]int) {
	var size, x, y, z int

	size = len(board)
	transpose(board)

	for y = 0; y < size; y++ {
		for x, z = 0, size-1; x < z; x, z = x+1, z-1 {
			board[y][x], board[y][z] = board[y][z], board[y][x]
		}
	}
}

func rotateNeg90(board [][]int) {
	var size, x, y, z int

	size = len(board)
	transpose(board)

	for x = 0; x < size; x++ {
		for y, z = 0, size-1; y < z; y, z = y+1, z-1 {
			board[y][x], board[z][x] = board[z][x], board[y][x]
		}
	}
}

func rotate180(board [][]int) {
	var size, x, y int

	size = len(board)

	for y = 0; y < ceilDiv2(size); y++ {
		for x = 0; x < size; x++ {
			board[y][x], board[size-1-y][size-1-x] = board[size-1-y][size-1-x], board[y][x]
		}
	}
}

func mergeLeft(board [][]int) bool {
	var (
		size, x, y, blockX int
		ok                 bool
	)

	size = len(board)

	for y = 0; y < size; y++ {
		blockX = -1

		for x = 0; x < size; x++ {
			if board[y][x] == 0 {
				continue
			}

			if blockX != -1 && board[y][blockX] == board[y][x] {
				board[y][blockX], board[y][x] = board[y][blockX]+1, 0
				ok = true
				continue
			}

			blockX++
			board[y][x], board[y][blockX] = 0, board[y][x]
			ok = ok || x != blockX
		}
	}

	return ok
}

func mergeRight(board [][]int) bool {
	var ok bool

	rotate180(board)
	ok = mergeLeft(board)
	rotate180(board)

	return ok
}

func mergeUp(board [][]int) bool {
	var ok bool

	rotateNeg90(board)
	ok = mergeLeft(board)
	rotate90(board)

	return ok
}

func mergeDown(board [][]int) bool {
	var ok bool

	rotate90(board)
	ok = mergeLeft(board)
	rotateNeg90(board)

	return ok
}

func randBlock(board [][]int) {
	var (
		emptyY     [][]int
		emptyX     []int
		size, x, y int
	)

	size = len(board)
	emptyY = make([][]int, 0, size)

	for y = 0; y < size; y++ {
		emptyX = make([]int, 0, size+1)

		for x = 0; x < size; x++ {
			if board[y][x] == 0 {
				emptyX = append(emptyX, x)
			}
		}

		if len(emptyX) != 0 {
			emptyX = append(emptyX, y)
			emptyY = append(emptyY, emptyX)
		}
	}

	emptyX = emptyY[rand.Intn(len(emptyY))]
	x = len(emptyX) - 1
	board[emptyX[x]][emptyX[rand.Intn(x)]] = 1
}

func moveTo(board [][]int, direction int) bool {
	var ok bool

	switch direction {
	case up:
		ok = mergeUp(board)
	case down:
		ok = mergeDown(board)
	case left:
		ok = mergeLeft(board)
	case right:
		ok = mergeRight(board)
	}

	return ok
}

func canMove(board [][]int) ([][]int, bool) {
	var (
		oldBoard [][]int
		size, y  int
	)

	size = len(board)

	oldBoard = make([][]int, size)
	for y = range oldBoard {
		oldBoard[y] = make([]int, size)
		copy(oldBoard[y], board[y])
	}

	for y = up; y <= right; y++ {
		if moveTo(board, y) {
			return oldBoard, true
		}
	}

	return oldBoard, false
}

func win(board [][]int, goal int) bool {
	var size, x, y int

	size = len(board)

	for y = 0; y < size; y++ {
		for x = 0; x < size; x++ {
			if board[y][x] == goal {
				return true
			}
		}
	}

	return false
}

func highest(board [][]int) int {
	var size, x, y, high int

	size = len(board)

	for y = 0; y < size; y++ {
		for x = 0; x < size; x++ {
			if board[y][x] > high {
				high = board[y][x]
			}
		}
	}

	return high
}

func border(side string, size int) string {
	return strings.Repeat(side, size) + string(side[0]) + "\r\n"
}

func row(g *game, sideY string, y int) string {
	var (
		builder strings.Builder
		x, n    int
	)

	for x = 0; x < len(g.board); x++ {
		n = g.board[y][x]

		if n == 0 {
			builder.WriteString(sideY)
			continue
		}

		builder.WriteString(fmt.Sprintf("| %-*d ", g.goalDigits, pow2(n)))
	}

	builder.WriteString("|\r\n")

	return builder.String()
}

func printBoard(g *game) {
	var (
		sideY, borderX, borderY string
		builder                 strings.Builder
		size, y                 int
	)

	size = len(g.board)
	sideY = "|" + strings.Repeat(" ", g.goalDigits+2)
	borderX = border(strings.Repeat("-", g.goalDigits+3), size)
	borderY = border(sideY, size)

	builder.WriteString("\033[2J\033[H")

	for y = 0; y < size; y++ {
		builder.WriteString(borderX)
		builder.WriteString(borderY)
		builder.WriteString(row(g, sideY, y))
		builder.WriteString(borderY)
	}

	builder.WriteString(borderX)
	builder.WriteString(fmt.Sprintf("quit: %2c\r\n", quitKey))
	builder.WriteString(fmt.Sprintf("up: %4c\r\n", upKey))
	builder.WriteString(fmt.Sprintf("left: %2c\r\n", leftKey))
	builder.WriteString(fmt.Sprintf("down: %2c\r\n", downKey))
	builder.WriteString(fmt.Sprintf("right: %c\r\n", rightKey))
	builder.WriteString(fmt.Sprintf("goal:  %d / %d\r\n", pow2(highest(g.board)), pow2(g.goal)))

	fmt.Print(builder.String())
}

func mkGame(goal, size int) *game {
	var (
		g     *game
		board [][]int
		y     int
	)

	board = make([][]int, size)
	for y = 0; y < size; y++ {
		board[y] = make([]int, size)
	}

	g = &game{
		board:      board,
		goal:       goal,
		goalDigits: digits(pow2(goal)),
	}

	randBlock(board)
	randBlock(board)

	return g
}

func main() {
	var (
		g         *game
		direction int
		ok        bool
	)

	g = mkGame(11, 4)
	defer rawTerm()()

start:
	for {
		printBoard(g)

		if win(g.board, g.goal) {
			fmt.Print("you win\r\n")
			break
		}

		g.board, ok = canMove(g.board)
		if !ok {
			fmt.Print("no more moves\r\n")
			break
		}

		for {
			direction = getDirection()
			if direction == quit {
				break start
			}

			if moveTo(g.board, direction) {
				randBlock(g.board)
				break
			}
		}
	}
}
