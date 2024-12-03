package main

func main() {

	s := InitWebServer()

	err := s.Run("127.0.0.1:8088")
	if err != nil {
		panic(err)
	}
}
