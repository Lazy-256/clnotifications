package clnotifications

import (
	"fmt"
	"io"
	"sync"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/registry"
)

const (
	COUNT_READ_KEYS = 200
	COUNT_SKIP_KEYS = 140 // !!!
	COUNT_IN_CHUNKS = 100
)

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

	fmt.Printf("Notifications key count: %d", keyInfo.ValueCount)
	log.Infof("Notifications key count: %d", keyInfo.ValueCount)
	return nil
}

func ClearValues(log *log.Entry) error {
	regKey, err := registry.OpenKey(registry.LOCAL_MACHINE,
		"SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Notifications", registry.ALL_ACCESS)
	if err != nil {
		log.Fatalf("Notifications key can't be opened: %v\n", err)
		return err
	}
	defer regKey.Close()

	keyInfo, err := regKey.Stat()
	if err != nil {
		log.Fatalf("Notifications key's properties can't be readd: %v\n", err)
		return err
	}

	// Skip them from deletion
	// Keep reading in portion of 300 keys => send to pipeline to goroutine to delete
	var wg sync.WaitGroup
	var deleted_values_count uint32

	count_read_keys := COUNT_READ_KEYS
	for {
		// Read first count_read_keys values
		values, err := regKey.ReadValueNames(count_read_keys)
		if err != nil && err != io.EOF {
			log.Fatalf("Notifications keys can't be readd: %v\n", err)
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
		fmt.Printf("%d values of %d have been read\n", len(values), keyInfo.ValueCount)
		log.Infof("%d values of %d have been read\n", len(values), keyInfo.ValueCount)

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
			fmt.Printf("chunk: %v\n", chunk)
			wg.Add(1)
			deleted_values_count += (uint32(len(chunk)))
			go func(key registry.Key, values []string, val_handled_count uint32, val_total_count uint32) {
				go delValues(key, values)
				wg.Done()
				fmt.Printf("%d values of %d\n", val_handled_count, val_total_count)
			}(regKey, chunk, deleted_values_count, keyInfo.ValueCount)
			log.Infof("%d values of %d\n", deleted_values_count, keyInfo.ValueCount)
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
			log.Fatalf("Notifications value can't be deleted: %v\n", err)
			return err
		}
		fmt.Printf("%s value have been deleted\n", val)
	}
	return nil
}
