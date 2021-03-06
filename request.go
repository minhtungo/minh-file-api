package cryp

import (
	"fmt"

	"net/http"

	"github.com/labstack/echo/v4"

	"io/ioutil"
	"log"
	"os"
)

var keyString []byte = []byte("the-key-has-to-be-32-bytes-long!")

type Data struct {
	Content string `json:"content" query:"content" form:"content" param:"content"`
}

type H map[string]interface{}

// Handlers
func GetData(c echo.Context) error {
	cid := c.QueryParam("cid")
	outdir := fmt.Sprintf("%v", cid)
	GetFileFromIPFS(cid, outdir)
	data, err := ioutil.ReadFile(outdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	return c.JSON(http.StatusOK, H{
		"cid": cid,
		"data": string(data),
	})
}

func Hello(c echo.Context) error {
	return c.JSON(http.StatusOK, H{
		"message": "Helloo",
	})
}

func AddData(c echo.Context) error {
	data := new(Data)
	err := c.Bind(&data)
	if err != nil {
		log.Fatalf("Failed reading the request body %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	encryptedString, err := Encrypt(data.Content, keyString)
	if err != nil {
		log.Fatalf("Failed encrypting the data%s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}

	content := string(encryptedString)
	cid := AddFileToIPFS(content)

	//log.Printf("Added data", cid)
	return c.JSON(http.StatusOK, H{
		"cid": cid,
	})
}
