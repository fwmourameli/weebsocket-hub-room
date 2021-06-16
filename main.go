package main

import (
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	go Hub.run()

	router := gin.New()
	router.LoadHTMLFiles("index.html")

	router.GET("/talk&listen/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		answer := make(chan string)
		m := message{
			waitMsg: true,
			data:    []byte("olar"),
			room:    roomId,
			answer:  answer,
		}
		Hub.broadcast <- m
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go readMsg(c, answer, wg)
		wg.Wait()
	})

	router.GET("/talk/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		m := message{
			waitMsg: false,
			data:    []byte("olar"),
			room:    roomId,
		}
		Hub.broadcast <- m
	})

	router.GET("/room/:roomId", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/ws/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		serveWs(c.Writer, c.Request, roomId)
	})

	router.Run("0.0.0.0:8080")
}

func readMsg(c *gin.Context, answer chan string, wg *sync.WaitGroup) {
	for {
		select {
		case str, ok := <-answer:
			if !ok {
				return
			}
			c.String(200, str)
			close(answer)
			wg.Done()
		}
	}
}
