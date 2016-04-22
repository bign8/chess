package chess

import "testing"

func BenchmarkStateApply(b *testing.B) {
	game := New()
	moves := game.Moves() // cache moves
	for i := 0; i < b.N; i++ {
		game.Apply(moves[0])
	}
}
