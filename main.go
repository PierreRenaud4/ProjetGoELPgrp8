package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strconv"
	//"time"
)

func main() {
	//hourOfDay := time.Now().Hour()
	//greeting := getGreeting(hourOfDay)
	//fmt.Println(greeting)
	file, err := os.Open("schema.txt")
  if err != nil {
    log.Fatal(err)
  }
	defer file.Close()
	nbLignes:=0
	scanner := bufio.NewScanner(file)
	premièreLigne:=true

	for scanner.Scan() {
      //fmt.Println(scanner.Text())
			if premièreLigne{
				n, err:=strconv.Atoi(scanner.Text())
				if err != nil {
        	// handle error
        	fmt.Println(err)
        	os.Exit(2)
    		}
				nbLignes=n


				premièreLigne=false
			}else{

			}
  }
	fmt.Println(nbLignes)

	var liens [nbLignes][4]string
	fmt.Println(liens)

  if err := scanner.Err(); err != nil {
      log.Fatal(err)
  }



	//nbLiens=

	//lien3 :=
	//m:=make(map[string][]string)

	//lien1
	//if lien1[0] == [] {
		//m[lien1[0]]=[1][]string{lien1}
	//}





	//m["A"] = [lien1]
	//m["B"] = [lien1, lien2]
	//m["C"] = [lien2]


}
