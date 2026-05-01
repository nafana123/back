package game

import (
	"back/internal/model"
	gamerepo "back/internal/repository/game"
)

type GameService struct {
	gameRepo *gamerepo.GameRepository
}

func NewGameService(gameRepo *gamerepo.GameRepository) *GameService {
	return &GameService{
		gameRepo: gameRepo,
	}
}

func (service *GameService) GetAllGames() ([]model.Game, error) {
	return service.gameRepo.GetAll()
}
