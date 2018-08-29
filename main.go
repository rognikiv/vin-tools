package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	InvalidVinLength = errors.New("Error: vin length is not equal 17 chars")
	InvalidVin       = errors.New("Error vin check failed")
)

type VIN [17]int32

func main() {
	fmt.Println("-- VIN check works on vehicles manufactured after 1980 --")
	for {
		var vin string
		fmt.Print("Enter VIN: ")
		fmt.Scan(&vin)
		v, err := parseVIN(vin)
		if err != nil {
			fmt.Printf("** ERROR: %s **", err)
		}
		fmt.Printf("\tVIN Check is OK! (%s)\n\tVIN Year: %d\n", v, v.ModelYear())
	}
}
func parseVIN(v string) (*VIN, error) {
	if len(v) != 17 {
		return nil, InvalidVinLength
	}
	vin := new(VIN)
	for i, c := range strings.ToUpper(v) {
		vin[i] = c
	}
	if pass := vin.checkDigit(); !pass {
		return nil, InvalidVin
	}
	return vin, nil
}
func (v *VIN) checkDigit() bool {
	//https://github.com/MartinThoma/vin_decoder/blob/master/vin_decoder/decode.py
	translate := make(map[int32]int)
	for i, c := range "0123456789" {
		translate[c] = i
	}
	for i, c := range "ABCDEFGH" {
		translate[c] = i + 1
	}
	for i, c := range "JKLMNOPQR" {
		if c == 'O' || c == 'Q' {
			continue
		}
		translate[c] = i + 1
	}
	for i, c := range "STUVWXYZ" {
		translate[c] = i + 2 //Wikipedia says S = 1 check this
	}
	weights := []int{8, 7, 6, 5, 4, 3, 2, 10, 0, 9, 8, 7, 6, 5, 4, 3, 2}
	var sum int
	for i, c := range v {
		sum += translate[c] * weights[i]
	}
	rem := sum % 11
	switch {
	case rem == 10 && v[8] == 'X':
		return true
	case rem == translate[v[8]]:
		return true
	}
	return false
}
func (v *VIN) String() string {
	return fmt.Sprintf("%c%c%c%c%c%c%c%c%c%c%c%c%c%c%c%c%c", v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7], v[8], v[9], v[10], v[11], v[12], v[13], v[14], v[15], v[16])
}
func (v *VIN) ModelYear() (year int) {
	translate := make(map[int32]int32)
	var delta int32
	for start := 'A'; start <= 'Z'; start++ {
		if strings.ContainsRune("IOQUZ", start) {
			//IOQUZ are not permissible for year identification
			delta += 1
		}
		translate[start] = start - 'A' + 2010 - delta
	}
	for start := '1'; start <= '9'; start++ {
		translate[start] = start - '1' + 2031
	}
	year = int(translate[v[9]])
	if yearLimit := time.Now().Year() + 3; year > yearLimit {
		delta = 0
		for start := 'A'; start <= 'Z'; start++ {
			if strings.ContainsRune("IOQUZ", start) {
				delta += 1
			}
			translate[start] = start - 'A' + 1980 - delta
		}
		for start := '1'; start <= '9'; start++ {
			translate[start] = start - '1' + 2001
		}
		year = int(translate[v[9]])
	}
	return
}
