package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var files []string

func VisitFile(fp string, fi os.FileInfo, err error) error {
	if err != nil {
		log.Println(err) // can't walk here,
		return nil       // but continue walking elsewhere
	}
	if strings.HasSuffix(fp, ".json") || strings.HasSuffix(fp, ".yml") || strings.HasSuffix(fp, ".properties") || strings.HasSuffix(fp, ".txt") {
		if strings.HasPrefix(fp, "server/world/stats") {
			return nil
		} else {
			files = append(files, fp)
		}
	}
	return nil
}

func loadConfig() {
	filepath.Walk("server", VisitFile)
}
