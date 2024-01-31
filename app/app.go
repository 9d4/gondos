package app

import "database/sql"

type Config struct {
	DB *sql.DB
}

func (c *Config) apply(app *App) {
	app.db = c.DB
}

type App struct {
	db *sql.DB
}

func New(c *Config) *App {
	app := &App{}
	c.apply(app)
	return app
}
