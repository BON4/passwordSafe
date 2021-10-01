package main

import (
	"flag"
	"fmt"
	"os"
	"passwordSafe/internal/cli"
)

var defaultFile = "store.txt"

const envFile = "PASSWORD_FILE"


func main() {
	//ONLY FOR PROD
	defer func(){
		if recoveryMessage := recover(); recoveryMessage != nil {
			fmt.Printf("ERROR: %s\n", recoveryMessage)
		}
	}()

	userFile := flag.Bool("f", false, "file that will store your passwords")

	getPassword := flag.String("get", "", "Provide a key by witch you saved your password")

	setPassword := flag.Bool("set", false, "Set new password by providing a new key and password")

	listPassword := flag.Bool("ls", false, "get all available keys")

	flag.Parse()

	if *userFile {
		fmt.Printf("Check your environment variable : %q\nIf it is not set file with passwords will be stored in current directory\n", envFile)
		return
	}

	file := os.Getenv(envFile)
	if file == "" {
		file = defaultFile
	}

	c := cli.NewCLI(file)

	if *getPassword == "" && !*setPassword && !*listPassword && *userFile {
		return
	}

	if !((((*getPassword != "") != (*setPassword)) != (*listPassword)) != (*userFile)){
		c.InvalidCommandHandler("Only one command at the time")
		return
	}

	//LIST
	if *getPassword == "" && !*setPassword && *listPassword {
		err := c.ListHandler(*getPassword)
		if err != nil { panic(err) }
		return
	}

	//GET
	if *getPassword != "" && !*setPassword {
		err := c.GetHandler(*getPassword)
		if err != nil { panic(err) }
		return
	}

	//SET
	if *getPassword == "" && *setPassword {
		err := c.SetHandler()
		if err != nil { panic(err) }
		return
	}

}

//func main() {
//	cred := map[string][]byte{"Google": []byte("Password: 90900")}
//
//	ecred, err := crypt.EncryptChip(cred["Google"], secret)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(string(ecred))
//
//	dcred, err := crypt.DecryptChip(ecred, secret)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(string(dcred))
//}