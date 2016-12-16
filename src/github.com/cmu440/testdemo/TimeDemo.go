package main





import (
	//_ "fmt"
	//_ "time"
	"time"
	"fmt"
)

func main() {
	tick := time.Tick(1e8)
	boom := time.After(5e8) //只是执行一次

	for{
		select {
		case <-tick:
			fmt.Println("tick.")
			fmt.Println(<-tick)
		case <-boom:
			fmt.Println("boom.")
			fmt.Println(<-boom) //由于里面只是填充了一个数值，所以这个时候拿不到
			return
		default:
			fmt.Println("     .")
			time.Sleep(5e7)
		}
	}
}
