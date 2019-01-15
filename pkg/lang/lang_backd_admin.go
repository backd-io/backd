package lang

import (
	"github.com/backd-io/backd/backd"
)

func (l *Lang) addBackdAdminCommands() {

	l.AddCommand(
		"app_id",
		"Sets the application ID of the application we want to operate to.",
		`Sets the application ID of the application we want to operate to.`,
		`app_id("idApplication")`,
		l.appID)

	l.AddCommand(
		"me",
		"Returns the information of the current user.",
		`Returns the information of the current user logged in.`,
		``,
		l.me)

	l.AddCommand(
		"app_get",
		"Gets the application object by ID.",
		`Gets the application object by ID`,
		`
app, err = app_get("idApplication")

println(app.name)
// returns "thisIsASample"
`,
		l.appGet)

	l.AddCommand(
		"app_create",
		"Creates an application",
		`Creates an application. Returns application ID and error if any.`,
		`app = {}
app.name = "applicationExample"
app.description = "this is an example"

id, err = app_create(app)
`,
		l.appCreate)

}

// appID - app_id
func (l *Lang) appID(id string) {
	l.currentAppID = id
}

// SetSessionID - set_session_id
func (l *Lang) SetSessionID(sessionID string) {
	l.b.SetSessionID(sessionID)
}

func (l *Lang) me() (user backd.User, err error) {
	return l.b.Me()
}

// appCreate - app_create
func (l *Lang) appCreate(data map[string]interface{}) (map[string]interface{}, error) {
	return l.b.Apps().Insert(data)
}

// appGet - app_get
func (l *Lang) appGet(id string) (data map[string]interface{}, err error) {
	err = l.b.Apps().GetByID(id, &data)
	return data, err
}
