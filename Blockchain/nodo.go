package main

import (
    "encoding/json"
    "fmt"
    //"math/rand"
    "net"
    "os"
    "encoding/csv"
    "strconv"
    //"time"
    "io"
)

type Msg struct {
    Command     string
    Hostname    string
    List        []string
}

var friends         []string
var local           string
var end             chan bool

var ready2listen    chan bool
var decisions       map[string]string
// var cont            int


type Info struct {
    Id float64
    Departamento float64
    Edad float64
    Sexo float64
    Peso float64
    Altura float64
    Temperatura float64
    TosSeca float64
    Cansancio float64
    Infectado float64
}

type Block struct {
    Hashanterior float64
    Hash float64
    Datosdelblock Info
}




type ProcessData struct {
    Uno     string `json:"Uno"`
    Dos     string `json:"Dos"`
}


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
    var blockchain []Block
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



func serv() {
    // fmt.Println("(", local, ")")
    ln, _ := net.Listen("tcp", local)
    fmt.Println(local)
    fmt.Println(ln)
    defer ln.Close()
    for {
        conn, _ := ln.Accept()
        fmt.Println("Entro")
        fmt.Println(conn)
        go handle(conn)
    }
}
func handle(conn net.Conn) {
    defer conn.Close()
    dec := json.NewDecoder(conn)
    fmt.Println(dec)
    var data ProcessData
    das := ProcessData{Uno:"a", Dos:"b"}
    ReadCsv()
    if err := dec.Decode(&data); err == nil {
        fmt.Println(err)
        fmt.Println(data)
        fmt.Println(das)
    }
    // var msg Msg
    // if err := dec.Decode(&msg); err == nil {
    //     switch msg.Command {
    //     case "hello":
    //         resp := Msg{"hey", local, friends}
    //         enc := json.NewEncoder(conn)
    //         if err := enc.Encode(&resp); err == nil {
    //             for _, friend := range friends {
    //                 // fmt.Println(local, friend, "meet", msg.Hostname)
    //                 send(friend, "meet new friend", []string{msg.Hostname})
    //             }
    //         }
    //         friends = append(friends, msg.Hostname)
    //         // fmt.Println(local, "updated list", friends)
    //     case "meet new friend":
    //         friends = append(friends, msg.List...)
    //         // fmt.Println(local, "new friend", msg.List)

    //     case "test consensus":
    //         if rand.Intn(100) % 2 == 0 {
    //             decisions[local] = "atacar"
    //         } else {
    //             decisions[local] = "retirada"
    //         }
    //         fmt.Println(local, decisions[local])
    //         cont = 0
    //         for _, friend := range friends {
    //             send(friend, "decision", []string{decisions[local]})
    //         }
    //         ready2listen<-true
    //     case "decision":
    //         <-ready2listen
    //         decisions[msg.Hostname] = msg.List[0]
    //         cont++
    //         if cont == len(friends) {
    //             contAtack := 0
    //             contFallb := 0
    //             for _, decision := range decisions {
    //                 if decision == "atacar" {
    //                     contAtack++
    //                 } else {
    //                     contFallb++
    //                 }
    //             }
    //             if contAtack < contFallb {
    //                 fmt.Println(local, "RETIRADA!!")
    //             } else {
    //                 fmt.Println(local, "ATACAR!!!!")
    //             }
    //             //end<-true
    //         } else {
    //             ready2listen<-true
    //         }

    //     case "finish":
    //         //end<-true
    //     }
    // }
}
func send(remote, command string, list []string) {
    conn, _ := net.Dial("tcp", remote)
    defer conn.Close()
    msg := Msg{command, local, list}
    enc := json.NewEncoder(conn)
    if err := enc.Encode(&msg); err == nil {
        fmt.Println(local, "sent", msg)
        
    }
}
func main() {
    ready2listen    = make(chan bool)
    end             = make(chan bool)
    local           = os.Args[1]
    decisions       = make(map[string]string)

    // rand.Seed(time.Now().UTC().UnixNano())
    go serv()

    if len(os.Args) == 3 {
        remote := os.Args[2]
        friends = append(friends, os.Args[2])
        send(remote, "hello", []string{})
    }

    <-end
    fmt.Println(local, "time to die")
    fmt.Println(local, friends)
}