package main

import (
	"fmt"
	"router"
//	"encoding/json"
//	"io/ioutil"
//	"os"
//	"util"
)

type (
)


func main() {
	fmt.Println("Start Server")
	e := router.Initialize()
    e.Start(":8080")
}