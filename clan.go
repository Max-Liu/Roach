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

func init() {

}

func NewClan(config *ClanConfig) *Clan {
	link := newLink(config.StartPoint)
	config.LinkConfig.Target = config.StartPoint
	config.LinkConfig.Log = config.Log
	config.LinkConfig.Host = config.Host
	link.setConfig(config.LinkConfig)
	return &Clan{
		startPoint: link,
		config:     config,
		workerChan: make(chan *Link, config.ConcurrentNumber),
	}
}

func (clan *Clan) Rush() {

	producer(clan.startPoint)
	clan.customer(clan.startPoint)
	<-exitChan
}

func (clan *Clan) SetConfig() {

}

func (clan *Clan) customer(l *Link) {
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
			//file.WriteString(l.Url + " " + l.Title + strconv.Itoa(i.StatusCode) + "\n")
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

func producer(l *Link) {
	go func() {
		l.GetPageUrls()
	}()
}
