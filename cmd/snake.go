package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	term "github.com/nsf/termbox-go"
)

//GridSize is the size of the grid
const GridSize = 20

//SnakeCapacity is the max snake size
const SnakeCapacity = GridSize*GridSize - 1

//RefreshRate is the frequency with which we redraw the grid
const RefreshRate = 100 * time.Millisecond

//IsPaused is a flag to determine if the game is paused or not
var IsPaused = false

//Location struct to make maintaining locations easier
type Location struct {
	X int
	Y int
}

//Board holds the grid
type Board struct {
	Grid [][]byte
}

//Snake holds the direction that a snake is going as well as the locations of all it's segments
type Snake struct {
	Segments  []Location
	Direction Location
}

func initializeBoard() Board {
	myGrid := make([][]byte, GridSize)
	for i := 0; i < GridSize; i++ {
		myGrid[i] = make([]byte, GridSize)
	}

	board := Board{
		Grid: myGrid,
	}
	board.addFood()
	return board
}

func initializeSnake() Snake {
	mySnake := Snake{
		Segments: make([]Location, 1, SnakeCapacity),
		Direction: Location{
			X: 1,
			Y: 0,
		},
	}
	for i := range mySnake.Segments {
		mySnake.Segments[i].X = 0
		mySnake.Segments[i].Y = 0
	}
	return mySnake
}

func (snake *Snake) updateDirection() {
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			if ev.Key == term.KeyCtrlC {
				endGame(snake)
			}
			if ev.Key == term.KeyCtrlP {
				pauseGame()
			}
			switch ev.Key {
			case term.KeyArrowRight:
				snake.Direction.Y = 1
				snake.Direction.X = 0
			case term.KeyArrowLeft:
				snake.Direction.Y = -1
				snake.Direction.X = 0
			case term.KeyArrowUp:
				snake.Direction.Y = 0
				snake.Direction.X = -1
			case term.KeyArrowDown:
				snake.Direction.Y = 0
				snake.Direction.X = 1
			}
		}
	}
}

func (snake *Snake) isInBounds(board *Board) bool {
	if (snake.Segments[0].X+snake.Direction.X < len(board.Grid) && snake.Segments[0].Y+snake.Direction.Y < len(board.Grid[0])) && (snake.Segments[0].X+snake.Direction.X >= 0 && snake.Segments[0].Y+snake.Direction.Y >= 0) {
		return true
	}
	return false
}

func (snake *Snake) updateHeadLocation(board *Board) {
	if snake.isInBounds(board) {
		snake.Segments[0].X += snake.Direction.X
		snake.Segments[0].Y += snake.Direction.Y
	}
}

func (snake *Snake) updateSegments(board *Board) {
	prev := snake.Segments[0]
	for i := 1; i < len(snake.Segments); i++ {
		curr := snake.Segments[i]
		snake.Segments[i] = prev
		snake.Segments[i].X = prev.X
		snake.Segments[i].Y = prev.Y
		prev = curr
	}
	snake.updateHeadLocation(board)
}

func (snake *Snake) clear(board *Board) {
	for i := 0; i < len(snake.Segments); i++ {
		board.Grid[snake.Segments[i].X][snake.Segments[i].Y] = ' '
	}
}

func (snake *Snake) redraw(board *Board) {
	for i := 0; i < len(snake.Segments); i++ {
		board.Grid[snake.Segments[i].X][snake.Segments[i].Y] = 's'
	}
}

func main() {
	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()
	rand.Seed(time.Now().UTC().UnixNano())
	board := initializeBoard()
	snake := initializeSnake()
	refreshScreen(RefreshRate, board, snake)
}

func endGame(s *Snake) {
	fmt.Println("Game Over")
        printScore(s)
	os.Exit(2)
}

func pauseGame() {
	IsPaused = !IsPaused
}

func (board *Board) addFood() {
	randX := rand.Intn(GridSize - 1)
	randY := rand.Intn(GridSize - 1)
	if board.Grid[randX][randY] == 's' {
		board.addFood()
	} else {
		board.Grid[randX][randY] = 'f'
	}
}

func (snake *Snake) addSegment() {
	lenSnake := len(snake.Segments) + 1
	if lenSnake < SnakeCapacity {
		snake.Segments = snake.Segments[:lenSnake]
		snake.Segments[len(snake.Segments)-1] = Location{}
	}
	if lenSnake == SnakeCapacity {
		endGame(snake)
	}
}

func (snake *Snake) eatFood(board *Board) {
	target := Location{
		X: snake.Segments[0].X + snake.Direction.X,
		Y: snake.Segments[0].Y + snake.Direction.Y,
	}
	if !snake.isInBounds(board) {
		endGame(snake)
	} else if board.Grid[target.X][target.Y] == 's' {
		endGame(snake)
	} else if board.Grid[target.X][target.Y] == 'f' {
		board.Grid[target.X][target.Y] = ' '
		snake.addSegment()
		board.addFood()
	}
}

func (snake *Snake) act(board *Board) {
	snake.eatFood(board)
	snake.clear(board)
	snake.updateSegments(board)
	snake.redraw(board)
}

func printScore(snake *Snake) {
	fmt.Printf("Score: %d\n", 100*(len(snake.Segments)-1))
}

func (board *Board) print() {
	for i := 0; i < GridSize; i++ {
		fmt.Print("##")
	}
	fmt.Println()
	for i := 0; i < len(board.Grid); i++ {
		for j := 0; j < len(board.Grid[0]); j++ {
			if j == 0 {
				fmt.Print("#")
			}
			if board.Grid[i][j] == 's' {
				fmt.Print(" o")
			} else if board.Grid[i][j] == 'f' {
				fmt.Print(" *")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Print("#")
		fmt.Println()
	}
	for i := 0; i < GridSize; i++ {
		fmt.Print("##")
	}
	fmt.Println()
}

func refreshScreen(refreshRate time.Duration, board Board, snake Snake) {
	go snake.updateDirection()
	for true {
		if !IsPaused {
			fmt.Printf("x:%d, y:%d\n", snake.Segments[0].X, snake.Segments[0].Y)
			snake.act(&board)
			printScore(&snake)
			board.print()
			time.Sleep(refreshRate)
			clear()
		}
	}
}

func clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
