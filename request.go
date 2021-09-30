package cryp

import (
	//"fmt"

	"net/http"

	"github.com/labstack/echo/v4"

	// "github.com/labstack/echo/v4/middleware"

	"encoding/json"
	//"io/ioutil"
	"log"
)

var keyString []byte = []byte("the-key-has-to-be-32-bytes-long!")

type Data struct {
	Content string `json:"content"`
}

type H map[string]interface{}

// Handlers
func GetData(c echo.Context) error {
	cid := c.QueryParam("cid")
	// outdir := fmt.Sprintf("%v", cid)
	// GetFileFromIPFS(cid, outdir)
	// data, err := ioutil.ReadFile(outdir)
	// if err != nil {
	// 	panic(err.Error())
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	// }
	return c.JSON(http.StatusOK, H{
		"cid": cid,
	})
}

func Hello(c echo.Context) error {
	return c.JSON(http.StatusOK, H{
		"message": "Hello",
	})
}

func AddData(c echo.Context) error {
	var data Data
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		log.Fatalf("Failed reading the request body %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}

	encryptedString, err := Encrypt(data.Content, keyString)
	if err != nil {
		log.Fatalf("Failed encrypting the data%s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	cid := AddFileToIPFS(string(encryptedString))

	log.Printf("Added data", cid)
	return c.JSON(http.StatusOK, H{
		"cid": cid,
	})
}
