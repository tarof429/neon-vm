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
	VM_STATE_DELETING    string = "Deleting"
)

var (
	reader *bufio.Reader
)

type VMManager struct {
	vmlist     []VM
	cpuUsed    int
	memoryUsed int
}

type VM struct {
	name   string
	memory int
	cpu    int
	status string
}

func newVMManager() (m *VMManager) {
	m = &VMManager{
		vmlist: make([]VM, 0),
	}

	return m
}

func (m *VMManager) list() {
	fmt.Println("***********************")
	fmt.Println("***** Listing VMs *****")
	fmt.Println("***********************")

	for _, vm := range m.vmlist {
		fmt.Printf("%v\n", vm)
	}
}

func (m *VMManager) start() {
	var s string

	fmt.Printf("Enter the VM name: ")
	s, _ = reader.ReadString('\n')
	name := strings.TrimSpace(s)

	index := m.find(name)

	if index != -1 {
		state := m.vmlist[index].status
		if state == VM_STATE_NOT_STARTED || state == VM_STATE_STOPPED {
			go m._start(index)
		} else if state == VM_STATE_PENDING {
			fmt.Println("VM is still being created")
		}
	} else {
		fmt.Printf("VM %v does not exists\n", name)
	}
}

func (m *VMManager) _start(index int) {
	defer m.updateStatus(index, VM_STATE_RUNNING)
	m.updateStatus(index, VM_STATE_STARTING)
	time.Sleep(10 * time.Second)

}

func (m *VMManager) stop() {
	var s string

	fmt.Printf("Enter the VM name: ")
	s, _ = reader.ReadString('\n')
	name := strings.TrimSpace(s)

	index := m.find(name)

	if index != -1 {
		go m._stop(index)
	} else {
		fmt.Printf("VM %v does not exists\n", name)
	}
}

func (m *VMManager) _stop(index int) {
	m.updateStatus(index, VM_STATE_STOPPING)
	defer m.updateStatus(index, VM_STATE_STOPPED)
	time.Sleep(10 * time.Second)
}

func (m *VMManager) validateVMSettings(name string, cpu int, memory int) bool {
	if cpu+m.cpuUsed > MAX_CPUS {
		fmt.Printf("Unable to allocated %v CPU to VM", cpu)
		return false
	}
	if memory+m.memoryUsed > MAX_MEMORY {
		fmt.Printf("Unable to allocate %v memory to VM", memory)
		return false
	}

	if m.find(name) != -1 {
		fmt.Printf("Unable to create a VM with the same name: %v\n", name)
		return false
	}

	return true
}

func (m *VMManager) create() {
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

	if !m.validateVMSettings(name, cpu, memory) {
		return
	}

	vm := VM{
		name:   name,
		memory: memory,
		cpu:    cpu,
		status: VM_STATE_PENDING}

	m.vmlist = append(m.vmlist, vm)

	index := m.find(name)

	if index != -1 {
		go m._create(index)
	}
}

func (m *VMManager) _create(index int) {
	defer m.updateStatus(index, VM_STATE_STOPPED)
	time.Sleep(10 * time.Second)
}

func (m *VMManager) updateStatus(index int, status string) int {

	if index != -1 {
		m.vmlist[index].status = status
	}
	return index
}

func (m *VMManager) getStatus(index int) string {

	return m.vmlist[index].status
}

func (m *VMManager) delete() {
	fmt.Printf("Enter the name of the VM: ")

	var s string

	s, _ = reader.ReadString('\n')
	name := strings.TrimSpace(s)

	index := m.find(name)

	if index != -1 {
		go m._delete(index)
	} else {
		fmt.Printf("VM %v does not exists\n", name)
	}
}

func (m *VMManager) _delete(index int) {
	if m.getStatus(index) == VM_STATE_RUNNING {
		m._stop(index)
	}
	m.updateStatus(index, VM_STATE_DELETING)
	time.Sleep(10 * time.Second)
	m.vmlist = slices.Delete(m.vmlist, index, index+1)
}

func (m *VMManager) find(name string) int {
	for index, vm := range m.vmlist {
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

func main() {
	var choice string

	m := newVMManager()

	reader = bufio.NewReader(os.Stdin)

	for {
		usage()

		fmt.Print("Enter choice: ")
		fmt.Scanf("%v", &choice)

		switch choice {
		case "1":
			m.list()
		case "2":
			m.start()
		case "3":
			m.stop()
		case "4":
			m.create()
		case "5":
			m.delete()
		case "6":
			exitNeonVM()
		default:
			usage()
		}
	}
}
