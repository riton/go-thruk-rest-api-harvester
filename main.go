package main

// Copyright Rémi Ferrand <remi.ferrand_at_cc.in2p3.fr> (2018)
//
// This software is governed by the CeCILL  license under French law and
// abiding by the rules of distribution of free software.  You can  use,
// modify and/ or redistribute the software under the terms of the CeCILL
// license as circulated by CEA, CNRS and INRIA at the following URL
// "http://www.cecill.info".
//
// As a counterpart to the access to the source code and  rights to copy,
// modify and redistribute granted by the license, users are provided only
// with a limited warranty  and the software's author,  the holder of the
// economic rights,  and the successive licensors  have only  limited
// liability.
//
// In this respect, the user's attention is drawn to the risks associated
// with loading,  using,  modifying and/or developing or reproducing the
// software by the user in light of its specific status of free software,
// that may mean  that it is complicated to manipulate,  and  that  also
// therefore means  that it is reserved for developers  and  experienced
// professionals having in-depth computer knowledge. Users are therefore
// encouraged to load and test the software's suitability as regards their
// requirements in conditions enabling the security of their systems and/or
// data to be ensured and,  more generally, to use and operate it in the
// same conditions as regards security.
//
// The fact that you are presently reading this means that you have had
// knowledge of the CeCILL license and that you accept its terms.

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	nbWorker, _ := strconv.Atoi(getEnv("WORKERS", "5"))
	workerMaxWait, _ := strconv.Atoi(getEnv("WORKER_MAX_WAIT", "5"))
	httpTimeout, _ := strconv.Atoi(getEnv("HTTP_TIMEOUT", "3"))
	user := getEnv("THRUK_USER", "thrukadmin")
	password := getEnv("THRUK_PASSWORD", "thrukadmin")
	thrukEndpoint := getEnv("THRUK_ENDPOINT", "http://localhost:8080/thruk")

	c := cli{
		nbWorker:      nbWorker,
		workerMaxWait: workerMaxWait,
		user:          user,
		password:      password,
		thrukEndpoint: thrukEndpoint,
		httpTimeout:   time.Duration(httpTimeout) * time.Second,
	}

	log.Fatal(c.Run())
}

type cli struct {
	nbWorker      int
	workerMaxWait int
	user          string
	password      string
	thrukEndpoint string
	httpTimeout   time.Duration
}

func (c cli) Run() error {

	var wg sync.WaitGroup
	wg.Add(c.nbWorker)

	for i := 0; i < c.nbWorker; i++ {
		go func(idx int) {
			c.queryRestAPI(idx)
		}(i)
	}

	wg.Wait()

	return nil
}

func (c cli) getThrukRestApiURL() string {
	return c.thrukEndpoint + "/r/"
}

func (c cli) getHttpRequest() *http.Request {

	hostReq := "hosts?columns=name"

	req, err := http.NewRequest("GET", c.getThrukRestApiURL()+hostReq, nil)
	if err != nil {
		log.Fatalf("building http request: %v", err)
	}

	req.SetBasicAuth(c.user, c.password)

	return req
}

func (c cli) getHttpClient() *http.Client {
	return &http.Client{
		Timeout: c.httpTimeout,
	}
}

type apiResponse []struct {
	Name string
}

func (c cli) queryRestAPI(workerIdx int) {
	hc := c.getHttpClient()

	wlog := log.WithFields(log.Fields{
		"worker-id": workerIdx,
	})

	randomInt := rand.Intn(c.workerMaxWait)
	time.Sleep(time.Duration(randomInt) * time.Second)

	for {
		req := c.getHttpRequest()

		resp, err := hc.Do(req)
		if err != nil {
			wlog.Errorf("error querying host: %v / %T", err, err)
			continue
		}

		var ar apiResponse

		jd := json.NewDecoder(resp.Body)
		if err := jd.Decode(&ar); err != nil {
			wlog.Errorf("error decoding API response: %v / %T", err, err)
			continue
		}

		var hosts []string
		for _, host := range ar {
			hosts = append(hosts, host.Name)
		}

		wlog.Infof("hosts = %s", strings.Join(hosts, ","))

		randomInt = rand.Intn(c.workerMaxWait)

		time.Sleep(time.Duration(randomInt) * time.Second)
	}
}
