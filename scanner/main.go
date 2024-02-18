package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
)

const colorRed = "\033[31m"
const colorReset = "\033[0m"

func worker(ports chan int, results chan int, target string) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", target, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			// port is closed or filtered
			results <- 0
			continue
		}

		conn.Close()
		results <- p
	}
}

func main() {

	amountPorts := flag.Int("amount-ports", 1024, "This option specifies wich amount ports you want scan and override default")
	target := flag.String("target", "127.0.0.1", "This option specifies wich target you will scan. You can use url or ip")

	flag.Parse()

	ports := make(chan int, *amountPorts)
	results := make(chan int)
	var openports []int
	fmt.Println(string(colorRed), "Start scanner...", string(colorReset))
	for i := 0; i <= cap(ports); i++ {
		go worker(ports, results, *target)
	}
	go func() {
		for i := 1; i <= *amountPorts; i++ {
			ports <- i
		}
	}()
	for i := 1; i <= *amountPorts; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}
	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
