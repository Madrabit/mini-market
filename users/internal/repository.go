package internal

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) BeginTransaction() (tx *sqlx.Tx, err error) {
	return r.db.Beginx()
}

func (r *Repository) CreateUser(ctx context.Context, tx *sqlx.Tx, user User) error {
	if _, err := tx.ExecContext(ctx, "INSERT INTO users (id, name, email, password_hash) VALUES  ($1, $2, $3, $4) ",
		user.Id, user.Name, user.Email, user.PasswordHash); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAllRoles(ctx context.Context) ([]Role, error) {
	var roles []Role
	if err := r.db.SelectContext(ctx, &roles, "SELECT id, name, created_at, updated_at FROM roles;"); err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *Repository) AddUserRoles(ctx context.Context, tx *sqlx.Tx, userId uuid.UUID, roles []uuid.UUID) error {
	values := make([]string, 0, len(roles))
	args := make([]interface{}, 0, 1+len(roles))
	args = append(args, userId)
	for i, roleID := range roles {
		placeholder := fmt.Sprintf("($1, $%d)", i+2)
		values = append(values, placeholder)
		args = append(args, roleID)
	}
	query := fmt.Sprintf(`INSERT INTO user_roles (user_id, role_id) 
		VALUES %s ON CONFLICT (user_id, role_id) 
        DO NOTHING;`, strings.Join(values, ","))
	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, user User) error {
	result, err := r.db.ExecContext(ctx, `UPDATE users SET name = $1, email = $2  WHERE id = $3 `,
		user.Name, user.Email, user.Id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID uuid.UUID) (User, error) {
	var user User
	if err := r.db.GetContext(ctx, &user, "SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE id=$1", userID); err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *Repository) GetUsersByIds(ctx context.Context, IDs []uuid.UUID) ([]User, error) {
	var users []User
	q, args, err := sqlx.In("SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE id IN (?)", IDs)
	if err != nil {
		return nil, err
	}
	q = r.db.Rebind(q)
	err = r.db.SelectContext(ctx, &users, q, args...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) CreateRole(ctx context.Context, role Role) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO roles (id, name) VALUES ($1, $2)",
		role.Id, role.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteRole(ctx context.Context, id uuid.UUID) (err error) {
	_, err = r.db.ExecContext(ctx, "DELETE FROM roles WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateRole(ctx context.Context, id uuid.UUID, role string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE roles SET name = $1 WHERE id = $2 `,
		role, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) GetUsersByRole(ctx context.Context, role string) ([]User, error) {
	var users []User
	err := r.db.SelectContext(ctx, &users, `SELECT * FROM 
             users INNER JOIN user_roles ON users.id = user_roles.user_id
        	 INNER JOIN roles r on r.id = user_roles.role_id	
             WHERE roles.name = $1`, role)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetRoleByName(ctx context.Context, name string) (Role, error) {
	var role Role
	err := r.db.GetContext(ctx, &role, `SELECT id, name FROM roles WHERE name = $1`, name)
	if err != nil {
		return Role{}, err
	}
	return role, nil
}

func (r *Repository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error) {
	var roles []Role
	err := r.db.SelectContext(ctx, &roles, `
	SELECT	roles.id, roles.name
	FROM user_roles
	INNER JOIN roles ON user_roles.role_id = roles.id
	WHERE user_roles.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	return roles, nil
}
