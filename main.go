package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	filename := "accounts.txt"
	createFile(filename)

	sections := promptInt("How many sections are there? (ex: 3)")
	students := promptInt("How many students are in each section? (ex: 40)")
	prefix := promptString("What prefix should be used? (ex: IS3030_FL_18)")
	password := promptString("What default password should be used? (ex: Bearcat1)")

	f, err := os.OpenFile(filename, os.O_RDWR, 0644) // 0644 specifies to overwrite the file
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Truncate(0) // Clears out any data in accounts.txt

	for sec := 1; sec <= sections; sec++ {
		for st := 1; st <= students; st++ {
			u := buildUser(prefix, sec, st)
			q := buildQuery(u, password)
			if _, err = f.WriteString(q); err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Println("Account script has been generated and can be found in", filename)
	fmt.Println("Copy and paste all of the file's text into your admin enabled SQL Plus account")
}

func createFile(filename string) {
	var _, err = os.Stat(filename)

	if os.IsNotExist(err) {
		var file, err = os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}
}

func buildQuery(user string, password string) string {
	s := `
DROP USER "%[1]s" CASCADE;
CREATE USER "%[1]s" PROFILE "DEFAULT" IDENTIFIED BY "%[2]s" DEFAULT TABLESPACE "USERS" TEMPORARY TABLESPACE "TEMP" ACCOUNT UNLOCK;
GRANT UNLIMITED TABLESPACE TO "%[1]s";
GRANT "CONNECT" TO "%[1]s";
GRANT "RESOURCE" TO "%[1]s";
`
	return fmt.Sprintf(s, user, password)
}

func buildUser(prefix string, section int, student int) string {
	return fmt.Sprintf("%s_%03d_%03d", prefix, section, student)
}

func promptInt(prompt string) int {
	var i int
	fmt.Println(prompt)
	_, err := fmt.Scan(&i)
	if err != nil || i < 1 {
		log.Fatal("Error: Please enter a valid number (greater than 0) next time!")
	}

	return i
}

func promptString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)
	s, err := reader.ReadString('\n')
	s = strings.TrimRight(s, "\n")
	if err != nil || s == "" {
		log.Fatal("Error: Please enter a non-empty string next time!")
	}
	return s
}
