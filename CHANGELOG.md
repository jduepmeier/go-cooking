## 0.1.12 (2025-09-03)

### Fix

- **deps**: update module github.com/mattn/go-sqlite3 to v1.14.32
- **deps**: update module golang.org/x/crypto to v0.41.0

## 0.1.11 (2025-07-06)

### Fix

- **deps**: update module golang.org/x/crypto to v0.39.0
- **deps**: update module github.com/mattn/go-sqlite3 to v1.14.28

## 0.1.10 (2025-04-09)

### Fix

- **deps**: update module github.com/fsnotify/fsnotify to v1.9.0
- **deps**: update module github.com/mattn/go-sqlite3 to v1.14.27
- **deps**: update module golang.org/x/crypto to v0.37.0

## 0.1.9 (2025-02-09)

### Fix

- **deps**: update module golang.org/x/crypto to v0.33.0

## 0.1.8 (2025-01-12)

### Fix

- **deps**: update module github.com/gorilla/sessions to v1.4.0
- **deps**: update module golang.org/x/crypto to v0.32.0
- **dep**: update golang => 1.23
- **deps**: update module github.com/mattn/go-sqlite3 to v1.14.24
- **deps**: update module golang.org/x/crypto to v0.31.0
- **deps**: update module github.com/fsnotify/fsnotify to v1.8.0

## 0.1.7 (2024-06-29)

### Fix

- **deps**: update module golang.org/x/crypto to v0.24.0
- **deps**: update module github.com/gorilla/sessions to v1.3.0
- **deps**: update module github.com/jessevdk/go-flags to v1.6.1
- **deps**: update module golang.org/x/crypto to v0.23.0

## 0.1.6 (2024-04-03)

### Fix

- **deps**: update module golang.org/x/crypto to v0.21.0
- **deps**: update module golang.org/x/crypto to v0.20.0
- **deps**: update module github.com/mattn/go-sqlite3 to v1.14.22
- **deps**: update module golang.org/x/crypto to v0.19.0
- **deps**: update module github.com/mattn/go-sqlite3 to v1.14.20
- **deps**: update module golang.org/x/crypto to v0.18.0

## 0.1.5 (2024-01-07)

### Fix

- **deps**: update module github.com/mattn/go-sqlite3 to v1.14.19
- **deps**: update module golang.org/x/crypto to v0.17.0
- **deps**: update module golang.org/x/crypto to v0.16.0

## 0.1.4 (2023-11-10)

### Fix

- **deps**: update module github.com/gorilla/mux to v1.8.1
- **deps**: update module github.com/gorilla/sessions to v1.2.2
- **deps**: update module golang.org/x/crypto to v0.15.0
- **deps**: update module github.com/gorilla/securecookie to v1.1.2
- **deps**: update module github.com/mattn/go-sqlite3 to v1.14.18
- **deps**: update module github.com/fsnotify/fsnotify to v1.7.0
- **deps**: update module golang.org/x/crypto to v0.14.0

## 0.1.3 (2023-09-10)

### Fix

- **deps**: update module golang.org/x/crypto to v0.13.0

## 0.1.2 (2023-08-27)

### Fix

- **deps**: update module golang.org/x/crypto to v0.12.0
- **deps**: update module golang.org/x/crypto to v0.11.0
- **deps**: update module golang.org/x/crypto to v0.10.0
- **deps**: update module github.com/mattn/go-sqlite3 to v1.14.17
- **deps**: update module github.com/sirupsen/logrus to v1.9.3

## 0.1.1 (2023-05-18)

### Fix

- **deps**: update module golang.org/x/crypto to v0.9.0
- **deps**: update module github.com/sirupsen/logrus to v1.9.2
- **deps**: update module golang.org/x/crypto to v0.8.0
- **deps**: update module golang.org/x/crypto to v0.7.0

## 0.1.0 (2023-02-25)

### Feat

- **go.mod**: rename module to github path
- update dependencies
- add collapsed navbar for small screens
- add search and css fixes
- add description field
- add favicon
- add cooker image
- persist session key in database
- add deletion of recipes
- add buttons to show/hide cards on print
- initial commit

### Fix

- **deps**: update transient dependencies pretty and check.v1
- **handler_static**: sanitize filename even more
- **database**: use smaller max int size for ParseInt
- **handler_static**: add more checks to filename
- **logrus**: fix call has possible formatting directive
- **deps**: update module github.com/sirupsen/logrus to v1.9.0
- **deps**: update module golang.org/x/crypto to v0.6.0
- **deps**: update module gopkg.in/yaml.v2 to v3
- **deps**: update module github.com/fsnotify/fsnotify to v1.6.0
- **deps**: update module github.com/jessevdk/go-flags to v1.5.0
- fix typo
- source shouldn't be a required field

### Docs

- add LICENSE
- add CHANGELOG

### Build

- **go.sum**: sync with go.mod
- **makefile**: add test target
- **goreleaser**: add config
- **makefile**: depend on go.mod and go.sum for build

### ci

- **github-actions**: add codeql, go and release
