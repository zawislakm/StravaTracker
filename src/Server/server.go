package Server

import (
	"app/src/Database"
	"app/src/StravaAPI"
	"app/src/Templates"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"sort"
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

	go GetActivities(serviceApi, serviceDb)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		//athletesData := serviceDb.GetAthletesData()
		//return Render(c, http.StatusOK, Templates.Table(athletesData))
		return Render(c, http.StatusOK, Templates.Index())
	})

	e.GET("/table/:sort", func(c echo.Context) error {
		sortField := c.Param("sort")
		athletesData := serviceDb.GetAthletesData()

		sort.Slice(athletesData, func(i, j int) bool {
			valI := reflect.ValueOf(athletesData[i]).FieldByName(sortField)
			valJ := reflect.ValueOf(athletesData[j]).FieldByName(sortField)
			if !valI.IsValid() || !valJ.IsValid() {
				return false
			}
			switch valI.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return valI.Int() > valJ.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return valI.Uint() > valJ.Uint()
			case reflect.Float32, reflect.Float64:
				return valI.Float() > valJ.Float()
			case reflect.String:
				return valI.String() > valJ.String()
			default:
				return false
			}
		})

		return Render(c, http.StatusOK, Templates.Table(athletesData))
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
