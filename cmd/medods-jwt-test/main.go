package main

import (
	"github.com/joho/godotenv"
	"github.com/pheezz/medods-jwt-test/internal/pkg/app"
	"log"
)

// type UserBaseSchema struct {
// 	GUID      string
// 	tokenPair token.TokenPair
// }

// var users = []UserBaseSchema{
// 	{
// 		GUID: "1",
// 		tokenPair: token.TokenPair{
// 			AccessToken:  "access_token_one",
// 			RefreshToken: "refresh_one"},
// 	},
// 	{
// 		GUID: "2",
// 		tokenPair: token.TokenPair{
// 			AccessToken:  "access_token_two",
// 			RefreshToken: "refresh_two"},
// 	}}

// var conf = config.Conf

// type TokenPairCookies struct {
// 	AccessToken  http.Cookie
// 	RefreshToken http.Cookie
// }

func init() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}

}
