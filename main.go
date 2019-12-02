package main

import (
	"flag"
	"fmt"
	"go-updater/scraper"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

var (
	opSys            string
	latestVersion    string
	goInstallPath    string
	installedVersion string
	arch             string
	goroot           string
	goArchiveName    string
	downloadAddrPage = "https://golang.org/dl/"
	fileDownloadURL  = "https://dl.google.com/go/" // append goArchiveName after this for complete archive download

)

func main() {

	opSys = runtime.GOOS
	arch = runtime.GOARCH
	goroot = runtime.GOROOT()
	goInstallPath = "/usr/local"
	installedVersion = runtime.Version()[2:]

	checkFlag := flag.Bool("c", false, "Checks latest available version, no other actions performed")
	flag.Parse()

	root := scraper.GetRootNode(downloadAddrPage)
	latestVersion = scraper.GetLatestVersionNumber(root)
	goArchiveName = fmt.Sprintf("go%s.%s-%s.tar.gz", latestVersion, opSys, arch)

	if *checkFlag {
		fmt.Println("Latest available version:", latestVersion)
		fmt.Println("Currently installed version:", runtime.Version()[2:])
		os.Exit(0)
	}

	downloadDest := os.Getenv("HOME") + "/Downloads/" + goArchiveName
	downloadFrom := fileDownloadURL + goArchiveName

	err := downloadFile(downloadFrom, downloadDest)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Encountered an error:", err)
		fmt.Println("Cleaning up after myself...")
		if err := os.Remove(downloadDest); err != nil {
			fmt.Fprintln(os.Stderr, "error removing artifacts:", err)
			fmt.Println("Manually remove the file " + downloadDest)
			os.Exit(-1)
		}
	}

	fmt.Println("Removing old Go installation. Enter your password if prompted.")

	delCmd := fmt.Sprintf("sudo rm -rf %s", goroot)
	_, err = exec.Command("/bin/sh", "-c", delCmd).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Encountered an error:", err)
		os.Exit(-1)
	}

	fmt.Println("Old Go installation removed.")

	fmt.Println("Extracting fresh Go package. Enter your password if prompted.")

	ec := fmt.Sprintf("sudo tar -C %s -xzf %s", goInstallPath, downloadDest)
	cmd := exec.Command("/bin/sh", "-c", ec)
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error extracting archive:", err)
		os.Exit(-1)
	}

	fmt.Println("\nCleaning up...")
	if err := os.Remove(downloadDest); err != nil {
		fmt.Fprintln(os.Stderr, "error removing artifacts:", err)
		fmt.Println("Manually remove the file" + downloadDest)
		os.Exit(-1)
	}

}

// downloadFile downloads from url and places download in dest
func downloadFile(url string, dest string) error {
	fmt.Println("Downloading archive...")
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
