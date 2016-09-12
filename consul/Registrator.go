package consul

import consulapi "github.com/hashicorp/consul/api"
import "time"
import "fmt"
import "math/rand"
import "github.com/zaunerc/cntinsight/system"
import "strconv"

/**
 * ScheduleRegistration return immediately after the
 * container registration job is scheduled.
 */
func ScheduleRegistration(consulUrl string) {
	fmt.Printf("Scheduling registration task using consul URL >%s<.\n", consulUrl)
	serviceId := RandStringBytesMaskImprSrc(8)
	go registerContainer(consulUrl, 5, serviceId)
}

func registerContainer(consulUrl string, sleepSeconds int, serviceId string) {
	for {
		fmt.Printf("Registering container...\n")

		config := consulapi.DefaultConfig()
		config.Address = consulUrl
		consul, err := consulapi.NewClient(config)

		if err != nil {
			fmt.Printf("Error while trying to register container: %s\n", err)
			return
		}

		kv := consul.KV()

		// cntrInfodUrl
		data := &consulapi.KVPair{Key: "containers/" + serviceId + "/cntrInfodUrl",
			Value: []byte("XXXXX")}
		_, err = kv.Put(data, nil)
		if err != nil {
			fmt.Printf("Error while trying to register container: %s\n", err)
			return
		}

		// MAC
		data = &consulapi.KVPair{Key: "containers/" + serviceId + "/macAdress",
			Value: []byte("XXXXX")}
		_, err = kv.Put(data, nil)
		if err != nil {
			fmt.Printf("Error while trying to register container: %s\n", err)
			return
		}

		// IP Adress
		data = &consulapi.KVPair{Key: "containers/" + serviceId + "/ipAdress",
			Value: []byte("XXXXX")}
		_, err = kv.Put(data, nil)
		if err != nil {
			fmt.Printf("Error while trying to register container: %s\n", err)
			return
		}

		// Unix Epoch Timestamp
		unixEpochTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
		data = &consulapi.KVPair{Key: "containers/" + serviceId + "/unixEpochTimestamp",
			Value: []byte(unixEpochTimestamp)}
		_, err = kv.Put(data, nil)
		if err != nil {
			fmt.Printf("Error while trying to register container: %s\n", err)
			return
		}

		// Hostname
		hostname := system.FetchContainerHostname()
		data = &consulapi.KVPair{Key: "containers/" + serviceId + "/hostname",
			Value: []byte(hostname)}
		_, err = kv.Put(data, nil)
		if err != nil {
			fmt.Printf("Error while trying to register container: %s\n", err)
			return
		}

		fmt.Printf("Successfully registered container. Next registration in %d seconds...\n", sleepSeconds)
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
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
