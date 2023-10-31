package main

import (
	"log"
	"newProject/cmd"
)

var Tag, CommitHash, CommitDate string

func main() {
	log.Printf("目前版本: %s, hash :%s , Update_Time: %s", Tag, CommitHash, CommitDate[:10])
	cmd.Run()

}
