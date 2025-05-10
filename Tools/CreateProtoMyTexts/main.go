package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var linebuf strings.Builder

func main() {
	//	prototechName := "ESPrototech"

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	grow := false
	flush := false
	keep := false
	skipEmptyLine := false
	line := ""
	for scanner.Scan() {
		line = scanner.Text()
		switch {

		case strings.Contains(line, "<data"):
			{
				grow = true
				if strings.Contains(line, "Description_HeavyArmor") {
					keep = true
					line = strings.Replace(line, "_Heavy", "_ESPrototech", 1)
					break
				}

				if strings.Contains(line, "DisplayName_Block") && strings.Contains(line, "Heavy") {
					keep = true
					line = strings.Replace(line, "Block_", "Block_ESPrototech", 1)
					line = strings.Replace(line, "BlockGroup_", "BlockGroup_ESPrototech", 1)
					line = strings.Replace(line, "Heavy", "", 1)
				}
			}
		case strings.Contains(line, "</data>"):
			{
				flush = true
			}
		case strings.Contains(line, "<value>") && strings.Contains(line, "Heavy"):
			{
				line = strings.Replace(line, "Heavy", "Prototech", 1)
			}
		}

		if grow {
			linebuf.WriteString(fmt.Sprintf("%s\r\n", line))
		} else {

			if len(strings.TrimSpace(line)) > 0 {
				skipEmptyLine = false
			}
			if !skipEmptyLine {
				fmt.Printf("%s\r\n", line)
			}

			if len(strings.TrimSpace(line)) == 0 {
				skipEmptyLine = true
			}
			continue
		}

		if keep && flush {
			fmt.Print(linebuf.String())
			linebuf.Reset()
			grow = false
			flush = false
			keep = false
		}

		if !keep && flush {
			linebuf.Reset()
			grow = false
			flush = false
			keep = false
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
