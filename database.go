package cooking

import (
	"cooking/database"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html"
	"image/color"
	"net/url"
	"sort"
	"strconv"
	"strings"

	// import sqlite3 sql driver. Must be blank import for loading the init register function.
	"github.com/gorilla/securecookie"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/skip2/go-qrcode"
)

const (
	// Fresh if the reciped needs fresh ingredients.
	Fresh Freshness = 1
	// MiddleFresh if the reciped needs middle fresh ingredients.
	MiddleFresh Freshness = 2
	// NotFresh if no fresh ingredients are needed.
	NotFresh Freshness = 3
	// MaxFreshness is the currently highest value for freshness.
	MaxFreshness int64 = int64(NotFresh)
)

// Freshness defines the freshness.
type Freshness int

// HTML return the freshness value as html formatted scala.
func (f Freshness) HTML() string {
	var builder strings.Builder

	for i := 0; i < int(MaxFreshness); i++ {
		if i < int(f) {
			builder.WriteString("&#9679;")
		} else {
			builder.WriteString("&#9675;")
		}
	}
	return builder.String()
}

const (
	// StatementAdd adds a recipe to database.
	StatementAdd = "INSERT INTO recipes(name, length, freshness, source) VALUES(:name, :length, :freshness, :source)"
	// StatementUpdate updates a recipe.
	StatementUpdate = "UPDATE recipes SET name = :name, length = :length, freshness = :freshness, source = :source WHERE id = :id"
	// StatementDelete deletes a recipe from database.
	StatementDelete = "DELETE FROM recipes WHERE id = :id"
	// StatementGetAll returns all recipes from database.
	StatementGetAll = "SELECT id, name, length, freshness, source FROM recipes"
	// StatementGet returns the recipe with the given id from database.
	StatementGet = "SELECT id, name, length, freshness, source FROM recipes WHERE id = :id"
	// StatementGetUser returns the user with the given username from database.
	StatementGetUser = "SELECT username, password FROM user WHERE username = :username"
	// StatementAddUser adds a user to the database.
	StatementAddUser = "INSERT INTO user (username, password) VALUES(:username, :password)"
	// StatementSetPassword sets the password for a user.
	StatementSetPassword = "UPDATE user SET password = :password WHERE username = :username"
	// StatementGetAllUsers returns all users.
	StatementGetAllUsers = "SELECT username, password FROM user"
)

// Storage handles access to the storage.
type Storage struct {
	Path       string
	db         *sql.DB
	Statements map[string]*sql.Stmt
}

// Recipe contains all infos about a recipe.
type Recipe struct {
	ID        uint64
	Name      string
	Length    string
	Freshness Freshness
	Source    string
}

// SourceHTML returns the source field as html code.
// If it is a url it will be converted to qrcode.
// Otherwise it will be returned as text.
func (recipe *Recipe) SourceHTML() string {
	_, err := url.ParseRequestURI(recipe.Source)
	if err != nil {
		return fmt.Sprintf("<p>%s</p>", recipe.Source)
	}
	size := 256
	htmlSource := html.EscapeString(recipe.Source)
	png, err := qrcode.New(recipe.Source, qrcode.Highest)
	if err != nil {
		logrus.Errorf("could not encode qr code for %s", recipe.Source)
		return fmt.Sprintf(`<a href=%q>%q<a />`, htmlSource, htmlSource)
	}
	png.BackgroundColor = color.Transparent
	pngEncoded, err := png.PNG(size)
	if err != nil {
		logrus.Errorf("could not encode qr code for %s", recipe.Source)
		return fmt.Sprintf(`<a href=%q>%q<a />`, htmlSource, htmlSource)
	}

	png64 := base64.StdEncoding.EncodeToString(pngEncoded)
	return fmt.Sprintf(`<a href=%q><img width="%d" height="%d" alt=%q src="data:image/png;base64,%s" /></a>`, htmlSource, size, size, htmlSource, png64)
}

// Validate validates the struct for missing or broken values.
func (recipe *Recipe) Validate() error {
	if recipe.Name == "" {
		return fmt.Errorf("missing key: name")
	}
	if recipe.Length == "" {
		return fmt.Errorf("missing key: length")
	}

	return nil
}

// ParseFreshness parses the freshness value and adds it to the recipe.
func (recipe *Recipe) ParseFreshness(freshness string) error {
	value, err := strconv.ParseInt(freshness, 10, 64)
	if err != nil {
		return err
	}
	if value > MaxFreshness && value < 1 {
		return fmt.Errorf("freshness must be between 1 and %d", MaxFreshness)
	}

	recipe.Freshness = Freshness(value)

	return nil
}

// Connect connects to the database.
func (storage *Storage) Connect() (err error) {
	if storage.db == nil {
		if storage.Path == "" {
			storage.Path = ":memory:"
		}
		storage.db, err = sql.Open("sqlite3", storage.Path)
		if err != nil {
			return err
		}

		err = storage.checkDB()
		if err != nil {
			return err
		}
		err = storage.prepareStatements()
		if err != nil {
			return err
		}
	}

	return nil
}

func (storage *Storage) prepareStatements() (err error) {
	statements := map[string]string{
		"get":         StatementGet,
		"getall":      StatementGetAll,
		"update":      StatementUpdate,
		"add":         StatementAdd,
		"delete":      StatementDelete,
		"setpassword": StatementSetPassword,
		"adduser":     StatementAddUser,
		"getuser":     StatementGetUser,
		"getallusers": StatementGetAllUsers,
	}
	storage.Statements = make(map[string]*sql.Stmt, len(statements))
	for identifier, statement := range statements {
		stmt, err := storage.db.Prepare(statement)
		if err != nil {
			return err
		}
		storage.Statements[identifier] = stmt
	}

	return nil
}

// Close closes the database connection.
func (storage *Storage) Close() {
	if storage.db != nil {
		storage.db.Close()
	}
}

func (storage *Storage) checkDB() (err error) {

	row := storage.db.QueryRow("SELECT value FROM meta WHERE key = 'version'")
	var currVersionString string
	var currVersion int64
	var maxVersion int64
	err = row.Scan(&currVersionString)
	logrus.Infof("error is %s and currVersionString=%s", err, currVersionString)
	if err == nil {
		currVersion, err = strconv.ParseInt(currVersionString, 10, 64)
		if err != nil {
			currVersion = -1
		}
	} else if strings.Contains(err.Error(), "no such table") {
		currVersion = -1
	} else {
		return err
	}

	maxVersion = currVersion
	logrus.Infof("current version is %d", currVersion)

	sort.Sort(database.ByVersion(database.Versions))
	for _, version := range database.Versions {
		if currVersion < int64(version.Version) {
			logrus.Infof("install version %d", version.Version)
			_, err := storage.db.Exec(version.Statement)
			if err != nil {
				return fmt.Errorf("cannot install version %d: %s", version.Version, err)
			}
			maxVersion = int64(version.Version)
		}
	}

	_, err = storage.db.Exec("UPDATE meta SET value = :version WHERE key = 'version';", sql.Named("version", &maxVersion))
	return err
}

// GetRecipes returns all saved recipes.
func (storage *Storage) GetRecipes() ([]Recipe, error) {
	var recipes []Recipe
	rows, err := storage.Statements["getall"].Query()
	if err != nil {
		return recipes, err
	}

	for rows.Next() {
		recipe := Recipe{}
		err = rows.Scan(&recipe.ID, &recipe.Name, &recipe.Length, &recipe.Freshness, &recipe.Source)
		if err != nil {
			return recipes, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, rows.Err()
}

// Add adds a recipe to storage.
func (storage *Storage) Add(recipe Recipe) error {
	logrus.Infof("add %v", recipe)
	_, err := storage.Statements["add"].Exec(
		sql.Named("name", &recipe.Name),
		sql.Named("length", &recipe.Length),
		sql.Named("freshness", &recipe.Freshness),
		sql.Named("source", &recipe.Source),
	)

	return err
}

// Update updates a recipe.
func (storage *Storage) Update(recipe Recipe) error {
	if recipe.ID == 0 {
		return fmt.Errorf("id must be defined")
	}
	logrus.Infof("update %v", recipe)

	_, err := storage.Statements["update"].Exec(
		sql.Named("id", &recipe.ID),
		sql.Named("name", &recipe.Name),
		sql.Named("length", &recipe.Length),
		sql.Named("freshness", &recipe.Freshness),
		sql.Named("source", &recipe.Source),
	)
	return err
}

// Get returns the recipe with the given id.
func (storage *Storage) Get(id uint64) (Recipe, error) {
	row := storage.Statements["get"].QueryRow(sql.Named("id", id))
	recipe := Recipe{}
	err := row.Scan(
		&recipe.ID,
		&recipe.Name,
		&recipe.Length,
		&recipe.Freshness,
		&recipe.Source,
	)
	return recipe, err
}

// Delete deletes a recipe.
func (storage *Storage) Delete(id uint64) error {
	_, err := storage.Statements["delete"].Exec(sql.Named("id", &id))
	return err
}

// GetUser returns the user with the given username.
func (storage *Storage) GetUser(username string) (User, error) {
	row := storage.Statements["getuser"].QueryRow(sql.Named("username", username))
	user := User{}
	err := row.Scan(
		&user.Username,
		&user.Password,
	)
	return user, err
}

// AddUser adds a user to the database.
func (storage *Storage) AddUser(user User) error {
	_, err := storage.Statements["adduser"].Exec(
		sql.Named("username", user.Username),
		sql.Named("password", user.Password),
	)
	return err
}

// SetPassword sets the password for a user.
func (storage *Storage) SetPassword(user User) error {
	_, err := storage.Statements["setpassword"].Exec(
		sql.Named("username", user.Username),
		sql.Named("password", user.Password),
	)
	return err
}

// GetAllUsers returns all users.
func (storage *Storage) GetAllUsers() ([]User, error) {
	var users []User
	rows, err := storage.Statements["getallusers"].Query()
	if err != nil {
		return users, err
	}

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Username, &user.Password)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

// GetOrCreateSessionCookieKey returns the session cookie key from database or creates it.
func (storage *Storage) GetOrCreateSessionCookieKey() ([]byte, error) {
	row := storage.db.QueryRow("SELECT value FROM meta WHERE key = 'sessionkey'")
	var base64Key string
	err := row.Scan(&base64Key)
	if err == nil {
		key, err := base64.StdEncoding.DecodeString(base64Key)
		if err == nil {
			return key, err
		}
	}
	logrus.Warn("cannot get session key form db: %s", err)

	key := securecookie.GenerateRandomKey(64)
	base64Key = base64.StdEncoding.EncodeToString(key)
	_, err = storage.db.Exec("INSERT INTO meta (key, value) VALUES ('sessionkey', :key)", sql.Named("key", base64Key))
	return key, err
}
