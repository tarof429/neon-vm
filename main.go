package main

import (
	"fmt"

	"neonvm/neonvm"
)

func usage() {
	fmt.Println("***********************")
	fmt.Println("***** Main Menu *******")
	fmt.Println("***********************")
	fmt.Println("1) List VMs")
	fmt.Println("2) Start VM")
	fmt.Println("3) Stop VM")
	fmt.Println("4) Create VM")
	fmt.Println("5) Delete VM")
	fmt.Println("6) Stop all VMs")
	fmt.Println("7) Exit Neon VM")
}

func main() {
	var choice string

	m := neonvm.NewVMManager()
	m.SimulatorMode = true

	m.ReadFromDisk()

	for {
		usage()

		fmt.Print("Enter choice: ")
		fmt.Scanf("%v", &choice)

		switch choice {
		case "1":
			m.List()
		case "2":
			m.Start()
		case "3":
			m.Stop()
		case "4":
			m.Create()
		case "5":
			m.Delete()
		case "6":
			m.StopAllVMs()
		case "7":
			m.ExitNeonVM()
		default:
			usage()
		}
	}
}
