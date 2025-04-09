package main

import (
	"fmt"

	"github.com/go_web/task1"
	"github.com/go_web/task2"
)

func print() {
	fmt.Println("---------------------------------------------------")
}
func main() {
	task1.SingleNum()
	print()
	task1.Rob()
	print()
	task1.RverseString()
	print()

	task1.Sqrt(10)
	print()

	task1.RemoveDuplicates()
	print()
	task1.Merge()
	print()
	cal := task1.MyCalendar{}
	cal.Book(10, 20)
	cal.Book(15, 25)
	cal.Book(20, 30)
	fmt.Println("Bookings ", cal.Bookings)
	print()

	p := 11
	task2.Add10(&p)
	print()

	num := []int{1, 2, 13, 4, 1, 6, 5, 28, 6, 10}
	task2.DoubleSlice(&num)
	print()

	task2.GoRoutine()
	print()
	task2.TaskScheduler()
	print()

	task2.Inf()
	print()
	task2.Employee{Person: task2.Person{Name: "Alice", Age: 25}, EmployeeID: 1001}.PrintInfo()
	print()
}
