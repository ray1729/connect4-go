package main

import (
    "bytes"
    "fmt"
//    "log"

//    "github.com/pkg/errors"
)

type Player int

const (
    Nobody Player = 0
    Red           = 1
    Yellow        = 2
)

const boardHeight = 6
const boardWidth = 7

type GameState struct {
  NextPlayer Player
  IsGameOver bool
  Winner     Player
  Columns    [boardWidth][]Player
}

func (s *GameState) Clone() *GameState {
  res := &GameState{NextPlayer: s.NextPlayer, IsGameOver: s.IsGameOver, Winner: s.Winner}
  for i := 0; i < boardWidth; i++ {
    res.Columns[i] = make([]Player, len(s.Columns[i]))
    copy(res.Columns[i], s.Columns[i])
  }
  return res
}

func (s *GameState) PlayerAt(row, col int) Player {
  if col < 0 || col >= boardWidth {
    return Nobody
  }
  if row < 0 || row >= len(s.Columns[col]) {
    return Nobody
  }
  return s.Columns[col][row]
}

func (s *GameState) CanPlay(col int) bool {
    if len(s.Columns[col]) < boardHeight {
        return true
    }
    return false
}

func (s *GameState) ValidMoves() []int {
    var res []int
    for i := 0; i < boardWidth; i++ {
        if s.CanPlay(i) {
            res = append(res, i)
        }
    }
    return res
}

func TogglePlayer(p Player) Player {
    if p == Red {
        return Yellow
    }
    return Red
}

func (s *GameState) MakeMove(col int) *GameState {
  res := s.Clone()
  res.Columns[col] = append(res.Columns[col], s.NextPlayer)
  if res.IsWinningMove(col) {
    res.Winner = s.NextPlayer
    res.IsGameOver = true
    res.NextPlayer = Nobody
  } else if len(res.ValidMoves()) == 0 {
    res.IsGameOver = true
    res.NextPlayer = Nobody
  } else {
    res.NextPlayer = TogglePlayer(s.NextPlayer)
  }
  return res
}

func (s *GameState) String() string {
  var buffer bytes.Buffer
  for r := boardHeight - 1; r >= 0; r-- {
    for c := 0; c < boardWidth; c++ {
      switch s.PlayerAt(r, c) {
      case Nobody:
        buffer.WriteString(".")
      case Red:
        buffer.WriteString("o")
      case Yellow:
        buffer.WriteString("x")
      }
      if c < boardWidth-1 {
        buffer.WriteString(" ")
      }
    }
    buffer.WriteString("\n")
  }
  return buffer.String()
}

func (s *GameState) IsWinningMove(col int) bool {
  return winningMoveInDirection(s, col, 1, 0) ||
    winningMoveInDirection(s, col, 0, 1) ||
    winningMoveInDirection(s, col, 1, 1) ||
    winningMoveInDirection(s, col, 1, -1)
}

func winningMoveInDirection(s *GameState, col, dx, dy int) bool {
  for n := -3; n <= 0; n++ {
    // (x,y) is the starting point of this possible winning run of 4
    x := col + n*dx
    y := len(s.Columns[col]) + n*dy - 1
    // If all 4 cells in the direction (dx, dy) from the starting position
    // are the next player, it's a win
    win := true
    for i := 0; i < 4; i++ {
      if s.PlayerAt(y, x) != s.NextPlayer {
        win = false
        break
      }
      x += dx
      y += dy
    }
    if win {
      return win
    }
  }
  return false
}

func main() {
  s := &GameState{NextPlayer: Red}
  fmt.Println(s)
  fmt.Println(s.ValidMoves())
  for i := 0; i < boardHeight; i++ {
    s  = s.MakeMove(0)
    fmt.Println(s)
  }
}
