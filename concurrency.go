package main

import (
	"log"
	"sync"
	"time"
)

func backup(server string, sem chan struct{}, inProgMap *map[string]bool, round int, wg *sync.WaitGroup, mutex *sync.RWMutex) {
	sem <- struct{}{}
	defer func() {
		<-sem

		wg.Done()
	}()
	mutex.RLock()
	if inProg, ok := (*inProgMap)[server]; ok && inProg {
		mutex.RUnlock()
		return
	}
	mutex.RUnlock()
	mutex.Lock()
	(*inProgMap)[server] = true
	mutex.Unlock()
	log.Printf("round %v: starting backup for server %v", round, server)
	time.Sleep(5 * time.Second)
	log.Printf("round %v: finish backup for server %v", round, server)
}

func metric(server string, sem chan struct{}, round int, wg *sync.WaitGroup) {
	sem <- struct{}{}
	defer func() {
		<-sem
		wg.Done()
	}()
	log.Printf("round %v: starting metric for server %v", round, server)
	time.Sleep(1 * time.Second)
	log.Printf("round %v: finish metric for server %v", round, server)

}

func main() {

	var (
		backupStartedMap      = map[string]bool{}
		backupStartedMapMutex = sync.RWMutex{}
		backupWg              sync.WaitGroup
	)
	semb := make(chan struct{}, 2)
	semm := make(chan struct{}, 4)
	defer close(semb)
	defer close(semm)
	for i := 0; i < 6; i++ {

		var wg sync.WaitGroup
		servers := []string{"a", "b", "c", "d", "e"}
		for _, server := range servers {
			backupStartedMapMutex.RLock()
			started, ok := backupStartedMap[server]
			backupStartedMapMutex.RUnlock()
			if !ok || !started {
				log.Printf("round %v: server: %v, started: %v ok: %v ", i, server, started, ok)
				backupWg.Add(1)
				go backup(server, semb, &backupStartedMap, i, &backupWg, &backupStartedMapMutex)
			}
			wg.Add(1)
			go metric(server, semm, i, &wg)
		}

		wg.Wait()

	}
	backupWg.Wait()
	log.Printf("finished!!!!!!")
}

