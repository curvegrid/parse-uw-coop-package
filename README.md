# Parse UW Co-op Package PDFs
Parse out relevant information from co-op student resume package PDFs provided by the University of Waterloo.

**NOTE: This program is not perfect, and the results (email address, etc.) will need to be manually inspected for errors and fixed.**

# Installation
Ensure that [Go](https://golang.org/) is installed and setup with a working [`$GOPATH`](https://golang.org/doc/code.html#GOPATH).

Installation and sanity check:

```sh
$ go get github.com/curvegrid/parse-uw-coop-package
$ parse-uw-coop-package -h
```

# Usage
This assumes you are an employer of University of Waterloo co-operative education (co-op, interns) students and have a valid Employer login on [WaterlooWorks](https://waterlooworks.uwaterloo.ca/home.htm).

1. Post a job on WaterlooWorks and wait for student applications to become available.
1. Login to WaterlooWorks, navigate to the applications list, and click the blue 'Application Options' button button near the top of the page to create a **custom** application bundle with **each application as a separate PDF**.
1. Download and unzip the consolidated package.
1. Install this utility, [`parse-uw-coop-package`](https://github.com/curvegrid/parse-uw-coop-package#installation).
1. From the directory where you unzipped the consolidated package of PDFs, run `parse-uw-coop-package` and pipe the output to a CSV file (e.g., `parse-uw-coop-package > applicants.csv`). You can tweak the options (try `parse-uw-coop-package -h`) as required.
1. Import into your spreadsheet of choice. As noted above, manual cleanup will be required. 

# Running and Command Line Options
By default, searches the current directory for all PDFs that fit a regular expression (`-fileregex`) and parse the text within for fields specific to UW co-op.

```sh
Usage of ./parse-uw-coop-package:
  -averagesRegex string
    	Regex for averages (default "Term Average:\\s*([0-9]{2}\\.*[0-9]*)")
  -concurrency int
    	Number of PDF parsing threads to run in parallel (default 4)
  -coverLetterRegex string
    	Regex for cover letter yes/no (default "[Ss]incerely|[Hh]iring [Mm]anager")
  -emailRegex string
    	Regex for email address (default "[A-Za-z0-9_.-]+\\@[A-Za-z0-9.-]+\\.[A-Za-z0-9]+")
  -fileregex string
    	Regex filter for filenames (default "([A-Za-z-]+) ([A-Za-z-]+) \\(([0-9]+)\\).pdf")
  -githubRegex string
    	Regex for Github (default "github.com/[A-Za-z0-9_.-]+")
  -linkedInRegex string
    	Regex for LinkedIn (default "linkedin.com/in/[A-Za-z0-9_.-]+")
  -pdftoascii string
    	PDF to ASCII converter (default "ps2ascii")
  -worktermEvalRegex string
    	Regex for work term evaluations (default "UNSATISFACTORY|MARGINAL|SATISFACTORY|VERY GOOD|EXCELLENT|OUTSTANDING")
```

# Sample Run
```sh
$ parse-uw-coop-package 
ID,First name,Last name,Email,Email with name,LinkedIn,Github,Included a cover letter,Work term evaluations,Term averages,Overall average
123456,Able,Baker,able@example.com,Able Baker <abel@example.com>,,,Yes,"OUTSTANDING,OUTSTANDING,OUTSTANDING,GOOD,OUTSTANDING","72,81,84.5,72,78",73.4
...
```

# Known Issues and Limitations
- This has only been tested on macOS.
- The PDF-to-text converter defaults to `pdftotext`, part of [Xpdf](https://www.xpdfreader.com/download.html), which may not be available on your system. On macOS via Homebrew it's part of [Poppler](https://poppler.freedesktop.org/): `brew install poppler`. See the command line options to adjust. Previously we used `ps2ascii`, a not-default part of Ghostscript, which unfortunately segfaults on a high proportion of modern PDF files.
- The PDF-to-text process is not perfect, especially with formatted PDFs. Email addresses seem to be especially problematic, with many of them mangled. For example, we've seen `jeff@example.com` turn into `.com      example       je   ef@` with `ps2ascii`, even in what seems like a fairly "standard" formatted PDF. Manual cleanup will be required.

# Future enhancements
- DRY up the whole program
- Switch from `pdftotext` to a native Go PDF-to-text solution
- Improve the parsing accuracy: better regexes, etc.
- Direct package download, and integration with tabular info, from WaterlooWorks
- Keyword extraction

# Contributing
[Pull requests welcome](https://github.com/curvegrid/parse-uw-coop-package/pulls).

# Development
Assuming `parse-uw-coop-package` was installed per the previous step, then change to the directory where `go get` downloaded the source:

```sh
$ cd $GOPATH/src/github.com/curvegrid/parse-uw-coop-package
$ go build parse-uw-coop-package.go
$ ./parse-uw-coop-package
```

Note that you will now have two copies of the `parse-uw-coop-package` binary on your system, the one in `$GOPATH/bin` via `go install`, and the one just built in `$GOPATH/src/curvegrid/parse-uw-coop-package` via `go build`.

# License and Copyright
Licensed under the MIT License. See the `LICENSE` file for details of the MIT License. Copyright 2018-2021 by Curvegrid Inc.
