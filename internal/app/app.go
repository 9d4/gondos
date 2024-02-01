package app

type Dependencies struct {
	UserStore UserStore
}

type App struct {
	Dependencies
}

func New(d Dependencies) *App {
	app := &App{
		Dependencies: d,
	}

	return app
}
