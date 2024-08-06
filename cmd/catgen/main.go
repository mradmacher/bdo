package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type CodeDescs map[string]string
type CodeTree map[string][]string
type FinalCodes []string

func readWasteCatalog() (CodeTree, CodeDescs, FinalCodes, error) {
	buffer, err := os.ReadFile("catalog")
	if err != nil {
		return nil, nil, nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(buffer))

	var codeDescs CodeDescs
	var codeTree CodeTree
	var finalCodes FinalCodes
	codeDescs = make(CodeDescs)
	codeTree = make(CodeTree)

	wasteRegexp := regexp.MustCompile(`(\d\d \d\d \d\d\*|\d\d \d\d \d\d|\d\d \d\d|\d\d)\s+(\S.+)`)
	totalLines := 0
	for sc.Scan() {
		totalLines++
		result := wasteRegexp.FindStringSubmatch(sc.Text())

		if result[1] == "" || result[2] == "" {
			return nil, nil, nil, fmt.Errorf("Something weird in line %d: %q: %q", totalLines, result[1], result[2])
		}

		code := strings.Join(strings.Split(result[1], " "), "")
		codeDescs[code] = result[2]
		if len(code) == 2 {
			codeTree["00"] = append(codeTree["00"], code)
		} else if len(code) == 4 {
			codeTree[code[:2]] = append(codeTree[code[:2]], code)
		} else {
			codeTree[code[:4]] = append(codeTree[code[:4]], code)
			finalCodes = append(finalCodes, code)
		}
	}
	if err = sc.Err(); err != nil {
		return nil, nil, nil, err
	}

	total := 0
	total += len(codeTree["00"])
	for _, code1 := range codeTree["00"] {
		total += len(codeTree[code1])
		for _, code2 := range codeTree[code1] {
			total += len(codeTree[code2])
		}
	}
	if totalLines != total {
		return nil, nil, nil, fmt.Errorf("Total lines %d in file different from discovered codes %d", totalLines, total)
	}

	return codeTree, codeDescs, finalCodes, nil
}

func writeCodeTree(f *os.File, codeTree CodeTree) error {
	fmt.Fprintln(f, "const wasteCodes = {")
	fmt.Fprintln(f, "  \"00\": [")
	for _, code1 := range codeTree["00"] {
		fmt.Fprintf(f, "    %q,\n", code1)
	}
	fmt.Fprintln(f, "  ],")
	for _, code1 := range codeTree["00"] {
		fmt.Fprintf(f, "  %q: [\n", code1)
		for _, code2 := range codeTree[code1] {
			fmt.Fprintf(f, "    %q,\n", code2)
		}
		fmt.Fprintln(f, "  ],")
	}
	for _, code1 := range codeTree["00"] {
		for _, code2 := range codeTree[code1] {
			fmt.Fprintf(f, "  %q: [\n", code2)
			for _, code3 := range codeTree[code2] {
				fmt.Fprintf(f, "    %q,\n", code3)
			}
			fmt.Fprintln(f, "  ],")
		}
	}
	fmt.Fprintf(f, "}\n")
	return nil
}

func writeCodeDescs(f *os.File, codeDescs CodeDescs) error {
	fmt.Fprintf(f, "const wasteCodeDescs = {\n")
	for code, desc := range codeDescs {
		normalizedCode, _ := strings.CutSuffix(code, "*")
		fmt.Fprintf(f, "  %q: %q,\n", normalizedCode, desc)
	}

	fmt.Fprintf(f, "}\n")
	return nil
}

func writeFinalCodes(f *os.File, finalCodes FinalCodes) error {
	fmt.Fprintf(f, "const wasteFinalCodes = {\n")
	for _, code := range finalCodes {
		fmt.Fprintf(f, "  {\n    \"code\": %q\n  },\n", code)
	}

	fmt.Fprintf(f, "}\n")
	return nil
}

func main() {
	var codeDescs CodeDescs
	var codeTree CodeTree
	var finalCodes FinalCodes
	var err error
	fileName := "test_catalog.js"

	codeTree, codeDescs, finalCodes, err = readWasteCatalog()
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writeCodeTree(file, codeTree)
	writeCodeDescs(file, codeDescs)
	writeFinalCodes(file, finalCodes)
}
