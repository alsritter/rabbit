package main

import (
	"fmt"

	"github.com/philchia/agollo/v4"
)

func main() {
	agollo.Start(&agollo.Conf{
		AppID:           "rabbit-apollo-test",
		Cluster:         "dev",
		NameSpaceNames:  []string{"application"},
		MetaAddr:        "http://localhost:28070/",
		AccesskeySecret: "2c92c614222b4f07b16b32e5b1c4a83a",
	})

	val := agollo.GetString("rabbit.mysql.password")
	fmt.Println("password: ", val)
}
