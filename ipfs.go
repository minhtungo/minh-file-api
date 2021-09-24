package cryp

import (
	"fmt"
	"strings"
    "os"

    shell "github.com/ipfs/go-ipfs-api"
)

var sh *shell.Shell
var url string = "localhost:8080"

func init() {
	sh = shell.NewShell(url)
}

func AddFileToIPFS(content string) string {
	cid, err := sh.Add(strings.NewReader(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("added %s\n", cid)

	err2 := sh.Pin(cid)
	if err2 != nil {
		fmt.Println("\nerr: ", err2)
	}
	return cid
}

func GetFileFromIPFS(cid string, outdir string) {
	err := sh.Get(cid, outdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
}
