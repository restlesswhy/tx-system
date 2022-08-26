package store

import (
	"context"
	"txsystem/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *store {
	return &store{pool: pool}
}

func (s *store) CreateUser(u *models.User) error {
	q := `--sql
		INSERT INTO users
		("id","balance")
		VALUES
		($1,$2);
	`

	_, err := s.pool.Exec(context.Background(), q, u.ID, u.Balance)
	if err != nil {
		return errors.Wrap(err, "create user error")
	}

	return nil
}

func (s *store) CreateTx(tx *models.Transaction) error {
	q := `--sql
		INSERT INTO txs
		("user_id","amount","action","status")
		VALUES
		($1,$2,$3,$4,$5)
		RETURNING "id", "create_at";
	`

	if err := s.pool.QueryRow(context.Background(), q, tx.UserID, tx.Amount, tx.Amount, tx.Status).Scan(&tx.ID, &tx.CreateAt); err != nil {
		return errors.Wrap(err, "create tx error")
	}

	return nil
}

func (s *store) UpdateTxStatusByID(status models.Status, id int) error {
	q := `--sql
		UPDATE txs
		SET "status"=$1
		WHERE "id"=$2;
	`

	_, err := s.pool.Exec(context.Background(), q, status, id)
	if err != nil {
		return errors.Wrap(err, "update status error")
	}

	return nil
}

func (s *store) GetBalanceByUserID(id int) (int, error) {
	balance := 0
	
	q := `--sql
		SELECT balance
		FROM "users"
		WHERE "id"=$1;
	`

	if err := s.pool.QueryRow(context.Background(), q, id).Scan(&balance); err != nil {
		return 0, errors.Wrap(err, "get balance error")
	}

	return balance, nil
}

func (s *store) UpdateBalanceByID(amount int) error {
	
}