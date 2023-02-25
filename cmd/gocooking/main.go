package main

import (
	"fmt"
	"os"

	"github.com/jduepmeier/go-cooking"

	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v3"
)

type opts struct {
	Config string `short:"c" long:"config" default:"config.yaml" description:"Path to config file."`
}

type userOpts struct {
	Username string `short:"u" long:"username" description:"username to add"`
	Password string `short:"p" long:"password" description:"password for the username"`
}

func loadConfig(configFile string) (*cooking.Config, error) {
	config := cooking.DefaultConfig()
	file, err := os.Open(configFile)
	if err != nil {
		logrus.Errorf("cannot open config file %s: %s", configFile, err)
		return config, err
	}

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		logrus.Errorf("cannot parse config file %s: %s", configFile, err)
	}
	return config, err
}

func main() {
	var opts opts
	var userOpts userOpts
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.AddCommand("serve", "serve the server", "Serves a http server.", &struct{}{})
	if err != nil {
		logrus.Errorf("error adding command serve: %s", err)
		return
	}
	_, err = parser.AddCommand("add-user", "add a user", "this adds a user to the database", &userOpts)
	if err != nil {
		logrus.Errorf("error adding command add-user: %s", err)
		return
	}
	_, err = parser.AddCommand("set-password", "sets the password", "this sets the password for a user in the database", &userOpts)
	if err != nil {
		logrus.Errorf("error adding command add-user: %s", err)
		return
	}
	_, err = parser.AddCommand("list-users", "lists all users", "list all users from the database", &struct{}{})
	if err != nil {
		logrus.Errorf("error adding command add-user: %s", err)
		return
	}
	_, err = parser.Parse()
	if err != nil {
		return
	}
	config, err := loadConfig(opts.Config)
	if err != nil {
		return
	}
	logrus.SetLevel(logrus.DebugLevel)

	server, err := cooking.NewServer(config)
	if err != nil {
		logrus.Error(err)
		return
	}
	cmd := "serve"
	if parser.Active != nil {
		cmd = parser.Active.Name
	}
	switch cmd {
	case "serve":
		err = server.Serve()
		if err != nil {
			logrus.Errorf("got error after serve: %s", err)
		}
	case "add-user":
		user, err := fillUser(server, userOpts)
		if err != nil {
			return
		}
		server.Storage.AddUser(user)
	case "set-password":
		user, err := fillUser(server, userOpts)
		if err != nil {
			return
		}
		server.Storage.SetPassword(user)
	case "list-users":
		users, err := server.Storage.GetAllUsers()
		if err != nil {
			logrus.Errorf("Cannot get users: %s", err)
			return
		}
		for _, user := range users {
			fmt.Printf("%s\n", user.Username)
		}
	default:
		logrus.Errorf("unkown command %s", cmd)
	}
}

func fillUser(server *cooking.Server, userOpts userOpts) (cooking.User, error) {
	user := cooking.User{}
	if userOpts.Password == "" {
		fmt.Fprintf(os.Stderr, "password: ")
		pw, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			logrus.Errorf("error getting password: %s", err)
			return user, err
		}
		userOpts.Password = string(pw)
	}
	hash, err := server.HashPassword(userOpts.Password)
	if err != nil {
		logrus.Errorf("cannot hash password: %s", err)
		return user, err
	}
	user = cooking.User{
		Username: userOpts.Username,
		Password: string(hash),
	}
	return user, nil
}
