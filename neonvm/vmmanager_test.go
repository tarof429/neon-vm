package neonvm

import (
	"os"
	"testing"
)

func TestCreateFindDelete(t *testing.T) {

	var m VMManager
	m.SimulatorMode = false

	m.internalCreate("abcd", 1, 4)
	m.internalCreate("efgh", 2, 5)
	m.internalCreate("test", 3, 6)

	var index int

	index = m.Find("test")

	if index == -1 {
		t.Fail()
	}

	vm := m.VMList[index]

	if vm.Cpu != 6 {
		t.Fail()
	}

	m.internalDelete(index)

	index = m.Find("test")

	if index != -1 {
		t.Fail()
	}

}

func TestVMState(t *testing.T) {
	var m VMManager
	m.SimulatorMode = false

	m.internalCreate("abcd", 1, 4)
	m.internalCreate("efgh", 2, 5)
	m.internalCreate("test", 3, 6)

	index := m.Find("efgh")

	if index == -1 {
		t.Fail()
	}

	m.StartWithIndex(index)
	status := m.GetStatus(index)

	if status != VM_STATE_RUNNING {
		t.Fail()
	}

	m.StopWithIndex(index)
	status = m.GetStatus(index)

	if status != VM_STATE_STOPPED {
		t.Fail()
	}

}

func TestPersistence(t *testing.T) {
	var m VMManager
	m.SimulatorMode = false

	m.internalCreate("fff", 1, 4)
	m.internalCreate("ggg", 2, 5)
	m.internalCreate("hhh", 3, 6)

	err := os.Remove("neon.json")

	if err != nil {
		t.Fail()
	}

	status := m.WriteToDisk()

	if !status {
		t.Fail()
	}

	status = m.ReadFromDisk()

	if !status {
		t.Fail()
	}

	if index := m.Find("hhh"); index == -1 {
		t.Fail()
	}

	if index := m.Find("fff"); index == -1 {
		t.Fail()
	}

	if index := m.Find("ggg"); index == -1 {
		t.Fail()
	}

}
