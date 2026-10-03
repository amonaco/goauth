package db

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5"
)

// User represents a user record in PostgreSQL.
type User struct {
    ID        int32
    Email     string
    FirstName string
    LastName  string
    Password  string
}

// FindUserByEmail loads a user by email.
func FindUserByEmail(ctx context.Context, email string) (User, error) {
    var user User
    err := QueryRow(ctx, `SELECT id, email, first_name, last_name, password_hash FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password)
    if err != nil {
        return User{}, err
    }
    return user, nil
}

// FindUserByID loads a user by database ID.
func FindUserByID(ctx context.Context, userID int32) (User, error) {
    var user User
    err := QueryRow(ctx, `SELECT id, email, first_name, last_name, password_hash FROM users WHERE id = $1`, userID).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password)
    if err != nil {
        return User{}, err
    }
    return user, nil
}

// FindCompanyRoles loads the roles for a user within a company.
func FindCompanyRoles(ctx context.Context, userID, companyID int32) ([]string, error) {
    rows, err := DB().Query(ctx, `SELECT r.name FROM user_roles ur JOIN roles r ON r.id = ur.role_id WHERE ur.user_id = $1 AND ur.company_id = $2`, userID, companyID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var roles []string
    for rows.Next() {
        var role string
        if err := rows.Scan(&role); err != nil {
            return nil, err
        }
        roles = append(roles, role)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return roles, nil
}

// Company represents a company record.
type Company struct {
    ID   int32
    Name string
    Slug string
}

// FindCompanyByID loads a company by ID.
func FindCompanyByID(ctx context.Context, companyID int32) (Company, error) {
    var company Company
    err := QueryRow(ctx, `SELECT id, name, slug FROM companies WHERE id = $1`, companyID).Scan(&company.ID, &company.Name, &company.Slug)
    if err != nil {
        return Company{}, err
    }
    return company, nil
}

// ValidateUserPassword is a placeholder for a password verification implementation.
func ValidateUserPassword(ctx context.Context, email, password string) (User, error) {
    user, err := FindUserByEmail(ctx, email)
    if err != nil {
        return User{}, err
    }
    if password == "" {
        return User{}, fmt.Errorf("empty password")
    }
    return user, nil
}

var _ pgx.Row
