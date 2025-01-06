package main
import ("fmt")

func myFunction(x int,y int) int {
	return x+y
}

func main()  {
	fmt.Println(myFunction(1,2))
}