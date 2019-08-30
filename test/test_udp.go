package main 

import(
	"github.com/ailumiyana/tools/udp"
	"fmt"
	"time"
	"net"
	"strconv"
)

func main() {

	pool := udp.NewPool(100, 6000, 10000)
	
	time.Sleep(10*time.Second)
	
	
	fmt.Println("actives : ", pool.ActiveNums)
	fmt.Println("Get:", pool.Get())
	fmt.Println("Get:", pool.Get())
	fmt.Println("Get:", pool.Get())
	time.Sleep(1*time.Second)
	fmt.Println("actives : ", pool.ActiveNums)
	
	fmt.Println("Put:9999")
	pool.Put(9999)
	fmt.Println("actives : ", pool.ActiveNums)
	
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:" + strconv.Itoa(pool.Get()))
	if err != nil {
		fmt.Println(":",err)
		return
	}

	_, err = net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("unavailable addr")
		return
	}

	fmt.Println("ListenUDP:", addr)
	time.Sleep(60*time.Second)

}