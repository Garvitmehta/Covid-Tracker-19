package routes

import (
	"echo-mongo-api/controllers" //add this

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {

	e.GET("/insert/:stateId", controllers.InsertTable) //add this
	e.GET("/state/:stateId", controllers.GetAState)    //add this
	// e.PUT("/update/:stateId", controllers.UpdateState) //add this
}
