package main

cosnt (
  port = ":50052"
)

func main() {
  lis, err := net.Listen("tcp", port)
}