package testing

import (
	"context"
	"sync"
	"testing"
)

func TestFull2P(t *testing.T) {
	ctx := context.Background()
	c, close := createBufferedServer(t, ctx)
	defer close()

	gameID := gameForPlayers(t, ctx, c, []string{"jt", "stef"})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		playMovesUntilDone(t, ctx, c, "jt", gameID)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		playMovesUntilDone(t, ctx, c, "stef", gameID)
	}()

	wg.Wait()
}
