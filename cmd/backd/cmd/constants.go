package cmd

// version
const version = "0.0.1"

// config
const (
	configCliNoColor           string = "cli.no_color"
	configCliQuiet             string = "cli.quiet"
	configCliSaveSession       string = "cli.login.save_session"
	configCliDontAskSession    string = "cli.login.dont_ask_session"
	configCliDontAskUserDomain string = "cli.login.dont_ask_user_domain"
	configCliDontAskPassword   string = "cli.login.dont_ask_password"
	configLoginUsername        string = "login.username"
	configLoginPassword        string = "login.password"
	configLoginDomain          string = "login.domain"
	configSessionID            string = "session.id"
	configSessionExpiresAt     string = "session.expires_at"
	configURLAuth              string = "url.auth"
	configURLObjects           string = "url.objects"
	configURLAdmin             string = "url.admin"
)

// Questions & answers
const (
	answerYes    = "Yes"
	answerNo     = "No"
	answerAlways = "Always"
	answerNever  = "Never"
	answerSave   = "Save"
	answerCancel = "Cancel"
)
