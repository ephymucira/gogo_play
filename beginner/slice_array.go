package main
import ("fmt")

func main() {
	arr1 := [6]int{10,11,12,13,14,15}
	myslice := arr1[2:4]

	fmt.Printf("myslice = %v\n",myslice)
	fmt.Printf("Length = %d\n",len(myslice))
	fmt.Printf("myslice = %d\n",cap(myslice))
}