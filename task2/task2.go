package task2

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// 编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
func Add10(p *int) {
	base := *p
	*p += 10
	fmt.Println(base, "Add10:", *p)
}

// 实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
func DoubleSlice(nums *[]int) {
	for i := range *nums {
		(*nums)[i] *= 2
	}
	fmt.Println("DoubleSlice", *nums)
}

// 编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
func GoRoutine() {
	fmt.Println("start go routine")
	go func() {
		for i := 1; i <= 10; i += 2 {
			fmt.Println("奇数：", i)
		}
	}()
	go func() {
		for i := 2; i <= 10; i += 2 {
			fmt.Println("偶数：", i)
		}
	}()
	time.Sleep(time.Second * 1)
}

// 设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
func TaskScheduler() {
	var wg sync.WaitGroup
	tasks := []func(){
		func() {
			defer wg.Done()
			time.Sleep(time.Second * 1)
			fmt.Println("Task 1")
		},
		func() {
			defer wg.Done()
			fmt.Println("Task 2")
		},
		func() {
			defer wg.Done()
			fmt.Println("Task 3")
		},
	}
	start := time.Now() //记录开始时间
	for _, task := range tasks {
		go task() //启动协程
		wg.Add(1)
	}
	wg.Wait()

	fmt.Println("All tasks done, time elapsed:", time.Since(start))
}

// 定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
// 在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radius float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}
func Inf() {
	rect := Rectangle{Width: 5, Height: 3}
	circle := Circle{Radius: 2}
	fmt.Println("Rectangle area:", rect.Area())
	fmt.Println("Rectangle perimeter:", rect.Perimeter())
	fmt.Println("Circle area:", circle.Area())
	fmt.Println("Circle perimeter:", circle.Perimeter())
}

// 使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。
// 为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
type Person struct {
	Name string
	Age  int
}
type Employee struct {
	Person
	EmployeeID int
}

func (e Employee) PrintInfo() {
	fmt.Printf("Employee ID: %d, Name: %s, Age: %d\n", e.EmployeeID, e.Name, e.Age)
}

// 编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
func Channel() {
	ch := make(chan int)
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	for num := range ch {
		fmt.Println("Received:", num)
	}

}

// 实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
func BufferChannel() {
	ch := make(chan int, 7)
	go func() {
		for i := 1; i <= cap(ch); i++ {
			ch <- i
		}
		close(ch)
	}()
	for num := range ch {
		fmt.Println("bufferChannel Received:", num)
	}
}

// 编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func Mutex() {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Println("mutexCounter:", counter)
}

// 使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func AtomicCounter() {
	var counter int64
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Println("atomicCounter:", counter)
}
