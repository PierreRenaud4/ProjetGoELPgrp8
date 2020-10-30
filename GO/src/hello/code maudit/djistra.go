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
	//"sync"
)
//********************************************************************************
//
//********************************************************************************
// Programme principale
func main() {
	//hourOfDay := time.Now().Hour()
	//greeting := getGreeting(hourOfDay)
	//fmt.Println(greeting)
	file, err := os.Open("schema.txt") // permet d'ouvir le fichier txt en appelant file
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()// ferme le fichier quand le main est fini 

	compteurLiens := 0

	scanner := bufio.NewScanner(file)

	scanner.Scan() //On lit la première ligne du doc pour avoir le nombre de liens avec .scan()
	//fmt.Println(scanner.Text())

	n, err := strconv.Atoi(scanner.Text()) //On le convertie en entier pour pouvoir creer nos slice de slice : "tableau de liens"
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	liens := make([][]string, n) // notre slice de slice
	
	//On remplit un slice avec les différents liens
	for scanner.Scan() { // on parcours tous le fichier ligne par ligne  
		
		lien := strings.Split(scanner.Text(), " ")// on lit chaque ligne avec la fonction scanner.text() et on créer un slice de string avec la fonction strings.Split avec comme indication " " comme séparateur  
		liens[compteurLiens] = lien // on associe pour chaque item de la "slice de slice" "un lien"(sous forme d'un slice) du fichier texte dans l'ordre chronologique 
		compteurLiens++

	}
	//fmt.Println(liens)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	
	m := make(map[string][][]string)// création d'une map
	//On remplit la map avec les élément de notre " slice de slice "  
	for _, lien := range liens { // on parcours chaque "lien" de "liens" (slice de slice) 
		//Idée : On parcours tout liens, dans chaque lien on lit les 2 premiers item qui sont les sommets
		//Ensuite on teste s'ils ont déja été créé dans dans la map si non on les crées, ces sommets seront les clés de la map   
		if m[lien[0]] == nil { 
			m[lien[0]] = make([][]string, 0)// on associe une clé à un tableau 
		}
		if m[lien[1]] == nil { 
			m[lien[1]] = make([][]string, 0)
		}
		//Enfin on rempli chaque tableau de la clé avec le lien associé(càd l'iteration en cours)  
		m[lien[0]] = append(m[lien[0]], lien)// append permet de pouvoir aggrandir la taille de ce tableau en lui mettant un lien   
		m[lien[1]] = append(m[lien[1]], lien)
		// exemple : le tableau de "B" va être creer à l'itération 1, ensuite il va être rempli avec le lien 1 ET à l'itération 2 et 3 on va l'agrendir en lui mettant le lien 2 et le lien 3
		// B :[[lien1][lien2][lien3]]
	}

	//fmt.Println(m)
	var Depart ="A"
	var Arrive ="S"	
	var p=40
	var k=0
	Q := make([][]string, 12)
	for i:= 0; i < len(Q); i++{
		Q[i]=[]string{"0", "0", "0","0"}
	}
	for {
		
		for key, element := range m{
			//fmt.Println(key)
			if key == Depart{
				p=40
				
				//println("nombre d'element :",len(element))
				for i:= 0; i < len(element); i++{
					var s = 0
					s, err := strconv.Atoi(element[i][3])
					if err != nil {
						fmt.Println(err)
						os.Exit(2)
					}
					if s<p  {
						if k == 0{	
							p = s
							Q[k][0]= element[i][2]
							Q[k][1]= element[i][1]
							Q[k][2]= element[i][3]
							Q[k][3]= element[i][0]
							Depart = element[i][1]
						}
						
						if k == 1{
					//		println("iteration i: ",i)
					//		println("iteration k: ",k)
					//		println("ancient p: ",p)
					//		println(element[i][0])
					//		println(Q[k-1][3])
							if element[i][0] != Q[k-1][3] {
								p = s
					//			println("nouveau p: ",p)								
								Q[k][0]= element[i][2]
								Q[k][1]= element[i][1]
								Q[k][2]= element[i][3]
								Q[k][3]= element[i][0]
								Depart = element[i][1]
							}
						}
						if k > 1{
							println("iteration i: ",i)
							println("iteration k: ",k)
							println("ancient p: ",p)
							println(element[i][0])
							println(Q[k-1][3])
							if element[i][0] != Q[k-1][3] && element[i][0] != Q[k-2][3]{
								p = s
								println("nouveau p: ",p)								
								Q[k][0]= element[i][2]
								Q[k][1]= element[i][1]
								Q[k][2]= element[i][3]
								Q[k][3]= element[i][0]
								Depart = element[i][1]
							}
						}								
					}

				}
				k++
				fmt.Println(Q)
				fmt.Println("noeud suivant:",Depart)	
			}
			if Depart == Arrive{
				fmt.Println(Q)
				println("vous êtes arrivé à destination")
				break
			}
				

		}
		
		if Depart == Arrive{
			break
		}
	
	}
}
//***********************************************************************************
//
//***********************************************************************************
//	Objectif
//go routine
//entree : sommet du graphe
//sortie : les chemins les plus courts vers les autres sommets
//on lance une routine par sommet dans m
