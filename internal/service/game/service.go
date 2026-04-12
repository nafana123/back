package game

import (
	"back/internal/model"
	"back/internal/repository"
)

type GameService struct {
	gameRepo *repository.GameRepository
}

func NewGameService(gameRepo *repository.GameRepository) *GameService {
	if gameRepo == nil {
		panic("gameRepo cannot be nil")
	}

	return &GameService{
		gameRepo: gameRepo,
	}
}

func (service *GameService) GetAllGames() ([]model.Game, error) {
	return service.gameRepo.GetAll()
}
