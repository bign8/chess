package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/bign8/games"
	"github.com/bign8/games/impl"
)

// TODO: allow this to support 3+ player games
func play(slug string, x, y io.ReadWriteCloser) {
	fmt.Fprintln(x, "sFound one! Say hi.")
	fmt.Fprintln(y, "sFound one! Say hi.")

	// Setup socket managers for sockets
	xMan := createManager(x)
	yMan := createManager(y)
	xChat := xMan.Room('u')
	yChat := yMan.Room('u')
	xGame := xMan.Room('g')
	yGame := yMan.Room('g')
	_, isBot := y.(bot)

	// Setup player chat-room // TODO: handle > 2 players
	errc := make(chan error, 1)
	go cp(xChat, yChat, errc)
	go cp(yChat, xChat, errc)

	// Convert actors to real players
	// TODO: support this for 3+ game players
	i := -1
	actors := []struct {
		msg   io.ReadWriteCloser
		isBot bool
	}{
		{xGame, false},
		{yGame, isBot},
	}
	builder := func(g games.Game, name string) games.Actor {
		i++
		if actors[i].isBot {
			return g.AI(g, name)
		}
		return newSocketActor(name, actors[i].msg, errc)
	}
	game, _ := impl.Get(slug)

	// Play the game
	data := game4client(games.Run(game, builder), true)
	xGame.Write(data) // Broadcast final game state
	yGame.Write(data)

	// Log errors if necessary
	if err := <-errc; err != nil {
		log.Println(err)
	}
	x.Close()
	y.Close()
}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(io.MultiWriter(w, chain), r) // copy chats to markov chain
	errc <- err
}

type actor struct {
	name  string
	s     *bufio.Scanner
	write io.Writer
}

func newSocketActor(name string, s io.ReadWriteCloser, errc chan<- error) *actor {
	a := &actor{
		name:  name,
		s:     bufio.NewScanner(s),
		write: s,
	}
	return a
}

func (a *actor) Name() string {
	return a.name
}

func (a *actor) Act(s games.State) games.Action {
	actions := s.Actions()
	a.write.Write(game4client(s, false))
	var chosen *games.Action
	for chosen == nil && a.s.Scan() {
		move := a.s.Text()
		for _, a := range actions {
			if a.String() == move {
				chosen = &a
				break
			}
		}
		if chosen == nil {
			a.write.Write([]byte("sInvalidMove... Try again!"))
		}
	}
	return *chosen
}

type gameMSG struct {
	SVG   string
	Moves []gameMoveMSG
}

type gameMoveMSG struct {
	Name string
	Type string
	SVG  string
}

func game4client(s games.State, done bool) []byte {
	moves := make([]gameMoveMSG, len(s.Actions()))

	for i, a := range s.Actions() {
		moves[i] = gameMoveMSG{
			Name: a.String(),
			Type: a.Type(),
			SVG:  s.Apply(a).SVG(false),
		}
	}

	data := gameMSG{
		SVG:   s.SVG(!done),
		Moves: moves,
	}
	js, _ := json.Marshal(data)
	return js
}
