package service

import (
	"back/internal/dto"
	"back/internal/model"
	"back/internal/repository"
	"fmt"
)

type TournamentService struct {
	tournamentRepo  *repository.TournamentRepository
	participantRepo *repository.ParticipantRepository
}

func NewTournamentService(tournamentRepo *repository.TournamentRepository, participantRepo *repository.ParticipantRepository) *TournamentService {
	if tournamentRepo == nil {
		panic("tournamentRepo cannot be nil")
	}
	if participantRepo == nil {
		panic("participantRepo cannot be nil")
	}

	return &TournamentService{
		tournamentRepo:  tournamentRepo,
		participantRepo: participantRepo,
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

func (s *TournamentService) ChangeStatus(id, status string) error {
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

func (s *TournamentService) JoinTournament(data *dto.JoinTournamentRequest) (*model.Participant, error) {
	tournament, err := s.tournamentRepo.GetById(data.TournamentID)
	if err != nil {
		return nil, err
	}

	if tournament.Status != model.StatusRegistration {
		return nil, fmt.Errorf("tournament is not open for registration")
	}

	if tournament.CurrentPlayers >= tournament.MaxPlayers {
		return nil, fmt.Errorf("tournament is full")
	}

	exists, err := s.participantRepo.Exists(data.TournamentID, data.UserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("already joined")
	}

	participant := data.ToModel()
	if err := s.participantRepo.Create(participant); err != nil {
		return nil, err
	}

	// todo тут нужна логика оплаты либо с кошелька юзера либо напряму

	tournament.CurrentPlayers++

	if err := s.tournamentRepo.Update(tournament); err != nil {
		return nil, err
	}

	return participant, nil
}

func (s *TournamentService) GetParticipantsByTournament(tournamentId string) ([]model.Participant, error) {
	return s.participantRepo.GetByTournament(tournamentId)
}
