package model

import (
	"fmt"
	"sync"
	"time"
)

type FerryManStatus int

const (
	AVAILABLE FerryManStatus = iota
	SLEEPING
	RIDING
	CHECKING
)

type FerryMan struct {
	sync.Mutex
	Id           uint
	Name         string
	Passenger    *Passenger
	Status       FerryManStatus
	RideDuration time.Duration
}

func NewFerryMan(id uint, name string, rideDuration int) FerryMan {
	return FerryMan{
		Id:           id,
		Name:         name,
		Status:       AVAILABLE,
		RideDuration: time.Duration(rideDuration) * time.Second,
	}
}

func (f *FerryMan) Awake(p *Passenger) {
	fmt.Printf("[%s] awake by passenger [%s].\n", f.Name, p.Name)
	f.Status = AVAILABLE
}

func (f *FerryMan) StartRiding(p *Passenger) {
	fmt.Printf("[%s] start riding, passenger [%s].\n", f.Name, p.Name)
	f.Status = RIDING
	f.Passenger = p
	time.Sleep(f.RideDuration)
	fmt.Printf("[%s] end riding, passenger [%s].\n", f.Name, p.Name)
}

func (f *FerryMan) StartChecking() {
	fmt.Printf("[%s] is checking the waiting room.\n", f.Name)
	f.Status = CHECKING
}

func (f *FerryMan) StartSleeping() {
	fmt.Printf("[%s] is going to sleep\n", f.Name)
	f.Status = SLEEPING
	f.Passenger = nil
}
