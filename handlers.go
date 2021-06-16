package main

import (
	"app/checksum"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

var (
	Status = struct {
		Version    string
		NewVersion *UpdateInfo
		Refresh    bool
	}{
		Version,
		nil,
		false,
	}

	// mutex to avoid data races when working with Status
	statusLock sync.RWMutex
)

func handleError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(500)
	fmt.Fprintf(w, "Internal Error")
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	statusLock.RLock()
	defer statusLock.RUnlock()
	must(PageTemplate.Execute(w, Status))
}

func handleCheck(w http.ResponseWriter, r *http.Request) {
	hasNew, info, err := checkNewVersion(UpdateServer, Version)
	if err != nil {
		handleError(w, err)
		return
	}
	if hasNew {
		statusLock.Lock()
		Status.NewVersion = &info
		statusLock.Unlock()
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func handleInstall(w http.ResponseWriter, r *http.Request) {
	statusLock.RLock()
	if Status.NewVersion == nil {
		// if no version, get back to the main page
		statusLock.RUnlock()
		http.Redirect(w, r, "/?no-updates=1", http.StatusTemporaryRedirect)
		return
	}
	// get new version file checksum
	ck := Status.NewVersion.Sha256
	statusLock.RUnlock()

	// exe will be used to get current executable filename
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	exePath := filepath.Dir(exe)
	exeName := filepath.Base(exe)

	// create a temporary file next to the currently running executable
	f, err := ioutil.TempFile(exePath, exeName+".*")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// download the executable
	err = download(f, UpdateServer+"/app.exe")
	if err != nil {
		log.Fatal(err)
	}

	//// the idea was to check PE checksum, but golang appears not to write it: it's always 00 00 00 00 :(
	// isValid, err := checksum.ValidatePEChecksum(f)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// compare the checksum of the downloaded file with the expected one
	isValid, err := checksum.ValidateSha256Checksum(f, ck)
	if err != nil {
		log.Fatal(err)
	}

	// If not valid, maybe we need to retry? But for the simplicity let's just redirect to the home page
	if !isValid {
		log.Println("downloaded file checksum does not match the expected one")
		http.Redirect(w, r, "/?failed=1", http.StatusTemporaryRedirect)
		return
	}

	// TODO: I could have also verified the signature of the downloaded file,
	// TODO: but to be able to do that I would need a valid certificate to sign the exe first

	// rename current executable to `app.exe.bak`
	err = os.Rename(exe, exe+".bak")
	if err != nil {
		log.Fatal(err)
	}

	// close temporary file to be able to rename it
	f.Close() // yes, we have defer, but calling Close twice doesn't panic, to make it cleaner we could have some internal var to check if we closed the file
	// rename our temp file to the current executable original name
	err = os.Rename(f.Name(), exe)
	if err != nil {
		// in case of any error rename the exe file back
		_ = os.Rename(exe+".bak", exe) // first restore the old file
		log.Fatal(err)
	}

	// note, using cmd.exe lets us running the program in foreground
	// if not a requirement, it can be run directly
	cmd := exec.Command("cmd.exe", "/C", "start", exe)
	if err := cmd.Run(); err != nil {
		// in case of any error rename the exe file back
		_ = os.Rename(exe+".bak", exe) // first restore the old file
		log.Fatal(err)
	}

	statusLock.Lock()
	Status.Refresh = true
	statusLock.Unlock()

	// finally, redirect to the home page and wait...
	http.Redirect(w, r, "/?updated=1", http.StatusTemporaryRedirect)
	go func() {
		time.Sleep(50 * time.Millisecond) // allow the redirect to finish
		os.Exit(0)
	}()
}
