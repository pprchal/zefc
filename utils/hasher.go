package utils

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"io"
	"sync"
	structs "zefc/structs"
)

func CalculateHashes(zipPath string, profile structs.Profile, commands chan *zip.File, workerWg *sync.WaitGroup) {
	archive, err := zip.OpenReader(zipPath)
	if err != nil {
		panic(err)
	}

	// Send all files to workers, then wait for them to finish before closing archive
	go func() {
		defer archive.Close()

		// First, send all files to workers
		for _, f := range archive.File {
			if !canBeFileHashed(f, profile) {
				continue
			}
			commands <- f
		}
		close(commands)

		// Wait for all workers to finish before closing archive
		workerWg.Wait()
	}()
}

func CalculateHashForFile(f *zip.File) structs.Etalon {
	fileInArchive, err := f.Open()
	if err != nil {
		panic("Cannot open file " + f.Name + " in archive: " + err.Error())
	}
	defer fileInArchive.Close()
	fmt.Printf("%s\n", f.Name)

	sha1Hasher := sha256.New()
	if _, err := io.Copy(sha1Hasher, fileInArchive); err != nil {
		panic("Cannot create SHA1 hash for file " + f.Name + ": " + err.Error())
	}

	sha1hash := sha1Hasher.Sum(nil)
	return structs.Etalon{
		FileName: f.Name,
		SHA1:     fmt.Sprintf("%x", sha1hash),
		Date:     f.Modified,
		Size:     int64(f.UncompressedSize64),
	}
}

func canBeFileHashed(f *zip.File, p structs.Profile) bool {
	// fail not-so-fast if no accept patterns matched
	// (give a chance to match reject patterns)
	for _, accept := range p.AcceptREs {
		if !accept.MatchString(f.Name) {
			return false
		}
	}

	// finally check reject patterns
	for _, reject := range p.RejectREs {
		if reject.MatchString(f.Name) {
			return false
		}
	}

	return true
}
