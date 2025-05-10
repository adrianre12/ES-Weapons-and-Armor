package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var linebuf strings.Builder

func main() {

	prototechName := "ESPrototech"

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

		case strings.Contains(line, "<Definition>"):
			{
				grow = true
			}
		case strings.Contains(line, "</Definition>"):
			{
				flush = true
			}
		case strings.Contains(line, "<SubtypeId>"):
			{
				line = strings.Replace(line, "Heavy", prototechName, 1)
			}
		case strings.Contains(line, "<DisplayName>") && strings.Contains(line, "Heavy"):
			{
				keep = true
				line = strings.Replace(line, "Heavy", prototechName, 1)
				line = ElementAddLoc(line)
			}
		case strings.Contains(line, "<Icon>"):
			{
				linebuf.WriteString(fmt.Sprintf("%s\r\n", line)) //Kludge to add extra icon rows
				s := strings.Index(line, "<")
				indent := line[:s]
				//linebuf.WriteString(fmt.Sprintf("%s%s\r\n", indent, `<Icon>Textures\GUI\Icons\Overlays\Weight.dds</Icon>`))
				line = fmt.Sprintf("%s%s", indent, `<Icon>Textures\GUI\Icons\Overlays\PrototechArmor.dds</Icon>`)
			}
		case strings.Contains(line, "<Description>"):
			{
				line = strings.Replace(line, "Heavy", prototechName, 1)
				line = ElementAddLoc(line)
			}
		case strings.Contains(line, "<Component "):
			{
				line = ConvertComponent(line, 0.633, 0.6)
			}
		case strings.Contains(line, "<CriticalComponent"):
			{
				line = strings.Replace(line, "SteelPlate", "PrototechPanel", 1)

			}
		case strings.Contains(line, "<BlockPairName>"):
			{
				line = strings.Replace(line, "Heavy", prototechName, 1)
			}
		case strings.Contains(line, "<BuildTimeSeconds>"):
			{
				s := strings.Index(line, ">") + 1
				e := strings.Index(line[s:], "<")
				i, err := strconv.ParseFloat(line[s:s+e], 64)
				if err != nil {
					log.Fatal("strconv failed %w", err)
				}
				line = fmt.Sprintf("%s%s%s", line[:s], strconv.FormatFloat(math.Ceil(float64(i)*1.5), 'f', -1, 64), line[s+e:])
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

func ElementAddLoc(line string) string {
	s := strings.Index(line, ">") + 1
	e := strings.Index(line[s:], "<") + s

	return fmt.Sprintf("%s{LOC:%s}%s", line[:s], line[s:e], line[e:])
}

func ElementPrefix(line string) string {
	s := strings.Index(line, ">") + 1

	return fmt.Sprintf("%s%s%s", line[:s], "Prototech", line[s:])
}

func ConvertComponent(line string, protoRatio float64, otherRatio float64) string {
	if strings.Contains(line, "SteelPlate") {
		line = strings.Replace(line, "SteelPlate", "PrototechPanel", 1)
		c := strings.Index(line, "Count=") + 7
		e := strings.Index(line[c:], "\"")
		i, err := strconv.ParseFloat(line[c:c+e], 64)
		if err != nil {
			log.Fatal("strconv failed %w", err)
		}
		line = fmt.Sprintf("%s%d%s", line[:c], int(math.Ceil(i*protoRatio)), line[c+e:])
	} else {
		c := strings.Index(line, "Count=") + 7
		e := strings.Index(line[c:], "\"")
		i, err := strconv.ParseFloat(line[c:c+e], 64)
		if err != nil {
			log.Fatal("strconv failed %w", err)
		}
		line = fmt.Sprintf("%s%d%s", line[:c], int(math.Ceil(i*otherRatio)), line[c+e:])
	}

	return line
}
