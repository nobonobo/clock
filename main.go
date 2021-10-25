package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/zserge/lorca"
)

type counter struct {
	sync.Mutex
	count int
}

func (c *counter) Add(n int) {
	c.Lock()
	defer c.Unlock()
	c.count = c.count + n
}

func (c *counter) Value() int {
	c.Lock()
	defer c.Unlock()
	return c.count
}

func main() {
	log.SetFlags(log.Llongfile)
	ui, err := lorca.New("", "", 480, 320)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	c := &counter{}
	ui.Bind("counterAdd", c.Add)
	ui.Bind("counterValue", c.Value)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, nil)
	ui.Load(fmt.Sprintf("http://%s/", ln.Addr()))

	ui.Eval(`console.log("Hello, world!");`)

	<-ui.Done()
}
