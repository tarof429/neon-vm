# NeonVM

NeonVM simulates a VM manager. This program does not actually create VMs. It mainly illustrates Go features such as:

- Variables
- Switch statements
- Reading from the console
- For loops
- Constants
- Structs
- Goroutines
- Functions
- Functions as variables
- Pointers
- Packages
- Reading/writing to files
- JSON
- Testing

See the Makefile for instructions on how to build neonvm.

```sh
$ make
help                           Display this help
build                          Build the code
test                           Run all tests
clean                          Remove the dist directory
```

Below is a sample run.

```sh
$ ./neonvm 
***********************
***** Main Menu *******
***********************
1) List VMs
2) Start VM
3) Stop VM
4) Create VM
5) Delete VM
6) Stop all VMs
7) Exit Neon VM
Enter choice: 1   
{vm1 1 1 Stopped}
{vm2 1 1 Stopped}
***********************
***** Main Menu *******
***********************
1) List VMs
2) Start VM
3) Stop VM
4) Create VM
5) Delete VM
6) Stop all VMs
7) Exit Neon VM
Enter choice: 2 vm1
***********************
***** Main Menu *******
***********************
1) List VMs
2) Start VM
3) Stop VM
4) Create VM
5) Delete VM
6) Stop all VMs
7) Exit Neon VM
Enter choice: 1
{vm1 1 1 Starting}
{vm2 1 1 Stopped}
***********************
***** Main Menu *******
***********************
1) List VMs
2) Start VM
3) Stop VM
4) Create VM
5) Delete VM
6) Stop all VMs
7) Exit Neon VM
Enter choice: 1
{vm1 1 1 Running}
{vm2 1 1 Stopped}
***********************
***** Main Menu *******
***********************
1) List VMs
2) Start VM
3) Stop VM
4) Create VM
5) Delete VM
6) Stop all VMs
7) Exit Neon VM
Enter choice: 6
***********************
***** Main Menu *******
***********************
1) List VMs
2) Start VM
3) Stop VM
4) Create VM
5) Delete VM
6) Stop all VMs
7) Exit Neon VM
Enter choice: 1
{vm1 1 1 Stopping}
{vm2 1 1 Stopped}
***********************
***** Main Menu *******
***********************
1) List VMs
2) Start VM
3) Stop VM
4) Create VM
5) Delete VM
6) Stop all VMs
7) Exit Neon VM
Enter choice: 1
{vm1 1 1 Stopping}
{vm2 1 1 Stopped}
***********************
***** Main Menu *******
***********************
1) List VMs
2) Start VM
3) Stop VM
4) Create VM
5) Delete VM
6) Stop all VMs
7) Exit Neon VM
Enter choice: 1
{vm1 1 1 Stopped}
{vm2 1 1 Stopped}
***********************
***** Main Menu *******
***********************
1) List VMs
2) Start VM
3) Stop VM
4) Create VM
5) Delete VM
6) Stop all VMs
7) Exit Neon VM
Enter choice: 7
Goodbye!
```