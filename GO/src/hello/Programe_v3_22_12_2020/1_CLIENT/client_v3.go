package main

import (

	"bufio"
	"fmt"
	"net"
	"os"
	"log"
	"strconv"
	"strings"
	"encoding/gob"
	"time"
)
// fonction gestionnaire d'erreur
func gestionErreur(err error) {
	if err != nil {
		panic(err)
	}
}
// creation d'un type pour l'encodage de la donné à envoyer au serveur 
type Maping struct {
    M map[string][][]string
}

// creation d'un type pour le decodage de la donné envoyé par le serveur 
type Caping struct {
    C map[string]map[string][]string
}

// Constante utiliser pour deffinir l'adress IP et le PORT 
const (
	IP   = "127.0.0.01" // IP local
	PORT = "3569"       // Port utilisé
)
var C = &Caping{}
func main() {

//**********************************************************************************************************************************
// Partie du propro qui va: -> 1er: lire un fichier	 -> 2nd: va prendre chaque ligne du fichier texte pour le mettre dans slice de slice 
// -> 3rd: transformer ce slice de slice en map avec key : sommet et element : liens qui sont connecté
  
	file, err := os.Open("schema.txt") // permet d'ouvir le fichier txt en appelant file
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // ferme le fichier quand le main est fini

	compteurLiens := 0

	scanner := bufio.NewScanner(file)

	scanner.Scan() //On lit la première ligne du doc pour avoir le nombre de liens avec .scan()


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
fmt.Println("*********ProproClient*********")
fmt.Println("\nmap : <<sommet et liens>> a envoyer : \n" ,m)	
	
//************************************************************************************************************************



//************************************************************************************************************************
	// Connexion au serveur
	//Initialisation de la connection
	fmt.Println("\nTentative de connection au serveur")
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	if err != nil {fmt.Println("Connection à échouer") 
	} else { 
		// s'il y a pas de problème
		fmt.Println("Connection établie") 
		
		// envoie des donné au server puis reçois réponse du server s'il y a eu aucune réponse du server le Client renvoit le message en boucle
		for{
			// Encodage de la variable à envoyer en tcp  
			encoder := gob.NewEncoder(conn)
			M := &Maping{m} // M est la variable encoder
			encoder.Encode(M)// on l'envoie au serveur 
			//on attend la réponse du serveur 
			buffer := make([]byte, 4096)       // taille maximum du message qui sera envoyé par le client
			length, er := conn.Read(buffer) 
			message := string(buffer[:length]) // réponse du serveur 
			
			if er == nil { 
				fmt.Print("serveur : " + message)
				dec := gob.NewDecoder(conn)
				
				dec.Decode(C)
				fmt.Printf("\n\nPaquet reçu du serveur : %+v\n", C.C);
				conn.Write([]byte("Copy that Server, Over..."))
				conn.Close()// on ferme la connexion
				break // on break la boucle si on a bien reçu ACK
			} else{
				time.Sleep(3000 * time.Millisecond) // on attend 3 sec 
				fmt.Print("\nfatal error : pas ACK\nRenvoie du message")// si on n'a de message du serveur 
			}
			gestionErreur(er) // erreur si le serveur n'envoie pas d'aquitement 
		}
	}
	gestionErreur(err)	
//***************************************************************************************************************************
	
//***************************************************************************************************************************
// traitement des informations reçues	
	fmt.Println("\n\nTraitement des résultats : \n");
	traitement_final(C.C)
	fmt.Println("\n\n*****Travaille terminé******");
//***************************************************************************************************************************	
}
//fonction qui affiche en beau le résultat reçu du server
func traitement_final ( C map[string]map[string][]string){
	for key, elem :=range C {
		fmt.Println("Depart : ", key)
			for k := range elem{
			fmt.Println("Arrivé : ", k , " , Chemin : ", C[key][k])
			}	
	}	
}

