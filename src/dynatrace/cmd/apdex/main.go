package main

import (
	"dynatrace"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	ring "github.com/zfjagann/golang-ring"
)

const metricName = "apdex"
const metricDisplayName = "Mobile Apdex"
const metricType = "MobileDevice"
const deviceID = "MobileApdexCalculator"
const deviceType = metricType
const deviceDisplayName = "Mobile Apdex Calculator"
const ipAddress = "10.0.0.10"

type experience int

const (
	SATISFIED experience = iota
	TOLERATING
	FRUSTRATED
)

type apdexAction struct {
	StartTime      uint64
	ResponseTime   uint64
	UserExperience experience
}

var actionsring = ring.Ring{}
var (
	targetTime        int
	applicationFilter string
	visitCount        int
	actionCount       int
)

func processRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "visitCount:", visitCount)
		fmt.Fprintln(w, "actionCount:", actionCount)
		bufferSize := getBufferSize()
		fmt.Fprintln(w, "bufferSize:", bufferSize)
		if bufferSize > 0 {
			fmt.Fprintln(w, "apdex:", calculateApdex())
		}
		return
	}

	visits := dynatrace.Parse(r.Body)
	for _, v := range visits {
		visitCount++
		for _, a := range v.UserActions {
			processUserAction(a)
		}
	}

	fmt.Fprintln(w, "OK")
}

func processUserAction(action dynatrace.UserAction) {
	if len(applicationFilter) > 0 && applicationFilter != action.Application {
		return
	}
	actionCount++
	a := apdexAction{
		StartTime:    action.StartTime,
		ResponseTime: action.EndTime - action.StartTime,
	}
	if a.ResponseTime < (uint64)(targetTime*1000) {
		a.UserExperience = SATISFIED
	} else if a.ResponseTime < (uint64)(targetTime*1000*4) {
		a.UserExperience = TOLERATING
	} else {
		a.UserExperience = FRUSTRATED
	}

	actionsring.Enqueue(a)
}

func calculateApdex() float32 {
	var satisfiedCount, toleratingCount, totalCount int
	for _, a := range actionsring.Values() {
		totalCount++
		action := a.(apdexAction)
		if action.UserExperience == SATISFIED {
			satisfiedCount++
		} else if action.UserExperience == TOLERATING {
			toleratingCount++
		}
	}

	if totalCount == 0 {
		return 0
	}

	return (float32(satisfiedCount) + (float32(toleratingCount) / 2)) / float32(totalCount)
}

func getBufferSize() int {
	return len(actionsring.Values())
}

func postData(api dynatrace.API) {
	dp := dynatrace.DataPoint{
		CustomDeviceID:    deviceID,
		IPAddress:         ipAddress,
		DeviceDisplayName: deviceDisplayName,
		MetricDisplayName: metricDisplayName,
		DeviceType:        deviceType,
		MetricName:        metricName,
	}
	for {
		time.Sleep(60 * time.Second)

		apdex := calculateApdex()
		dp.Timestamp = time.Now().Unix() * 1000
		dp.Value = apdex
		respBody, err := api.ReportDataPoint(dp)
		if err != nil {
			log.Println("Error while reporting data point:", err)
			continue
		}
		log.Println(respBody)
		log.Println("Successfully POSTed data point with apdex", apdex)
	}
}

func main() {
	var (
		port           int
		ringBufferSize int
		serverURL      string
		apiToken       string
	)

	flag.IntVar(&port, "port", 8080, "HTTP listener port")
	flag.IntVar(&targetTime, "targettime", 2, "Satisfied threshold in seconds")
	flag.IntVar(&ringBufferSize, "ringbuffersize", 1000, "Ring buffer size")
	flag.StringVar(&applicationFilter, "application", "", "Only calculate apdex for this particular application - calculate for all applications if blank")
	flag.StringVar(&serverURL, "serverurl", "", "Dynatrace server URL")
	flag.StringVar(&apiToken, "apitoken", "", "API token")

	flag.Parse()

	if len(serverURL) == 0 {
		missingParameter("serverurl")
	}

	log.Println("Server URL:", serverURL)

	if len(apiToken) == 0 {
		missingParameter("apitoken")
	}

	log.Println("API Token:", apiToken)

	actionsring.SetCapacity(ringBufferSize)

	if len(applicationFilter) > 0 {
		log.Println("Filtering on application:", applicationFilter)
	} else {
		log.Println("No application filter - will calculate apdex for all applications")
	}

	api := dynatrace.NewAPI(serverURL, apiToken)
	log.Println("Checking if metric", metricName, "exists...")
	ok, err := api.CustomMetricExists(metricName)
	if err != nil {
		log.Fatalln(err)
	}

	if !ok {
		log.Println("Metric", metricName, "does not exist.")
		log.Println("Creating custom metric", metricName, "...")
		err := api.RegisterCustomMetric(metricName, metricDisplayName, "Ratio", metricType)
		if err != nil {
			log.Fatalln("Could not create custom metric", err)
		}
		log.Println("Custom metric successfully created")
	} else {
		log.Println("Custom metric exists")
	}

	go postData(api)

	log.Println("Listening on port", port)
	http.HandleFunc("/", processRequest)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal("Could not bind to listener port", err)
	}
}

func missingParameter(name string) {
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "Mandatory parameter", name, "is missing.")
	os.Exit(1)
}
