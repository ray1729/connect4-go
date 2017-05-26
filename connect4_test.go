package main

import (
  //"fmt"
  "testing"
  "github.com/stretchr/testify/assert"
)

func playMoves(firstPlayer Player, moves []int) *GameState {
  s := &GameState{NextPlayer: firstPlayer}
  for _, m := range moves {
    if s.IsGameOver {
      return s
    }
    s = s.MakeMove(m)
  }
  return s
}

func assertRedWins(t *testing.T, moves []int) {
  s := playMoves(Red, moves)
  assert := assert.New(t)
  assert.True(s.IsGameOver)
  assert.Equal(Player(Nobody), s.NextPlayer)
  assert.Equal(Player(Red), s.Winner)
}

func Test_WinningMove(t *testing.T) {
  assertRedWins(t, []int{0,0,1,0,2,0,3})
  assertRedWins(t, []int{0,1,0,2,0,3,0})
  assertRedWins(t, []int{0,1,1,2,3,2,2,3,3,4,3})
  assertRedWins(t, []int{1,0,0,0,0,1,1,2,2,4,3})
}
