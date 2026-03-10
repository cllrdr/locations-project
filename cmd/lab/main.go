package main

import (
    "fmt"
    "log"
    "locations-project/internal/api"
)

func main() {
    fmt.Println("Привет, Go!")
    log.Println("App Start!")
    api.StartServer()
    log.Println("App Terminated!")
}
