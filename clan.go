package roach

import (
	"log"
	"time"
)

type Clan struct {
	startPoint *Link
	log        Logger
	config     *ClanConfig
	workerChan chan *Link
}

func NewClan(config *ClanConfig) *Clan {
	link := newLink(config.StartPoint)
	config.LinkConfig.Target = config.StartPoint
	config.LinkConfig.Log = config.Log
	config.LinkConfig.Host = config.Host
	config.LinkConfig.header = config.Header
	link.setConfig(config.LinkConfig)
	return &Clan{
		startPoint: link,
		config:     config,
		workerChan: make(chan *Link, config.ConcurrentNumber),
	}
}

func (clan *Clan) Rush() {

	clan.producer()
	clan.customer()
	<-exitChan
}

func (clan *Clan) customer() {
	go func() {
		for {
			select {
			case linkStr := <-LinkStrChan:
				if _, ok := PureStack[linkStr]; ok == true {
					continue
				} else {
					PureStack[linkStr] = newLink(linkStr)
					PureStack[linkStr].setConfig(clan.config.LinkConfig)

					clan.workerChan <- PureStack[linkStr]
				}

			case <-time.After(10 * time.Second):
				log.Println("break")
				exitChan <- true
				break
			}
		}
	}()

	go func() {
		for {
			i := <-clan.workerChan
			go func(*Link) {
				err := i.GetPageUrls()
				if err != nil {
					log.Println(err)
				}
			}(i)
			<-time.After(clan.config.Rate)
		}
	}()

	//for {
	//<-time.After(time.Second)
	//go func() {
	//persent := (requestCounter / float64(len(PureStack))) * 100
	//log.Println(persent, "%")
	//}()
	//}
}

func (clan *Clan) producer() {
	go func() {
		clan.startPoint.GetPageUrls()
	}()
}
