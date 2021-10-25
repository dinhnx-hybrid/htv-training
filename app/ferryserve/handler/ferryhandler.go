package handler

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dinhnx-hybrid/htv-training/app/ferryserve/model"
)

type FerryHandlerInterface interface {
	Execute(passengers []*model.Passenger) (string, error)
}

type ferryHandler struct {
	FerryMan        *model.FerryMan
	WaitingRoom     chan *model.Passenger
	Wakers          chan *model.Passenger
	MaxPassengerNum uint
	WG              sync.WaitGroup
}

func NewFerryHandler(
	ferryMan *model.FerryMan,
	waitingRoomCapacity uint,
	wakersCapacity uint,
	maxPassengers uint,
) FerryHandlerInterface {
	return &ferryHandler{
		FerryMan:        ferryMan,
		WaitingRoom:     make(chan *model.Passenger, waitingRoomCapacity),
		Wakers:          make(chan *model.Passenger, wakersCapacity),
		MaxPassengerNum: maxPassengers,
	}
}

func (h *ferryHandler) Execute(passengers []*model.Passenger) (string, error) {
	if h.FerryMan == nil {
		return "Serving failed", errors.New("no ferry man")
	}
	go h.StartFerryManChecking()

	h.WG.Add(int(h.MaxPassengerNum))
	fmt.Printf("[%s] has started his job.\n", h.FerryMan.Name)
	for _, p := range passengers {
		time.Sleep(p.ComingTime)
		go h.StartPassengerWelcome(p)
	}
	h.WG.Wait()
	fmt.Printf("[%s] has end his job.\n", h.FerryMan.Name)
	return "Serving done", nil
}

func (h *ferryHandler) StartPassengerWelcome(p *model.Passenger) {
	fmt.Printf("Passenger [%s] came after %f second(s)\n", p.Name, p.ComingTime.Seconds())
	switch h.FerryMan.Status {
	case model.AVAILABLE, model.CHECKING:
		h.Wakers <- p
	case model.SLEEPING:
		select {
		case h.Wakers <- p:
		default:
			select {
			case h.WaitingRoom <- p:
				fmt.Printf("Passenger [%s] is going to waiting room\n", p.Name)
			default:
				fmt.Printf("Passenger [%s] is leaving\n", p.Name)
				h.WG.Done()
			}
		}
	case model.RIDING:
		select {
		case h.WaitingRoom <- p:
			fmt.Printf("Passenger [%s] is going to waiting room\n", p.Name)
		default:
			fmt.Printf("Passenger [%s] is leaving\n", p.Name)
			h.WG.Done()
		}
	}
}

func (h *ferryHandler) StartFerryManChecking() {
	for {
		h.FerryMan.Lock()
		h.FerryMan.StartChecking()
		fmt.Printf("There are %d passenger(s) in the waiting room.\n", len(h.WaitingRoom))
		select {
		case p := <-h.WaitingRoom:
			h.FerryMan.StartRiding(p)
			h.FerryMan.Unlock()
			h.WG.Done()
		default:
			h.FerryMan.StartSleeping()
			h.FerryMan.Unlock()
			p := <-h.Wakers
			h.FerryMan.Lock()
			h.FerryMan.Awake(p)
			h.FerryMan.StartRiding(p)
			h.FerryMan.Unlock()
			h.WG.Done()
		}
	}

}
