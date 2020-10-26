// Projet Go grTC1_8:  Dijkstra
//Tom POUPARD
//Pierre RENAUD
//Enzo ZATTARIN
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

//*******************************************************************************
//
//*******************************************************************************
// fonction de l'alogorithme dijkstra appeler par une go routine
func dijkstra(sommet string, liens map[string][][]string, wg *sync.WaitGroup) {
	defer wg.Done()

	/*fmt.Print("sommet :")
	fmt.Print(sommet)
	fmt.Println(liens[sommet])*/
}

//********************************************************************************
//
//********************************************************************************
// Programme principal
func main() {
	//hourOfDay := time.Now().Hour()
	//greeting := getGreeting(hourOfDay)
	//fmt.Println(greeting)
	file, err := os.Open("schema.txt") // permet d'ouvir le fichier txt en appelant file
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // ferme le fichier quand le main est fini

	compteurLiens := 0

	scanner := bufio.NewScanner(file)

	scanner.Scan() //On lit la première ligne du doc pour avoir le nombre de liens avec .scan()
	fmt.Println(scanner.Text())

	n, err := strconv.Atoi(scanner.Text()) //On le convertit en entier pour pouvoir créer nos slice de slice : "tableau de liens"
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	liens := make([][]string, n) // notre slice de slice

	//On remplit un slice avec les différents liens
	for scanner.Scan() { // on parcoure tous le fichier ligne par ligne

		lien := strings.Split(scanner.Text(), " ") // on lit chaque ligne avec la fonction scanner.text() et on créer un slice de string avec la fonction strings.Split avec comme indication " " comme séparateur
		liens[compteurLiens] = lien                // on associe pour chaque item de la "slice de slice" "un lien"(sous forme d'un slice) du fichier texte dans l'ordre chronologique
		compteurLiens++

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	m := make(map[string][][]string) // création d'une map
	//On remplit la map avec les élément de notre " slice de slice "
	for _, lien := range liens { // on parcours chaque "lien" de "liens" (slice de slice)
		//Idée : On parcours tout liens, dans chaque lien on lit les 2 premiers item qui sont les sommets
		//Ensuite on teste s'ils ont déja été créé dans dans la map si non on les crées, ces sommets seront les clés de la map
		if m[lien[0]] == nil {
			m[lien[0]] = make([][]string, 0) // on associe une clé à un tableau
		}
		if m[lien[1]] == nil {
			m[lien[1]] = make([][]string, 0)
		}
		//Enfin on rempli chaque tableau de la clé avec le lien associé(càd l'iteration en cours)
		m[lien[0]] = append(m[lien[0]], lien) // append permet de pouvoir aggrandir la taille de ce tableau en lui mettant un lien
		m[lien[1]] = append(m[lien[1]], lien)
		// exemple : le tableau de "B" va être creer à l'itération 1, ensuite il va être rempli avec le lien 1 ET à l'itération 2 et 3 on va l'agrendir en lui mettant le lien 2 et le lien 3
		// B :[[lien1][lien2][lien3]]
	}

	// Utilisation du principe de go routine en executant simultanément plusieur fonction tout en evitant que l'execution des fonction se chevauche
	var wg sync.WaitGroup
	//fmt.Println(m)

	for sommet, _ := range m {
		wg.Add(1)
		//fmt.Println(sommet)
		go dijkstra(sommet, m, &wg)
	}
	wg.Wait()

}

//***********************************************************************************
//
//***********************************************************************************
//	Objectif
//go routine
//entree : sommet du graphe
//sortie : les chemins les plus courts vers les autres sommets
//on lance une routine par sommet dans m
