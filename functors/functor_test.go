package functors

import (
	"log"
	"testing"
)

func TestBasics(t *testing.T) {
	mp := Transform(func(data interface{}) interface{} {
		log.Println("data:", data)
		return data
	})

	mp.Call("sucks")
}
