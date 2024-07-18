package main

import "log"

func main() {
	//embeds main as an encrypted resource
	log.Println(string(cats))
}

//go:generate go run github.com/dbeleon/encembed -i main.go -decvarname cats -srcname encembed.go -scramble
