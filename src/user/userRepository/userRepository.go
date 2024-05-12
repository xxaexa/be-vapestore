package userRepository

import (
	"clean-architecture/model/dto/userDto"
	"clean-architecture/model/entity"
	"clean-architecture/src/user"
	"database/sql"
	"fmt"
	"strings"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.UserRepository {
	return &userRepository{db}
}

func (repo *userRepository) CreateUser(user *userDto.CreateUserRequest) error {
	sqlQuery := `INSERT INTO users (email,fullname, password) VALUES ($1, $2,$3)`
	_, err := repo.db.Exec(sqlQuery, user.Email, user.FullName, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) GetUserByEmail(email string) (*entity.User, error) {
	sqlQuery := `SELECT id, email, fullname, password FROM users WHERE email = $1`
	row := repo.db.QueryRow(sqlQuery, email)
	u := new(entity.User)
	err := row.Scan(&u.ID, &u.Email, &u.FullName, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return u, nil
}

func (repo *userRepository) GetUserByID(id string) (*entity.User, error) {
	sqlQuery := `SELECT id, email, email, password FROM users WHERE id = $1`
	rows, err := repo.db.Query(sqlQuery, id)
	if err != nil {
		return nil, err
	}
	u := new(entity.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (repo *userRepository) GetUsers(page, limit int, email, fullName string) ([]*entity.User, int, error) {
	offset := (page - 1) * limit
	baseQuery := "SELECT id, fullname, email, password FROM users WHERE deleted_at IS NULL"

	var conditions []string
	var args []interface{}

	if email != "" {
		conditions = append(conditions, "email = $1")
		args = append(args, email)
	}
	if fullName != "" {
		if email != "" {
			conditions = append(conditions, "fullname = $2")
			args = append(args, fullName)
		} else {
			conditions = append(conditions, "fullname = $1")
			args = append(args, fullName)
		}
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	orderLimitOffset := fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	baseQuery += orderLimitOffset

	args = append(args, limit, offset)

	rows, err := repo.db.Query(baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := new(entity.User)
		if err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Password); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	count, err := GetTotalUsers(repo.db)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (repo *userRepository) UpdateUser(user *userDto.UpdateUserRequest) error {
	query := "UPDATE users SET fullname = $2, password = $3 WHERE id = $1"
	_, err := repo.db.Exec(query, user.ID, user.FullName, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) DeleteUser(id string) error {
	query := "UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL"

	_, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func GetTotalUsers(db *sql.DB) (int, error) {
	count := 0
	sqlQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	row := db.QueryRow(sqlQuery)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func scanRowsIntoUser(rows *sql.Rows) (*entity.User, error) {
	user := new(entity.User)

	err := rows.Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
