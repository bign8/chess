// Package games borrows terminology from "AI - A Modern Approach" Chapter 5
package games

import "fmt"

// SVGChooseMove is the JS function that needs to be called within a SVG for a move to be chosen
const SVGChooseMove = `parent.N8.games.chooseMove`

// Starter is a function used to create a game's initial state
type Starter func(...Actor) State

// ActorBuilder is a builder of actors
type ActorBuilder func(g Game, name string) Actor

// Action is the base type for a game move
type Action interface {
	fmt.Stringer
	Type() string // allows types of moves to be grouped
}

// Actor is a method that choose an Action given a particular State
type Actor interface {
	Name() string
	Act(State) Action
}

// State is the base type for the state of a game
type State interface {
	fmt.Stringer
	Actors() []Actor    // List of active actors for a game
	Player() int        // index of the active player given a State (also index in Utility array)
	Apply(Action) State // Applying an action to a game
	Actions() []Action  // List of available actions in a State
	Utility() []int     // If the game is in a terminal state return the utility for each Actor, else nil
	SVG(bool) string    // Browser representation of a state (bool: editable)
}

// Game is contains all the meta-data surrounding a game so it can be played
type Game struct {
	Name    string       // Name of the game
	Slug    string       // Short name of game
	Board   string       // SVG of board state
	Players []string     // List of Player names
	Counts  []uint8      // Possible number of players to play a game (if nil assume == len(Players))
	Start   Starter      `json:"-"`
	AI      ActorBuilder `json:"-"` // TODO: use smart enough ai that this can be removed
}

// Run is the primary game runner
func Run(g Game, ab ActorBuilder) (final State) {
	actors := make([]Actor, len(g.Players))
	for i, name := range g.Players {
		actors[i] = ab(g, name)
	}
	game := g.Start(actors...)
	for game.Utility() == nil {
		game = Play(game)
	}
	return game
}

// Play takes the game through the next phase
func Play(g State) State { return g.Apply(g.Actors()[g.Player()].Act(g)) }
