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
	"github.com/golang/protobuf/proto"
	"qnetwork.net/fbcrawl/fbcolly"
	"unsafe"
)

const logPath = "parse.log"

var verbose = flag.Bool("verbose", true, "print info level logs to stdout")
var email = flag.String("email", "change_me@gmail.com", "facebook email")
var password = flag.String("password", "change_me", "facebook password")
var otp = flag.String("otp", "123456", "facebook otp")
var groupId = flag.String("groupId", "334294967318328", "facebook group id, default is 334294967318328")

var allInstances = map[uintptr]*fbcolly.Fbcolly{}

//export Init
func Init() uintptr {
	instance := fbcolly.New()
	ptr := (uintptr)(unsafe.Pointer(instance))
	allInstances[ptr] = instance
	return ptr
}

//export FreeColly
func FreeColly(pointer unsafe.Pointer) {
	delete(allInstances, uintptr(pointer))
}

//export Login
func Login(pointer unsafe.Pointer, email *C.char, password *C.char) {
	p := (*fbcolly.Fbcolly)(pointer)
	//print(p.E)
	p.Login(C.GoString(email), C.GoString(password), "")
}

//export FetchGroupFeed
func FetchGroupFeed(pointer unsafe.Pointer, groupId *C.char) unsafe.Pointer {
	p := (*fbcolly.Fbcolly)(pointer)
	_, postsList := p.FetchGroupFeed(C.GoString(groupId))
	marshaledPostsList, _ := proto.Marshal(postsList)
	return C.CBytes(append(marshaledPostsList, 0))
}

//export FetchPost
func FetchPost(pointer unsafe.Pointer, groupId *C.char, postId *C.char) unsafe.Pointer {
	p := (*fbcolly.Fbcolly)(pointer)
	_, post := p.FetchPost(C.GoString(groupId), C.GoString(postId))
	marshaledPost, _ := proto.Marshal(post)
	return C.CBytes(append(marshaledPost, 0))
}

func main() {
	//r := regexp.MustCompile("/([\\d\\w.]+)").FindStringSubmatch()[1]
	//print(r.FindStringSubmatch("/liem.phamthanh.161?refid=18&__tn__=R")[1])
	//flag.Parse()
	//
	////post := fbcrawl.FacebookPost{}
	////bPost, err := proto.Marshal(&post)
	//
	//lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	//if err != nil {
	//	logger.Fatalf("Failed to open log file: %v", err)
	//}
	//defer lf.Close()
	//defer logger.Init("fb-colly", *verbose, false, lf).Close()
	//f := fbcolly.New()
	//err = f.Login(*email, *password, *otp)
	//f.FetchGroupFeed(*groupId)
}
