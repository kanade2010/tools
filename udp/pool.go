package udp

import(
	"net"
	"strconv"
	"time"
	"sync"
	"math/rand"
	
	//"fmt"
	//"sync/atomic"
)

const MaxReapCnts = 65535

type portPool struct {
	MaxNums     int
	ActiveNums  int
	baseNumber  int
	MaxNumber   int
	cMaps 	    map[int]chan struct{}
	ch          chan struct{}
	get         chan struct{}
	sync.Mutex
}

func NewPool(max_nums, base_number, max_number int) *portPool{
	pool := &portPool{
		MaxNums:    max_nums,
		baseNumber: base_number,
		MaxNumber:  max_number,
		cMaps:		make(map[int]chan struct{}),
		ch:         make(chan struct{}),
		get:        make(chan struct{}),
	}
	
	pool.reapPort()

	return pool
}

func (p *portPool)reapPort() {

	if p.baseNumber > p.MaxNumber {
		panic("error args")
	}

	i := 0
	p.ActiveNums = 0

	for p.ActiveNums < p.MaxNums {
		if i > MaxReapCnts {
			break
		}
		//rand port
		seed := rand.New(rand.NewSource(time.Now().UnixNano()))
		randNum := seed.Intn(p.MaxNumber-p.baseNumber)
		port := p.baseNumber + randNum
		
		//port := p.baseNumber + i
		
		go p.listenUdp(port)
		
		<- p.ch

		i++
	}

	p.MaxNums = p.ActiveNums

	//fmt.Println("end-")
}

func (p *portPool)listenUdp(port int) {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:" + strconv.Itoa(port))
	if err != nil {
		//fmt.Println(":",err)
		p.ch <- struct{}{}
		return
	}

	ln, err := net.ListenUDP("udp", addr)
	if err != nil {
		//fmt.Println("unavailable addr")
		p.ch <- struct{}{}
		return
	}

	p.ActiveNums++
	
	//fmt.Printf(" : %d\n", port)

	p.ch <- struct{}{}
	p.cMaps[port] = make(chan struct{})
	<- p.cMaps[port]
	
	//fmt.Printf("Get : %d\n", port)
	
	delete(p.cMaps, port)

	p.ActiveNums--

	ln.Close()

	p.get <- struct{}{}
}

func (p *portPool)Get() int {

	p.Lock()
	defer p.Unlock()

	if p.ActiveNums < 1 {
		return 0
	}

	for k, _ := range p.cMaps {
		p.cMaps[k] <- struct{}{}
		<- p.get
		return k
	}

	return 0
}

//user to ensure that port is available
func (p *portPool)Put(port int) {

	p.Lock()
	defer p.Unlock()

	if p.ActiveNums < p.MaxNums {
		go p.listenUdp(port)
		<- p.ch
	}
}