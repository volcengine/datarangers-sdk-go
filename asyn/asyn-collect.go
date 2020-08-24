package asyn

import (
	"code.byted.org/data/datarangers-sdk-go/pb_event"
	"encoding/json"
	"os"
	"strconv"
)

type executePool struct {
	max int
	tickets chan *ticket
	writer []*os.File
	//register *register
}

func NewExecutorPool(x int) (*executePool,error){
	p:=&executePool{}
	p.max = x
	p.tickets = make(chan *ticket, p.max)
	for i := 0; i < p.max; i++ {
		p.tickets <- &ticket{ id : i}
		logFile, err := os.OpenFile("./bak/test"+strconv.Itoa(i), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err!=nil { return nil,err}
		p.writer = append(p.writer, logFile)
	}
	return p,nil
}

func (this *executePool) execute(user *pb_event.User, header *pb_event.Header, event *pb_event.Event)error{

	t := <-this.tickets
	go func() error {
		iolog:= this.writer[t.id]
		message:=map[string]interface{}{
			"user": *user,
			"header":*header,
			"event": *event,
		}
		data, err:=json.Marshal(message)
		if err!=nil {return err}
		iolog.Write(data)
		iolog.WriteString("\n")
		this.tickets <- t
		return nil
	}()
	return nil
}

func (this *executePool) isComplete() bool{
	return len(this.tickets) == this.max
}


type ticket struct{
	id int
}

