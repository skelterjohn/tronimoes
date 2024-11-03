package game

import "errors"

var (
	ErrGameAlreadyStarted       = errors.New("game already started")
	ErrGameTooManyPlayers       = errors.New("game already has 6 players")
	ErrGameNotEnoughPlayers     = errors.New("not enough players")
	ErrGamePreviousRoundNotDone = errors.New("previous round not done")
	ErrPlayerAlreadyInGame      = errors.New("player already in game")
	ErrRoundAlreadyDone         = errors.New("round already done")
	ErrRoundNotStarted          = errors.New("round not started")
	ErrNotYourTurn              = errors.New("not your turn")
	ErrNotYourGame              = errors.New("not your game")
	ErrNotYou                   = errors.New("not you")
	ErrMissingToken             = errors.New("missing token")
	ErrNoSuchGame               = errors.New("no such game")
	ErrPlayerNotFound           = errors.New("player not found")
	ErrTileNotFound             = errors.New("tile not found")
	ErrNoTile                   = errors.New("no tile provided")
	ErrNotRoundLeader           = errors.New("not round leader")
	ErrTileOccluded             = errors.New("tile occluded")
	ErrNoLine                   = errors.New("no line to play on")
	ErrTileOutOfBounds          = errors.New("tile out of bounds")
	ErrGameOver                 = errors.New("game is over")
)
