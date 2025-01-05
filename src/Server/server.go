package Server

import (
	"app/src/Database"
	"app/src/Models"
	"app/src/StravaAPI"
	"app/src/Templates"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func SetupServer() {
	serviceDb := Database.GetDbClient()
	serviceApi := StravaAPI.GetStravaClient()
	cache := newDataCache(serviceDb)
	go GetActivities(serviceApi, serviceDb, cache)

	e := echo.New()

	e.Static("/static", "static")
	e.GET("/", func(c echo.Context) error {
		tableLabels := []string{"Name", "Distance", "AverageTime", "AverageSpeed", "AverageLength", "LongestActivity", "ElevationGain", "TotalActivities"}
		return Render(c, http.StatusOK, Templates.Index(tableLabels))
	})

	e.GET("/table", func(c echo.Context) error {
		yearFilter := c.QueryParams().Get("year")
		if yearFilter == "" {
			yearFilter = cache.Year
		}
		sortField := c.QueryParams().Get("year")
		if sortField == "" {
			sortField = "Distance"
		}
		athletesData := cache.GetActivities(serviceDb, yearFilter)
		Models.SortAthletesData(athletesData, sortField)

		return Render(c, http.StatusOK, Templates.Table(athletesData))
	})

	e.GET("/years", func(c echo.Context) error {
		// TODO rebuild cache to store data years and add this for frontend
		years, err := serviceDb.GetUniqueYears()
		if err != nil {
			return err
		}
		return Render(c, http.StatusOK, Templates.Years(years))
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
