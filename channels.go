package main

import (
	"fmt"
	"sync"
)

func worker1(ch1,ch2 chan bool,wg *sync.WaitGroup){
defer wg.Done()
for i:=1;i<=20;i+=2{
	<-ch1
	fmt.Println(i)
    ch2<-true
}
}
func worker2(ch1,ch2 chan bool,wg *sync.WaitGroup){
defer wg.Done()
for i:=2;i<=20;i+=2{
	<-ch2
	fmt.Println(i)
	if i<20{
    ch1<-true
	}
}
}

func main2(){
ch1:=make(chan bool)
ch2:=make(chan bool)
var wg sync.WaitGroup
wg.Add(2)
go worker1(ch1,ch2,&wg)
go worker2(ch1,ch2,&wg)
ch1<-true
wg.Wait()
}

func main() {
	var x int = 0
	ch1 := make(chan int)
	ch2 := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)
	go func(x int) {
		defer wg.Done()
		ch1 <- x + 1
		// fmt.Println(<-ch1)
	}(x)
	go func(x int) {
		defer wg.Done()
		ch2 <- x + 1
		// fmt.Println(<-ch2)
	}(x)
	fmt.Println(<-ch1)
	fmt.Println(<-ch2)
	wg.Wait()
	close(ch1)
	close(ch2)
}
