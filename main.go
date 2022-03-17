// Live API Assignment
// by JustWatch
// for Muhammed S. Baldeh github@alagiesellu
// Date: 17th March 2022

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
)

const GhibliApi string = "https://ghibliapi.herokuapp.com"

type film struct {
	Data []interface{}
}

func main() {
	r := gin.Default()
	r.GET("/movies", func(c *gin.Context) {

		films := GetFilmsOfSpecies(c.Query("species"))

		c.JSON(200, films)
	})

	// listen and serve on 0.0.0.0:8080
	err := r.Run()
	CatchError(err)
}

// GetFilmsOfSpecies returns list of films of the speciesId provided.
func GetFilmsOfSpecies(speciesId string) []interface{} {

	// MakeRequest returns species records
	seriesResponse := MakeRequest(GhibliApi + "/species/" + speciesId)

	// store films from species records
	films := seriesResponse["films"].([]interface{})

	filmRecords := film{}

	for i := range films {

		// MakeRequest of each film
		filmResponse := MakeRequest(films[i].(string))

		// append film records filmRecords
		filmRecords.Data = append(filmRecords.Data, filmResponse)
	}

	return filmRecords.Data
}

// MakeRequest to API server and return map[string]interface{}
func MakeRequest(link string) map[string]interface{} {

	client := &http.Client{}

	request, err := http.NewRequest("GET", link, nil)
	CatchError(err)

	response, err := client.Do(request)
	CatchError(err)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		CatchError(err)
	}(response.Body)

	responseBody, err := ioutil.ReadAll(response.Body)
	CatchError(err)

	jsonResponse := make(map[string]interface{})

	err = json.Unmarshal(responseBody, &jsonResponse)
	CatchError(err)

	return jsonResponse
}

// CatchError and print it, if err is not nil
func CatchError(err error) {
	if err != nil {
		fmt.Println("Error => ", err)
	}
}
