package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
)

const GhibliApi string = "https://ghibliapi.herokuapp.com"

type ResponseRecord struct {
	Data []interface{}
}

func GetMoviesOfSpecies(c *gin.Context) {

	speciesQuery := c.Query("species")

	// if species ID not provided, response BadResponse
	if speciesQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Species ID not provided."})
		return
	}

	// MakeRequest returns species records
	seriesResponse := MakeRequest(GhibliApi+"/species/"+speciesQuery, c)

	// store films from species records
	films := seriesResponse["films"].([]interface{})

	filmRecords := ResponseRecord{}

	for i := range films {

		// MakeRequest of each film
		filmResponse := MakeRequest(films[i].(string), c)

		// append film records filmRecords
		filmRecords.Data = append(filmRecords.Data, filmResponse)
	}

	// response successful request
	c.JSON(http.StatusOK, filmRecords.Data)

	return
}

// MakeRequest to API server and return json map
func MakeRequest(link string, c *gin.Context) map[string]interface{} {

	client := &http.Client{}

	request, err := http.NewRequest("GET", link, nil)
	CatchError(err, c)

	response, err := client.Do(request)
	CatchError(err, c)

	// if response is not http.StatusOK, throw error.
	if response.StatusCode != http.StatusOK {
		err := errors.New(response.Status)
		CatchError(err, c)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		CatchError(err, c)
	}(response.Body)

	responseBody, err := ioutil.ReadAll(response.Body)
	CatchError(err, c)

	jsonResponse := make(map[string]interface{})

	err = json.Unmarshal(responseBody, &jsonResponse)
	CatchError(err, c)

	return jsonResponse
}

// CatchError and print it, if err is not nil
func CatchError(err error, c *gin.Context) {

	// if there is error, response error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
