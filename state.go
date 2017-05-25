package main

import (
    "bytes"
    "fmt"
    "log"

    "github.com/pkg/errors"
)

type Player int

const (
    Red    Player = 0
    Yellow        = 1
)

const boardHeight = 6
const boardWidth = 7

type GameState struct {
    NextPlayer Player
    Columns    [boardWidth][]Player
}

func NewGameState(p Player) *GameState {
    s := new(GameState)
    s.NextPlayer = p
    return s
}

func (s *GameState) Clone() *GameState {
    res := NewGameState(s.NextPlayer)
    for i := 0; i < boardWidth; i++ {
        res.Columns[i] = make([]Player, len(s.Columns[i]))
        copy(res.Columns[i], s.Columns[i])
    }
    return res
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

func (s *GameState) MakeMove(i int) (*GameState, error) {
    if !s.CanPlay(i) {
        return s, errors.New("Illegal move")
    }
    res := s.Clone()
    res.NextPlayer = TogglePlayer(s.NextPlayer)
    res.Columns[i] = append(res.Columns[i], s.NextPlayer)
    return res, nil
}

func (s *GameState) String() string {
    var buffer bytes.Buffer
    for r := boardHeight - 1; r >= 0; r-- {
        for c := 0; c < boardWidth; c++ {
            col := s.Columns[c]
            if len(col) <= r {
                buffer.WriteString(".")
            } else if col[r] == Red {
                buffer.WriteString("o")
            } else {
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

func main() {
    s := NewGameState(Red)
    var err error
    fmt.Println(s)
    fmt.Println(s.ValidMoves())
    for i := 0; i < boardHeight; i++ {
        s, err = s.MakeMove(0)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(s)
    }

}
