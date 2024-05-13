package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"sync"
)

/*
 * Améliorations possibles:
 * - Mettre une option pour rendre insensible à la casse
 * - Rendre le nb de workers dynamique
 * - colorer le pattern recherché dans l'output
 */

func grepInFilePart(filePath string, pattern string, index int, start, end int64, wg *sync.WaitGroup, result chan<- Result) {
	defer wg.Done()
	var matchedLines []string

	file, err := os.Open(filePath)
	if err != nil {
		panic(err) //should not happen
	}
	defer file.Close()

	if start > 0 {
		file.Seek(int64(start-1), 0)
	}
	scanner := bufio.NewScanner(file)
	totalChar := 0
	isFirstLineSkipped := false
	for scanner.Scan() {
		line := scanner.Bytes()
		totalChar += len(line) + 1

		//fmt.Printf(" ---------- %d %d : %s\n", index, totalChar, line)

		// If we are at the beginning of a line, prev caracter is '\n'
		if index > 0 && len(line) == 1 {
			isFirstLineSkipped = true
			continue
		}

		// If line + prev caracter has a len != 1 => it is an incomplete line
		if index > 0 && !isFirstLineSkipped {
			isFirstLineSkipped = true
			continue
		}

		//fmt.Printf(" ++++++++++ %d : %s\n", index, line)

		if matched, _ := regexp.MatchString(pattern, string(line)); matched {
			matchedLines = append(matchedLines, string(line))
		}

		if int64(totalChar)+start >= end {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err) //should not happen
	}

	result <- Result{
		index:   index,
		resList: matchedLines,
	}
}

type Result struct {
	index   int
	resList []string
}

func grepFile(filePath string, pattern string, numWorkers int) (res []string, err error) {
	result := make(chan Result)
	var wg sync.WaitGroup

	// Divide file in parts to be processed by workers
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	partSize := fileSize / int64(numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := int64(i) * partSize
		end := int64(i+1) * partSize
		if i == numWorkers-1 {
			end = fileSize
		}
		wg.Add(1)
		go grepInFilePart(filePath, pattern, i, start, end, &wg, result)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	results, _ := sortResults(result)
	matchedLines := concatResults(results)

	return matchedLines, nil
}

func concatResults(results []Result) []string {
	var res []string
	for _, r := range results {
		res = append(res, r.resList...)
	}
	return res
}

func sortResults(result chan Result) ([]Result, error) {
	var results []Result
	for res := range result {
		results = append(results, res)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})
	return results, nil
}

func main() {
	numWorkers := 4 // Nb of workers to use for the grep
	var err error

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Utilisation: %s <pattern> <fichier>\n", os.Args[0])
		os.Exit(1)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Fprintf(os.Stderr, "Utilisation: %s <pattern> <fichier>\n", os.Args[0])
		os.Exit(0)
	}

	pattern := os.Args[1]
	filePath := os.Args[2]

	var matchedLines []string

	if filePath == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if matched, _ := regexp.MatchString(pattern, line); matched {
				matchedLines = append(matchedLines, line)
			}
		}
		if err = scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de la lecture de l'entrée standard: %v\n", err)
		}
	} else {
		matchedLines, err = grepFile(filePath, pattern, numWorkers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors du grep du fichier %v\n", err)
		}
	}

	// Résultat
	for _, line := range matchedLines {
		fmt.Println(line)
	}
}
