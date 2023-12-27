package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"log"

	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Data struct {
	NumOfWords      int  `json:"numOfWordsHeader"`
	NumbersWanted   bool `json:"numbersWantedHeader"`
	UppercaseWanted bool `json:"uppercaseWantedHeader"`
	SymbolWanted    bool `json:"symbolWantedHeader"`
}

type passwordStruct struct {
	Password string `json:"Password"`
}

var words []string

func readInput(r *http.Request) Data {
	// Parse the input from the front end and return a json object in the format of the Data struct
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var d Data
	json.Unmarshal(body, &d)

	return d

}

func readFile() []string {
	// Read the wordlist and extract the lines into the array words
	f, err := os.Open("wordlist.txt")

	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return words
}

func genRandomNumber(upperLimit int64) int64 {
	// generate a random number [1, upperlimit]
	rng, err := rand.Int(rand.Reader, big.NewInt(1).Add(big.NewInt(1), big.NewInt(upperLimit)))
	if err != nil {
		panic(err)
	}
	return rng.Int64()
}

func sendPassword(w http.ResponseWriter, passwordJson []uint8) {

	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(passwordJson)

}

func genPassword(w http.ResponseWriter, r *http.Request) {

	// var d Data
	d := readInput(r)
	wordList := readFile()

	numOfWords := d.NumOfWords
	numbersWanted := d.NumbersWanted
	uppercaseWanted := d.UppercaseWanted
	symbolWanted := d.SymbolWanted

	// fmt.Println(numOfWords, numbersWanted, uppercaseWanted, symbolWanted)

	var password string
	i := 0
	for i < numOfWords {
		randomWord := wordList[genRandomNumber(7776)]

		if uppercaseWanted {
			randomWord = strings.Title(randomWord)
		}

		password = password + randomWord

		if numbersWanted {
			password = password + strconv.FormatInt(genRandomNumber(9), 16)
		}
		i++
		// ADD SYMBOLS :DDDDD
	}
	if symbolWanted {
		const specialCharacters = "!$%^*(){}:@~<>?[];#,./"
		password = password + string(specialCharacters[genRandomNumber(int64(len(specialCharacters)))%int64(len(specialCharacters))])
	}

	passwordJson, err := json.Marshal(passwordStruct{password})
	if err != nil {
		panic(err)
	}

	fmt.Println("password", password, "passwordjson", passwordJson)

	sendPassword(w, passwordJson)
}

func main() {

	router := NewRouter()

	router.HandleFunc("/", rootMain).Methods("GET")

	router.HandleFunc("/", genPassword).Methods("POST")
	router.HandleFunc("/passGen/", genPassword).Methods("POST")

	http.ListenAndServe(":8000", router)

}

func rootMain(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("passphrase.html")
	if err != nil {
		log.Print(err)
		panic(err)
	}
	t.Execute(w, 0)

}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	return router
}