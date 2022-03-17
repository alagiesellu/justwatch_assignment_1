// Live API Assignment
// by JustWatch
// for Muhammed S. Baldeh github@alagiesellu
// Date: 17th March 2022

package main

import (
	"JustWatch/controllers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.GET("/movies", controllers.GetMoviesOfSpecies)

	// listen and serve on 0.0.0.0:8080
	err := r.Run()
	log.Fatal(err)
}
