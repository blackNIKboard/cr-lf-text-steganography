package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
)

type coder struct {
	text    []byte
	message []position
}

type position struct {
	indexA int
	indexB int
	value  bool
}

func main() {
	text, err := ioutil.ReadFile("raw.txt")
	if err != nil {
		panic(err)
	}

	stegoSystem, err := NewCoder(text)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nRaw text")
	spew.Dump(text)

	fmt.Println("\nRaw message")
	fmt.Println(stegoSystem.DecodeMessage())

	err = stegoSystem.EncodeMessage([]int{1, 0, 1, 0, 1}) // Example message
	if err != nil {
		panic(err)
	}

	fmt.Println("\nEncoded message")
	fmt.Println(stegoSystem.DecodeMessage())

	err = stegoSystem.WriteFile("encoded.txt")
	if err != nil {
		panic(err)
	}

	// test and verify
	encodedText, err := ioutil.ReadFile("encoded.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println("\nEncoded text")
	spew.Dump(encodedText)

	stegoSystem, _ = NewCoder(encodedText)

	fmt.Println("\nDecoded message")
	fmt.Println(stegoSystem.DecodeMessage())
}

func NewCoder(text []byte) (*coder, error) {
	result := coder{
		text:    text,
		message: []position{},
	}

	// Find CR-LF LF-CR sequences
	for i := 0; i < len(result.text)-1; i++ {
		if result.text[i] == 0x0D && result.text[i+1] == 0x0A {
			result.message = append(result.message, position{
				indexA: i,
				indexB: i + 1,
				value:  false,
			})

			i++
		} else if result.text[i] == 0x0A && result.text[i+1] == 0x0D {
			result.message = append(result.message, position{
				indexA: i,
				indexB: i + 1,
				value:  true,
			})

			i++
		}
	}

	if len(result.message) == 0 {
		return nil, fmt.Errorf("not detected escape symbols")
	}

	return &result, nil
}

func (receiver *coder) DecodeMessage() ([]int, error) {
	var result []int

	if !receiver.ContainsMessage() {
		return nil, fmt.Errorf("coder doesn't contain a message")
	}

	for _, p := range receiver.message {
		result = append(result, b2i(p.value))
	}

	return result, nil
}

func (receiver *coder) ContainsMessage() bool {
	return receiver.message[0].value // Check mark
}

func b2i(b bool) int {
	if b {
		return 1
	}

	return 0
}

func i2b(i int) bool {
	return i == 1
}

func (receiver *coder) EncodeMessage(message []int) error {
	receiver.message[0].value = true

	for i := 0; i < len(message); i++ {
		receiver.message[i+1].value = i2b(message[i])
	}

	return nil
}

func (receiver *coder) WriteFile(filename string) error {
	for _, p := range receiver.message { // Parse (encode) message in text
		if p.value {
			receiver.text[p.indexA] = 0x0A
			receiver.text[p.indexB] = 0x0D
		} else {
			receiver.text[p.indexA] = 0x0D
			receiver.text[p.indexB] = 0x0A
		}
	}

	err := ioutil.WriteFile(filename, receiver.text, 0644)
	if err != nil {
		return err
	}

	return nil
}
