package postgresql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	db_dto "main/internal/adapters/dto"
	"main/internal/apperror"
	"main/internal/domain/entity"
	"main/pkg/client/postgresql"
	"main/pkg/logging"
)

const (
	Product = "\"product\""
)

type productStorage struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewProductStorage(db postgresql.Client) *productStorage {
	return &productStorage{client: db}
}

func (p *productStorage) GetProducts(ctx context.Context) ([]entity.Product, error) {
	var productsReturn []entity.Product
	sql, args, err := sq.Select("*").From(Product).ToSql()
	if err != nil {
		return nil, apperror.NewAppError(err, "incorrect data entered", err.Error(), "PR-000004")
	}
	rows, err := p.client.Query(ctx, sql, args...)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			p.logger.Errorln(newErr)
			return nil, apperror.NewAppError(err, "SQL error", pgErr.Message, "PR-000005")
		}
		return nil, apperror.NewAppError(err, "SQL error", err.Error(), "PR-000005")
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Discount, &product.Seller, &product.Rating)
		if err != nil {
			return nil, apperror.NewAppError(err, "can`t scan products", err.Error(), "PR-000002")
		}
		productsReturn = append(productsReturn, product)
	}
	return productsReturn, nil
}

func (u *productStorage) Add(ctx context.Context, product db_dto.CreateProductDTO) (string, error) {
	sql, args, err := sq.Insert(Product).Columns("name", "description", "price", "seller").Values(product.Name, product.Description, product.Price, product.Seller).Suffix("RETURNING \"id\"").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return "", apperror.NewAppError(err, "incorrect data entered", err.Error(), "PR-000004")
	}
	var productID string
	if err = u.client.QueryRow(ctx, sql, args...).Scan(&productID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			u.logger.Errorln(newErr)
			return "", apperror.NewAppError(err, "SQL error", pgErr.Message, "PR-000005")
		}
		return "", apperror.NewAppError(err, "SQL error", err.Error(), "PR-000005")
	}
	return productID, nil
}
