package chess

// Utility generates the game state utility for the game.
func (s State) Utility() []int {
	val := ValueUtility(s)
	res := make([]int, 2)
	for i := 0; i < 2; i++ {
		if s.Player() == 0 {
			res[i] = -val
		} else {
			res[i] = val
		}
	}
	return res
}

// ValueUtility is a uses the common standard value of pieces to rate a state.
// https://en.wikipedia.org/wiki/Chess_piece_relative_value
func ValueUtility(s State) int { // TODO: convert to games.State
	ctr := 0
	for _, square := range s.board {
		switch square {
		case 'P':
			ctr++
		case 'p':
			ctr--
		case 'N':
			fallthrough
		case 'B':
			ctr += 3
		case 'n':
			fallthrough
		case 'b':
			ctr -= 3
		case 'R':
			ctr += 5
		case 'r':
			ctr -= 5
		case 'Q':
			ctr += 9
		case 'q':
			ctr -= 9
		}
	}
	return ctr
}
