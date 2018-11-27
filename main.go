package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type InputData struct {
	String string
	NumSpl bool
	SecSpl bool
}

var (
	inputData         string
	numSplitter       rune
	sequencesSplitter rune
)

func SplitAll(r rune) bool {
	return r == numSplitter || r == sequencesSplitter
}

func SplitArg(r rune) bool {
	return r == '='
}

func SplitNum(r rune) bool {
	return r == numSplitter
}

func SplitSecq(r rune) bool {
	return r == sequencesSplitter
}

func ReturnHelp() {
	msg := `Usage:	
	main {numbers} [Parameters]
Parameters:
	{numbers} Number for build array, required
	[-i | --itemSep] Separator for items in string, default ','
	[-s | --secSep] Separator for sequence in string, default '|'
Example: 
	main '1|5,6,7|10'
	main '1%5.6.7%10' -i='.' -s='%'`

	fmt.Println(msg)
}

func CheckArgs(args []string) {
	var (
		NumSpl    bool = false
		SecSpl    bool = false
		NumSplSet bool = false
		SecSplSet bool = false
	)
	if len(args) == 1 {
		getError("Number for build array not specified")
	}

	inputData = args[1]
	for _, argument := range args[2:] {

		a := strings.FieldsFunc(argument, SplitArg)
		arg, val := a[0], a[1]

		switch arg {
		case "-i", "--itemSep":
			NumSpl = true
		case "-s", "--secSep":
			SecSpl = true
		}

		if NumSpl {
			runeVal := []rune(val)
			numSplitter = runeVal[0]
			NumSpl = false
			NumSplSet = true
			continue
		}
		if SecSpl {
			runeVal := []rune(val)
			sequencesSplitter = runeVal[0]
			SecSpl = false
			SecSplSet = true
			continue
		}
	}

	if !NumSplSet {
		numSplitter = ','
	}
	if !SecSplSet {
		sequencesSplitter = '|'
	}

}

func getError(err string) {
	ReturnHelp()
	fmt.Println("Error:", err)
	os.Exit(1)
}

func AppendInt(slice []int, data ...int) []int {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) {
		newSlice := make([]int, (n+1)*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:n]
	copy(slice[m:n], data)
	return slice
}

func AppendStr(slice []string, data ...string) []string {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) {
		newSlice := make([]string, (n+1)*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:n]
	copy(slice[m:n], data)
	return slice
}

func Check(numbers string) InputData {
	data := InputData{}

	NumSpl := strings.ContainsAny(inputData, string(numSplitter))
	SecSpl := strings.ContainsAny(inputData, string(sequencesSplitter))

	if NumSpl || SecSpl {
		a := strings.FieldsFunc(inputData, SplitAll)
		for i := 0; i < len(a); i++ {
			if _, err := strconv.Atoi(a[i]); err != nil {
				fmt.Printf("%q is a NAN \n", a[i])
				os.Exit(1)
			}

		}
	} else {
		inputStringLen := len(inputData)
		if inputStringLen == 0 {
			getError("Number for build array not specified")
		} else {
			e := []string{}
			for i := 0; i < inputStringLen; i++ {
				num := string(inputData[i])
				if _, err := strconv.Atoi(num); err != nil {
					e = AppendStr(e, string(num))
				}
			}
			fmt.Printf("%s wrong sequnse ", inputData)
			fmt.Printf("and '%s' is a NAN \n", strings.Join(e, ","))
			os.Exit(1)
		}

		os.Exit(1)
	}

	data.String = numbers
	data.NumSpl = NumSpl
	data.SecSpl = SecSpl

	return data
}

func Process(numbers string) []int {
	r := []int{}
	data := Check(numbers)
	if data.NumSpl {
		a := strings.FieldsFunc(data.String, SplitNum)
		for i := 0; i < len(a); i++ {
			SecSpl := strings.ContainsAny(a[i], string(sequencesSplitter))
			if SecSpl {
				b := strings.FieldsFunc(a[i], SplitSecq)
				start, _ := strconv.Atoi(b[0])
				end, _ := strconv.Atoi(b[1])
				for i := start; i <= end; i++ {
					r = AppendInt(r, i)
				}
			} else {
				num, _ := strconv.Atoi(a[i])
				r = AppendInt(r, num)
			}
		}
	} else {
		a := strings.FieldsFunc(data.String, SplitSecq)
		for i := 0; i < len(a); i++ {
			num, _ := strconv.Atoi(a[i])
			r = AppendInt(r, num)
		}
	}

	return r
}

func main() {
	CheckArgs(os.Args)
	fmt.Println(Process(inputData))
}
