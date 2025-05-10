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
					line = strings.Replace(line, "_Heavy", "_PrototechHeavy", 1)
					break
				}

				if strings.Contains(line, "DisplayName_Block") && strings.Contains(line, "Heavy") {
					keep = true
					line = strings.Replace(line, "Block_", "Block_Prototech", 1)
				}
			}
		case strings.Contains(line, "</data>"):
			{
				flush = true
				//fmt.Printf("EndDef\n")
			}
		case strings.Contains(line, "<value>") && strings.Contains(line, "Heavy"):
			{
				//keep = true
				line = strings.Replace(line, ">Heavy", ">Prototech Heavy", 1)
				//fmt.Println(line)
				//fmt.Printf("Found\n")
			}

			// case strings.Contains(line, "<Description>"):
			// 	{
			// 		line = strings.Replace(line, "Description_", "Description_Prototech", 1)
			// 	}
			// case strings.Contains(line, "<Component "):
			// 	{
			// 		line = ConvertComponent(line, 0.6, 0.6)
			// 	}
			// case strings.Contains(line, "<CriticalComponent"):
			// 	{
			// 		line = strings.Replace(line, "SteelPlate", "PrototechPanel", 1)

			// 	}
			// case strings.Contains(line, "<BlockPairName>"):
			// 	{
			// 		line = ElementPrefix(line)
			// 	}
			// case strings.Contains(line, "<BuildTimeSeconds>"):
			// 	{
			// 		s := strings.Index(line, ">") + 1
			// 		e := strings.Index(line[s:], "<")
			// 		i, err := strconv.Atoi(line[s : s+e])
			// 		if err != nil {
			// 			log.Fatal("strconv failed %w", err)
			// 		}
			// 		line = fmt.Sprintf("%s%d%s", line[:s], int(math.Ceil(float64(i)*1.5)), line[s+e:])
			// 	}
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

// func ElementPrefix(line string) string {
// 	s := strings.Index(line, ">") + 1

// 	return fmt.Sprintf("%s%s%s", line[:s], "Prototech", line[s:])
// }

// func ConvertComponent(line string, protoRatio float64, otherRatio float64) string {
// 	if strings.Contains(line, "SteelPlate") {
// 		line = strings.Replace(line, "SteelPlate", "PrototechPanel", 1)
// 		c := strings.Index(line, "Count=") + 7
// 		e := strings.Index(line[c:], "\"")
// 		i, err := strconv.Atoi(line[c : c+e])
// 		if err != nil {
// 			log.Fatal("strconv failed %w", err)
// 		}
// 		line = fmt.Sprintf("%s%d%s", line[:c], int(math.Ceil(float64(i)*protoRatio)), line[c+e:])
// 	} else {
// 		c := strings.Index(line, "Count=") + 7
// 		e := strings.Index(line[c:], "\"")
// 		i, err := strconv.Atoi(line[c : c+e])
// 		if err != nil {
// 			log.Fatal("strconv failed %w", err)
// 		}
// 		line = fmt.Sprintf("%s%d%s", line[:c], int(math.Ceil(float64(i)*otherRatio)), line[c+e:])
// 	}

// 	return line
// }
