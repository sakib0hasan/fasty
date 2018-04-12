package main

import (
	"os"
	"bufio"
	"net/http"
	"io"
	"strings"
)

func main(){
	pool := NewPool(100)

	inFile, _ := os.Open("images.txt")
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "/")
		pool.Exec(DownloadTask{"images/"+parts[len(parts)-1], scanner.Text()})
	}

	pool.Close()
	pool.Wait()
}

type DownloadTask struct {
	filepath string
	url string
}
func (e DownloadTask) Execute() {
	downloadFile(e.filepath, e.url)
}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil  {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}

	return nil
}