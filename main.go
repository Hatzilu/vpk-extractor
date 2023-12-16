package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/NublyBR/go-vpk"
)

func main() {
    rootPath := `E:\SteamLibrary\steamapps\common\Team Fortress 2\tf\tf2_sound_misc_dir.vpk`	
	logfile, err  := os.Create("vpk-extractor-"+time.Now().Format("2006-01-02-15-04-05")+".log")

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	defer logfile.Close()
	logger.SetOutput(logfile)

	if err != nil {
		logger.Fatal(err)
	}

	pak, err := vpk.OpenDir(rootPath)

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
		logger.Println("Writing file",path)
		fmt.Println("Writing file",path)
		writeErr := WriteVpkFile(entry, path)
		if writeErr != nil {
			logger.Fatal(err)
		}	
	}
	logger.Println("Successfully extracted content to "+rootPath)
	fmt.Println("Successfully extracted content to "+rootPath)
}



func WriteVpkFile(file vpk.FileReader, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

    for {
		buf := make([]byte, 1024)
        n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			content := buf[:n]
			// Print the actual content of the file, TODO: uncomment this when we add a --verbose flag
			// fmt.Println("Buf: ",string(content))
			f.Write(content)
		}
    }
	f.Sync()
	return nil
}