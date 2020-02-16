package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/charts"
)

// https://github.com/go-echarts/go-echarts/blob/master/docs/docs/line.md

const (
	host   = "http://127.0.0.1:8081"
	maxNum = 50
)

var (
	nameItems = []string{"0", "1", "2", "3", "4", "5"}
)

var seed = rand.NewSource(time.Now().UnixNano())

func randInt() []int {
	cnt := len(nameItems)
	r := make([]int, 0)
	for i := 0; i < cnt; i++ {
		r = append(r, int(seed.Int63())%maxNum)
	}
	return r
}

func main() {
	http.HandleFunc("/", lineHandler)
	log.Printf("Please open URL: %v", host)
	http.ListenAndServe(":8081", nil)
}

func lineHandler(w http.ResponseWriter, _ *http.Request) {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.TitleOpts{Title: "Equity"})
	line.AddXAxis(nameItems).AddYAxis("Broker0", randInt())
	f, err := os.Create("./result.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	line.Render(f, w)
}
