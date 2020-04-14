package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func init() {

	configRoot, _ := os.Getwd()
	var configPath = flag.String("config", configRoot, "defines path to config file")

	flag.Parse()

	// setup config file
	viper.SetConfigName("config")    // name of config file (without extension)
	viper.AddConfigPath(*configPath) // path to look for the config file in
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Couldn't find config.toml!")
		log.Fatal(err)
	}

}

func main() {

	// create db instance to deprecated db
	dbOld, err := sql.Open("mysql", viper.GetString("connectionDeprecatedDB"))
	if err != nil {
		log.Fatal(err)
	}
	err = dbOld.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer dbOld.Close()

	// create db instance to actual db
	db, err := sql.Open("mysql", viper.GetString("connectionDB"))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	results, err := dbOld.Query(`
	SELECT
			id, 
			COALESCE(username, '') as username,
			email, lastname, firstname, mobile, need_sms,
			street, zip, city, 
			COALESCE(date_of_birth, '') as date_of_birth,			
			COALESCE(date_of_entry, '') as date_of_entry,
			COALESCE(date_of_exit, '') as date_of_exit,
			state, credit, credit_date, credit_comment,
			COALESCE(iban, '') as iban,
			COALESCE(bic, '') as bic,
			COALESCE(sepa, '') as sepa,
			COALESCE(additionals, '') as additionals, 
			COALESCE(comment, '') as comment,
			COALESCE(group_comment, '') as group_comment,
			created_at, updated_at,
			COALESCE(last_login, '') as last_login
		FROM users`)

	if err != nil {
		fmt.Println(err)
	}

	defer results.Close()

	var users []UserDeprecated
	var newUsers []User
	for results.Next() {
		var user UserDeprecated
		err = results.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Lastname,
			&user.Firstname,
			&user.Mobile,
			&user.NeedSMS,
			&user.Street,
			&user.ZIP,
			&user.City,
			&user.DateOfBirth,
			&user.DateOfEntry,
			&user.DateOfExit,
			&user.State,
			&user.Credit,
			&user.CreditDate,
			&user.CreditComment,
			&user.IBAN,
			&user.BIC,
			&user.SEPA,
			&user.Additionals,
			&user.Comment,
			&user.GroupComment,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.LastLogin,
		)
		if err != nil {
			fmt.Println(err)
		}
		users = append(users, user)
	}
	for _, usr := range users {

		var state string
		if usr.State == 2 {
			state = "active"
		}
		newUser := User{
			ID:           usr.ID,
			FirstName:    usr.Firstname,
			LastName:     usr.Lastname,
			Email:        usr.Email,
			Phone:        usr.Mobile,
			Street:       usr.Street,
			StreetNumber: "0815", // TODO: split Street to get Number ?
			Zip:          usr.ZIP,
			Country:      "Germany",
			Birthday:     usr.DateOfBirth,
			EntranceDate: usr.DateOfEntry,
			// LeavingDate:     usr.DateOfExit,
			AdditionalInfos: usr.Comment,
			LastActivityAt:  usr.LastLogin,
			CreatedAt:       usr.CreatedAt,
			UpdatedAt:       usr.UpdatedAt,
			CreatedBy:       174,
			Password:        "2020_b!ok!ste_2020",
			State:           state,
		}

		newUsers = append(newUsers, newUser)
	}
	addUser := newUsers[143]
	id, err := AddUserReq(addUser)

	fmt.Println(id)
}