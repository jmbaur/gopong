package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	build := exec.Command("go", "build", "-o", "assets/main.wasm", "pong/main.go")
	build.Env = append(os.Environ(), []string{"GOOS=js", "GOARCH=wasm"}...)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir("assets")))
	log.Println("running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
