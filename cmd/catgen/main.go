package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strings"
)

type Material struct {
	Code  string
	Name  string
	Desc  string
	Items []Material
}

func readMaterials(fileName string) ([]Material, error) {
	var materials []Material
	buffer, err := os.ReadFile(fileName)
	if err != nil {
		return nil, errors.Join(errors.New("Problem reading materials file"), err)
	}
	err = yaml.Unmarshal(buffer, &materials)
	if err != nil {
		return nil, errors.Join(errors.New("Problem unmarshaling materials file"), err)
	}
	fmt.Printf("%v\n", materials)

	return materials, nil
}

func genMaterials(fileName string, materials []Material) error {
	f, err := os.Create(fileName)
	if err != nil {
		return errors.Join(errors.New("Problem creating materials output file"), err)
	}
	defer f.Close()

	fmt.Fprintln(f, "export const materialCodes = {")
	fmt.Fprintln(f, "  \"00\": [")
	for _, m := range materials {
		fmt.Fprintf(f, "    %q,\n", m.Code)
	}
	fmt.Fprintln(f, "  ],")
	for _, m1 := range materials {
		fmt.Fprintf(f, "  %q: [\n", m1.Code)
		for _, m2 := range m1.Items {
			fmt.Fprintf(f, "    %q,\n", m2.Code)
		}
		fmt.Fprintln(f, "  ],")
	}
	for _, m1 := range materials {
		for _, m2 := range m1.Items {
			fmt.Fprintf(f, "  %q: [\n", m2.Code)
			for _, m3 := range m2.Items {
				fmt.Fprintf(f, "    %q,\n", m3.Code)
			}
			fmt.Fprintln(f, "  ],")
		}
	}
	fmt.Fprintf(f, "}\n")

	fmt.Fprintf(f, "export const materialNames = {\n")
	for _, m1 := range materials {
		fmt.Fprintf(f, "  %q: %q,\n", m1.Code, m1.Name)
		for _, m2 := range m1.Items {
			fmt.Fprintf(f, "  %q: %q,\n", m2.Code, m2.Name)
			for _, m3 := range m2.Items {
				fmt.Fprintf(f, "  %q: %q,\n", m3.Code, m3.Name)
			}
		}
	}
	fmt.Fprintf(f, "}\n")

	return nil
}

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

func makeWasteCodes(fileName string) {
	var codeDescs CodeDescs
	var codeTree CodeTree
	var finalCodes FinalCodes
	var err error

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

func main() {
	//makeWasteCodes("test_catalog.js")
	var materials []Material
	var err error
	materials, err = readMaterials("internal/seeds/materials.yaml")
	if err != nil {
		panic(err)
	}

	err = genMaterials("js/material_catalog.js", materials)
	if err != nil {
		panic(err)
	}
}
