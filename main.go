package main

import (
	"fmt"

	"net/http"

	"github.com/labstack/echo/v4"

	// "github.com/labstack/echo/v4/middleware"

	"encoding/json"
	"io/ioutil"
	"log"

	"crypto/rand"
	"encoding/hex"

	encryption "github.com/minhtungo/minh-file-api/encryption"
	ipfs "github.com/minhtungo/minh-file-api/ipfs-api"
)

var keyString string

type Data struct {
	content string `json:"content"`
}

func main() {
	// create a new echo instance
	e := echo.New()

	// Routes
	e.GET("/", hello)
	e.GET("/add/:cid", getData)

	e.POST("/add", addData)

	// generate a random 32 byte key for AES-256
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}
	keyString = hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault

	// Start server
	e.Logger.Fatal(e.Start(":8000"))

}

func getData(c echo.Context) error {
	cid := c.Param("cid")
	outdir := fmt.Sprintf("%v", cid)
	ipfs.GetFileFromIPFS(cid, outdir)
	data, err := ioutil.ReadFile(outdir)
	if err != nil {
		panic(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	return c.String(http.StatusOK, fmt.Sprintf("Data: %v", data))
}

// Handlers
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, Taubyte!")
}

func addData(c echo.Context) error {
	data := Data{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		log.Fatalf("Failed reading the request body %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	dataString := fmt.Sprintf("%v", data)

	encryptedString := encryption.Encrypt(dataString, keyString)
	cid := ipfs.AddFileToIPFS(encryptedString)

	log.Printf("Added data", cid)
	return c.String(http.StatusOK, fmt.Sprintf("Data added: %v", data))
}
