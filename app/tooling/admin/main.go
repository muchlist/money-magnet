package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/muchlist/moneymagnet/business/user/model"
	urrepo "github.com/muchlist/moneymagnet/business/user/repo"
	urserv "github.com/muchlist/moneymagnet/business/user/service"
	"github.com/muchlist/moneymagnet/cfg"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/mcrypto"
	"github.com/muchlist/moneymagnet/pkg/mjwt"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/validate"
)

var grantType = []string{
	"user",
	"admin",
}

func main() {
	// config
	config := cfg.Load()

	// init log
	log := mlogger.New(mlogger.Options{
		Level:  mlogger.LevelInfo,
		Output: "stdout",
	})

	// dependency
	// init database
	database, err := db.OpenDB(db.Config{
		DSN:          config.DB.DSN,
		MaxOpenConns: config.DB.MaxOpenCons,
		MinOpenConns: config.DB.MinOpenCons,
	})
	if err != nil {
		log.Error("connection to database", err)
		panic(err.Error())
	}
	defer database.Close()

	jwt := mjwt.New(config.App.Secret)
	bcrypt := mcrypto.New()

	// middleware
	userRepo := urrepo.NewRepo(database, log)

	userService := urserv.NewCore(log, userRepo, bcrypt, jwt)

	inputHint := []string{"name", "email", "password", "roles"}
	inputValue := make([]string, len(inputHint))

	fmt.Println("=====================================")
	fmt.Println("Tools to inject user")
	fmt.Println("Fill in the data below!")
	fmt.Println("=====================================")

	// colect input from terminal
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < len(inputHint); i++ {
		fmt.Printf("input %s\t:", inputHint[i])
		input, err := reader.ReadString('\n') // read input until hit enter
		if err != nil {
			panic(fmt.Sprintf("error read input %v", err))
		}
		input = strings.TrimSpace(input)
		inputValue[i] = input
	}

	// validating input
	req := model.UserRegisterReq{
		Name:     inputValue[0],
		Email:    inputValue[1],
		Password: inputValue[2],
		Roles:    []string{inputValue[3]},
	}

	if err = validateRoles(inputValue[3]); err != nil {
		fmt.Println("=====================================")
		fmt.Println("input not valid: ", err)
		return
	}

	validator := validate.New()
	_, err = validator.Struct(req)
	if err != nil {
		fmt.Println("=====================================")
		fmt.Println("input not valid: ", err)
		return
	}

	result, err := userService.InsertUser(context.Background(), req)
	if err != nil {
		fmt.Println("=====================================")
		fmt.Println("error register user: ", err)
		return
	}
	fmt.Println("=====================================")
	fmt.Printf("Success register user : %v", result)
	fmt.Println("=====================================")
}

func validateRoles(role string) error {
	for _, v := range grantType {
		if role == v {
			return nil
		}
	}
	return fmt.Errorf("role must be one of: %v", grantType)
}
