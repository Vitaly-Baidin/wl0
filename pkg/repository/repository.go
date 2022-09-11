package repository

import "github.com/jackc/pgx/v4/pgxpool"

type Repository struct {
	database *pgxpool.Pool
}
