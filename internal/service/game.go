package service

import "back/internal/repository"

type GameService struct {
	tournamentRepo *repository.TournamentRepository
	gameRepo       *repository.GameRepository
}

func NewGameService(tournamentRepo *repository.TournamentRepository, gameRepo *repository.GameRepository) *GameService {
	if tournamentRepo == nil {
		panic("tournamentRepo cannot be nil")
	}
	if gameRepo == nil {
		panic("gameRepo cannot be nil")
	}

	return &GameService{
		tournamentRepo: tournamentRepo,
		gameRepo:       gameRepo,
	}
}
