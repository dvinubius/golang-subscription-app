package data

import (
	"database/sql"
	"fmt"
	"time"
)

// TestNew is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application. This
// is only used when running tests.
func TestNew(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: DBUsersTest{},
		Plan: DBPlansTest{},
	}
}

type DBUsersTest struct{}
type DBPlansTest struct{}

// GetAll returns a slice of all users, sorted by last name
func (dbuTest DBUsersTest) GetAll() ([]*User, error) {
	var users []*User

	user := User{
		ID:        1,
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "Admin",
		Password:  "abc",
		Active:    1,
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	users = append(users, &user)

	return users, nil
}

// GetByEmail returns one user by email
func (dbuTest DBUsersTest) GetByEmail(email string) (*User, error) {
	user := User{
		ID:        1,
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "Admin",
		Password:  "abc",
		Active:    1,
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

// GetOne returns one user by id
func (dbuTest DBUsersTest) GetOne(id int) (*User, error) {
	return dbuTest.GetByEmail("")
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (dbuTest DBUsersTest) Update(user *User) error {
	return nil
}

// Delete deletes one user from the database, by User.ID
func (dbuTest DBUsersTest) Delete(user *User) error {
	return nil
}

// DeleteByID deletes one user from the database, by ID
func (dbuTest DBUsersTest) DeleteByID(id int) error {
	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (dbuTest DBUsersTest) Insert(user *User) (int, error) {
	return 2, nil
}

// ResetPassword is the method we will use to change a user's password.
func (dbuTest DBUsersTest) ResetPassword(user *User, password string) error {
	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (dbuTest DBUsersTest) PasswordMatches(user *User, plainText string) (bool, error) {
	return true, nil
}

func (dbpTest DBPlansTest) GetAll() ([]*Plan, error) {
	var plans []*Plan

	plan := Plan{
		ID:         1,
		PlanName:   "Bronze Plan",
		PlanAmount: 1000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	plans = append(plans, &plan)

	return plans, nil
}

// GetOne returns one plan by id
func (dbpTest DBPlansTest) GetOne(id int) (*Plan, error) {
	plan := Plan{
		ID:         1,
		PlanName:   "Bronze Plan",
		PlanAmount: 1000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return &plan, nil
}

// SubscribeUserToPlan subscribes a user to one plan by insert
// values into user_plans table
func (dbpTest DBPlansTest) SubscribeUserToPlan(user User, plan Plan) error {
	return nil
}

// AmountForDisplay formats the price we have in the DB as a currency string
func (dbpTest DBPlansTest) AmountForDisplay(plan *Plan) string {
	amount := float64(plan.PlanAmount) / 100.0
	return fmt.Sprintf("$%.2f", amount)
}
