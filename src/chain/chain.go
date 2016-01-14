package chain

import (
	"fmt"
	"led"
)

type token struct {
	chain      int
	revolution int
}

func make_token_chan() chan token {
	return make(chan token, 1)
}

type chanpair struct {
	in  chan token
	out chan token
}

type Chain struct {
	revolutions int
	procs       []chanpair
	led         led.Led
	done        chan bool
}

// Pass the token around and toggle LED state if current chain number
// is equal to toggle_chain.
func tokenpass(chans chanpair, maxrevs int, done chan bool, toggle_chain int, led led.Led) {

	for {
		tok, more := <-chans.in

		if more == false {
			close(chans.out)
			return
		}
		// fmt.Printf("token: %s\n", tok)

		new_chain := tok.chain + 1
		if new_chain != toggle_chain {
			chans.out <- token{new_chain, tok.revolution}
		} else {
			// already finishing?
			if tok.revolution < maxrevs {
				chans.out <- token{0, tok.revolution + 1}
				led.Toggle()
			} else {
				fmt.Printf("set %d %d done\n", maxrevs, toggle_chain)
				// close output
				close(chans.out)

				done <- true
				return
			}
		}
	}
}

func New(revs int, count int, led led.Led) *Chain {
	chain := &Chain{
		revs,
		make([]chanpair, count),
		led,
		make(chan bool),
	}

	first := make_token_chan()
	for i := range chain.procs {
		if i == 0 {
			// setup first
			chain.procs[i] = chanpair{
				first,
				make_token_chan(),
			}
		} else if i == (len(chain.procs) - 1) {
			// setup last
			chain.procs[i] = chanpair{
				chain.procs[i-1].out,
				first,
			}
		} else {
			chain.procs[i] = chanpair{
				chain.procs[i-1].out,
				make_token_chan(),
			}
		}
	}

	return chain
}

func (chain *Chain) Spawn() {
	for i := range chain.procs {
		go tokenpass(chain.procs[i], chain.revolutions, chain.done,
			len(chain.procs), chain.led)
	}
}

func (chain *Chain) Start() {
	chain.procs[0].in <- token{0, 0}
}

func (chain *Chain) Wait() {
	<-chain.done
	close(chain.done)
}
