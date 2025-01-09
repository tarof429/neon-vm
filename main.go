package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	MAX_CPUS   int = 25
	MAX_MEMORY int = 128

	VM_STATE_PENDING     string = "Pending"
	VM_STATE_NOT_STARTED string = "Not started"
	VM_STATE_STARTING    string = "Starting"
	VM_STATE_RUNNING     string = "Running"
	VM_STATE_STOPPING    string = "Stopping"
	VM_STATE_STOPPED     string = "Stopped"
)

var (
	reader     *bufio.Reader
	vmlist     = make([]VM, 0)
	cpuUsed    int
	memoryUsed int
)

type VM struct {
	name   string
	memory int
	cpu    int
	status string
}

func usage() {
	fmt.Println("***********************")
	fmt.Println("***** Main Menu *******")
	fmt.Println("***********************")
	fmt.Println("1) List VMs")
	fmt.Println("2) Start VM")
	fmt.Println("3) Stop VM")
	fmt.Println("4) Create VM")
	fmt.Println("5) Delete VM")
	fmt.Println("6) Exit Neon VM")
}

func listVMs() {
	fmt.Println("***********************")
	fmt.Println("***** Listing VMs *****")
	fmt.Println("***********************")

	for _, vm := range vmlist {
		fmt.Printf("%v\n", vm)
	}
}

func startVM() {
	var s string

	fmt.Printf("Enter the VM name: ")
	s, _ = reader.ReadString('\n')
	name := strings.TrimSpace(s)

	index := findVM(name)

	if index != -1 {
		state := vmlist[index].status
		if state == VM_STATE_NOT_STARTED || state == VM_STATE_STOPPED {
			go _startVM(index)
		} else if state == VM_STATE_PENDING {
			fmt.Println("VM is still being created")
		}
	} else {
		fmt.Printf("VM %v does not exists\n", name)
	}
}

func _startVM(index int) {
	defer updateStatus(index, VM_STATE_RUNNING)
	updateStatus(index, VM_STATE_STARTING)
	time.Sleep(10 * time.Second)

}

func stopVM() {
	var s string

	fmt.Printf("Enter the VM name: ")
	s, _ = reader.ReadString('\n')
	name := strings.TrimSpace(s)

	index := findVM(name)

	if index != -1 {
		go _stopVM(index)
	} else {
		fmt.Printf("VM %v does not exists\n", name)
	}
}

func _stopVM(index int) {
	updateStatus(index, VM_STATE_STOPPING)
	defer updateStatus(index, VM_STATE_STOPPED)
	time.Sleep(10 * time.Second)
}

func validateVMSettings(name string, cpu int, memory int) bool {
	if cpu+cpuUsed > MAX_CPUS {
		fmt.Printf("Unable to allocated %v CPU to VM", cpu)
		return false
	}
	if memory+memoryUsed > MAX_MEMORY {
		fmt.Printf("Unable to allocate %v memory to VM", memory)
		return false
	}

	if findVM(name) != -1 {
		fmt.Printf("Unable to create a VM with the same name: %v\n", name)
		return false
	}

	return true
}

func createVM() {
	var name string
	var memory int
	var cpu int

	var s string

	fmt.Printf("Enter the VM name: ")
	s, _ = reader.ReadString('\n')
	name = strings.TrimSpace(s)

	fmt.Printf("Enter the VM memory(GB): ")
	s, _ = reader.ReadString('\n')
	memory, err := strconv.Atoi(strings.TrimSpace(s))

	if err != nil {
		fmt.Printf("Invalid memory given: %v\n", s)
		return
	}

	fmt.Printf("Enter the VM CPU: ")
	s, _ = reader.ReadString('\n')
	cpu, err = strconv.Atoi(strings.TrimSpace(s))

	if err != nil {
		fmt.Printf("Invalid cpu given: %v\n", s)
		return
	}

	if !validateVMSettings(name, cpu, memory) {
		return
	}

	vm := VM{
		name:   name,
		memory: memory,
		cpu:    cpu,
		status: VM_STATE_PENDING}

	vmlist = append(vmlist, vm)

	index := findVM(name)

	if index != -1 {
		go _createVM(index)
	}
}

func _createVM(index int) {
	defer updateStatus(index, VM_STATE_STOPPED)
	time.Sleep(10 * time.Second)

}

func updateStatus(index int, status string) int {

	if index != -1 {
		vmlist[index].status = status
	}
	return index
}

func deleteVM() {
	fmt.Printf("Enter the name of the VM: ")

	var s string

	s, _ = reader.ReadString('\n')
	name := strings.TrimSpace(s)

	index := findVM(name)

	if index != -1 {
		go _deleteVM(index)
	} else {
		fmt.Printf("VM %v does not exists\n", name)
	}
}

func _deleteVM(index int) {
	time.Sleep(10 * time.Second)
	updateStatus(index, VM_STATE_STOPPED)
	time.Sleep(10 * time.Second)
	vmlist = slices.Delete(vmlist, index, index+1)
}

func findVM(name string) int {
	for index, vm := range vmlist {
		if vm.name == name {
			return index
		}
	}
	return -1
}

func exitNeonVM() {
	fmt.Println("Exiting Neon VM. Goodbye!")
	os.Exit(0)
}

func main() {
	var choice string

	reader = bufio.NewReader(os.Stdin)

	for {
		usage()

		fmt.Print("Enter choice: ")
		fmt.Scanf("%v", &choice)

		fmt.Printf("Your choice was %v\n", choice)

		switch choice {
		case "1":
			listVMs()
		case "2":
			startVM()
		case "3":
			stopVM()
		case "4":
			createVM()
		case "5":
			deleteVM()
		case "6":
			exitNeonVM()
		default:
			usage()
		}
	}
}
