package tournament

import (
	tournamentdto "back/internal/dto/tournament"
	"back/internal/model"
	participantrepo "back/internal/repository/participant"
	tournamentrepo "back/internal/repository/tournament"
	"fmt"
)

type TournamentService struct {
	tournamentRepo  *tournamentrepo.TournamentRepository
	participantRepo *participantrepo.ParticipantRepository
}

func NewTournamentService(tournamentRepo *tournamentrepo.TournamentRepository, participantRepo *participantrepo.ParticipantRepository) *TournamentService {
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

func (s *TournamentService) CreateTournament(data *tournamentdto.CreateTournamentRequest) (*model.Tournament, error) {
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

func (s *TournamentService) JoinTournament(tournamentID string, userID int) (*model.Participant, error) {
	tournament, err := s.tournamentRepo.GetById(tournamentID)
	if err != nil {
		return nil, err
	}

	if tournament.Status != model.StatusRegistration {
		return nil, fmt.Errorf("tournament is not open for registration")
	}

	if tournament.CurrentPlayers >= tournament.MaxPlayers {
		return nil, fmt.Errorf("tournament is full")
	}

	exists, err := s.participantRepo.Exists(tournamentID, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("already joined")
	}

	participant := &model.Participant{
		TournamentID: tournamentID,
		UserID:       userID,
	}

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
