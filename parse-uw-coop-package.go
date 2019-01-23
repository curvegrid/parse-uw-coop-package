// Licensed under the MIT License
// Copyright (c) 2018 Curvegrid Inc.

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func uniqueStrings(squashCase bool, inStrs []string) (outStrs []string) {
	strsMap := make(map[string]bool)
	for _, inStr := range inStrs {
		if squashCase {
			inStr = strings.ToLower(inStr)
		}

		strsMap[inStr] = true
	}

	for k, _ := range strsMap {
		outStrs = append(outStrs, k)
	}

	return outStrs
}

func main() {
	options := struct {
		fileRegex         string
		pdfToAscii        string
		pathToParse       string
		concurrency       int
		emailRegex        string
		linkedInRegex     string
		githubRegex       string
		coverLetterRegex  string
		worktermEvalRegex string
		averagesRegex     string
	}{
		fileRegex:         `([A-Za-z -]+) ([A-Za-z-]+) \(([0-9]+)\).pdf`,
		pdfToAscii:        "ps2ascii",
		pathToParse:       ".",
		concurrency:       4, // 4 seems to be a sweet spot
		emailRegex:        `[A-Za-z0-9_.-]+\@[A-Za-z0-9.-]+\.[A-Za-z0-9]+`,
		linkedInRegex:     `linkedin.com/in/[A-Za-z0-9_.-]+`,
		githubRegex:       `github.com/[A-Za-z0-9_.-]+`,
		coverLetterRegex:  `[Ss]incerely|[Hh]iring [Mm]anager`,
		worktermEvalRegex: `UNSATISFACTORY|MARGINAL|SATISFACTORY|VERY GOOD|EXCELLENT|OUTSTANDING`,
		averagesRegex:     `Term Average:\s*([0-9]{2}\.*[0-9]*)`,
	}

	flag.StringVar(&options.fileRegex, "fileregex", options.fileRegex, "Regex filter for filenames")
	flag.StringVar(&options.pdfToAscii, "pdftoascii", options.pdfToAscii, "PDF to ASCII converter")
	flag.IntVar(&options.concurrency, "concurrency", options.concurrency, "Number of PDF parsing threads to run in parallel")
	flag.StringVar(&options.emailRegex, "emailRegex", options.emailRegex, "Regex for email address")
	flag.StringVar(&options.linkedInRegex, "linkedInRegex", options.linkedInRegex, "Regex for LinkedIn")
	flag.StringVar(&options.githubRegex, "githubRegex", options.githubRegex, "Regex for Github")
	flag.StringVar(&options.coverLetterRegex, "coverLetterRegex", options.coverLetterRegex, "Regex for cover letter yes/no")
	flag.StringVar(&options.worktermEvalRegex, "worktermEvalRegex", options.worktermEvalRegex, "Regex for work term evaluations")
	flag.StringVar(&options.averagesRegex, "averagesRegex", options.averagesRegex, "Regex for averages")
	flag.Parse()

	if flag.Arg(0) != "" {
		options.pathToParse = flag.Arg(0)
	}

	fileRe := regexp.MustCompile(options.fileRegex)

	files, err := ioutil.ReadDir(options.pathToParse)
	if err != nil {
		log.Fatal(err)
	}

	// collect filenames to parse
	filenames := []string{}
	for _, file := range files {
		// ensure filename matches criteria
		if file.IsDir() {
			continue
		}

		filenameComponents := fileRe.FindStringSubmatch(file.Name())
		if len(filenameComponents) != 4 {
			continue
		}

		// filename matches criteria: save for later processing
		filenames = append(filenames, file.Name())
	}

	// setup channel
	filenameChan := make(chan string, len(filenames))
	recordsChan := make(chan []string, len(filenames))
	var wg sync.WaitGroup

	// Compile regexes
	emailRe := regexp.MustCompile(options.emailRegex)
	linkedInRe := regexp.MustCompile(options.linkedInRegex)
	githubRe := regexp.MustCompile(options.githubRegex)
	coverLetterRe := regexp.MustCompile(options.coverLetterRegex)
	worktermEvalRe := regexp.MustCompile(options.worktermEvalRegex)
	averagesRe := regexp.MustCompile(options.averagesRegex)

	// spin up processing goroutines
	for i := 0; i < options.concurrency; i++ {
		wg.Add(1)
		go func() {
			for filename := range filenameChan {
				filenameComponents := fileRe.FindStringSubmatch(filename)
				if len(filenameComponents) != 4 {
					continue
				}

				// key fields from filename
				firstName := filenameComponents[1]
				lastName := filenameComponents[2]
				id := filenameComponents[3]

				// extract text from PDF
				cmd := exec.Command(options.pdfToAscii, filename)
				pdfText, err := cmd.Output()
				if err != nil {
					log.Fatal(err)
				}
				pdfTextStr := string(pdfText)

				// parse additional information

				// email
				emails := emailRe.FindAllString(pdfTextStr, -1)
				uniqueEmails := uniqueStrings(true, emails)
				emailsFiltered := strings.Join(uniqueEmails, ",")

				emailsFull := []string{}
				for _, uniqueEmail := range uniqueEmails {
					emailsFull = append(emailsFull, firstName+" "+lastName+" <"+uniqueEmail+">")
				}
				emailsFullFiltered := strings.Join(emailsFull, ",")

				// linkedin
				linkedIn := linkedInRe.FindString(pdfTextStr)

				// github
				github := githubRe.FindString(pdfTextStr)

				// included a cover letter?
				coverLetter := coverLetterRe.MatchString(pdfTextStr)
				coverLetterStr := "No"
				if coverLetter {
					coverLetterStr = "Yes"
				}

				// work term evaluation
				worktermEvals := strings.Join(worktermEvalRe.FindAllString(pdfTextStr, -1), ",")

				// grades and overall average
				averagesMatch := averagesRe.FindAllStringSubmatch(pdfTextStr, -1)
				averages := []string{}
				var overallAverage float64
				for _, averageMatch := range averagesMatch {
					averages = append(averages, averageMatch[1])
					termAverage, err := strconv.ParseFloat(averageMatch[1], 64)
					if err != nil {
						fmt.Printf("Error parsing average '%v' for id '%v'", averageMatch[1], id)
						break
					}

					overallAverage += termAverage
				}

				overallAverageFiltered := "Unknown"
				if len(averagesMatch) > 0 {
					overallAverage /= float64(len(averagesMatch))
					overallAverageFiltered = fmt.Sprintf("%.1f", overallAverage)
				}
				averagesFiltered := strings.Join(averages, ",")

				// append to records slice
				recordsChan <- []string{id, firstName, lastName, emailsFiltered, emailsFullFiltered, linkedIn, github, coverLetterStr, worktermEvals, averagesFiltered, overallAverageFiltered}
			}
			wg.Done()
		}()
	}

	// seed filenames channel
	for _, filename := range filenames {
		filenameChan <- filename
	}
	close(filenameChan)

	// wait for goroutines to finish and collect results
	wg.Wait()

	// TODO: perhaps this should be within a collection goroutine? Although it is protected by wg.Wait()
	close(recordsChan)

	records := [][]string{}
	headers := []string{
		"ID",
		"First name",
		"Last name",
		"Email",
		"Email with name",
		"LinkedIn",
		"Github",
		"Included a cover letter",
		"Work term evaluations",
		"Term averages",
		"Overall average",
	}

	for record := range recordsChan {
		records = append(records, record)
	}

	// sort the records for consistency
	sort.Slice(records, func(i, j int) bool { return records[i][0] < records[j][0] })

	// output in a parseable format
	w := csv.NewWriter(os.Stdout)
	w.Write(headers)
	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}
