## Unreleased

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
- **makefile**: add test target
- **goreleaser**: add config
- **makefile**: depend on go.mod and go.sum for build

### ci
- **github-actions**: add codeql, go and release
