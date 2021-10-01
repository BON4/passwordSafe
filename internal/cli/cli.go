package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
	"passwordSafe/internal"
	"passwordSafe/internal/store"
	"passwordSafe/pkg/crypt"
	"passwordSafe/utils"
	"strings"
	"syscall"
)

type CLI struct {
	fileName string
	store    internal.Store
	in       *bufio.Reader
}

func checkFile(fName string) bool {
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		return false
	}
	return true
}

func openFile(fName string) *os.File {
	f, err := os.OpenFile(fName, os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	return f
}

func openCreate(fName string) *os.File {
	f, err := os.OpenFile(fName, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	return f
}

func NewCLI(fName string) CLI {
	if checkFile(fName) {
		return CLI{fileName: fName, store: store.NewTxtStore(openFile(fName)), in: bufio.NewReader(os.Stdin)}
	}
	return CLI{fileName: fName, store: store.NewTxtStore(openCreate(fName)), in: bufio.NewReader(os.Stdin)}
}

func (c *CLI) NewDumpFile(fName string) {
	c.store = store.NewTxtStore(openFile(fName))
}

func (c CLI) getByKey(key string, secret []byte, writer io.Writer) error {
	cred, err := c.store.Get()
	if err != nil {
		return err
	}

	if v, ok := cred[key]; ok {
		dVal, err := crypt.DecryptChip(v, secret)
		if err != nil {
			if err.Error() == "cipher: message authentication failed" {
				return utils.NewWrongSecretKeyError()
			}
			return err
		}

		_, err = fmt.Fprintln(writer, string(dVal))
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Fprintln(writer, "key does not exist")
		if err != nil {
			return err
		}
	}
	return nil
}

func (c CLI) setByKey(key string, val string, secret []byte, writer io.Writer) error {
	eVal, err := crypt.EncryptChip([]byte(val), secret)
	if err != nil {
		return err
	}

	return c.store.Save(internal.Credentials{key: eVal})
}

func (c CLI) getSecret() ([]byte, error) {
	fmt.Print("Provide a SECRET KEY:")

	key, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}

	fmt.Println("")

	//TODO IF WINDOWS
	return crypt.ValidateSecretKey(bytes.TrimSuffix(key, []byte("\r\n"))), nil
}

func (c CLI) GetHandler(key string) error {
	secret, err := c.getSecret()
	if err != nil {
		return err
	}

	err = c.getByKey(strings.TrimSuffix(key, "\n"), secret, os.Stdout)
	if err != nil {
		return err
	}
	return nil
}

func (c CLI) SetHandler() error {
	secret, err := c.getSecret()
	if err != nil {
		return err
	}

	fmt.Println("New note.\nProvide a KEY:")

	key, err := c.in.ReadString('\n')
	if err != nil {
		return err
	}

	fmt.Println("\nProvide a VALUE:")
	val, err := c.in.ReadString('\n')
	if err != nil {
		return err
	}

	//TODO IF WINDOWS
	k := strings.TrimSuffix(key, "\r\n")
	v := strings.TrimSuffix(val, "\r\n")

	err = c.setByKey(k, v, secret, os.Stdout)
	if err != nil {
		return err
	}
	return nil
}

func (c CLI) ListHandler(msg string) error {
	cred, err := c.store.Get()
	if err != nil {
		return err
	}

	for k, _ := range cred {
		fmt.Println(k)
	}
	return nil
}

func (c CLI) InvalidCommandHandler(msg string) {
	fmt.Println(msg)
	return
}