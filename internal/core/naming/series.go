package naming

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SeriesService struct {
	DB *pgxpool.Pool
}

func NewSeriesService(db *pgxpool.Pool) *SeriesService {
	return &SeriesService{DB: db}
}

func (s *SeriesService) NextSeries(seriesKey string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	seriesKey = strings.TrimSpace(seriesKey)
	if seriesKey == "" {
		return "", fmt.Errorf("series key is required")
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	var (
		prefix       string
		currentValue int64
		digits       int
	)

	err = tx.QueryRow(ctx, `
		SELECT prefix, current_value, digits
		FROM "tabNaming Series"
		WHERE series_key = $1
		FOR UPDATE
	`, seriesKey).Scan(&prefix, &currentValue, &digits)
	if err != nil {
		return "", fmt.Errorf("naming series not found: %s", seriesKey)
	}

	nextValue := currentValue + 1
	nextName := fmt.Sprintf("%s%0*d", prefix, digits, nextValue)

	_, err = tx.Exec(ctx, `
		UPDATE "tabNaming Series"
		SET
			current_value = $1,
			modified = NOW(),
			modified_by = 'Administrator'
		WHERE series_key = $2
	`, nextValue, seriesKey)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return nextName, nil
}

func (s *SeriesService) GenerateName(doctype string, rule string, doc Document) (string, error) {
	return GenerateName(doctype, rule, doc, func(prefix string, digits int) (string, error) {
		seriesKey := prefix + strings.Repeat("#", digits)
		return s.NextSeries(seriesKey)
	})
}
