package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

func newUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New User Endpoint Hit")

	db, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/testdb?parseTime=true")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	fmt.Println(name)
	fmt.Println(email)

	db.Create(&User{Name: name, Email: email})
	fmt.Fprintf(w, "New User Successfully Created")
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/testdb?parseTime=true")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var users []User
	db.Find(&users)
	fmt.Println("{}", users)

	json.NewEncoder(w).Encode(users)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "JUst hello..........")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/testdb?parseTime=true")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var user User
	db.Where("name = ?", name).Find(&user)
	db.Delete(&user)

	fmt.Fprintf(w, "Successfully Deleted User")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/testdb?parseTime=true")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	var user User
	db.Where("name = ?", name).Find(&user)

	user.Email = email

	db.Save(&user)
	fmt.Fprintf(w, "Successfully Updated User")
}

func main() {
	fmt.Println("Go ORM Tutorial")

	db, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/testdb?parseTime=true")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()
	db.AutoMigrate(&User{})

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", hello).Methods("GET")
	myRouter.HandleFunc("/user/{name}/{email}", newUser).Methods("POST")
	myRouter.HandleFunc("/user", allUsers).Methods("GET")
	myRouter.HandleFunc("/user/{name}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{name}/{email}", updateUser).Methods("PUT")
	log.Fatal(http.ListenAndServe(":9000", myRouter))

}
