package postgresql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"log"
	db_dto "main/internal/adapters/dto"
	"main/internal/apperror"
	"main/internal/domain/entity"
	"main/pkg/client/postgresql"
	"main/pkg/logging"
)

const (
	User = "\"user\""
)

type userStorage struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewUserStorage(db postgresql.Client) *userStorage {
	return &userStorage{client: db}
}

func (u *userStorage) Create(ctx context.Context, user db_dto.CreateUserDTO) (string, error) {
	sql, args, err := sq.Insert(User).Columns("email", "username", "passwordHash", "session").Values(user.Email, user.Username, user.PasswordHash, user.Session).Suffix("RETURNING \"id\"").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return "", apperror.NewAppError(err, "incorrect data entered", err.Error(), "US-000004")
	}
	if err = u.client.QueryRow(ctx, sql, args...).Scan(&user.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			u.logger.Error(newErr)
			return "", apperror.NewAppError(err, "SQL error", pgErr.Message, "US-000005")
		}
		return "", apperror.NewAppError(err, "SQL error", err.Error(), "US-000005")
	}
	userID := user.ID
	return userID, nil
}

func (u *userStorage) IsUserCreated(ctx context.Context, user db_dto.IsUserExistsDTO) (bool, error) {
	sql, args, err := sq.Select("username", "email").From(User).Prefix("SELECT NOT EXISTS(").Where("(username = $1 OR email = $2)", user.Username, user.Email).Suffix(")").ToSql()
	if err != nil {
		return false, apperror.NewAppError(err, "incorrect data entered", err.Error(), "US-000004")
	}
	var isCreated bool
	if err = u.client.QueryRow(ctx, sql, args...).Scan(&isCreated); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			u.logger.Error(newErr)
			return false, apperror.NewAppError(err, "SQL error", pgErr.Message, "US-000005")
		}
		return false, apperror.NewAppError(err, "SQL error", err.Error(), "US-000005")
	}
	return isCreated, nil
}

func (u *userStorage) AuthUser(ctx context.Context, user db_dto.AuthUserDTO) (string, error) {
	sql, args, err := sq.Select("id").From(User).Where("(email = $1 OR username = $2) AND passwordHash = $3", user.Email, user.Username, user.PasswordHash).ToSql()
	if err != nil {
		return "", apperror.NewAppError(err, "incorrect data entered", err.Error(), "US-000004")
	}
	var id string
	if err = u.client.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		if err == pgx.ErrNoRows {
			return "", apperror.NewAppError(err, "invalid login or password entered", "", "US-000009")
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			u.logger.Error(newErr)
			return "", apperror.NewAppError(err, "SQL error", pgErr.Message, "US-000005")
		}
		return "", apperror.NewAppError(err, "SQL error", err.Error(), "US-000005")
	}
	return id, nil
}

func (u *userStorage) GetUser(ctx context.Context, user db_dto.GetUserDTO) (entity.User, error) {
	var userReturn entity.User
	sql, args, err := sq.Select("id", "username", "email", "session").From(User).Where("id = $1", user.ID).ToSql()
	if err != nil {
		return userReturn, apperror.NewAppError(err, "incorrect data entered", err.Error(), "US-000004")
	}
	if err = u.client.QueryRow(ctx, sql, args...).Scan(&userReturn.ID, &userReturn.Username, &userReturn.Email, &userReturn.Session); err != nil {
		if err == pgx.ErrNoRows {
			return userReturn, apperror.NewAppError(err, "invalid token", "", "TJ-000004")
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			u.logger.Error(newErr)
			return userReturn, apperror.NewAppError(err, "SQL error", pgErr.Message, "US-000005")
		}
		return userReturn, apperror.NewAppError(err, "SQL error", err.Error(), "US-000005")
	}
	return userReturn, nil
}

func (u *userStorage) GetUserByRefreshToken(ctx context.Context, user db_dto.GetUserByRefreshTokenDTO) (string, error) {
	sql, args, err := sq.Select("id").From(User).Where("session->>'RefreshToken' = $1", user.Token).ToSql()
	if err != nil {
		return "", apperror.NewAppError(err, "incorrect data entered", err.Error(), "US-000004")
	}
	var id string
	if err = u.client.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		log.Println(sql, args)
		if err == pgx.ErrNoRows {
			return "", apperror.NewAppError(err, "invalid token", "", "TJ-000004")
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			u.logger.Error(newErr)
			return "", apperror.NewAppError(err, "SQL error", pgErr.Message, "US-000005")
		}
		return "", apperror.NewAppError(err, "SQL error", err.Error(), "US-000005")
	}
	return id, nil
}

func (u *userStorage) Update(ctx context.Context, user entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *userStorage) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
