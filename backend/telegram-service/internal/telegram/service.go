package telegram

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/alina965/pLaNtS/telegram-service/internal/client"
	"github.com/alina965/pLaNtS/telegram-service/internal/domain"
	"github.com/alina965/pLaNtS/telegram-service/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrNotLinked       = errors.New("telegram is not linked")
	ErrInvalidStart    = errors.New("invalid start payload")
	ErrLinkCodeMissing = errors.New("link code not found")
)

type Service struct {
	repo        *repository.TelegramLinksRepository
	bot         *client.BotClient
	botUsername string
}

func NewService(repo *repository.TelegramLinksRepository, bot *client.BotClient, botUsername string) *Service {
	return &Service{
		repo:        repo,
		bot:         bot,
		botUsername: strings.TrimPrefix(botUsername, "@"),
	}
}

func (s *Service) CreateLink(ctx context.Context, userID uuid.UUID) (*domain.CreateLinkResponse, error) {
	code, err := randomCode(16)
	if err != nil {
		return nil, err
	}

	if _, err := s.repo.UpsertLinkCode(ctx, userID, code); err != nil {
		return nil, err
	}

	return &domain.CreateLinkResponse{
		LinkCode: code,
		LinkURL:  fmt.Sprintf("https://t.me/%s?start=%s", s.botUsername, code),
	}, nil
}

func (s *Service) HandleStart(ctx context.Context, code string, chatID int64) error {
	code = strings.TrimSpace(code)
	if code == "" {
		return ErrInvalidStart
	}

	link, err := s.repo.FindByLinkCode(ctx, code)
	if err != nil {
		return err
	}
	if link == nil {
		return ErrLinkCodeMissing
	}

	if err := s.repo.AttachChat(ctx, code, chatID); err != nil {
		return err
	}

	_ = s.bot.SendMessage(ctx, chatID, "Telegram успешно привязан к аккаунту pLaNtS.")
	return nil
}

func (s *Service) Notify(ctx context.Context, userID uuid.UUID, text string) error {
	if strings.TrimSpace(text) == "" {
		return errors.New("text is required")
	}

	link, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if link == nil || link.ChatID == nil {
		return ErrNotLinked
	}

	return s.bot.SendMessage(ctx, *link.ChatID, text)
}

func (s *Service) StartPolling(ctx context.Context) {
	var offset int64
	log.Println("telegram polling started")

	for {
		select {
		case <-ctx.Done():
			log.Println("telegram polling stopped")
			return
		default:
		}

		updates, err := s.bot.GetUpdates(ctx, offset, 30)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("getUpdates: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, update := range updates {
			offset = update.UpdateID + 1
			if update.Message == nil {
				continue
			}

			code, ok := parseStartCode(update.Message.Text)
			if !ok {
				continue
			}

			if err := s.HandleStart(ctx, code, update.Message.Chat.ID); err != nil {
				log.Printf("handle start: %v", err)
			}
		}
	}
}

func parseStartCode(text string) (string, bool) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", false
	}

	parts := strings.Fields(text)
	if len(parts) == 0 {
		return "", false
	}

	cmd := parts[0]
	if i := strings.IndexByte(cmd, '@'); i >= 0 {
		cmd = cmd[:i]
	}
	if cmd != "/start" {
		return "", false
	}
	if len(parts) < 2 {
		return "", false
	}
	return parts[1], true
}

func randomCode(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
