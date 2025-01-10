package neonvm

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	MAX_CPUS       int    = 25
	MAX_MEMORY     int    = 128
	NEON_DATA_FILE string = "neon.json"

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
	VMList        []VM
	cpuUsed       int
	memoryUsed    int
	SimulatorMode bool
}

type VM struct {
	Name   string
	Memory int
	Cpu    int
	status string
}

func A() {

}

func NewVMManager() (m *VMManager) {
	m = &VMManager{
		VMList: make([]VM, 0),
	}

	reader = bufio.NewReader(os.Stdin)

	return m
}

func (m *VMManager) List() {
	for _, vm := range m.VMList {
		fmt.Printf("%v\n", vm)
	}
}

func (m *VMManager) Start() {
	var s string

	fmt.Printf("Enter the VM name: ")
	s, _ = reader.ReadString('\n')
	name := strings.TrimSpace(s)

	index := m.Find(name)

	if index != -1 {
		state := m.VMList[index].status
		if state == VM_STATE_NOT_STARTED || state == VM_STATE_STOPPED {
			go m.StartWithIndex(index)
		} else if state == VM_STATE_PENDING {
			fmt.Println("VM is still being created")
		}
	} else {
		fmt.Printf("VM %v does not exists\n", name)
	}
}

func (m *VMManager) StartWithIndex(index int) {
	defer m.UpdateStatus(index, VM_STATE_RUNNING)
	m.UpdateStatus(index, VM_STATE_STARTING)
	if m.SimulatorMode {
		time.Sleep(10 * time.Second)
	}
}

func (m *VMManager) Stop() {
	var s string

	fmt.Printf("Enter the VM name: ")
	s, _ = reader.ReadString('\n')
	name := strings.TrimSpace(s)

	index := m.Find(name)

	if index != -1 {
		go m.StopWithIndex(index)
	} else {
		fmt.Printf("VM %v does not exists\n", name)
	}
}

func (m *VMManager) StopWithIndex(index int) {
	if m.VMList[index].status == VM_STATE_STOPPED {
		return
	}
	m.UpdateStatus(index, VM_STATE_STOPPING)
	defer m.UpdateStatus(index, VM_STATE_STOPPED)
	if m.SimulatorMode {
		time.Sleep(10 * time.Second)
	}
}

func (m *VMManager) StopAllVMs() {
	for i := 0; i < len(m.VMList); i++ {
		go m.StopWithIndex(i)
	}
}

func (m *VMManager) ValidateVMSettings(name string, cpu int, memory int) bool {
	if cpu+m.cpuUsed > MAX_CPUS {
		fmt.Printf("Unable to allocated %v CPU to VM", cpu)
		return false
	}
	if memory+m.memoryUsed > MAX_MEMORY {
		fmt.Printf("Unable to allocate %v memory to VM", memory)
		return false
	}

	if m.Find(name) != -1 {
		fmt.Printf("Unable to create a VM with the same name: %v\n", name)
		return false
	}

	return true
}

func (m *VMManager) Create() {
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

	if !m.ValidateVMSettings(name, cpu, memory) {
		return
	}

	index := m.Find(name)

	if index == -1 {
		m.internalCreate(name, memory, cpu)
	}
}

func (m *VMManager) internalCreate(name string, memory int, cpu int) {
	vm := VM{
		Name:   name,
		Memory: memory,
		Cpu:    cpu,
		status: VM_STATE_PENDING}

	m.VMList = append(m.VMList, vm)

	actuallyCreate := func() {
		index := m.Find(name)
		if m.SimulatorMode {
			time.Sleep(10 * time.Second)
		}
		m.UpdateStatus(index, VM_STATE_STOPPED)
	}
	go actuallyCreate()
}

func (m *VMManager) UpdateStatus(index int, status string) int {

	if index != -1 {
		m.VMList[index].status = status
	}
	return index
}

func (m *VMManager) GetStatus(index int) string {

	return m.VMList[index].status
}

func (m *VMManager) Delete() {
	fmt.Printf("Enter the name of the VM: ")

	var s string

	s, _ = reader.ReadString('\n')
	name := strings.TrimSpace(s)

	index := m.Find(name)

	if index != -1 {
		go m.internalDelete(index)
	} else {
		fmt.Printf("VM %v does not exists\n", name)
	}
}

func (m *VMManager) internalDelete(index int) {
	if m.GetStatus(index) == VM_STATE_RUNNING {
		m.StopWithIndex(index)
	}
	m.UpdateStatus(index, VM_STATE_DELETING)
	if m.SimulatorMode {
		time.Sleep(10 * time.Second)
	}
	m.VMList = slices.Delete(m.VMList, index, index+1)
}

func (m *VMManager) Find(name string) int {
	for index, vm := range m.VMList {
		if vm.Name == name {
			return index
		}
	}
	return -1
}

func (m *VMManager) WriteToDisk() bool {

	data, err := json.MarshalIndent(m.VMList, "", "\t")

	if err != nil {
		fmt.Println("Error when writing data")
		return false
	}

	err = os.WriteFile(
		filepath.Join(".", NEON_DATA_FILE), data, 0644)

	if err != nil {
		fmt.Println(err.Error())
	}

	return err == nil
}

func (m *VMManager) ReadFromDisk() bool {
	data, err := os.ReadFile(
		filepath.Join(".", NEON_DATA_FILE))

	if err != nil {
		fmt.Printf("Error while loading %v\n", NEON_DATA_FILE)
		return false
	}

	err = json.Unmarshal(data, &m.VMList)

	if err != nil {
		return false
	} else {
		for i := 0; i < len(m.VMList); i++ {
			m.VMList[i].status = VM_STATE_STOPPED
		}
		return true
	}
}

func (m *VMManager) ExitNeonVM() {
	fmt.Println("Goodbye!")

	status := m.WriteToDisk()

	if status {
		os.Exit(0)
	} else {
		os.Exit(1)
	}

}
