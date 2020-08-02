package main

/*
#include <stdio.h>
#include <stdlib.h>

static void myprint(char* s) {
  printf("%s\n", s);
}


*/
import "C"

import (
	"flag"
	"github.com/google/logger"
	"os"
	"qnetwork.net/fbcrawl/fbcolly"
)

const logPath = "parse.log"

var verbose = flag.Bool("verbose", true, "print info level logs to stdout")
var email = flag.String("email", "change_me@gmail.com", "facebook email")
var password = flag.String("password", "change_me", "facebook password")
var otp = flag.String("otp", "123456", "facebook otp")
var groupId = flag.String("groupId", "334294967318328", "facebook group id, default is 334294967318328")

func main() {
	flag.Parse()

	//post := fbcrawl.FacebookPost{}
	//bPost, err := proto.Marshal(&post)

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()
	defer logger.Init("fb-colly", *verbose, false, lf).Close()
	f := fbcolly.New()
	f.Login(*email, *password, *otp)
	f.FetchGroupFeed(*groupId)
}
