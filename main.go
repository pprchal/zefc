package main

import (
	"archive/zip"
	"fmt"
	"os"
	"runtime"
	"slices"
	"sync"
	structs "zefc/structs"
	utils "zefc/utils"

	"github.com/fatih/color"
)

type ResultState int

const (
	OK ResultState = iota
	Missing
	Differs
)

type Result struct {
	Eta   structs.Etalon
	Cur   structs.Etalon
	State ResultState
}

func compareCurrentWithEtalon(current chan structs.Etalon, etalon []structs.Etalon) []Result {
	results := make([]Result, len(etalon))

	// collect all current files
	currentFiles := make(map[string]structs.Etalon, len(etalon))
	for cur := range current {
		currentFiles[cur.FileName] = cur
	}

	// check each etalon file
	for _, eta := range etalon {
		if cur, exists := currentFiles[eta.FileName]; exists {
			// File exists in current
			var state ResultState
			if cur.SHA1 == eta.SHA1 {
				state = OK
			} else {
				state = Differs
			}
			results = append(results, Result{
				Eta:   eta,
				Cur:   cur,
				State: state,
			})
			delete(currentFiles, eta.FileName) // Remove from map
		} else {
			// File missing in current
			results = append(results, Result{
				Eta:   eta,
				Cur:   structs.Etalon{},
				State: Missing,
			})
		}
	}
	return results
}

// Prepare worker pool for hashing files
// fully automatized number of workers based on available CPUs
// returns channels for commands and results and waitgroup to sync workers
func preparePool() (chan *zip.File, chan structs.Etalon, *sync.WaitGroup) {
	commands := make(chan *zip.File, runtime.NumCPU())
	results := make(chan structs.Etalon, runtime.NumCPU()*2)

	// prepare ,,pool'' of workers
	var wg sync.WaitGroup
	for n := 0; n < runtime.NumCPU(); n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for command := range commands {
				hash := utils.CalculateHashForFile(command)
				results <- hash
			}
		}()
	}

	// wait for all workers to finish and then close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	return commands, results, &wg
}

func main() {
	if len(os.Args) < 2 {
		utils.PrintHelp()
		os.Exit(1)
	}

	showGUI := false
	if slices.Contains(os.Args, "--gui") {
		showGUI = true
	}

	zipFile := os.Args[1]
	if !utils.FileExists(zipFile) {
		utils.PrintGui(showGUI, fmt.Sprintf("The provided ZIP file: [%s] does not exist\n", zipFile))
		os.Exit(2)
	}

	config := utils.LoadConfig()

	is_handled, module, profile := utils.IsHandledFile(zipFile, config)
	if !is_handled {
		utils.PrintGui(showGUI, fmt.Sprintf("The provided ZIP file: [%s] is not a handled (check patterns)\n", zipFile))
		os.Exit(3)
	}
	fmt.Printf("ZIP_FILE:\t%s\n", zipFile)
	fmt.Printf("MODULE:\t\t%s\n", module)

	// load etalon
	etalon, etalonPath := utils.LoadEtalon(module)
	fmt.Printf("ETA_FILE:\t%s\n", etalonPath)
	fmt.Printf("Using cores:\t%d\n\n", runtime.NumCPU())

	// prepare cpu pools
	commands, results, workerWg := preparePool()

	if etalon == nil {
		// no etalon found, create one
		utils.CalculateHashes(zipFile, profile, commands, workerWg)
		utils.SaveEtalon(module, results, zipFile)
		os.Exit(0)
	}

	// compare with etalon
	utils.CalculateHashes(zipFile, profile, commands, workerWg)

	// compare and print results
	errors := 0
	for _, result := range compareCurrentWithEtalon(results, etalon) {
		switch result.State {
		case Differs:
			color.Red("! %s eta.%s\n", result.Eta.FileName, result.Eta.SHA1)
			color.Red("! %s cur.%s\n", result.Eta.FileName, result.Cur.SHA1)
			println()
			errors++

		case Missing:
			color.Yellow("- %s eta.%s\n", result.Eta.FileName, result.Eta.SHA1)
			color.Yellow("- %s cur.%s\n", result.Eta.FileName, result.Cur.SHA1)
			println()
			errors++
		}
	}

	if showGUI {
		utils.ShowGUI(zipFile, errors, len(etalon), profile)
	}

	if errors == 0 {
		color.Green("\nOK\n")
	} else {
		color.Red("\nERRORS: %d\n", errors)
		os.Exit(4)
	}
}
