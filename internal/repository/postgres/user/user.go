package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/internal/entity"
	"main/internal/services/user"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r Repository) Create(ctx context.Context, data user.Create, adminId int64, birthTime time.Time) (int64, error) {
	var id int64
	log.Println(data.FirstName)
	query := `
		INSERT INTO users
			(avatar, first_name, last_name, phone_number, password, login, birth_date, email, role, region_id, district_id, created_by, created_at)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		RETURNING id
	`

	err := r.DB.QueryRowContext(
		ctx,
		query,
		data.Avatar,
		data.FirstName,
		data.LastName,
		data.PhoneNumber,
		data.Password,
		data.Login,
		birthTime,
		data.Email,
		data.Role,
		data.RegionID,
		data.DistrictID,
		adminId,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r Repository) GetAll(ctx context.Context, filter entity.Filter) ([]user.Get, int, error) {
	var users []user.Get
	var limitQuery, offsetQuery string

	whereQuery := "WHERE u.deleted_at IS NULL AND u.id != 1"

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := "ORDER BY u.id DESC"

	if filter.Order != nil && *filter.Order != "" {
		parts := strings.Fields(*filter.Order)
		if len(parts) == 2 {
			column := parts[0]
			direction := strings.ToUpper(parts[1])
			if direction != "ASC" && direction != "DESC" {
				direction = "ASC"
			}
			orderQuery = fmt.Sprintf("ORDER BY %s %s", column, direction)
		}
	}

	query := fmt.Sprintf(
		`
			SELECT id, avatar, first_name, last_name, phone_number, login, birth_date, email, role, region_id, district_id, created_at
			FROM users u
			%s
			%s
			%s
			%s
		`,
		whereQuery,
		orderQuery,
		limitQuery,
		offsetQuery,
	)

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var u user.Get
		if err := rows.Scan(
			&u.ID,
			&u.Avatar,
			&u.FirstName,
			&u.LastName,
			&u.PhoneNumber,
			&u.Login,
			&u.BirthDate,
			&u.Email,
			&u.Role,
			&u.RegionID,
			&u.DistrictID,
			&u.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	countQuery := `SELECT COUNT(u.id) FROM users u WHERE u.deleted_at IS NULL AND u.id > 1`
	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0

	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select user count: %w", err)
	}

	return users, count, nil
}

func (r Repository) GetById(ctx context.Context, id int64) (user.Get, error) {
	if id <= 1 {
		return user.Get{}, errors.New("invalid user id")
	}

	var u user.Get
	query := `
		SELECT id, avatar, first_name, last_name, phone_number, login, birth_date, email, role, region_id, district_id, created_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
		LIMIT 1
	`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Avatar,
		&u.FirstName,
		&u.LastName,
		&u.PhoneNumber,
		&u.Login,
		&u.BirthDate,
		&u.Email,
		&u.Role,
		&u.RegionID,
		&u.DistrictID,
		&u.CreatedAt,
	)
	if err != nil {
		return user.Get{}, err
	}

	return u, nil
}

func (r Repository) Update(ctx context.Context, id int64, data user.Update, adminId int64) error {
	setParts := []string{}
	args := []any{}

	if id <= 1 {
		return errors.New("invalid user id")
	}

	if data.Avatar != nil {
		setParts = append(setParts, "avatar = ?")
		args = append(args, *data.Avatar)
	}

	if data.FirstName != nil {
		setParts = append(setParts, "first_name = ?")
		args = append(args, *data.FirstName)
	}

	if data.LastName != nil {
		setParts = append(setParts, "last_name = ?")
		args = append(args, *data.LastName)
	}

	if data.PhoneNumber != nil {
		setParts = append(setParts, "phone_number = ?")
		args = append(args, *data.PhoneNumber)
	}

	if data.Password != nil {
		setParts = append(setParts, "password = ?")
		args = append(args, *data.Password)
	}

	if data.Login != nil {
		setParts = append(setParts, "login = ?")
		args = append(args, *data.Login)
	}

	if data.Status != nil {
		setParts = append(setParts, "status = ?")
		args = append(args, *data.Status)
	}

	if data.BirthDate != nil {
		layout := "2006-01-02"
		parsedTime, err := time.Parse(layout, *data.BirthDate)
		if err != nil {
			return fmt.Errorf("invalid birth_date format: %w", err)
		}
		setParts = append(setParts, "birth_date = ?")
		args = append(args, parsedTime.Format("2006-01-02"))
	}

	if data.Email != nil {
		setParts = append(setParts, "email = ?")
		args = append(args, *data.Email)
	}

	if data.Role != nil {
		setParts = append(setParts, "role = ?")
		args = append(args, *data.Role)
	}

	if data.RegionId != nil {
		setParts = append(setParts, "region_id = ?")
		args = append(args, *data.RegionId)
	}

	if data.DistrictId != nil {
		setParts = append(setParts, "district_id = ?")
		args = append(args, *data.DistrictId)
	}

	if len(setParts) == 0 {
		return nil
	}

	setParts = append(setParts, "updated_at = NOW()")

	setParts = append(setParts, "updated_by = ?")
	args = append(args, adminId)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ? AND deleted_at IS NULL", strings.Join(setParts, ", "))
	args = append(args, id)

	_, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id, adminId int64) error {

	if id <= 1 {
		return errors.New("invalid user id")
	}

	query := "UPDATE users SET deleted_at = NOW(), deleted_by = ? WHERE id = ? AND deleted_at IS NULL"
	_, err := r.DB.ExecContext(ctx, query, adminId, id)
	if err != nil {
		return err
	}

	return nil
}
