package main

import (
	"encoding/json"
	"encoding/csv"
    "fmt"
    //"math/rand"
    "net"
    "os"
	//"time"
	"io"
	"strconv"
)

type Msg struct {
    Command     string	`json:"Command"`
    Hostname    string	`json:"Hostname"`
	List        []string `json:"List"`
	Informacion Info	 `json:"Informacion"`
	UltHa		UltimoHash `json:"UltHa"`
}

type Info struct {
    Id float64  `json:"Id"`
    Departamento float64    `json:"Departamento"`
    Edad float64    `json:"Edad"`
    Sexo float64    `json:"Sexo"`
    Peso float64    `json:"Peso"`
    Altura float64  `json:"Altura"`
    Temperatura float64 `json:"Temperatura"`
    TosSeca float64 `json:"TosSeca"`
    Cansancio float64   `json:"Cansancio"`
    Infectado float64   `json:"Infectado"`
}


type Block struct {
    Hashanterior float64
    Hash float64
    Datosdelblock Info
}

type UltimoHash struct {
    UltHash     float64  `json:"UltHash"`
}


var blockchain []Block




func ReadCsv() {
    f, err := os.Open("coviddata.csv")
    if err != nil {
        fmt.Println("Error al abrir el csv")
    }
    defer f.Close()

    r := csv.NewReader(f)
    r.Comma = ','
    r.Comment = '#'
    r.FieldsPerRecord = 10

    var informacion []Info
    //var blockchain []Block
    cont := 0
    for {
        record , err := r.Read()
        if err == io.EOF {
            break
        }

        f1, _ := strconv.ParseFloat(record[0], 64)
        f2, _ := strconv.ParseFloat(record[1], 64)
        f3, _ := strconv.ParseFloat(record[2], 64)
        f4, _ := strconv.ParseFloat(record[3], 64)
        f5, _ := strconv.ParseFloat(record[4], 64)
        f6, _ := strconv.ParseFloat(record[5], 64)
        f7, _ := strconv.ParseFloat(record[6], 64)
        f8, _ := strconv.ParseFloat(record[7], 64)
        f9, _ := strconv.ParseFloat(record[8], 64)
        f10, _ := strconv.ParseFloat(record[9], 64)

        i := Info {
            Id: f1,
            Departamento: f2,
            Edad: f3,
            Sexo: f4,
            Peso: f5,
            Altura: f6,
            Temperatura: f7,
            TosSeca: f8,
            Cansancio: f9,
            Infectado: f10,
        }

        informacion = append(informacion, i)
        


        if cont == 0 {
            hashaux := i.Id + i.Departamento + i.Edad +i.Sexo +i.Peso +i.Altura +i.Temperatura +i.TosSeca+ i.Cansancio + i.Infectado

            bloque := Block {
                Hashanterior: 321.123,
                Hash: hashaux,
                Datosdelblock: i,

            }
            blockchain = append(blockchain,bloque)
        } else {
            hashaux := i.Id + i.Departamento + i.Edad +i.Sexo +i.Peso +i.Altura +i.Temperatura +i.TosSeca+ i.Cansancio + i.Infectado

            bloque := Block {
                Hashanterior: blockchain[cont-1].Hash,
                Hash: hashaux,
                Datosdelblock: i,

            }
            blockchain = append(blockchain,bloque)
        }

        
       
        



        cont++
    }

    //fmt.Println(informacion)
    //fmt.Println(blockchain)

}





var friends         []string
var local           string
var end             chan bool

var ready2listen    chan bool
//var espera    chan bool
var decisions       map[string]string
var cont            int

func serv() {
    // fmt.Println("(", local, ")")
    ln, _ := net.Listen("tcp", local)
    defer ln.Close()
    for {
        conn, _ := ln.Accept()
        go handle(conn)
    }
}
func handle(conn net.Conn) {
    defer conn.Close()
    dec := json.NewDecoder(conn)
    var msg Msg
    if err := dec.Decode(&msg); err == nil {
        switch msg.Command {
        case "hello":
            resp := Msg{"hey", local, friends, Info{}, UltimoHash{}}
            enc := json.NewEncoder(conn)
            //fmt.Println("ASFSDFASDFASD")
            if err := enc.Encode(&resp); err == nil {
                for _, friend := range friends {
                    // fmt.Println(local, friend, "meet", msg.Hostname)
                    send(friend, "meet new friend", []string{msg.Hostname}, Info{}, UltimoHash{})
                }
            }
            friends = append(friends, msg.Hostname)
            // fmt.Println(local, "updated list", friends)
        case "meet new friend":
            friends = append(friends, msg.List...)
            // fmt.Println(local, "new friend", msg.List)
        case "Aviso":
            var Nuevainfo Info
			Nuevainfo = msg.Informacion
			Nuevainfo.Id = float64(len(blockchain)+1)
			

			hashaux := Nuevainfo.Id + Nuevainfo.Departamento + Nuevainfo.Edad +Nuevainfo.Sexo +Nuevainfo.Peso +Nuevainfo.Altura +Nuevainfo.Temperatura +Nuevainfo.TosSeca+ Nuevainfo.Cansancio + Nuevainfo.Infectado
            fmt.Println(hashaux)
			bloque := Block {
				Hashanterior: blockchain[len(blockchain)-1].Hash,
				Hash: hashaux,
				Datosdelblock: Nuevainfo,

            }
            fmt.Println(bloque)

			blockchain = append(blockchain, bloque)

			//fmt.Println(blockchain)

            for _, friend := range friends {
                send(friend, "Agregar", []string{}, msg.Informacion,UltimoHash{})
                //send(friend, "Decision", []string{}, msg.Informacion, decihash)
            }

		case "Agregar":
			var Nuevainfo Info
			Nuevainfo = msg.Informacion
			Nuevainfo.Id = float64(len(blockchain)+1)
			

			hashaux := Nuevainfo.Id + Nuevainfo.Departamento + Nuevainfo.Edad +Nuevainfo.Sexo +Nuevainfo.Peso +Nuevainfo.Altura +Nuevainfo.Temperatura +Nuevainfo.TosSeca+ Nuevainfo.Cansancio + Nuevainfo.Infectado
            fmt.Println(hashaux)
			bloque := Block {
				Hashanterior: blockchain[len(blockchain)-1].Hash,
				Hash: hashaux,
				Datosdelblock: Nuevainfo,

			}
            fmt.Println(bloque)
			blockchain = append(blockchain, bloque)

			//fmt.Println(blockchain)

			//fmt.Println(friends)

			//fmt.Println(conn)

			decihash := UltimoHash {
				UltHash: bloque.Hash,
			}

			cont = 0
			//fmt.Println(friends)

			// for _, friend := range friends {
			// 	send(friend, "Agregar", []string{}, msg.Informacion)
			// }

			// if len(os.Args) == 3 {
			// 	remote := os.Args[2]
			// 	friends = append(friends, os.Args[2])
			// 	//send(remote, "hello", []string{}, Info{})
			// 	send(remote, "Agregar", []string{}, msg.Informacion,UltimoHash{})
			// 	send(remote, "Decision", []string{}, msg.Informacion, decihash)
            // }
            
            for _, friend := range friends {
                //send(friend, "Agregar", []string{}, msg.Informacion,UltimoHash{})
                send(friend, "Decision", []string{}, msg.Informacion, decihash)
            }

			//ready2listen<-true

		case "Decision":
			//<-ready2listen
			// fmt.Println("ENTROOOO A DECIDIR")
            cont++
            // fmt.Println(friends)
            fmt.Println(cont)
			
			
			// fmt.Println(cont)
			// fmt.Println(msg.UltHa)
			if cont == len(friends) {

                fmt.Println(local, "AGREGADOOO!!!!")
                for i, friend := range friends {
                    if i == 0 {
                        send(friend, "Decision", []string{}, msg.Informacion, msg.UltHa)
                    }
                }
                cont = 0
                //espera<-true

                
            } else {
                //ready2listen<-true
            }
			

		
        // case "test consensus":
        //     if rand.Intn(100) % 2 == 0 {
        //         decisions[local] = "atacar"
        //     } else {
        //         decisions[local] = "retirada"
        //     }
        //     fmt.Println(local, decisions[local])
        //     cont = 0
        //     for _, friend := range friends {
        //         send(friend, "decision", []string{decisions[local]})
        //     }
        //     ready2listen<-true
        // case "decision":
        //     <-ready2listen
        //     decisions[msg.Hostname] = msg.List[0]
        //     cont++
        //     if cont == len(friends) {
        //         contAtack := 0
        //         contFallb := 0
        //         for _, decision := range decisions {
        //             if decision == "atacar" {
        //                 contAtack++
        //             } else {
        //                 contFallb++
        //             }
        //         }
        //         if contAtack < contFallb {
        //             fmt.Println(local, "RETIRADA!!")
        //         } else {
        //             fmt.Println(local, "ATACAR!!!!")
        //         }
        //         end<-true
        //     } else {
        //         ready2listen<-true
        //     }

        // case "finish":
        //     end<-true
        }
    }
}
func send(remote, command string, list []string, inf Info, ulha UltimoHash) {
    conn, _ := net.Dial("tcp", remote)
    defer conn.Close()
    msg := Msg{command, local, list, inf, ulha}
    enc := json.NewEncoder(conn)
    if err := enc.Encode(&msg); err == nil {
        // fmt.Println(local, "sent", msg)
        if command == "hello" {
            dec := json.NewDecoder(conn)
            var resp Msg
            if err := dec.Decode(&resp); err == nil {
                fmt.Println(resp.List)
                friends = append(friends, resp.List...)
                // fmt.Println(local, "recibí", resp.List)
            }
        }
    }
}
func main() {
    ready2listen    = make(chan bool)
    //espera          = make(chan bool)
    end             = make(chan bool)
    local           = os.Args[1]
    decisions       = make(map[string]string)

	//rand.Seed(time.Now().UTC().UnixNano())
	ReadCsv()
    go serv()

    if len(os.Args) == 3 {
        remote := os.Args[2]
        friends = append(friends, os.Args[2])
        send(remote, "hello", []string{}, Info{},UltimoHash{})
    }

    <-end
    // fmt.Println(local, "time to die")
    // fmt.Println(local, friends)
}