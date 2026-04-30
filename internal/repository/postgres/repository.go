package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"subscriptions_rest/internal/config"
	"subscriptions_rest/internal/domain"
	"time"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New(cfg *config.Config) (*Repository, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.Postgres.Username,
		cfg.Postgres.Password, cfg.Postgres.URL, cfg.Postgres.Database))
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetSubscribtions(ctx context.Context) ([]domain.Subscribtion, error) {
	rows, err := r.db.QueryContext(ctx, `
	SELECT 
		user_id, service_name, price, start_date, end_date
	FROM 
		subscriptions`)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	res := make([]domain.Subscribtion, 0, 100)
	for rows.Next() {
		subscribtion := domain.Subscribtion{}
		var startTime time.Time
		var endTime *time.Time
		err := rows.Scan(&subscribtion.UserId, &subscribtion.ServiceName, &subscribtion.Price, &startTime, &endTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscribtion: %w", err)
		}
		subscribtion.Start = &domain.Date{
			Year:  startTime.Year(),
			Month: int(startTime.Month()),
		}
		if endTime != nil {
			subscribtion.End = &domain.Date{
				Year:  endTime.Year(),
				Month: int(endTime.Month()),
			}
		}

		res = append(res, subscribtion)
	}

	return res, nil
}

func (r *Repository) GetSubscribtion(ctx context.Context, userId, seviceName string) (domain.Subscribtion, error) {
	row := r.db.QueryRowContext(ctx, `
	SELECT 
		price, start_date, end_date
	FROM 
		subscriptions
	WHERE 
		user_id = $1 AND service_name = $2;`, userId, seviceName)

	res := domain.Subscribtion{
		UserId:      userId,
		ServiceName: seviceName,
	}
	var startTime time.Time
	var endTime *time.Time
	err := row.Scan(&res.Price, &startTime, &endTime)
	if err != nil {
		return domain.Subscribtion{}, fmt.Errorf("failed to get subscribtion: %w", err)
	}
	res.Start = &domain.Date{
		Year:  startTime.Year(),
		Month: int(startTime.Month()),
	}
	if endTime != nil {
		res.End = &domain.Date{
			Year:  endTime.Year(),
			Month: int(endTime.Month()),
		}
	}

	return res, nil
}

func (r *Repository) GetSubscribtionsPrice(ctx context.Context, requestStruct domain.Subscribtion) (int, error) {
	end := "(CASE WHEN end_date IS NOT NULL THEN end_date ELSE current_date END)" //"current_date"
	start := "start_date"
	if requestStruct.End.IsValid() {
		dateEnd := requestStruct.End.DateString()
		end = fmt.Sprintf("(CASE WHEN end_date IS NOT NULL THEN LEAST(end_date, '%s'::date) ELSE '%s'::date END)", dateEnd, dateEnd)
	}
	if requestStruct.Start.IsValid() {
		start = fmt.Sprintf("GREATEST(start_date, '%s'::date)", requestStruct.Start.DateString())
	}

	row := r.db.QueryRowContext(ctx, fmt.Sprintf(`
	SELECT 
		price * ((%s-%s) / 30)
	FROM 
		subscriptions
	WHERE 
		user_id = $1 AND service_name = $2;`, end, start), requestStruct.UserId, requestStruct.ServiceName)

	res := 0
	err := row.Scan(&res)
	if err != nil {
		return 0, fmt.Errorf("failed to get price: %w", err)
	}

	return res, nil
}
func (r *Repository) DeleteSubscribtion(ctx context.Context, userId, seviceName string) error {
	_, err := r.db.ExecContext(ctx, `
	DELETE 
	FROM 
		subscriptions 
	WHERE
		user_id = $1 AND service_name = $2;`, userId, seviceName)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}
func (r *Repository) UpdateSubscribtion(ctx context.Context, subscribtion domain.Subscribtion) error {
	add := ""
	if subscribtion.End != nil {
		add = fmt.Sprintf(", end_date = '%s'", subscribtion.End.DateString())
	}

	_, err := r.db.ExecContext(ctx, fmt.Sprintf(`
	UPDATE  
		subscriptions
	SET price = $1, start_date = $2%s 
	WHERE
		user_id = $3 AND service_name = $4;`, add),
		subscribtion.Price, subscribtion.Start.DateString(),
		subscribtion.UserId, subscribtion.ServiceName)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}
func (r *Repository) AddSubscribtion(ctx context.Context, subscribtion domain.Subscribtion) error {
	endDateColAdd := ""
	endDateAdd := ""
	if subscribtion.End != nil {
		endDateColAdd = ", end_date"
		endDateAdd = fmt.Sprintf(", '%s'", subscribtion.End.DateString())
	}

	_, err := r.db.ExecContext(ctx, fmt.Sprintf(`
	INSERT INTO  
		subscriptions
	(user_id, service_name, price, start_date%s)
	VALUES ($1, $2, $3, $4%s) 
	`, endDateColAdd, endDateAdd),
		subscribtion.UserId, subscribtion.ServiceName,
		subscribtion.Price, subscribtion.Start.DateString())
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}
