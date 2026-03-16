package game

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/skelterjohn/tronimoes/tronserv/clog"
	"google.golang.org/api/idtoken"
)

var fbApp *firebase.App

func init() {
	ctx := context.Background()
	var err error
	fbApp, err = firebase.NewApp(ctx, &firebase.Config{
		ProjectID: "tronimoes",
	})
	if err != nil {
		clog.Fatal(ctx, "Error initializing Firebase app", "error", err.Error())
	}
}

func writeErr(w http.ResponseWriter, err error, code int) {
	if errors.Is(err, ErrVersionConflict) {
		code = http.StatusConflict
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func RegisterHandlers(r chi.Router, gs *GameServer) {
	r.Get("/game/{code}", gs.HandleGetGame)
	r.Put("/game/{code}", gs.HandlePutGame)
	r.Post("/game/{code}/report", gs.HandleReportIssue)
	r.Post("/game/{code}/start", gs.HandleStartRound)
	r.Post("/game/{code}/tile", gs.HandleLayTile)
	r.Post("/game/{code}/spacer", gs.HandleLaySpacer)
	r.Post("/game/{code}/draw", gs.HandleDrawTile)
	r.Post("/game/{code}/pass", gs.HandlePass)
	r.Post("/game/{code}/leave", gs.HandleLeaveOrQuit)
	r.Post("/game/{code}/foot", gs.HandleChickenFoot)
	r.Post("/game/{code}/react", gs.HandleReact)
	r.Post("/players", gs.HandleRegisterPlayerName)
	r.Get("/players/{playerID}", gs.HandleGetPlayer)
	r.Put("/players/{playerID}/config", gs.HandleUpdatePlayerConfig)
}

func RandomString(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var reservedInitials = map[string]bool{
	"RC": true,
	"GP": true,
}

func isBotName(name string) bool {
	tokens := strings.Split(name, " ")
	var initials string
	for _, t := range tokens {
		initials += t[0:1]
	}

	return reservedInitials[initials]
}
func validatePlayerName(name string) error {
	if len(name) > 32 {
		return ErrPlayerNameTooLong
	}

	tokens := strings.Split(name, " ")
	var initials string
	for _, t := range tokens {
		initials += t
	}

	if reservedInitials[initials] {
		return ErrPlayerInitialsReserved
	}

	return nil
}

type GameOptions struct {
	AgentRoundOut int `json:"agent_round_out"`
}

// BotTokenAudience is the required audience claim for bot (service account) ID tokens.
const BotTokenAudience = "https://games.tronimoes.com"

type GameServer struct {
	Store                     Store
	AgentSpawner              AgentSpawner
	CheckBotTokens            bool
	AllowedBotServiceAccounts []string // if non-empty, token's email must be in this list
}

// validateBotServiceAccountToken checks the request for a Bearer token, validates it as a
// Google-issued ID token with audience BotTokenAudience, and optionally requires the token's
// service account email to be in AllowedBotServiceAccounts. Call when registering a bot player.
func (s *GameServer) validateBotServiceAccountToken(ctx context.Context, r *http.Request) error {
	token := r.Header.Get("Authorization")
	if token == "" {
		return ErrYouAreNotABot
	}
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return ErrYouAreNotABot
	}

	payload, err := idtoken.Validate(ctx, token, BotTokenAudience)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if len(s.AllowedBotServiceAccounts) > 0 {
		email, _ := payload.Claims["email"].(string)
		allowed := false
		for _, a := range s.AllowedBotServiceAccounts {
			if a == email {
				allowed = true
				break
			}
		}
		if !allowed {
			clog.Info(ctx, fmt.Sprintf("Bot service account not allowed: %s", email))
			return ErrYouAreNotABot
		}
	}

	return nil
}

func (s *GameServer) validateToken(ctx context.Context, r *http.Request) error {
	token := r.Header.Get("Authorization")
	if token == "" {
		return ErrMissingToken
	}
	userID := r.Header.Get("X-Player-Id")
	if userID == "" {
		return ErrMissingUserID
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Verify the Firebase token
	client, err := fbApp.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error getting Auth client: %v", err)
	}

	decodedToken, err := client.VerifyIDToken(ctx, token)
	if err != nil {
		return fmt.Errorf("error verifying ID token: %v", err)
	}

	// Verify that the token's UID matches the X-Player-Id
	if decodedToken.UID != userID {
		return ErrInvalidToken
	}

	return nil
}

func (s *GameServer) validatePlayerID(ctx context.Context, playerID string, r *http.Request) error {
	userID := r.Header.Get("X-Player-Id")
	if userID == "" {
		return ErrMissingUserID
	}
	if userID != playerID {
		return ErrNotYourPlayer
	}
	if err := s.validateToken(ctx, r); err != nil {
		return err
	}
	return nil
}

func (s *GameServer) getName(ctx context.Context, r *http.Request) (string, error) {
	userID := r.Header.Get("X-Player-Id")
	if userID != "" {
		if err := s.validateToken(ctx, r); err != nil {
			return "", err
		}

		pi, err := s.Store.GetPlayer(ctx, userID)
		if err == nil {
			return pi.Name, nil
		}
		return "", err
	}

	tempName := r.Header.Get("X-Player-Name")
	_, err := s.Store.GetPlayerByName(ctx, tempName)
	if err == ErrNoRegisteredPlayer {
		// anonymous play is ok with unregistered names.
		if err := validatePlayerName(tempName); err != nil {
			return "", err
		}
		return tempName, nil
	}
	if err == nil {
		return "", ErrNotYourPlayer
	}
	return "", err
}

func (s *GameServer) encodeFilteredGame(ctx context.Context, w http.ResponseWriter, name string, g *Game) {
	for _, p := range g.Players {
		if p.Name == name {
			continue
		}
		// Hide the hands of other players, though we still send the tile counts.
		for i := range p.Hand {
			p.Hand[i].PipsA = 0
			p.Hand[i].PipsB = 0
		}
	}
	// Hide the bag from everyone.
	for i := range g.Bag {
		g.Bag[i].PipsA = 0
		g.Bag[i].PipsB = 0
	}

	// Add legal moves for this player to see.
	r := g.CurrentRound(ctx)
	if len(g.Players) > g.Turn {
		p := g.Players[g.Turn]
		if !g.Done && r != nil && p.Name == name {
			r.FindHints(ctx, g, p)
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(g); err != nil {
		clog.Error(ctx, "Error encoding game", "error", err.Error(), "code", g.Code)
		writeErr(w, err, http.StatusInternalServerError)
	}
}

func (s *GameServer) HandleLeaveOrQuit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		clog.Error(ctx, "Error getting name", "error", err.Error())
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.Store.ReadGame(ctx, code)
	if err != nil {
		clog.Error(ctx, "Error reading game", "error", err.Error(), "code", code)
		if err == ErrNoSuchGame {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	if !g.LeaveOrQuit(ctx, name) {
		clog.Info(ctx, "Player cannot leave game", "name", name, "code", code)
		writeErr(w, ErrNotYourGame, http.StatusBadRequest)
		return
	}

	if err := s.Store.WriteGame(ctx, g); err != nil {
		clog.Error(ctx, "Error writing game", "error", err.Error(), "code", code)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleDrawTile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		clog.Error(ctx, "Error getting name", "error", err.Error())
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.Store.ReadGame(ctx, code)
	if err != nil {
		clog.Error(ctx, "Error reading game", "error", err.Error(), "code", code)
		if err == ErrNoSuchGame {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	player := g.Players[g.Turn]
	if player.Name != name {
		clog.Info(ctx, "Player not in turn for game", "name", name, "code", code)
		writeErr(w, ErrNotYourTurn, http.StatusBadRequest)
		return
	}

	if len(g.Rounds) == 0 {
		clog.Info(ctx, "Player tried to play but game not started", "name", name, "code", code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}
	round := g.Rounds[len(g.Rounds)-1]
	if round.Done {
		clog.Info(ctx, "Player tried to play but round is done", "name", name, "code", code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}

	if !g.DrawTile(ctx, name) {
		clog.Info(ctx, "Player tried to play but game not started", "name", name, "code", code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}

	if err := s.Store.WriteGame(ctx, g); err != nil {
		clog.Error(ctx, "Error writing game", "error", err.Error(), "code", code)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	s.encodeFilteredGame(ctx, w, name, g)
}

type ChickenFoot struct {
	SelectedX int `json:"selected_x"`
	SelectedY int `json:"selected_y"`
}

func (s *GameServer) HandlePass(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		clog.Error(ctx, "Error getting name", "error", err.Error())
		writeErr(w, err, http.StatusForbidden)
		return
	}

	chickenFoot := &ChickenFoot{}
	if err := json.NewDecoder(r.Body).Decode(chickenFoot); err != nil {
		clog.Error(ctx, "Error decoding chicken-foot placement", "error", err.Error(), "name", name, "code", code)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	g, err := s.Store.ReadGame(ctx, code)
	if err != nil {
		clog.Error(ctx, "Error reading game", "error", err.Error(), "code", code)
		if err == ErrNoSuchGame {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	player := g.Players[g.Turn]
	if player.Name != name {
		clog.Info(ctx, "Player not in turn for game", "name", name, "code", code)
		writeErr(w, ErrNotYourTurn, http.StatusBadRequest)
		return
	}

	if len(g.Rounds) == 0 {
		clog.Info(ctx, "Player tried to play but game not started", "name", name, "code", code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}
	round := g.Rounds[len(g.Rounds)-1]
	if round.Done {
		clog.Info(ctx, "Player tried to play but round is done", "name", name, "code", code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}

	if err := g.Pass(ctx, name, chickenFoot.SelectedX, chickenFoot.SelectedY); err != nil {
		clog.Error(ctx, "Could not pass", "error", err.Error())
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.Store.WriteGame(ctx, g); err != nil {
		clog.Error(ctx, "Error writing game", "error", err.Error(), "code", code)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleLayTile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		clog.Error(ctx, "Error getting name", "error", err.Error())
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.Store.ReadGame(ctx, code)
	if err != nil {
		clog.Error(ctx, "Error reading game", "error", err.Error(), "code", code)
		if err == ErrNoSuchGame {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	player := g.Players[g.Turn]
	if player.Name != name {
		clog.Info(ctx, "Player not in turn for game", "name", name, "code", code)
		writeErr(w, ErrNotYourTurn, http.StatusBadRequest)
		return
	}

	if len(g.Rounds) == 0 {
		clog.Info(ctx, "Player tried to play but game not started", "name", name, "code", code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}
	round := g.Rounds[len(g.Rounds)-1]
	if round.Done {
		clog.Info(ctx, "Player tried to play but round is done", "name", name, "code", code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}

	lt := &LaidTile{}
	if err := json.NewDecoder(r.Body).Decode(lt); err != nil {
		clog.Error(ctx, "Error decoding tile", "error", err.Error(), "name", name, "code", code)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	lt.PlayerName = player.Name

	if err := g.LayTile(ctx, name, lt); err != nil {
		clog.Error(ctx, "Error laying tile", "error", err.Error(), "name", name, "code", code)
		tileErr := fmt.Errorf("%s %v", lt, err)
		writeErr(w, tileErr, http.StatusBadRequest)
		return
	}

	if err := s.Store.WriteGame(ctx, g); err != nil {
		clog.Error(ctx, "Error writing game", "error", err.Error(), "code", code)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleLaySpacer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		clog.Error(ctx, "Error getting name", "error", err.Error())
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.Store.ReadGame(ctx, code)
	if err != nil {
		clog.Error(ctx, "Error reading game", "error", err.Error(), "code", code)
		if err == ErrNoSuchGame {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	player := g.Players[g.Turn]
	if player.Name != name {
		clog.Info(ctx, "Player not in turn for game", "name", name, "code", code)
		writeErr(w, ErrNotYourTurn, http.StatusBadRequest)
		return
	}

	if len(g.Rounds) == 0 {
		clog.Info(ctx, "Player tried to play but game not started", "name", name, "code", code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}
	round := g.Rounds[len(g.Rounds)-1]
	if round.Done {
		clog.Info(ctx, "Player tried to play but round is done", "name", name, "code", code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}

	sp := &Spacer{}
	if err := json.NewDecoder(r.Body).Decode(sp); err != nil {
		clog.Error(ctx, "Error decoding spacer", "error", err.Error(), "name", name, "code", code)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := g.LaySpacer(ctx, name, sp); err != nil {
		clog.Error(ctx, "Error laying spacer", "error", err.Error(), "name", name, "code", code)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.Store.WriteGame(ctx, g); err != nil {
		clog.Error(ctx, "Error writing game", "error", err.Error(), "code", code)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleGetGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Four minutes before timing out. Client must re-request within one minute
	// or risk being booted.
	ctx, cancel := context.WithTimeout(ctx, 200*time.Second)
	defer cancel()

	code := chi.URLParam(r, "code")
	versionStr := r.URL.Query().Get("version")
	var version int64
	if versionStr != "" {
		var err error
		version, err = strconv.ParseInt(versionStr, 10, 64)
		if err != nil {
			clog.Error(ctx, "Error parsing version", "error", err.Error(), "version", versionStr)
			writeErr(w, err, http.StatusBadRequest)
			return
		}
	}
	name, err := s.getName(ctx, r)
	if err != nil {
		clog.Error(ctx, "Error getting name", "error", err.Error())
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.Store.ReadGame(ctx, code)
	if err != nil {
		clog.Error(ctx, "Error reading game", "error", err.Error(), "code", code)
		if err == ErrNoSuchGame {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	for _, p := range g.Players {
		if p.Name != name {
			continue
		}
		if err := s.Store.RecordPlayerActive(ctx, code, name, time.Now().Unix()); err != nil {
			clog.Error(ctx, "Error setting player active", "error", err.Error(), "name", name, "code", code)
		}
	}

	// if the game hasn't begun, cull inactive players or cull the game if it's
	// been waiting for too long.
	if len(g.Rounds) == 0 {
		now := time.Now().Unix()
		if now-g.Created > 1800 {
			clog.Info(ctx, "Culling game (waiting too long)", "code", code)
			if err := s.Store.DeleteGame(ctx, code); err != nil {
				clog.Error(ctx, "Error deleting game", "error", err.Error(), "code", code)
			}
			writeErr(w, errors.New("this game took too long to start"), http.StatusNotFound)
			return
		}

		anyBooted := false
		for _, p := range g.Players {
			if p.Name == name {
				continue
			}
			lastActive, err := s.Store.PlayerLastActive(ctx, code, p.Name)
			if err != nil {
				clog.Error(ctx, "Error getting last active", "error", err.Error(), "player", p.Name, "code", code)
				continue
			}
			idleSeconds := now - lastActive
			if idleSeconds > 300 {
				clog.Info(ctx, "last active", "player", p.Name, "code", code, "lastActive", fmt.Sprint(lastActive), "idleSeconds", fmt.Sprint(idleSeconds))
				if !g.LeaveOrQuit(ctx, p.Name) {
					clog.Info(ctx, "Could not boot inactive player", "player", p.Name, "code", code)
				} else {
					anyBooted = true
				}
			}
		}
		if anyBooted {
			clog.Info(ctx, "Booted players", "code", code)
			if err := s.Store.WriteGame(ctx, g); err != nil {
				clog.Error(ctx, "Could not store game after booting players", "error", err.Error())
			}
		}
	}

	// We aleady have something newer.
	if g.Version > version {
		s.encodeFilteredGame(ctx, w, name, g)
		return
	}

	// Otherwise, wait for am update.
	select {
	case <-ctx.Done():
		err := ctx.Err()
		if err == context.DeadlineExceeded {
			writeErr(w, err, http.StatusRequestTimeout)
			return
		}
		clog.Error(ctx, "broke connection", "error", err.Error(), "name", name, "code", code)
		if err := s.Store.RecordPlayerActive(ctx, code, name, 0); err != nil {
			clog.Error(ctx, "Error setting player inactive", "error", err.Error(), "name", name, "code", code)
		}
		if !g.LeaveOrQuit(ctx, name) {
			clog.Info(ctx, "Could not boot when connection broke", "name", name, "code", code)
		}
		return
	case g := <-s.Store.WatchGame(ctx, code, version):
		s.encodeFilteredGame(ctx, w, name, g)
	}
}

func (s *GameServer) HandlePutGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		clog.Error(ctx, "Error getting name", "error", err.Error())
		writeErr(w, err, http.StatusForbidden)
		return
	}

	options := &GameOptions{}
	if err := json.NewDecoder(r.Body).Decode(options); err != nil {
		clog.Error(ctx, "Error decoding game options", "error", err.Error(), "name", name, "code", code)
		writeErr(w, err, http.StatusBadRequest)
		return
	}
	clog.Info(ctx, "Game options", "name", name, "code", code)

	var g *Game

	pickup := code == "PICKUP"
	if pickup {
		g, err = s.Store.FindPickupGame(ctx)
		if err != nil && err != ErrNoSuchGame {
			clog.Error(ctx, "Error finding pickup game", "error", err.Error())
			writeErr(w, err, http.StatusInternalServerError)
			return
		}
		if g != nil {
			clog.Info(ctx, "Found a pickup game", "code", g.Code)
		}
	} else {
		prefix := code
		if len(code) > 6 {
			prefix = code[:6]
		}

		g, err = s.Store.FindGameAlreadyPlaying(ctx, prefix, name)
		if err != nil && err != ErrNoSuchGame {
			clog.Error(ctx, "Error reading game", "error", err.Error(), "prefix", prefix)
			writeErr(w, err, http.StatusInternalServerError)
			return
		}

		if g == nil {
			g, err = s.Store.FindOpenGame(ctx, prefix)
			if err != nil && err != ErrNoSuchGame {
				clog.Error(ctx, "Error reading game", "error", err.Error(), "prefix", prefix)
				writeErr(w, err, http.StatusInternalServerError)
				return
			}
			clog.Info(ctx, "Found open game with prefix", "prefix", prefix)
		}

		if g != nil {
			if len(code) > 6 && code != g.Code {
				clog.Info(ctx, "Player joined game but it already exists", "name", name, "code", code, "existingCode", g.Code)
				writeErr(w, ErrGameOver, http.StatusConflict)
				return
			}
		}
	}
	createdNewGame := false
	if g == nil {
		createdNewGame = true
		if pickup {
			code = fmt.Sprintf("%s-%s", RandomString(6), RandomString(6))
		} else {
			if len(code) != 6 {
				clog.Info(ctx, "Code is the wrong length", "code", code)
				writeErr(w, ErrBadCode, http.StatusBadRequest)
				return
			}
			for _, c := range code {
				if !unicode.IsLetter(c) || !unicode.IsUpper(c) {
					clog.Info(ctx, "Code is not a capital letter", "code", code)
					writeErr(w, ErrBadCode, http.StatusBadRequest)
					return
				}
			}
			code = fmt.Sprintf("%s-%s", code, RandomString(6))
		}
		g = NewGame(ctx, code)
		g.Pickup = pickup
	}

	clog.Info(ctx, "Attempting to add player", "name", name, "code", code)

	inGame := false
	for _, p := range g.Players {
		if p.Name == name {
			inGame = true
			clog.Info(ctx, "Player already in game", "name", name, "code", code)
		}
	}

	player := &Player{Name: name}

	if isBotName(name) {
		if s.CheckBotTokens {
			if err := s.validateBotServiceAccountToken(ctx, r); err != nil {
				writeErr(w, err, http.StatusUnauthorized)
				return
			}
		}
		player.Bot = true
	}

	if !inGame {
		if err := g.AddPlayer(ctx, player); err != nil {
			clog.Error(ctx, "Error adding player to game", "error", err.Error(), "name", name, "code", code)
			if err == ErrGameTooManyPlayers {
				writeErr(w, err, http.StatusUnprocessableEntity)
				return
			}
			if err == ErrGameAlreadyStarted {
				writeErr(w, err, http.StatusUnprocessableEntity)
				return
			}
			if err == ErrPlayerAlreadyInGame {
				writeErr(w, err, http.StatusConflict)
				return
			}
			writeErr(w, err, http.StatusInternalServerError)
			return
		}

		if err := s.Store.WriteGame(ctx, g); err != nil {
			clog.Error(ctx, "Error writing game", "error", err.Error(), "code", code)
			writeErr(w, err, http.StatusInternalServerError)
			return
		}
	}

	s.encodeFilteredGame(ctx, w, name, g)

	if createdNewGame && options.AgentRoundOut > 1 && s.AgentSpawner != nil {
		if err := s.AgentSpawner.NewAgent(ctx, "gibbs", code, options.AgentRoundOut); err != nil {
			clog.Error(ctx, "Error spawning agent", "error", err.Error(), "code", code)
		}
	}
}

func (s *GameServer) HandleStartRound(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		clog.Error(ctx, "Error getting name", "error", err.Error())
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.Store.ReadGame(ctx, code)
	if err != nil {
		clog.Error(ctx, "Error reading game", "error", err.Error(), "code", code)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	if err := g.Start(ctx, name); err != nil {
		clog.Error(ctx, "Error starting round for game", "error", err.Error(), "code", code)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.Store.WriteGame(ctx, g); err != nil {
		clog.Error(ctx, "Error writing game", "error", err.Error(), "code", code)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleChickenFoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		clog.Error(ctx, "Error getting name", "error", err.Error())
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.Store.ReadGame(ctx, code)
	if err != nil {
		clog.Error(ctx, "Error reading game", "error", err.Error(), "code", code)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	reqBody := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		clog.Error(ctx, "Error decoding chickenfoot", "error", err.Error(), "name", name, "code", code)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	url, ok := reqBody["url"]
	if !ok {
		clog.Info(ctx, "No url provided", "name", name, "code", code)
		writeErr(w, ErrNoURL, http.StatusBadRequest)
		return
	}

	player := g.GetPlayer(ctx, name)
	if player == nil {
		clog.Info(ctx, "Player not found in game", "name", name, "code", code)
		writeErr(w, ErrPlayerNotFound, http.StatusNotFound)
		return
	}

	player.ChickenFootURL = url

	if err := s.Store.WriteGame(ctx, g); err != nil {
		clog.Error(ctx, "Error writing game", "error", err.Error(), "code", code)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleReact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		clog.Error(ctx, "Error getting name", "error", err.Error())
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.Store.ReadGame(ctx, code)
	if err != nil {
		clog.Error(ctx, "Error reading game", "error", err.Error(), "code", code)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	reqBody := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		clog.Error(ctx, "Error decoding chickenfoot", "error", err.Error(), "name", name, "code", code)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	url, ok := reqBody["url"]
	if !ok {
		clog.Info(ctx, "No url provided", "name", name, "code", code)
		writeErr(w, ErrNoURL, http.StatusBadRequest)
		return
	}

	player := g.GetPlayer(ctx, name)
	if player == nil {
		clog.Info(ctx, "Player not found in game", "name", name, "code", code)
		writeErr(w, ErrPlayerNotFound, http.StatusNotFound)
		return
	}

	player.ReactURL = url

	if err := s.Store.WriteGame(ctx, g); err != nil {
		clog.Error(ctx, "Error writing game", "error", err.Error(), "code", code)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleRegisterPlayerName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()
	playerID := r.Header.Get("X-Player-ID")

	pi := &PlayerInfo{}
	if err := json.NewDecoder(r.Body).Decode(pi); err != nil {
		clog.Error(ctx, "Error decoding player info", "error", err.Error())
		writeErr(w, err, http.StatusBadRequest)
		return
	}
	pi.Id = playerID

	if err := validatePlayerName(pi.Name); err != nil {
		clog.Error(ctx, "Error validating player name", "error", err.Error(), "name", pi.Name)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if isBotName(pi.Name) {
		clog.Info(ctx, "Error trying to register name with reserved initials", "name", pi.Name)
		writeErr(w, ErrPlayerInitialsReserved, http.StatusForbidden)
		return
	}

	if rpi, err := s.Store.GetPlayerByName(r.Context(), pi.Name); err == nil {
		if rpi.Id != playerID {
			clog.Info(ctx, "Player already registered", "name", pi.Name, "id", rpi.Id)
			writeErr(w, ErrPlayerAlreadyRegistered, http.StatusConflict)
			return
		}
	}

	if playerID != "" {
		if err := s.Store.RegisterPlayerName(r.Context(), playerID, pi.Name); err != nil {
			clog.Error(ctx, "Error registering player", "error", err.Error(), "name", pi.Name)
			writeErr(w, err, http.StatusBadRequest)
			return
		}
		clog.Info(ctx, "Registered player", "playerID", playerID)
	} else {
		clog.Info(ctx, "Anonymous player", "name", pi.Name)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pi)
}

func (s *GameServer) HandleGetPlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	playerID := chi.URLParam(r, "playerID")
	pi, err := s.Store.GetPlayer(ctx, playerID)

	if err == ErrNoRegisteredPlayer {
		// If this isn't a valid player ID, try the player by name.
		pi, err = s.Store.GetPlayerByName(ctx, playerID)
	}

	if err != nil {
		if err == ErrNoRegisteredPlayer {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		clog.Error(ctx, "Error getting player name", "error", err.Error(), "playerID", playerID)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pi)
}

func (s *GameServer) HandleUpdatePlayerConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()
	playerID := chi.URLParam(r, "playerID")
	if err := s.validatePlayerID(ctx, playerID, r); err != nil {
		clog.Error(ctx, "Error validating player ID", "error", err.Error())
		writeErr(w, err, http.StatusForbidden)
		return
	}

	cfg := PlayerConfig{}
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		clog.Error(ctx, "Error decoding player info", "error", err.Error())
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.Store.UpdatePlayerConfig(ctx, playerID, cfg); err != nil {
		clog.Error(ctx, "Error updating player config", "error", err.Error())
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cfg)
}

type ReportIssueRequest struct {
	Summary          string `json:"summary"`
	WhatHappened     string `json:"what_happened"`
	WhatShouldHappen string `json:"what_should_happen"`
	ErrorMessage     string `json:"error_message"`
}

func (s *GameServer) HandleReportIssue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		clog.Error(ctx, "Error getting name", "error", err.Error())
		writeErr(w, err, http.StatusForbidden)
		return
	}

	reqBody := ReportIssueRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		clog.Error(ctx, "Error decoding report issue request", "error", err.Error())
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	g, err := s.Store.ReadGame(ctx, code)
	if err != nil {
		clog.Error(ctx, "Error reading game", "error", err.Error(), "code", code)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	s.Store.ReportIssue(ctx, name, g, reqBody.Summary, reqBody.WhatHappened, reqBody.WhatShouldHappen, reqBody.ErrorMessage)
}
