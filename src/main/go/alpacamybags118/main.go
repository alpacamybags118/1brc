package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type LocationData struct {
	min   float64
	sum   float64
	max   float64
	count int
}

func main() {
	start := time.Now()
	file, err := os.Open("../../../../measurements.txt")
	//file, err := os.Open("sample.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	cache := map[string]LocationData{}

	readStart := time.Now()
	for {
		buffer := make([]byte, 240000)
		readBytes, err := file.Read(buffer)
		if readBytes > 0 {
			buffer = buffer[:readBytes]
			// need to read up to the next newline to process records
			if buffer[readBytes-1] != '\n' {
				for {
					charBuff := make([]byte, 1)
					dataRead, err := file.Read(charBuff)

					if dataRead > 0 {
						buffer = append(buffer, charBuff...)
					} else {
						break
					}

					if err != nil {
						if err != io.EOF {
							fmt.Println(err)
						}
						fmt.Println("EOF")
						break
					}

					if charBuff[0] == '\n' {
						break
					}
				}
			}

			lineList := bytes.Split(buffer, []byte("\n"))
			lineList = lineList[:len(lineList)-1]
			// process each line into our cache
			for _, line := range lineList {
				data := strings.Split(string(line), ";")
				temp, err := strconv.ParseFloat(data[1], 32)

				if err != nil {
					log.Fatal(err)
				}

				current, exist := cache[data[0]]
				var location LocationData

				if exist {
					location = LocationData{
						min:   current.min,
						sum:   current.sum + temp,
						max:   current.max,
						count: current.count + 1,
					}

					if current.min > temp {
						location.min = temp
					}

					if current.max < temp {
						location.max = temp
					}
				} else {
					location = LocationData{
						min:   temp,
						sum:   temp,
						max:   temp,
						count: 1,
					}
				}

				cache[data[0]] = location
			}
		}

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}

			fmt.Println("EOF")
			break
		}

	}
	readDuration := time.Since(readStart)

	fmt.Printf("read time: %v\n", readDuration.Seconds())

	calcTime := time.Now()

	keys := make([]string, 0)

	for k := range cache {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := cache[k]
		fmt.Printf("%s: min=%.1f  mean=%.1f  max=%.1f \n", k, v.min, v.sum/float64(v.count), v.max)
	}

	calcDuration := time.Since(calcTime)

	fmt.Println()
	fmt.Printf("Time to calculate: %v\n", calcDuration.Seconds())

	duration := time.Since(start)
	fmt.Println()

	fmt.Printf("Total elapsed time: %v\n", duration.Seconds())
}

func Mean(values []float64) float64 {
	sum := 0.0

	for _, v := range values {
		sum += v
	}

	return sum / float64(len(values))
}
