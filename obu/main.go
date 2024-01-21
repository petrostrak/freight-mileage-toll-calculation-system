// The package Obu (on-board unit) simulates coordinates transmission.
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var sendInterval = time.Second

type OBUData struct {
	OBUID int     `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}

	return ids
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	obuids := generateOBUIDS(20)
	for {
		for i := 0; i < len(obuids); i++ {
			lat, long := genLatLong()
			data := OBUData{
				OBUID: obuids[i],
				Lat:   lat,
				Long:  long,
			}
			fmt.Printf("%+v\n", data)
		}
		time.Sleep(sendInterval)
	}
}
