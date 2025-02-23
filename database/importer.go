package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"tango/model"
	"time"
)

func (di *Database) ImportFromJSON(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var source model.JMdict
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&source); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	entriesChan := make(chan model.JMdictWord, di.batchSize)
	errorsChan := make(chan error, 1)
	var wg sync.WaitGroup

	numWorkers := 3
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go di.importJmdictEntries(entriesChan, errorsChan, &wg)
	}

	startTime := time.Now()
	for _, entry := range source.Words {
		select {
		case err := <-errorsChan:
			close(entriesChan)
			return fmt.Errorf("worker error: %v", err)
		default:
			entriesChan <- entry
		}
	}

	close(entriesChan)
	wg.Wait()

	log.Printf("Import completed. Processed %d entries in %v", len(source.Words), time.Since(startTime))
	return nil
}
