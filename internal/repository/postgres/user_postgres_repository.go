package postgres

import (
	"context"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/repository/postgres/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userPostgresRepository struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

var _ entity.UserPostgresRepository = (*userPostgresRepository)(nil)

func NewUserPostgresRepository(
	db *pgxpool.Pool,
) entity.UserPostgresRepository {
	return &userPostgresRepository{queries: sqlc.New(db), pool: db}
}

func (r *userPostgresRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	row, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:          user.ID,
		Email:       user.Email,
		Password:    user.Password,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhonePrefix: user.PhonePrefix,
		PhoneNumber: user.PhoneNumber,
		RoleID:      user.RoleID,
	})

	if err != nil {
		return nil, ParseError(err, "UserRepository.Create")
	}

	return r.toEntity(row, nil), nil
}
func (r *userPostgresRepository) Delete(ctx context.Context, id string) error {
	result, err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return ParseError(err, "UserRepository.Delete")
	}
	if result.RowsAffected() == 0 {
		return entity.ErrNotFound.WithOperation("UserRepository.Delete")
	}
	return nil
}
func (r *userPostgresRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ParseError(err, "UserRepository.GetByEmail")
	}

	return r.toEntity(row, nil), nil
}
func (r *userPostgresRepository) List(ctx context.Context, query entity.UserQuery) ([]*entity.User, error) {
	params := r.mapUserQueryToParams(query)
	rows, err := r.queries.ListUsers(ctx, params)
	if err != nil {
		return nil, ParseError(err, "UserRepository.List")
	}

	users := make([]*entity.User, len(rows))
	for i, row := range rows {
		users[i] = r.toEntity(sqlc.User{
			ID:          row.User.ID,
			Email:       row.User.Email,
			FirstName:   row.User.FirstName,
			LastName:    row.User.LastName,
			PhonePrefix: row.User.PhonePrefix,
			PhoneNumber: row.User.PhoneNumber,
			RoleID:      row.User.RoleID,
			CreatedAt:   row.User.CreatedAt,
			UpdatedAt:   row.User.UpdatedAt,
			DeletedAt:   row.User.DeletedAt,
		}, nil)
	}

	return users, nil
}
func (r *userPostgresRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, ParseError(err, "UserRepository.GetByID")
	}

	return r.toEntity(row.User, &row.Role), nil
}
func (r *userPostgresRepository) Update(ctx context.Context, u *entity.User) (*entity.User, error) {
	rows, err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:          u.ID,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		PhonePrefix: u.PhonePrefix,
		PhoneNumber: u.PhoneNumber,
	})
	if err != nil {
		return nil, ParseError(err, "UserRepository.Update")
	}
	return r.toEntity(rows, nil), nil
}

func (r *userPostgresRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountUsers(ctx)
	if err != nil {
		return 0, ParseError(err, "UserRepository.Count")
	}
	return count, nil
}

func (r *userPostgresRepository) GetUserForAuth(ctx context.Context, email string) (*entity.User, error) {
	row, err := r.queries.GetUserForAuth(ctx, email)
	if err != nil {
		return nil, ParseError(err, "UserRepository.GetUserForAuth")
	}

	return &entity.User{
		ID:       row.ID,
		Email:    row.Email,
		Password: row.Password,
	}, nil
}

func (r *userPostgresRepository) getQueries(ctx context.Context) *sqlc.Queries {
	if tx := extractTx(ctx); tx != nil {
		return sqlc.New(tx)
	}
	return r.queries
}

func (r *userPostgresRepository) toEntity(row sqlc.User, role *sqlc.Role) *entity.User {
	user := &entity.User{
		ID:          row.ID,
		Email:       row.Email,
		FirstName:   row.FirstName,
		LastName:    row.LastName,
		PhonePrefix: row.PhonePrefix,
		PhoneNumber: row.PhoneNumber,
		RoleID:      row.RoleID,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	if role != nil {
		user.RoleID = role.ID
		user.Role = &entity.Role{
			ID:        role.ID,
			Name:      role.Name,
			CreatedAt: role.CreatedAt,
			UpdatedAt: role.UpdatedAt,
			DeletedAt: role.DeletedAt,
		}
	}

	if row.DeletedAt.Valid {
		user.DeletedAt = &row.DeletedAt.Time
	} else {
		user.DeletedAt = nil
	}

	return user
}

func (r *userPostgresRepository) mapUserQueryToParams(q entity.UserQuery) sqlc.ListUsersParams {
	offset := (q.Page - 1) * q.Limit

	return sqlc.ListUsersParams{
		Limit:   int32(q.Limit),
		Offset:  int32(offset),
		Search:  pgtype.Text{String: q.Search, Valid: q.Search != ""},
		RoleID:  toUUID(q.RoleID),
		SortBy:  q.SortBy,
		SortDir: q.SortDir,
	}
}

func toUUID(s string) pgtype.UUID {
	if s == "" {
		return pgtype.UUID{Valid: false}
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: id, Valid: true}
}
