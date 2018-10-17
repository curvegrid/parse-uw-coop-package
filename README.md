# Parse UW Co-op Package PDFs
Parse out relevant information from co-op student resume package PDFs provided by the University of Waterloo.

**NOTE: This program is not perfect, and the results (email address, etc.) will need to be manually inspected for errors and fixed.**

# Installation
Ensure that [Go](https://golang.org/) is installed and setup with a working [`$GOPATH`](https://golang.org/doc/code.html#GOPATH).

Installation and sanity check:

```sh
$ go install github.com/curvegrid/parse-uw-coop-package
$ parse-uw-coop-package -h
```

# Development
Assuming `parse-uw-coop-package` was installed per the previous step, then change to the directory where `go get` downloaded the source:

```sh
$ cd $GOPATH/src/github.com/curvegrid/parse-uw-coop-package
$ go build parse-uw-coop-package.go
$ ./parse-uw-coop-package
```

Note that you will now have two copies of the `parse-uw-coop-package` binary on your system, the one in `$GOPATH/bin` via `go install`, and the one just built in `$GOPATH/src/curvegrid/parse-uw-coop-package` via `go build`.

# Usage
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

# Known Issues
- The PDF-to-text process is not perfect, especially with formatted PDFs. Email addresses seem to be especially problematic, with many of them mangled. For example, we've seen `jeff@example.com` turn into `.com      example       je   ef@` with `ps2ascii`, even in what seems like a fairly "standard" formatted PDF. Manual cleanup will be required. 

# License and Copyright
Licensed under the MIT License. See the `LICENSE` file for details of the MIT License. Copyright 2018 by Curvegrid Inc.
