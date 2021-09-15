# Text Steganography
This example based on old CR-LF swap method (of course, it's not working anymore due to modern LF standard)

## Principles
Text contains escape seq's of carriage return (CR-LF). 
Embedded message is hidden in swapped sequences (LF-CR).

# Running
- `go mod download`  
- `go run main.go`