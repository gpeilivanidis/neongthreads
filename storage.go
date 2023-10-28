package main

import (
	"database/sql"
	"log/slog"
)

type Storage interface {
	CreateUser(User) (*User, error)
	GetAllUsers() ([]User, error)
	GetUserById(int) (*User, error)
	GetUserByUsername(string) (*User, error)
	UpdateUser(User) error
	DeleteUserById(int) error

	CreateProduct(Product) (*Product, error)
	GetAllProducts() ([]Product, error)
	GetProductsByType(string) ([]Product, error)
	GetProductById(int) (*Product, error)
	GetProductByTitle(string) (*Product, error)
	UpdateProduct(Product) error
	DeleteProductById(int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=neongthreads sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	if err := s.createUserTable(); err != nil {
		return err
	}
	if err := s.createProductTable(); err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) createUserTable() error {
	query := `create table if not exists users (
		id serial primary key,
		username text,
		passwordHashed text,
		level integer
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createProductTable() error {
	query := `create table if not exists products (
		id serial primary key,
		type text
		title text,
		description text,
		price money,
		gender text,
		color text,
		small integer,
		medium integer,
		large integer,
		imageUrl text,
		imageAlt text,
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateUser(u User) (*User, error) {
	query := `insert into users (username, passwordHashed, level)
	values ($1, $2, $3) returning id`

	row := s.db.QueryRow(query, u.Username, u.PasswordHashed, u.Level)

	err := row.Scan(&u.Id)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return &u, nil
}

func (s *PostgresStore) GetAllUsers() ([]User, error) {
	query := `select * from users`

	rows, err := s.db.Query(query)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	users := []User{}
	for rows.Next() {
		user := User{}

		err := rows.Scan(&user.Id, &user.Username, &user.PasswordHashed, &user.Level)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return users, nil
}

func (s *PostgresStore) GetUserById(id int) (*User, error) {
	query := `select * from users where id = $1 limit 1`
	row := s.db.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.Id, &user.Username, &user.PasswordHashed, &user.Level)

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (s *PostgresStore) GetUserByUsername(username string) (*User, error) {
	query := `select * from users where username = $1 limit 1`
	row := s.db.QueryRow(query, username)

	user := &User{}
	err := row.Scan(&user.Id, &user.Username, &user.PasswordHashed, &user.Level)

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (s *PostgresStore) UpdateUser(u User) error {
	query := `update users set username=$1, passwordHashed=$2, level=$3 where id=$4`

	_, err := s.db.Exec(query, u.Username, u.PasswordHashed, u.Level, u.Id)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteUserById(id int) error {
	query := `delete from users where id=$1`

	_, err := s.db.Exec(query, id)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func (s *PostgresStore) CreateProduct(p Product) (*Product, error) {
	query := `insert into products
	(type, title, description, price, gender, color, small, medium,large, imageUrl, imageAlt)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	returning id`

	row := s.db.QueryRow(query, p.Type, p.Title, p.Description, p.Price, p.Gender, p.Color, p.Small, p.Medium, p.Large, p.ImageUrl, p.ImageAlt)

	err := row.Scan(&p.Id)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return &p, nil
}

func (s *PostgresStore) GetAllProducts() ([]Product, error) {
	query := `select * from products`
	rows, err := s.db.Query(query)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	products := []Product{}
	for rows.Next() {
		p := Product{}

		err := rows.Scan(&p.Id, &p.Type, &p.Title, &p.Description, &p.Price, &p.Gender, &p.Color, &p.Small, &p.Medium, &p.Large, &p.ImageUrl, &p.ImageAlt)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return products, nil
}

func (s *PostgresStore) GetProductsByType(t string) ([]Product, error) {
	query := `select * from products where type=$1`
	rows, err := s.db.Query(query, t)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	products := []Product{}
	for rows.Next() {
		p := Product{}

		err := rows.Scan(&p.Id, &p.Type, &p.Title, &p.Description, &p.Price, &p.Gender, &p.Color, &p.Small, &p.Medium, &p.Large, &p.ImageUrl, &p.ImageAlt)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return products, nil
}

func (s *PostgresStore) GetProductById(id int) (*Product, error) {
	query := `select * from products where id = $1 limit 1`
	row := s.db.QueryRow(query, id)

	p := &Product{}
	err := row.Scan(&p.Id, &p.Type, &p.Title, &p.Description, &p.Price, &p.Gender, &p.Color, &p.Small, &p.Medium, &p.Large, &p.ImageUrl, &p.ImageAlt)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return p, nil
}

func (s *PostgresStore) GetProductByTitle(title string) (*Product, error) {
	query := `select * from products where title = $1 limit 1`
	row := s.db.QueryRow(query, title)

	p := &Product{}
	err := row.Scan(&p.Id, &p.Type, &p.Title, &p.Description, &p.Price, &p.Gender, &p.Color, &p.Small, &p.Medium, &p.Large, &p.ImageUrl, &p.ImageAlt)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return p, nil
}

func (s *PostgresStore) UpdateProduct(p Product) error {
	query := `update products 
	set type=$1,
	title=$2,
	description=$3,
	price=$4,
	gender=$5,
	color=$6,
	small=$7,
	medium=$8,
	large=$9,
	imageUrl=$10,
	imageAlt=$11 where id=$12`

	_, err := s.db.Exec(query, p.Type, p.Title, p.Description, p.Price, p.Gender, p.Color, p.Small, p.Medium, p.Large, p.ImageUrl, p.ImageAlt)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteProductById(id int) error {
	query := `delete from products where id=$1`
	_, err := s.db.Exec(query)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
