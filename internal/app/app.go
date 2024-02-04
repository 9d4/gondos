package app

type Dependencies struct {
	UserStore UserStore
	ListStore ListStore
}

type App struct {
	d Dependencies
}

func New(d Dependencies) *App {
	app := &App{
		d: d,
	}

	return app
}
