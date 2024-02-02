package app

type Dependencies struct {
	UserStore UserStore
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
