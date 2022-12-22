package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/user"

	// "github.com/nkh361/go-todo-list/pkg/users"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type Ticket struct {
	ID       int64
	username string
	title    string
	priority float32
}

// https://go.dev/doc/tutorial/database-access left off at add data

func main() {
	// initial config
	cfg := mysql.Config{
		User: os.Getenv("DBUSER"),
		// Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "ticketing",
	}
	// get a database handle
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected")

	// add new user
	createNewUser, err := newUser(user.User{
		username: "nho21",
		password: "amazing",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("new user id: %d", createNewUser)

	// test check password
	check_password, err := checkPassword("nate", "test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("success: %v", check_password)

	// query by username
	tickets, err := ticketByUsername("nkh361")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tickets found: %v\n", tickets)

	// query by priority
	ticketPriority, errPriority := ticketByPriority(5)
	if errPriority != nil {
		log.Fatal(errPriority)
	}
	fmt.Printf("tickets found by priority 5: %v\n", ticketPriority)

	// add ticket
	// ticID, errAdd := addTicket(Ticket{
	// 	username: "nkh361",
	// 	title:    "apply for jobs",
	// 	priority: 30,
	// })
	// if errAdd != nil {
	// 	log.Fatal(errAdd)
	// }
	// fmt.Printf("ID of added ticket: %v\n", ticID)

	// query by id
	ticketID, errID := ticketByID(4)
	if errID != nil {
		log.Fatal(errID)
	}
	fmt.Printf("tickets found by id: %v\n", ticketID)
}

func ticketByUsername(username string) ([]Ticket, error) {
	// a tickets slice to hold data from returned rows
	var tickets []Ticket

	rows, err := db.Query("SELECT * FROM tickets WHERE username = ?", username)
	if err != nil {
		return nil, fmt.Errorf("ticketByUsername %q: %v", username, err)
	}
	defer rows.Close()
	// loop through rows, using scan to assign column data for rows
	for rows.Next() {
		var tic Ticket
		if err := rows.Scan(&tic.ID, &tic.username, &tic.title, &tic.priority); err != nil {
			return nil, fmt.Errorf("ticketByUsername %q: %v", username, err)
		}
		tickets = append(tickets, tic)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ticketByUsername %q: %v", username, err)
	}
	return tickets, nil
}

func ticketByPriority(priority int) ([]Ticket, error) {
	var tickets []Ticket

	rows, err := db.Query("SELECT * FROM tickets WHERE priority = ?", priority)
	if err != nil {
		return nil, fmt.Errorf("ticketByPriority %q: %v", priority, err)
	}
	defer rows.Close()
	for rows.Next() {
		var tic Ticket
		if err := rows.Scan(&tic.ID, &tic.username, &tic.title, &tic.priority); err != nil {
			return nil, fmt.Errorf("ticketByPriority %q: %v", priority, err)
		}
		tickets = append(tickets, tic)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ticketByPriority %q: %v", priority, err)
	}
	return tickets, nil
}

func ticketByID(id int64) (Ticket, error) {
	var tic Ticket

	row := db.QueryRow("SELECT * FROM tickets WHERE id = ?", id)
	if err := row.Scan(&tic.ID, &tic.username, &tic.title, &tic.priority); err != nil {
		if err == sql.ErrNoRows {
			return tic, fmt.Errorf("ticketByID %d: no such ticket", id)
		}
		return tic, fmt.Errorf("ticketByID %d: %v", id, err)
	}
	return tic, nil
}

func addTicket(tic Ticket) (int64, error) {
	result, err := db.Exec("INSERT INTO tickets (username, title, priority) VALUES (?,?,?)", tic.username, tic.title, tic.priority)
	if err != nil {
		return 0, fmt.Errorf("addTicket: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addTicket: %v", err)
	}
	return id, nil
}

func newUser(usr user.User) (int64, error) {
	hashedPassword := user.getHashedPassword(usr.password)
	result, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", usr.username, hashedPassword)
	if err != nil {
		return 0, fmt.Errorf("newUser: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("newUser: %v", err)
	}
	return id, nil
}

func checkPassword(username string, password string) (bool, error) {
	var comparison string
	if err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&comparison); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("username does not exist: %q", username)
		}
		return false, fmt.Errorf("check password: %v", err)
	}

	err := bcrypt.CompareHashAndPassword([]byte(comparison), []byte(password))
	if err != nil {
		log.Fatal(err)
	} else {
		return true, nil
	}
	return false, err
}
