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

type Chemin struct {
	poids            int
	sommetsEmpruntes []string
	fini             bool
}

type Graphe map[string]*Chemin

//*******************************************************************************
//
//*******************************************************************************
// fonction de l'alogorithme dijkstra appeler par une go routine
func dijkstra(sommetDepart string, liens map[string][][]string, wg *sync.WaitGroup) {

	defer wg.Done()
	const inf = 999

	//Initialisation
	tableCouts := make(Graphe)

	for key, _ := range liens {

		tableCouts[key] = &Chemin{inf, []string{}, false}

		if key == sommetDepart {
			tableCouts[key].poids = 0
		}
	}
	//fmt.Println(tableCouts)
	var prochainSommet string
	var continuer bool
	var cheminLeger int
	var poidsCheminActuel int
	continuer = true
	cheminLeger = inf
	poidsCheminActuel = 0

	itineraire := make([]string, 0)
	//Condition d'arrêt : Tous les chemins sont finaux
	//On en profite pour choisir le prochain sommet a étudier

	for continuer == true {
		//On vérifie si on continue et on en profite pour déterminer le lien à étudier
		continuer = false
		cheminLeger = inf
		for key, _ := range tableCouts {
			if tableCouts[key].fini == false {
				continuer = true
			}
			if tableCouts[key].poids < cheminLeger && !tableCouts[key].fini {
				cheminLeger = tableCouts[key].poids
				prochainSommet = key

			}
		}
		//Le prochain sommet à traiter a été choisi
		fmt.Println("Prochain sommet : ", prochainSommet)
		//L'itinéraire est mis à jour ainsi que le poids de ce-dernier
		//Attention ! l'itinéraire peut régresser à un état précédent !!!

		itineraire = append(itineraire, prochainSommet)

		poidsCheminActuel += tableCouts[prochainSommet].poids
		//le sommet atteint est fixé : son chemin le plus court est compris dans la variable itineraire
		tableCouts[prochainSommet].fini = true
		//tableCouts[prochainSommet].sommetsEmpruntes = itineraire

		//fmt.Println(liens[prochainSommet])
		for _, lien := range liens[prochainSommet] {
			fmt.Println("lien : ", lien)
			//On distingue le voisin du sommet étudié (l'ordre des sommets n'est pas toujours le même)
			var voisin string
			if voisin = lien[0]; lien[0] == prochainSommet {
				voisin = lien[1]
			}

			if !tableCouts[voisin].fini {
				fmt.Println("voisin : ", voisin)
				//On change le poids dans la table des coûts
				nouveauPoids, err := strconv.Atoi(lien[3])
				if err != nil {
					// handle error
					fmt.Println(err)
					os.Exit(2)
				}
				//On détermine si le chemin étudié est plus avantageux que celui actuel !!!
				if poidsCheminActuel+nouveauPoids < tableCouts[voisin].poids {
					tableCouts[voisin].poids = poidsCheminActuel + nouveauPoids
				}

			}

		}
		fmt.Println("Sommet traité ! : ", prochainSommet)
	}

	for key, elem := range tableCouts {
		fmt.Println("Sommet : ", key, "plus court chemin : ", elem)
	}

	//tant que

	//fmt.Println(tableCouts)
	//fmt.Println(liens[prochainSommet])

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

	n, err := strconv.Atoi(scanner.Text()) //On le convertie en entier pour pouvoir creer nos slice de slice : "tableau de liens"
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	liens := make([][]string, n) // notre slice de slice

	//On remplit un slice avec les différents liens
	for scanner.Scan() { // on parcours tous le fichier ligne par ligne

		lien := strings.Split(scanner.Text(), " ") // on lit chaque ligne avec la fonction scanner.text() et on créer un slice de string avec la fonction strings.Split avec comme indication " " comme séparateur
		liens[compteurLiens] = lien                // on associe pour chaque item de la "slice de slice" "un lien"(sous forme d'un slice) du fichier texte dans l'ordre chronologique
		compteurLiens++

	}
	//fmt.Println(liens)

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
	//fmt.Println(m)
	// Utilisation du principe de go routine en executant simultanément plusieur fonction tout en evitant que l'execution des fonction se chevauche
	var wg sync.WaitGroup
	//fmt.Println(m)
	wg.Add(1)
	dijkstra("A", m, &wg)
	//for sommet, _ := range m {
	//wg.Add(1)
	//fmt.Println(sommet)
	//go dijkstra(sommet, m, &wg)
	//}
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
