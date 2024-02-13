package main

import "github.com/0xAckerMan/Savanah/internal/data"

func (app *Application) migrations() {
	app.DB.AutoMigrate(
		&data.Product{},
	)
}