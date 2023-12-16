package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/NublyBR/go-vpk"
)

func main() {
    rootPath := `E:\SteamLibrary\steamapps\common\Team Fortress 2\tf\tf2_misc_dir.vpk`	
	logfile, err  := os.Create("write.log")


	defer logfile.Close()
	log.SetOutput(logfile)

	if err != nil {
		log.Fatal(err)
	}

	pak, err := vpk.OpenDir(rootPath)

	defer pak.Close()

	// Iterate through all files in the VPK
	for _, file := range pak.Entries() {
		entry, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		
		path := filepath.Join("output",file.Filename())

		// Ensure the directories exist by using os.MkdirAll
		dir := filepath.Dir(path)
		if dir_err := os.MkdirAll(dir, 0755); dir_err != nil {
			log.Fatal(dir_err)
		}

		err, fileBytes := WriteVpkFile(entry, path)
		if err != nil {
			log.Fatal(err)
		}	
		fmt.Println(fileBytes)
	}
	
}



func WriteVpkFile(file vpk.FileReader, path string) (error, bool) {
	f, err := os.Create(path)
	if err != nil {
		return err, false
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
			fmt.Println("Buf: ",string(content))
			fmt.Println("Bytes: ",content)
			f.Write(content)
		}
    }
	f.Sync()
	return nil, true
}