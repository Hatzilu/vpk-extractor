package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/NublyBR/go-vpk"
)

func main() {
    rootPath, logPath := getFlags()	
	enableLogging := len(logPath) > 0
	
	logger := CustomLogger(os.Stdout, "INFO: ", log.Ldate|log.Ltime, enableLogging)
	
	
	if enableLogging {
		logfile, err  := os.Create(filepath.Join(logPath,"vpk-extractor-"+time.Now().Format("2006-01-02-15-04-05")+".log"))
	
		if err != nil {
			fmt.Println(err.Error())
		}

		defer logfile.Close()
		logger.SetOutput(logfile)
	} else {
		logger.SetOutput(nil)
	}

	
	pak, err := vpk.OpenDir(rootPath)

	if err != nil {
		logger.Fatal(err)
	}

	defer pak.Close()

	// Iterate through all files in the VPK
	for _, file := range pak.Entries() {
		entry, err := file.Open()
		if err != nil {
			logger.Fatal(err)
		}
		
		path := filepath.Join("output",file.Filename())

		// Ensure the directories exist by using os.MkdirAll
		dir := filepath.Dir(path)
		if dir_err := os.MkdirAll(dir, 0755); dir_err != nil {
			logger.Fatal(dir_err)
		}
		logger.Println("Extract file",path)
		writeErr := ExtractVpkFile(entry, path)
		if writeErr != nil {
			logger.Fatal(err)
		}	
	}
	logger.Println("Successfully extracted content to "+rootPath)
}




func ExtractVpkFile(file vpk.FileReader, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	io.Copy(f, file)

	f.Sync()
	closeErr := f.Close()
	if closeErr != nil {
		return closeErr
	}
	return nil
}

func getFlags() (string, string) {
	path := flag.String("p", "", "the full path to the vpk file.")
	logPath := flag.String("logPath", "", "the path where the log will be written to, leave blank to disable the creation of a logfile.")
	
	flag.Parse()


	if *path == "" {
		panic("No .vpk file path was specified.")
	}

	_, err := os.Stat(*path)

	if err != nil {
		panic(err)
	}

	if filepath.Ext(*path) != ".vpk" {
		panic("File is not vpk format")
	}

	return *path, *logPath
}