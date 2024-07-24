package main

func main() {
	server := InitWebServer()
	for _, consumer := range server.consumers {
		err := consumer.Start()
		if err != nil {
			panic(err)
		}
	}
	server.engine.Run(":8081")
}