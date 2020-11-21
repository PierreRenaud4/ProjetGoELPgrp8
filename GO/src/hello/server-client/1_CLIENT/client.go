package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
	//"os"
)

func gestionErreur(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	IP   = "127.0.0.01" // IP local
	PORT = "3569"       // Port utilisé
)

func main() {

	// Connexion au serveur
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	gestionErreur(err)

	for {
		// entrée utilisateur
		//reader := bufio.NewReader(os.Stdin)
		fmt.Print("client: ")
		//text, err := reader.ReadString('\n')
		gestionErreur(err)

		// On envoie le message au serveur
		conn.Write([]byte("CouCOU"))
		time.Sleep(1000 * time.Millisecond)

		// On écoute tous les messages émis par le serveur et on rajouter un retour à la ligne
		message, err := bufio.NewReader(conn).ReadString('\n')
		gestionErreur(err)

		// on affiche le message utilisateur
		fmt.Print("serveur : " + message)
	}
}
