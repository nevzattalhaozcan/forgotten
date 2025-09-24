package services

import (
	"errors"
	"time"

	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"gorm.io/gorm"
)

type ReadingService struct {
	cfg          *config.Config
	userRepo     repository.UserRepository
	bookRepo     repository.BookRepository
	clubRepo     repository.ClubRepository
	readRepo     repository.ReadingRepository
	clubReadRepo repository.ClubReadingRepository
}

func NewReadingService(cfg *config.Config, userRepo repository.UserRepository, bookRepo repository.BookRepository, clubRepo repository.ClubRepository, readRepo repository.ReadingRepository, clubReadRepo repository.ClubReadingRepository) *ReadingService {
	return &ReadingService{
		cfg:          cfg,
		userRepo:     userRepo,
		bookRepo:     bookRepo,
		clubRepo:     clubRepo,
		readRepo:     readRepo,
		clubReadRepo: clubReadRepo,
	}
}

func (s *ReadingService) StartReading(userID, bookID uint) (*models.UserBookProgressResponse, error) {
	if _, err := s.userRepo.GetByID(userID); err != nil {
		return nil, err
	}
	book, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	p := &models.UserBookProgress{
		UserID:    userID,
		BookID:    bookID,
		Status:    models.ReadingActive,
		StartedAt: now,
	}
	if err := s.readRepo.UpsertUserProgress(p); err != nil {
		return nil, err
	}
	resp := p.ToResponse(book)
	return &resp, nil
}

func (s *ReadingService) UpdateProgress(userID, bookID uint, req *models.UpdateReadingProgressRequest) (*models.UserBookProgressResponse, error) {
	book, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		return nil, err
	}

	p, err := s.readRepo.GetUserBookProgress(userID, bookID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		p = &models.UserBookProgress{
			UserID: userID, BookID: bookID, Status: models.ReadingActive, StartedAt: time.Now(),
		}
	} else if err != nil {
		return nil, err
	}

	if req.CurrentPage != nil {
		p.CurrentPage = req.CurrentPage
	}
	if req.Percent != nil {
		p.Percent = req.Percent
	}
	if p.Status == models.ReadingNotStarted {
		p.Status = models.ReadingActive
	}

	log := &models.ReadingLog{
		UserID:     userID,
		BookID:     bookID,
		PagesDelta: req.PagesDelta,
		Minutes:    req.Minutes,
		Note:       req.Note,
	}
	if req.CurrentPage != nil {
		log.ToPage = req.CurrentPage
	}

	if err := s.readRepo.UpsertUserProgress(p); err != nil {
		return nil, err
	}
	if err := s.readRepo.AppendLog(log); err != nil {
		return nil, err
	}

	resp := p.ToResponse(book)
	return &resp, nil
}

func (s *ReadingService) CompleteReading(userID, bookID uint, note *string) (*models.UserBookProgressResponse, error) {
	book, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		return nil, err
	}

	p, err := s.readRepo.GetUserBookProgress(userID, bookID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("reading not started")
	} else if err != nil {
		return nil, err
	}
	now := time.Now()
	p.Status = models.ReadingFinished
	p.FinishedAt = &now

	if book.Pages != nil {
		lastPage := *book.Pages
		p.CurrentPage = &lastPage
		h := float32(100.0)
		p.Percent = &h
	}

	if err := s.readRepo.UpsertUserProgress(p); err != nil {
		return nil, err
	}

	_ = s.readRepo.AppendLog(&models.ReadingLog{
		UserID: userID, BookID: bookID, Note: note,
	})

	resp := p.ToResponse(book)
	return &resp, nil
}

func (s *ReadingService) ListUserProgress(userID uint) ([]*models.UserBookProgressResponse, error) {
	entries, err := s.readRepo.ListUserProgress(userID)
	if err != nil {
		return nil, err
	}
	var out []*models.UserBookProgressResponse
	for _, e := range entries {
		book, berr := s.bookRepo.GetByID(e.BookID)
		if berr != nil {
			continue
		}
		r := e.ToResponse(book)
		out = append(out, &r)
	}
	return out, nil
}

func (s *ReadingService) UserReadingHistory(userID uint) ([]models.UserReadingHistoryItem, error) {
	finished, err := s.readRepo.ListUserFinished(userID)
	if err != nil {
		return nil, err
	}
	var out []models.UserReadingHistoryItem
	for _, p := range finished {
		book, berr := s.bookRepo.GetByID(p.BookID)
		if berr != nil {
			continue
		}
		logs, _ := s.readRepo.ListLogsByUserAndBook(userID, p.BookID)
		item := models.UserReadingHistoryItem{
			Book: *book,
			FinishedAt: func() time.Time {
				if p.FinishedAt != nil {
					return *p.FinishedAt
				}
				return p.UpdatedAt
			}(),
			Logs: logs,
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *ReadingService) AssignBookToClub(clubID, bookID uint, req *models.AssignBookRequest) (*models.ClubAssignmentResponse, error) {
	if _, err := s.clubRepo.GetByID(clubID); err != nil {
		return nil, err
	}
	book, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		return nil, err
	}

	if active, err := s.clubReadRepo.GetActiveAssignment(clubID); err == nil && active != nil {
		return nil, errors.New("club already has an active assignment")
	}

	a := &models.ClubBookAssignment{
		ClubID:     clubID,
		BookID:     bookID,
		Status:     models.ClubAssignmentActive,
		StartDate:  req.StartDate,
		DueDate:    req.DueDate,
		TargetPage: req.TargetPage,
		Checkpoint: req.Checkpoint,
	}

	if err := s.clubReadRepo.CreateAssignment(a); err != nil {
		return nil, err
	}

	resp := &models.ClubAssignmentResponse{
		ID: a.ID, ClubID: clubID, Book: *book, Status: string(a.Status),
		StartDate: a.StartDate, DueDate: a.DueDate, TargetPage: a.TargetPage, Checkpoint: a.Checkpoint,
	}
	return resp, nil
}

func (s *ReadingService) UpdateClubCheckpoint(clubID uint, req *models.UpdateClubCheckpointRequest) (*models.ClubAssignmentResponse, error) {
	a, err := s.clubReadRepo.GetActiveAssignment(clubID)
	if err != nil {
		return nil, err
	}
	if req.TargetPage != nil {
		a.TargetPage = req.TargetPage
	}
	if req.Checkpoint != nil {
		a.Checkpoint = req.Checkpoint
	}
	if err := s.clubReadRepo.UpdateAssignment(a); err != nil {
		return nil, err
	}
	book, _ := s.bookRepo.GetByID(a.BookID)
	resp := &models.ClubAssignmentResponse{
		ID: a.ID, ClubID: a.ClubID, Book: *book, Status: string(a.Status),
		StartDate: a.StartDate, DueDate: a.DueDate, TargetPage: a.TargetPage, Checkpoint: a.Checkpoint,
	}
	return resp, nil
}

func (s *ReadingService) CompleteClubAssignment(clubID uint) (*models.ClubAssignmentResponse, error) {
	a, err := s.clubReadRepo.GetActiveAssignment(clubID)
	if err != nil {
		return nil, err
	}
	if err := s.clubReadRepo.CompleteAssignment(a.ID); err != nil {
		return nil, err
	}
	a.Status = models.ClubAssignmentCompleted
	now := time.Now()
	a.CompletedAt = &now
	book, _ := s.bookRepo.GetByID(a.BookID)
	resp := &models.ClubAssignmentResponse{
		ID: a.ID, ClubID: a.ClubID, Book: *book, Status: string(a.Status),
		StartDate: a.StartDate, DueDate: a.DueDate, TargetPage: a.TargetPage, Checkpoint: a.Checkpoint,
	}
	return resp, nil
}

func (s *ReadingService) ListClubAssignments(clubID uint) ([]models.ClubAssignmentResponse, error) {
	as, err := s.clubReadRepo.ListAssignments(clubID)
	if err != nil {
		return nil, err
	}
	var out []models.ClubAssignmentResponse
	for _, a := range as {
		book, berr := s.bookRepo.GetByID(a.BookID)
		if berr != nil {
			continue
		}
		out = append(out, models.ClubAssignmentResponse{
			ID: a.ID, ClubID: a.ClubID, Book: *book, Status: string(a.Status),
			StartDate: a.StartDate, DueDate: a.DueDate, TargetPage: a.TargetPage, Checkpoint: a.Checkpoint,
		})
	}
	return out, nil
}
