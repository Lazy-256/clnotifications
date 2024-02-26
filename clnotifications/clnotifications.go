package clnotifications

import (
	"flag"
	"io"
	"sync"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/registry"
)

// const (
var COUNT_READ_KEYS int //= 200
var COUNT_SKIP_KEYS int //= 140 // !!!
var COUNT_IN_CHUNKS int //= 100
//)

func init() {
	flag.IntVar(&COUNT_READ_KEYS, "count-values-to-read", 1000, "number of values to read in one iteration")
	flag.IntVar(&COUNT_SKIP_KEYS, "count-values-to-skip-key", 500, "number of values to skip deletion")
	flag.IntVar(&COUNT_IN_CHUNKS, "count-values-in-chunks", 100, "number of values to delete in one chunk")

	//flag.Parse()
}

func GetKeys(log *log.Entry) error {
	subRegKey, err := registry.OpenKey(registry.LOCAL_MACHINE,
		"SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Notifications", registry.READ)
	if err != nil {
		log.Fatalf("Notifications key can't be opened: %v", err)
		return err
	}

	keyInfo, err := subRegKey.Stat()
	if err != nil {
		log.Fatalf("Notifications key's properties can't be readd: %v", err)
		return err
	}

	// fmt.Printf("Notifications key count: %d", keyInfo.ValueCount)
	log.Infof("Notifications key count: %d", keyInfo.ValueCount)
	return nil
}

func ClearValues(log *log.Entry) error {
	log.Infof("number of values to read in one iteration: %d", COUNT_READ_KEYS)
	log.Infof("number of values to skip deletion: %d", COUNT_SKIP_KEYS)
	log.Infof("number of values to delete in one chunk: %d", COUNT_IN_CHUNKS)

	regKey, err := registry.OpenKey(registry.LOCAL_MACHINE,
		"SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Notifications", registry.ALL_ACCESS)
	if err != nil {
		log.Fatalf("Notifications key can't be opened: %v", err)
		return err
	}
	defer regKey.Close()

	keyInfo, err := regKey.Stat()
	if err != nil {
		log.Fatalf("Notifications key's properties can't be readd: %v", err)
		return err
	}
	log.Infof("number of values: %d", keyInfo.ValueCount)

	// Skip them from deletion
	// Keep reading in portion of 300 keys => send to pipeline to goroutine to delete
	var wg sync.WaitGroup
	var deleted_values_count uint32

	count_read_keys := COUNT_READ_KEYS
	for {
		// Read first count_read_keys values
		values, err := regKey.ReadValueNames(count_read_keys)
		if err != nil && err != io.EOF {
			log.Fatalf("Notifications keys can't be readd: %v", err)
			return err
		}
		if err == io.EOF {
			if count_read_keys-1 > COUNT_SKIP_KEYS {
				// decrement amount of read keys until it will be equal to COUNT_SKIP_KEYS
				count_read_keys = count_read_keys - 1
				continue
			} else {
				return nil
			}
		}

		// Skip COUNT_SKIP_KEYS and delete the rest in portions of COUNT_IN_CHUNKS
		values_to_delete := values[COUNT_SKIP_KEYS:len(values)]
		if len(values_to_delete) < 1 {
			break
		}

		for i := 0; i < len(values_to_delete); i += COUNT_IN_CHUNKS {
			end := i + COUNT_IN_CHUNKS
			if end > len(values_to_delete) {
				end = len(values_to_delete)
			}
			chunk := values_to_delete[i:end]
			// fmt.Printf("chunk: %v", chunk)
			wg.Add(1)
			deleted_values_count += (uint32(len(chunk)))
			go func(key registry.Key, values []string) {
				go delValues(key, values)
				wg.Done()
			}(regKey, chunk)
			log.Infof("%d values of %d handled", deleted_values_count, keyInfo.ValueCount)
		}
	}
	wg.Wait()
	return err
}

func delValues(key registry.Key, values []string) error {
	for _, val := range values {
		err := key.DeleteValue(val)
		//_, _, err := key.GetBinaryValue(val)
		if err != nil {
			log.Errorf("Notifications value can't be deleted: %v", err)
			return nil //err
		}
		// fmt.Printf("%s value have been deleted", val)
	}
	return nil
}
