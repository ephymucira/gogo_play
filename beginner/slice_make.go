package main
import ("fmt")

func main(){
	myslice1 := make([]int,5,10)
	fmt.Printf("Type = %T\n",myslice1[0])
	fmt.Printf("myslice1 = %v\n",myslice1)
	fmt.Printf("length = %d\n",len(myslice1))
	fmt.Printf("capacity = %d\n",cap(myslice1))

	// with omitted capacity

	myslice2 := make([] int,5)
	fmt.Printf("Type = %T\n",myslice2)
	fmt.Printf("myslice2 = %v\n",myslice2)
	fmt.Printf("length = %d\n",len(myslice2))
	fmt.Printf("Capacity = %d\n",cap(myslice2))
}