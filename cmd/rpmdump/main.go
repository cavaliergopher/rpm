/*
rpmdump displays all headers and tags in rpm packages as a YAML document.

	usage: rpmdump [package ...]

Example:

	$ rpmdump golang-1.6.3-2.el7.x86_64.rpm
	---
	- path: golang-1.6.3-2.el7.x86_64.rpm
	  signature:
		version: 1
		tags:
		- tag: 1000
			type: INT32
			value: [13140]
		- tag: 1002
			type: BIN
			value: |
			00000000  89 02 15 03 05 00 54 89  c9 9c 24 c6 a8 a7 f4 a8  |......T...$.....|
			00000010  0e b5 01 08 a9 dd 0f fe  3b 48 ad dc 55 d5 ec a8  |........;H..U...|
			00000020  a7 e3 64 13 4b ee 3c d4  fb 6e fa b4 c0 bf a9 57  |..d.K.<..n.....W|
			00000030  b9 4a d9 ac 01 3c 34 3a  c9 18 99 08 d8 c2 25 12  |.J...<4:......%.|
			00000040  4c 7a 5e 4b 05 41 94 4d  d4 85 8f 8a a5 12 60 67  |Lz^K.A.M......`g|
			00000050  15 67 fd 8c 7a 26 a0 24  26 81 a2 d9 f0 c5 fb 8d  |.g..z&.$&.......|
			00000060  e3 47 15 f7 9f 81 cf 84  3d 40 a2 37 31 c7 e1 b0  |.G......=@.71...|
			00000070  1e b7 3f 5e fd 9c 63 06  83 44 3e 84 39 f2 8e c6  |..?^..c..D>.9...|
			00000080  a6 ae 47 47 83 0c 16 62  42 07 4d 94 fe 1d 9c f1  |..GG...bB.M.....|
			00000090  24 cc 35 55 07 76 2f a7  9f e6 ed 94 39 7c 3f b6  |$.5U.v/.....9|?.|
			000000a0  27 82 22 e9 83 79 6b 6e  74 ac 72 38 db ea 65 e4  |'."..yknt.r8..e.|
			000000b0  14 78 cc bd 37 b5 ef 35  c0 17 04 3e 2c b6 f7 fd  |.x..7..5...>,...|
			000000c0  90 e5 12 1f 69 bd 1c 3e  31 83 cd 44 6b d1 c7 37  |....i..>1..Dk..7|
			000000d0  b6 4a 5e 5d fa fa f2 04  c9 51 9a 56 26 8e fb 0e  |.J^].....Q.V&...|
			000000e0  2c b4 d3 f4 a4 10 39 97  d0 be 99 6d 24 00 6d 59  |,.....9....m$.mY|
			000000f0  4e fc 58 0e 7c 8f 6f 88  a9 cb c2 03 80 37 b3 c0  |N.X.|.o......7..|
			00000100  ab c0 46 e6 29 85 4b 8b  5b f5 18 de b4 c4 77 b7  |..F.).K.[.....w.|
			00000110  61 43 5e 2c f2 f3 ea c6  1d f6 36 11 46 16 50 e2  |aC^,......6.F.P.|
			00000120  b5 b3 9c 9f 81 2f f7 03  b3 f3 83 9d d4 53 48 07  |...../.......SH.|
			00000130  49 cd 8d 3a f8 2d 24 50  db 69 3c 99 e0 37 4e 5d  |I..:.-$P.i<..7N]|
			00000140  38 19 96 69 ea d4 30 2f  4b 61 d6 69 8c ee 06 ee  |8..i..0/Ka.i....|
			00000150  ee 78 af 9a 34 70 0d b8  4a 86 30 a2 31 48 de 98  |.x..4p..J.0.1H..|
			00000160  55 5b 57 cb 6b 4b 81 52  1a ab d5 2c 0e b0 e0 03  |U[W.kK.R...,....|
			00000170  88 46 3e 9d 2d 71 98 27  8a 40 b3 81 4b 29 4f 98  |.F>.-q.'.@..K)O.|
			00000180  9a ea b8 c9 ec 6d f6 09  15 62 8b c2 72 73 87 2d  |.....m...b..rs.-|
			00000190  8b af 52 de b0 a4 c8 5d  59 f3 5c 7a de 98 d0 7f  |..R....]Y.\z....|
			000001a0  a6 c6 96 89 ca 85 12 35  90 c5 fe 73 67 28 c1 65  |.......5...sg(.e|
			000001b0  36 15 db ce 50 f4 fc 74  f8 77 92 6a 65 2a cf fe  |6...P..t.w.je*..|
			000001c0  eb 1e 22 81 fa 87 f2 32  fa fd 10 6d 23 86 92 c5  |.."....2...m#...|
			000001d0  c8 25 a6 51 51 24 11 57  0c 4d bf 4d 38 e9 59 ff  |.%.QQ$.W.M.M8.Y.|
			000001e0  66 73 d5 5b 63 c6 89 1c  b0 ba ec bd be d8 ff 51  |fs.[c..........Q|
			000001f0  80 18 3d ea e3 b9 87 b1  d8 28 58 1e eb 6b ee 03  |..=......(X..k..|
			00000200  5d 3d 37 9b 92 3f c1 58  55 8a c9 a9 34 46 3a df  |]=7..?.XU...4F:.|
			00000210  e3 3f c3 97 6f 21 37 ff                           |.?..o!7.........|
		- tag: 1004
			type: BIN
			value: [74 e3 cd 32 88 e6 9c 33 fb e4 75 ba df ac 0e 7c]
		- tag: 1007
			type: INT32
			value: [26088]
		- tag: 62
			type: BIN
			value: [00 00 00 3e 00 00 00 07 ff ff ff 90 00 00 00 10
	...

*/
package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/cavaliergopher/rpm"
)

func main() {
	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-") {
		os.Exit(usage(1))
	}
	fmt.Printf("---\n")
	for i, name := range os.Args[1:] {
		if i > 0 {
			fmt.Printf("\n")
		}
		printPackage(name)
	}
}

func printPackage(name string) {
	fmt.Printf("- path: %v\n", name)
	p, err := rpm.Open(name)
	if err != nil {
		fmt.Printf("  error: %v\n", err)
		return
	}
	fmt.Printf("  signature:\n")
	printHeader(&p.Signature)
	fmt.Println()
	fmt.Printf("  header:\n")
	printHeader(&p.Header)
}

func printHeader(h *rpm.Header) {
	fmt.Printf("    version: %v\n", h.Version)
	fmt.Printf("    tags:\n")
	for _, tag := range h.Tags {
		fmt.Printf("      - tag: %v\n", tag.ID)
		fmt.Printf("        type: %v\n", tag.Type)
		switch tag.Value.(type) {
		case []string:
			ss := tag.Value.([]string)
			if len(ss) == 1 && !strings.Contains(ss[0], "\n") {
				fmt.Printf("        value: [\"%v\"]\n", ss[0])
			} else {
				fmt.Printf("        value:\n")
				for _, s := range ss {
					if !strings.Contains(s, "\n") {
						fmt.Printf("          - \"%v\"\n", s)
					} else {
						fmt.Printf("          - |\n")
						lines := strings.Split(s, "\n")
						for _, line := range lines {
							fmt.Printf("            %v\n", line)
						}
					}
				}
			}

		case []byte:
			b := tag.Value.([]byte)
			if len(b) <= 16 {
				fmt.Print("        value: [")
				for i, x := range b {
					if i > 0 {
						fmt.Print(" ")
					}
					fmt.Printf("%02x", x)
				}
				fmt.Println("]")
			} else {
				fmt.Println("        value: |")
				for i := 0; i < len(b); i += 16 {
					fmt.Printf("          %08x  ", i)
					l := int(math.Min(16, float64(len(b)-i)))
					for j := 0; j < l; j++ {
						fmt.Printf("%02x ", b[i+j])
						if j == 7 {
							fmt.Print(" ")
						}
					}

					for j := 0; j < 16-l; j++ {
						fmt.Print("   ")
					}
					if l < 8 {
						fmt.Print(" ")
					}

					s := [16]byte{}
					copy(s[:], b[i:])
					for j := 0; j < 16; j++ {
						// print '.' if char is not printable ascii
						if s[j] < 32 || s[j] > 126 {
							s[j] = 46
						}
					}
					fmt.Printf(" |%s|\n", s)
				}
			}

		default:
			fmt.Printf("        value: %v\n", tag.Value)
		}
	}
}

func usage(exitCode int) int {
	w := os.Stdout
	if exitCode != 0 {
		w = os.Stderr
	}

	fmt.Fprintf(w, "usage: %v [package ...]\n", os.Args[0])
	return exitCode
}
