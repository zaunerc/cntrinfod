package consul

import consulapi "github.com/hashicorp/consul/api"
import "time"
import "fmt"
import "math/rand"
import "github.com/zaunerc/cntrinfod/system"
import "github.com/zaunerc/cntrinfod/docker"
import "strconv"

/**
 * ScheduleRegistration return immediately after the
 * container registration job is scheduled.
 */
func ScheduleRegistration(consulUrl string, cntrInfodHttpPort int) {
	serviceId := RandStringBytesMaskImprSrc(8)
	fmt.Printf("Scheduling registration task using consul URL >%s< and service id >%s<.\n", consulUrl, serviceId)
	go registerContainer(consulUrl, cntrInfodHttpPort, 5, serviceId)
}

func registerContainer(consulUrl string, cntrInfodHttpPort int, sleepSeconds int, serviceId string) {
	firstIteration := true
	for {
		if firstIteration {
			firstIteration = false
		} else {
			time.Sleep(time.Duration(sleepSeconds) * time.Second)
		}

		fmt.Printf("Registering container...\n")

		config := consulapi.DefaultConfig()
		config.Address = consulUrl
		consul, err := consulapi.NewClient(config)

		if err != nil {
			fmt.Printf("Error while trying to register container: %s\n", err)
			continue
		}

		kv := consul.KV()

		// cntrInfodUrl
		cntrInfodHttpUrl := "http://" + system.FetchContainerHostname() + ":" + strconv.Itoa(cntrInfodHttpPort)
		data := &consulapi.KVPair{Key: "containers/" + serviceId + "/cntrInfodHttpUrl",
			Value: []byte(cntrInfodHttpUrl)}
		_, err = kv.Put(data, nil)
		if err != nil {
			fmt.Printf("Error while trying to register container: %s\n", err)
			continue
		}

		// MAC
		macAdress := system.FetchFirstMac()
		data = &consulapi.KVPair{Key: "containers/" + serviceId + "/macAdress",
			Value: []byte(macAdress)}
		_, err = kv.Put(data, nil)
		if err != nil {
			fmt.Printf("Error while trying to register container: %s\n", err)
			continue
		}

		// IP Adress
		ipAdress := system.FetchFirstIp()
		data = &consulapi.KVPair{Key: "containers/" + serviceId + "/ipAdress",
			Value: []byte(ipAdress)}
		_, err = kv.Put(data, nil)
		if err != nil {
			fmt.Printf("Error while trying to register container: %s\n", err)
			continue
		}

		// Unix Epoch Timestamp
		unixEpochTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
		data = &consulapi.KVPair{Key: "containers/" + serviceId + "/unixEpochTimestamp",
			Value: []byte(unixEpochTimestamp)}
		_, err = kv.Put(data, nil)
		if err != nil {
			fmt.Printf("Error while trying to register container: %s\n", err)
			continue
		}

		// Hostname
		hostname := system.FetchContainerHostname()
		data = &consulapi.KVPair{Key: "containers/" + serviceId + "/hostname",
			Value: []byte(hostname)}
		_, err = kv.Put(data, nil)
		if err != nil {
			fmt.Printf("Error while trying to register container: %s\n", err)
			continue
		}

		// HostHostname
		hostHostname := docker.FetchHostHostname()
		data = &consulapi.KVPair{Key: "containers/" + serviceId + "/hostinfo/hostname",
			Value: []byte(hostHostname)}
		_, err = kv.Put(data, nil)
		if err != nil {
			fmt.Printf("Error while trying to register container: %s\n", err)
			continue
		}

		fmt.Printf("Successfully registered container. Next registration in %d seconds...\n", sleepSeconds)
	}
}

/*
 * RandStringBytesMaskImprSrc acquires a new seed each time it is called.
 * Inspired by http://stackoverflow.com/a/31832326/6551760.
 */
func RandStringBytesMaskImprSrc(n int) string {

	var src = rand.NewSource(time.Now().UnixNano())

	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
