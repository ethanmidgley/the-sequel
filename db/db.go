package db

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethanmidgley/the-sequel/in-memory/handlers"
	"github.com/ethanmidgley/the-sequel/in-memory/pkg/resp"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	go func() {
		for {
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func InstantiateAof() *Aof {

	// TODO: READ FROM FILE HERE
	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println("FAILED TO ACCESS: database.aof")
		os.Exit(1)
	}

	// aof.file.Read(func(value resp.Value) {
	//
	// })
	err = aof.Read()
	if err != nil {
		fmt.Println("FAILED TO PARSE DATABASE")
		os.Exit(1)
	}

	return aof

}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

// TODO:: TURN THIS IN TO A ROUTINE WHICH WILL UPDATE BASED ON A CHANNEL OF UPDATES
// WILL BOTTLE NECK EACH TIME WE SAVE A COMMAND
func (aof *Aof) Write(value resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) Read() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	aof.file.Seek(0, io.SeekStart)

	reader := resp.New(aof.file)

	for {
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := handlers.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
		}
		handler(args)

	}
	return nil
}

var AOF *Aof = InstantiateAof()
