package main

import (
    "fmt"
    "net"
    "encoding/gob"
    "sync"
    "strconv"
    "os"
    "time"
)
//*****************************************************************************************
//Creation var, type, const utiliser dans le programe
 
// Constante utiliser pour deffinir l'adress IP et le PORT 
const (
    IP   = "127.0.0.01" // IP local
    PORT = "3569"       // Port utilisé
)

// creation d'un type pour le decodage de la donné envoyé par le client 
type Maping struct {
    M map[string][][]string
}

// creation d'un type pour le decodage de la donné envoyé par le client 
type Caping struct {
    C map[string]map[string][]string
}

// Creation du type Chemin qui va être utiliser dans la fonction dijkstra
type Chemin struct {
	poids           int      //poids du meilleur chemin pour l'instant
	sommetPrecedent string   //sommet par lequel on vient
	fini            bool     //true si on est sur que le chemin est optimal
	path            []string //itinéraire sommet par sommet
}

// Creation du type graphe qui va être utiliser dans la fonction dijstra 
type Graphe map[string]*Chemin

//création de la map(var global) qui accueillera les résultats :
// cette map va prendre les resultat ap la fonction dikjstra 
var mPlusCourtsChemins map[string]Graphe
//cette map va prendre le resultat final après filtrage de l'information de Graphe  
var resultat map[string]map[string][]string

// var utilisé pour le wait Group des go routines :
//le wait group pour gestionnaire server-client
var wo sync.WaitGroup
//le wait group pour la fonction dijkstra
var wg sync.WaitGroup
//************************************************************************************************


//****************************************************************************************************************************************
//fonction qui permet d'initialiser mPlusCourtsChemins et resultat en fonction de la taille de la map envoyer par le client
func initialisation ( m map[string][][]string){
	mPlusCourtsChemins = make(map[string]Graphe)
	resultat = make(map[string]map[string][]string)
	for key := range m {
		resultat[key] = make(map[string][]string, 0)
		mPlusCourtsChemins[key] = make(Graphe, 1)
			
	}
	
}

// fonction gestionnaire d'erreur
func gestionErreur(err error) {
    if err != nil {
        panic(err)
    }
}
//*********************************************************************************************************************************



//*********************************************************************************************************************************************
// fonction de l'alogorithme dijkstra appeler par une go routine
func dijkstra(sommetDepart string, liens map[string][][]string, wg *sync.WaitGroup) {
	

	defer wg.Done()
	const inf = 999 //On suppose que les poids seront inférieurs à cette valeur

	//Initialisation
	tableCouts := make(Graphe)

	for key, _ := range liens {

		tableCouts[key] = &Chemin{inf, "", false, make([]string, 0)}

		if key == sommetDepart {
			tableCouts[key].poids = 0
		}
	}
	var prochainSommet string
	var continuer bool
	var cheminLeger int
	var poidsCheminActuel int
	continuer = true
	poidsCheminActuel = 0

	//Condition d'arrêt : Tous les chemins sont finaux

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
		

		poidsCheminActuel = tableCouts[prochainSommet].poids
		//le sommet atteint est fixé : son chemin le plus court est compris dans la variable itineraire
		tableCouts[prochainSommet].fini = true

		
		for _, lien := range liens[prochainSommet] {
			//On distingue le voisin du sommet étudié (l'ordre des sommets n'est pas toujours le même)
			var voisin string
			if voisin = lien[0]; lien[0] == prochainSommet {
				voisin = lien[1]
			}

			if !tableCouts[voisin].fini {
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
					tableCouts[voisin].sommetPrecedent = prochainSommet
				}

			}

		}
		
	}
	//Calcul du path et affichage des résultats
	for key, elem := range tableCouts {
		tableCouts[key].path = append(tableCouts[key].path, key)
		step := elem.sommetPrecedent
		for step != "" {
			tableCouts[key].path = append(tableCouts[key].path, step)
			step = tableCouts[step].sommetPrecedent
		}
	}
	mPlusCourtsChemins[sommetDepart] = tableCouts
	
}
//***************************************************************************************************************************************



//***************************************************************************************************************************************
// notre fonction qui va remplir notre map de map de string : "resultat", en prenant les paths comme element à interieur de la second map 
func remettre_les_choses_normalement ( m map[string]Graphe){
	for key, elem :=range m {
			for k,e := range elem{
			resultat[key][k]=e.path
			}	
	}	
}
//****************************************************************************************************************************************


//****************************************************************************************************************************************
// fonction  qui va permmetre de revoir le message encoder du client et le renvoyer
func handleConnection(conn net.Conn) {
	fmt.Println("\nUn client est connecté depuis", conn.RemoteAddr())
    dec := gob.NewDecoder(conn)
    M := &Maping{}
    dec.Decode(M)
    fmt.Printf("Paquet reçu du client : %+v\n", M.M);
    conn.Write([]byte("Copy that Client, Over..."))
    initialisation(M.M)
    //***************************************
    //fonction dijkstra en go routine
    // Utilisation du principe de go routine en executant simultanément plusieur fonction tout en evitant que l'execution des fonction se chevauche
	
	for sommet, _ := range M.M {
		wg.Add(1)
	//nous traitons l'ago dijkstra en go routine qui va nous donner en sortie une map de map de Graphe
		go dijkstra(sommet, M.M, &wg)
	}
	wg.Wait()
	fmt.Println("\nFonction dijkstra terminée")
	//*****************************************
	// en sortie des go routine nous avons une map de map de graphe dans lequel on va uniquement prendre l'info path de graph pour le mettre dans une map de map de string
	remettre_les_choses_normalement(mPlusCourtsChemins)
	//notre map de map de string 
	fmt.Println("\nNouvelle Map à envoyer, Resultat : \n",resultat)
	//*****************************************
	
	//*****************************************
	// envoie du résultat à client
	// envoie des donné au client puis reçois un ACK du client s'il y a eu aucune réponse du client le server renvoit le message en boucle
	for{
		encoder := gob.NewEncoder(conn)
		C := &Caping{resultat} // R est la variable encoder
		encoder.Encode(C)// on l'envoie au client
		//on attend la réponse du client 
		buffer := make([]byte, 4096)       // taille maximum du message qui sera envoyé par le client
		length, er := conn.Read(buffer) 
		message := string(buffer[:length]) // réponse du serveur
		
		if er == nil {
			fmt.Print("\nMessage envoyé")
			fmt.Print("\nClient : " + message)
			conn.Close()// on ferme la connexion
			fmt.Println("\n Le client: ",conn.RemoteAddr()," s'est déconnecté")
			fmt.Println("\n\n*****Travaille terminé******")
			break // on break la boucle si on a bien reçu ACK
		} else{
			time.Sleep(3000 * time.Millisecond) // on attend 3 sec 
			fmt.Print("\nfatal error : pas ACK\nRenvoie du message")// si on n'a de message du serveur 
		}
		gestionErreur(er) // erreur si le serveur n'envoie pas d'aquitement 
	}
	//******************************************
wo.Done()
}
//*****************************************************************************************************************************************


func main() {
	fmt.Println("*********ProproServer*********")
    fmt.Println("Lancement du serveur ...")

    // on écoute sur le port 3569
    ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, PORT))
    gestionErreur(err)
    
    for {
        conn, err := ln.Accept() // this blocks until connection or error
        if err != nil {
            continue
        }
        
        gestionErreur(err)
        wo.Add(1)
        go handleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
    }
}
