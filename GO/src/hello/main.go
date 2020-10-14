package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func dijkstra(sommet string, liens map[string][][]string, wg *sync.WaitGroup) {
	defer wg.Done()

	/*fmt.Print("sommet :")
	fmt.Print(sommet)
	fmt.Println(liens[sommet])*/
}

func main() {
	//hourOfDay := time.Now().Hour()
	//greeting := getGreeting(hourOfDay)
	//fmt.Println(greeting)
	file, err := os.Open("schema.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	compteurLiens := 0

	scanner := bufio.NewScanner(file)

	scanner.Scan() //On lit la première ligne du doc pour avoir le nombre de liens
	//fmt.Println(scanner.Text())

	n, err := strconv.Atoi(scanner.Text())
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	liens := make([][]string, n)
	//fmt.Println(liens)

	for scanner.Scan() { //On remplit un slice avec les différents liens
		//fmt.Println(scanner.Text())

		lien := strings.Split(scanner.Text(), " ")
		//fmt.Println(lien)
		liens[compteurLiens] = lien
		compteurLiens++

	}
	//fmt.Println(liens)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	m := make(map[string][][]string)
	for _, lien := range liens { //On remplit la map

		if m[lien[0]] == nil { //On ajoute le sommet si inexistant
			m[lien[0]] = make([][]string, 0)
		}
		if m[lien[1]] == nil {
			m[lien[1]] = make([][]string, 0)
		}
		m[lien[0]] = append(m[lien[0]], lien)
		m[lien[1]] = append(m[lien[1]], lien)
	}

	var wg sync.WaitGroup

	for sommet, _ := range m {
		wg.Add(1)
		//fmt.Println(sommet)
		go dijkstra(sommet, m, &wg)
	}
	wg.Wait()

	//go routine
	//entree : sommet du graphe
	//sortie : les chemins les plus courts vers les autres sommets

	//on lance une routine par sommet dans m

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
