// Copyright (c) 2015 Open-RnD Sp. z o.o.

// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use, copy,
// modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package chain

import (
	"fmt"
	"led"
	"runtime"
	"sync"
)

var (
	enableYield = false // use voluntary Gosched()
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
	wg          sync.WaitGroup
}

func (c *chanpair) getToken() (token, bool) {
	tok, more := <-c.in
	return tok, more
}

// Pass token to channel
func (c *chanpair) passToken(t token) {
	c.out <- t
}

// Close channel pair output
func (c *chanpair) closeOut() {
	// close output
	close(c.out)

}

// Pass the token around and toggle LED state if current chain number
// is equal to toggle_chain.
func tokenpass(chans *chanpair, maxrevs int, toggle_chain int, wg *sync.WaitGroup, led led.Led) {

	defer wg.Done()

	for {
		tok, more := chans.getToken()

		if more == false {
			chans.closeOut()
			return
		}

		// yield CPU
		if enableYield {
			runtime.Gosched()
		}

		// fmt.Printf("token: %s\n", tok)

		new_chain := tok.chain + 1
		if new_chain != toggle_chain {
			chans.passToken(token{new_chain, tok.revolution})
		} else {
			// already finishing?
			if tok.revolution < maxrevs {
				chans.passToken(token{0, tok.revolution + 1})
				led.Toggle()
			} else {
				fmt.Printf("set %d %d done\n", maxrevs, toggle_chain)
				chans.closeOut()

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
		sync.WaitGroup{},
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
		go tokenpass(&chain.procs[i], chain.revolutions,
			len(chain.procs), &chain.wg, chain.led)
	}
	chain.wg.Add(len(chain.procs))
}

func (chain *Chain) Start() {
	chain.procs[0].in <- token{0, 0}
}

func (chain *Chain) Wait() {
	chain.wg.Wait()
}

// Enable voluntary CPU yield
func SetYield(yield bool) {
	enableYield = yield
}
