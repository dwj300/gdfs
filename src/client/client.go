package client

import (
    "flag";
    "fmt"
)

func main() {
    var ip = flag.String("IP Address", "127.0.0.1", "IP Address of the server")
    fmt.Printf("Hello, World.\n")
}
