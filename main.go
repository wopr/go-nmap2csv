package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"nmap2csv/decodeXML"
	"os"
	"strings"
)

type Hosts struct {
	XMLName xml.Name `xml:"nmaprun"`
	Hosts   []Host   `xml:"host"`
}

type Host struct {
	XMLName xml.Name `xml:"host"`
	HStatus Status   `xml:"status"`
	IPs     []IP     `xml:"address"`
	Ports   []Port   `xml:"ports>port"`
}

type Status struct {
	XMLName xml.Name `xml:"status"`
	State   string   `xml:"state,attr"`
}

type IP struct {
	XMLName xml.Name `xml:"address"`
	IP      string   `xml:"addr,attr"`
	Type    string   `xml:"addrtype,attr"`
}

type Port struct {
	XMLName xml.Name `xml:"port"`
	Proto   string   `xml:"protocol,attr"`
	ID      string   `xml:"portid,attr"`
	State   State    `xml:"state"`
	Service Service  `xml:"service"`
	Scripts []Script `xml:"script"`
}

type State struct {
	XMLName xml.Name `xml:"state"`
	State   string   `xml:"state,attr"`
}

type Service struct {
	XMLName xml.Name `xml:"service"`
	Name    string   `xml:"name,attr"`
	Product string   `xml:"product,attr"`
	Version string   `xml:"version,attr"`
}

type Script struct {
	XMLName xml.Name `xml:"script"`
	ID      string   `xml:"id,attr"`
	Output  string   `xml:"output,attr"`
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -i <src filename> -o <dest filename> [-s=<true|false>]\n", os.Args[0])
		flag.PrintDefaults()
	}
	infile := flag.String("i", "", "Enter a source XML file to be converted.")
	outfile := flag.String("o", "", "Enter a destination file.")
	scripts := flag.Bool("s", true, "Display script data [true|false].")
	flag.Parse()

	if *infile == "" || *outfile == "" {
		flag.Usage()
	}

	hostsStruct := new(Hosts)
	decodeXML.DecodeXML(hostsStruct, infile)

	file, err := os.Create(*outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	wb := bufio.NewWriter(file)

	for _, host := range hostsStruct.Hosts {
		if host.HStatus.State == "up" {
			wb.WriteString(fmt.Sprintln(host.IPs[0].IP))

			for _, port := range host.Ports {
				wb.WriteString(fmt.Sprintf(",%s,%s,%s,%s,%s,%s\n", port.Proto, port.State.State, port.ID, port.Service.Name, port.Service.Product, port.Service.Version))

				if *scripts == true {
					for _, script := range port.Scripts {
						if script.ID != "fingerprint-strings" {
							str := strings.Replace(script.Output, ",", "", -1)
							wb.WriteString(fmt.Sprintf(",,,,,,script: %s,%s\n", script.ID, strings.TrimSpace(str)))
						}
					}
				}
			}
		}
		wb.WriteString(fmt.Sprintln())
	}
	wb.Flush()
}
