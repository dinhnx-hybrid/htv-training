package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dinhnx-hybrid/htv-training/app/ferryserve/handler"
	"github.com/dinhnx-hybrid/htv-training/app/ferryserve/model"
)

const (
	waitingRoomCapacity uint = 7
	wakersCapacity      uint = 1
	maxPassengers       int  = 15
	rideDuration        int  = 10
)

func main() {
	var ferryMan = model.NewFerryMan(uint(1), "Ferry Man", rideDuration)
	var passengers []*model.Passenger
	for i := 0; i < maxPassengers; i++ {
		p := &model.Passenger{
			Id:         uint(i + 1),
			Name:       fmt.Sprintf("passenger_%d", (i + 1)),
			ComingTime: time.Duration(rand.Intn(6)+1) * time.Second,
		}
		passengers = append(passengers, p)
	}
	var handler = handler.NewFerryHandler(&ferryMan, waitingRoomCapacity, wakersCapacity, uint(maxPassengers))
	res, err := handler.Execute(passengers)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res)
}
