package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

var (
	server string
	protocol string
	beginPort  int
	endPort  int
	IsWrite  bool
	help   bool
	openPort Port
	waitGroup sync.WaitGroup
	dialTimeout int
)

type Port struct {
	temp
	mutex sync.Mutex
}

type temp []int

func (p temp) Swap(i,j int) {
	p[i],p[j] = p[j],p[i]
}

func (p temp) Less(i,j int) bool {
	return p[i] < p[j]
}

func (p temp) Len() int {
	return len(p)
}

func init() {
	flag.StringVar(&server,"server","127.0.0.1","Server IP Address Or Domain Name")
	flag.StringVar(&protocol,"protocol","tcp","Protocol TCP")
	flag.IntVar(&beginPort,"beginPort",1,"Begin Port")
	flag.IntVar(&endPort,"endPort",65535,"End Port")
	flag.BoolVar(&IsWrite,"isWrite",false,"Write File In Dir")
	flag.BoolVar(&help,"h",false,"this help")
	flag.IntVar(&dialTimeout,"dialTimeOut",0,"Is Timeout")
}

func checkPort(port int) {
	address := server +":"+ fmt.Sprintf("%d",port)
	if conn,err:=net.DialTimeout(protocol,address,time.Duration(dialTimeout));err == nil {
		openPort.mutex.Lock()
		openPort.temp = append(openPort.temp,port)
		openPort.mutex.Unlock()
		conn.Close()
	}
	waitGroup.Done()
}

func writeFile(d time.Duration) {
  file,fErr:=os.OpenFile("port_scan.txt",os.O_APPEND|os.O_CREATE,0666)
  if fErr != nil {
  	fmt.Println("Open File Error:",fErr)
  }
  defer file.Close()
  bufio.NewReader(file)
  content := " server:%s\n" +
  			 " protocol:%s\n" +
  	         " openPort:%d\n" +
  	         " length:%d\n" +
  			 " spendTime:%s\n" +
  			 " createdTime:%s\n"
  wTitle := "------------------------------\n"
  wContent:= fmt.Sprintf(content,server,protocol,openPort.temp,len(openPort.temp),d,time.Now().Format("2006-01-02 15:04:05"))
  wEnd   := "-------------------------------\n"
  txt :=wTitle +wContent+wEnd
  file.WriteString(txt)
}

func Bar() {
	str :="端口扫描中"
	for range time.Tick(time.Second * 1){
		str += "="
		fmt.Printf("\r%s",str+">")
	}
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	go Bar()
	beginTime:=time.Now()
	for port:=beginPort;port<=endPort;port++ {
		waitGroup.Add(1)
		go checkPort(port)
	}
	waitGroup.Wait()
	sort.Sort(openPort.temp)
	endTime := time.Now()
	fmt.Println("\nopenPort:",openPort.temp)
	if IsWrite {
		writeFile(endTime.Sub(beginTime))
	}
}
