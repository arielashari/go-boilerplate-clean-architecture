package postgres

import (
	"context"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/repository/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type rolePostgresRepository struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

var _ entity.RolePostgresRepository = (*rolePostgresRepository)(nil)

func NewRolePostgresRepository(
	db *pgxpool.Pool,
) entity.RolePostgresRepository {
	return &rolePostgresRepository{queries: sqlc.New(db), pool: db}
}

func (r *rolePostgresRepository) Create(ctx context.Context, role *entity.Role) (*entity.Role, error) {
	row, err := r.queries.CreateRole(ctx, sqlc.CreateRoleParams{
		ID:   role.ID,
		Name: role.Name,
	})

	if err != nil {
		return nil, ParseError(err, "RoleRepository.Create")
	}

	return r.toEntity(row), nil
}

func (r *rolePostgresRepository) Delete(ctx context.Context, id string) error {
	result, err := r.queries.DeleteRole(ctx, id)
	if err != nil {
		return ParseError(err, "RoleRepository.Delete")
	}
	if result.RowsAffected() == 0 {
		return entity.ErrNotFound.WithOperation("RoleRepository.Delete")
	}
	return nil
}

func (r *rolePostgresRepository) GetByID(ctx context.Context, id string) (*entity.Role, error) {
	row, err := r.queries.GetRoleByID(ctx, id)
	if err != nil {
		return nil, ParseError(err, "RoleRepository.GetByID")
	}

	return r.toEntity(row), nil
}

func (r *rolePostgresRepository) List(ctx context.Context, limit, offset int) ([]*entity.Role, error) {
	rows, err := r.queries.ListRoles(ctx, sqlc.ListRolesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, ParseError(err, "RoleRepository.List")
	}

	roles := make([]*entity.Role, len(rows))
	for i, row := range rows {
		roles[i] = r.toEntity(row)
	}
	return roles, nil
}

func (r *rolePostgresRepository) Update(ctx context.Context, role *entity.Role) (*entity.Role, error) {
	row, err := r.queries.UpdateRole(ctx, sqlc.UpdateRoleParams{
		ID:   role.ID,
		Name: role.Name,
	})

	if err != nil {
		return nil, ParseError(err, "RoleRepository.Update")
	}

	return r.toEntity(row), nil
}

func (r *rolePostgresRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountRoles(ctx)
	if err != nil {
		return 0, ParseError(err, "RoleRepository.Count")
	}
	return count, nil
}

func (r *rolePostgresRepository) toEntity(row sqlc.Role) *entity.Role {
	return &entity.Role{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: row.DeletedAt,
	}
}
