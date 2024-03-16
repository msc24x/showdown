package main

import (
	"flag"
	"fmt"
)

func printHeadline() {
	fmt.Printf("Showdown! https://github.com/msc24x/showdown\nPortable server to execute and judge code.\n")
}

func printHelp() {
	printHeadline()
	fmt.Println("\nUsage:\n\t./showdown [...options] -start")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}
