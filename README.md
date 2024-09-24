# Strava Tracker

The Strava Tracker application is designed to monitor and track activities of a specific Strava Club.

## Technology 

- Golang 1.22 with [Echo](https://echo.labstack.com/), [templ](https://github.com/a-h/templ)
- [htmx](https://htmx.org/)
- MongoDB
- Docker
- Oracle Cloud -

Application fetch each data from [Strava API](https://developers.strava.com/docs/reference/) and store it in MongoDB. 
Data is displayed in the web application. To display the data, the application uses htmx to fetch the data from the server and 
templ to render HTML templates. Docker image of the application is hosted on Oracle Cloud VM.

Website overview (frontend still under construction):
[Strava Tracker](http://130.61.63.141) 
