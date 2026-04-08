package service

import (
	"back/internal/dto"
	"back/internal/model"
	"back/internal/repository"
)

type TournamentService struct {
	tournamentRepo *repository.TournamentRepository
	gameRepo       *repository.GameRepository
}

func NewTournamentService(tournamentRepo *repository.TournamentRepository, gameRepo *repository.GameRepository) *TournamentService {
	if tournamentRepo == nil {
		panic("tournamentRepo cannot be nil")
	}
	if gameRepo == nil {
		panic("gameRepo cannot be nil")
	}

	return &TournamentService{
		tournamentRepo: tournamentRepo,
		gameRepo:       gameRepo,
	}
}

func (s *TournamentService) GetAllTournaments() ([]model.Tournament, error) {
	return s.tournamentRepo.GetALl()
}

func (s *TournamentService) GetTournament(id string) (*model.Tournament, error) {
	return s.tournamentRepo.GetById(id)
}

func (s *TournamentService) CreateTournament(data *dto.CreateTournamentRequest) (*model.Tournament, error) {
	tournament := data.ToModel()

	if err := s.tournamentRepo.CreateTournament(tournament); err != nil {
		return nil, err
	}

	return tournament, nil
}

func (s *TournamentService) ChangeStatus(id string, status string) error {
	tournament, err := s.tournamentRepo.GetById(id)

	if err != nil {
		return err
	}

	// TODO в будущем добавить реализации в случае разных статусов
	//switch status {
	//case model.StatusInProcess:
	//case model.StatusFailed:
	//case model.StatusCompleted:
	// например в случае если турнир из статуса в процессе перешел в статус завершен, то нужно определить победителей и выплатить им призы
	//}

	tournament.Status = status

	if err := s.tournamentRepo.Update(tournament); err != nil {
		return err
	}

	return nil
}
