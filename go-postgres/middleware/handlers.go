package middleware

import (
	"github.com/joho/godotenv"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-postgres/models"
	"log"
	"net/http"
	"os"
	"strconv"
	
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type response struct {
	ID int64 `json:"id,omitempty"`
	Message string `json"message,omitempty"`
}


func creatConnection() *sql.DB {

	err := godotenv.Load(".env")


	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to DB")
	return db
}

// CreateUser create a user in database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header.Set("Context-type", "applicatopn/x-www-form-urlencoded")
	w.Header.Set("Access-Control-Allow-Origin", "*")
	w.Header.Set("Access-Control-Allow-Methods", "POST")
	w.Header.Set("Access-Control-Allow-Headers", "Content-Type")

	var user models.User

	err := json.NewDecoder(r.body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode request body. %v", err)
	}

	insertID := insertUser(user)

	res := response {
		ID: insertID,
		Message: "User created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

// GetUser will return a single user by id
func GetUser(w http.ResponseWriter,  r *http.Request) {
	w.Header.Set("Context-type", "applicatopn/x-www-form-urlencoded")
	w.Header.Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)

	id, err := strconv(params["id"])
	if err != nil {
		log.Fatalf("unable to convert the string into int. %v", err)
	}

	user, err := getUser(int64(id))

	if err != null {
		log.Fatalf("Unable to get user. %v", err)
	}

	json.NewEncoder(w).Encode(user)
}

// GetAllUsers returns all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	users, err := getAllUsers()

	if err != nil {
		log.Fataf("Unable to get all users %v", err)
	}
	
	json.NewEncoder(w).Encode(users)
}

// UpdateUser update user's detail in database
func UpdateUser(w http.ResponseWriter, r *http.Request) {	
	w.Header.Set("Context-type", "applicatopn/x-www-form-urlencoded")
	w.Header.Set("Access-Control-Allow-Origin", "*")
	w.Header.Set("Access-Control-Allow-Methods", "PUT")
	w.Header.Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	updatedRows := updateUser(int64(id), user)

	msg := fmt.Sprintf("User has been updated, Total records affected %v", updatedRows)

	res := response{
		ID: int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

}

// DeleteUser delete user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header.Set("Context-type", "applicatopn/x-www-form-urlencoded")
	w.Header.Set("Access-Control-Allow-Origin", "*")
	w.Header.Set("Access-Control-Allow-Methods", "PUT")
	w.Header.Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	deletedRows := deleteUser(int64(id))

	msg := fmt.Sprintf("User deleted sucessfully. Total records affected %v", deletedRows)

	res := response{
		ID: int64(id),
		Message: msg,

	}

	json.NewEncoder(w).Encode(res)
}

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//																HANDLER FUNCTIONS 																//
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////											

func insertUser(user models.User) int64 {
	db := creatConnection()
	defer db.Close()

	sqlQuery := `INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid`

	var id int64

	err := db.Query.Row(sqlQuery, user.Name, user.Location, user.Age).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	return id
}

func getUser(id int64) (models.User, error) {
	db := creatConnection()
	defer db.Close()

	var user models.User

	sqlQuery := `SELECT * FROM users WHERE userid=$1`
	row := db.QueryRow(sqlQuery, id)

	err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, nil
	case nil:
		return user, nil
	default: 
		log.Fatalf("Unable to scan the row. %v", err)		
	}

	return user, err
}

func getAllUsers() ([]models.User, error) {
	db := creatConnection()
	defer db.Close()

	var users []models.User

	sqlQuery := `SELECT * FROM users`

	rows, err := db.Query(sqlQuery)

	if err != nil {
		log.Fatalf("Unable to execute the Query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		users = append(users, user)
	}

	return users, err
}


func updateUser(id int64, user models.User) int64 {
	db := creatConnection()

	defer db.Close()

	sqlQuery := `UPDATE users SET name=$2, location=$3, age=$4 WHERE userid=$1`

	res, err := db.Exec(sqlQuery, id, user.Name, user.Location, user.Age)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v",  err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the effected Rows. %v", err)
	}
	fmt.Printf("Total records affected %v", rowsAffected)

	return rowsAffected
}

func deleteUser(id int64) int64 {
	db := creatConnection()

	defer db.Close()

	sqlQuery := `DELETE FROM users WHERE userid=$1`
	 res, err := db.Exec(sqlQuery, id)

	 if err != nil {
		 log.Fatalf("unable to execute query. %v", err)

	 }

	 rowsAffected, err := res.RowsAffected()

	 if err != nil {
		 log.Fatalf("Error while checking effected rows. %v", err)
	 }
	 
	 fmt.Printf("Total records affected. %v", rowsAffected)

	 return rowsAffected
}
