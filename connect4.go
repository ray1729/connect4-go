package main

import (
  "bytes"
  "fmt"
  "math/rand"
  "time"
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
  return col >= 0 && col < boardWidth && len(s.Columns[col]) < boardHeight
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

func token(p Player) string {
  if p == Red {
    return "o"
  }
  if p == Yellow {
    return "x"
  }
  return "."
}

func (s *GameState) String() string {
  var buffer bytes.Buffer
  for r := boardHeight - 1; r >= 0; r-- {
    for c := 0; c < boardWidth; c++ {
      p := s.PlayerAt(r, c)
      buffer.WriteString(token(p))
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

type Mover interface {
  NextMove(s *GameState) int
}

type RandomMover struct{}

func (m *RandomMover) NextMove(s *GameState) int {
  var winning, losing, neutral []int
MOVE:
  for _, move := range s.ValidMoves() {
    t := s.MakeMove(move)
    if t.Winner != Nobody {
      winning = append(winning, move)
      continue MOVE
    }
    for _, opponent_move := range t.ValidMoves() {
      u := t.MakeMove(opponent_move)
      if u.Winner != Nobody {
        losing = append(losing, move)
        continue MOVE
      }
    }
    neutral = append(neutral, move)
  }
  if len(winning) > 0 {
    i := rand.Intn(len(winning))
    return winning[i]
  }
  if len(neutral) > 0 {
    i := rand.Intn(len(neutral))
    return neutral[i]
  }
  i := rand.Intn(len(losing))
  return losing[i]
}

type ConsoleMover struct{}

func (m *ConsoleMover) NextMove(s *GameState) int {
  fmt.Print(s)
  fmt.Println("0 1 2 3 4 5 6")
  fmt.Printf("%s to play. Enter your move: ", token(s.NextPlayer))
  var i int
  for {
    _, err := fmt.Scan(&i)
    if err == nil && s.CanPlay(i) {
      break
    }
    fmt.Print("Please enter a valid move: ")
  }
  return i
}

type MonteCarloMover struct {
  level int
}

func (m *MonteCarloMover) NextMove(s *GameState) int {
  scores := make(map[int]int)
  randomMover := new(RandomMover)
  for _, move := range s.ValidMoves() {
    score := 0
    for i := 0; i < m.level; i++ {
      t := GameLoop(randomMover, randomMover, s.MakeMove(move))
      if t.Winner == s.NextPlayer {
        score += 1
      } else if t.Winner == TogglePlayer(s.NextPlayer) {
        score -= 1
      }
    }
    scores[move] = score
  }
  first := true
  var bestMove, bestScore int
  for move, score := range scores {
    if first {
      bestMove = move
      bestScore = score
      first = false
      continue
    }
    if score > bestScore {
      bestMove = move
      bestScore = score
    }
  }
  return bestMove
}

func GameLoop(red, yellow Mover, s *GameState) *GameState {
  for !s.IsGameOver {
    if s.NextPlayer == Red {
      s = s.MakeMove(red.NextMove(s))
    } else {
      s = s.MakeMove(yellow.NextMove(s))
    }
  }
  return s
}

func main() {
  rand.Seed(time.Now().UnixNano())
  p1 := new(ConsoleMover)
  p2 := &MonteCarloMover{level: 50}
  s := GameLoop(p1, p2, &GameState{NextPlayer: Red})
  fmt.Println(s)
  if s.Winner == Nobody {
    fmt.Println("Game over! It's a draw.")
  } else {
    fmt.Printf("Game over! %s won\n", token(s.Winner))
  }
}
