package main

/*import (
	"fmt"
	"github.com/d2r2/go-dht"
	"log"
)*/

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

/*func main2() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	s := "gopher"
	fmt.Printf("Hello and welcome, %s!\n", s)

	for i := 1; i <= 5; i++ {
		//TIP <p>To start your debugging session, right-click your code in the editor and select the Debug option.</p> <p>We have set one <icon src="AllIcons.Debugger.Db_set_breakpoint"/> breakpoint
		// for you, but you can always add more by pressing <shortcut actionId="ToggleLineBreakpoint"/>.</p>
		fmt.Println("i =", 100/i)
	}
}

func main2() {
	// Read DHT11 sensor data from pin 4, retrying 10 times in case of failure.
	// You may enable "boost GPIO performance" parameter, if your device is old
	// as Raspberry PI 1 (this will require root privileges). You can switch off
	// "boost GPIO performance" parameter for old devices, but it may increase
	// retry attempts. Play with this parameter.
	// Note: "boost GPIO performance" parameter is not work anymore from some
	// specific Go release. Never put true value here.
	temperature, humidity, retried, err :=
		dht.ReadDHTxxWithRetry(dht.DHT11, 4, false, 10)
	if err != nil {
		log.Fatal(err)
	}
	// Print temperature and humidity
	fmt.Printf("Temperature = %v*C, Humidity = %v%% (retried %d times)\n",
		temperature, humidity, retried)
}*/

import (
	"fmt"
	"os"
	"strconv"
	"time"

	//"github.com/d2r2/go-dht"

	"github.com/MichaelS11/go-dht"

	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*func main1() {
	fmt.Println("init ...")
	err := dht.HostInit()
	if err != nil {
		fmt.Println("HostInit error:", err)
		return
	}
	fmt.Println("init ok")

	fmt.Println("dht ...")
	//dht2, err := dht.NewDHT("GPIO19", dht.Fahrenheit, "")
	dht2, err := dht.NewDHT("GPIO4", dht.Celsius, "")
	if err != nil {
		fmt.Println("NewDHT error:", err)
		return
	}

	fmt.Println("temp ...")
	humidity, temperature, err := dht2.ReadRetry(11)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Printf("humidity: %v\n", humidity)
	fmt.Printf("temperature: %v\n", temperature)
}*/

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
	temp = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "temperature",
		Subsystem: "temperature",
		Name:      "temperature",
		Help:      "La temperature",
	})
	humidite = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "humidite",
		Subsystem: "humidite",
		Name:      "humidite",
		Help:      "L'humidite",
	})
)

func temperatureMetrics(sleepTime int, humidity *float64, temperature *float64) {
	go func() {
		for {

			// should have at least read the sensor twice after 30 seconds
			time.Sleep(time.Duration(sleepTime) * time.Second)

			fmt.Printf("humidity: %v\n", *humidity)
			fmt.Printf("temperature: %v\n", *temperature)

			temp.Set(*temperature)

			humidite.Set(*humidity)

		}

	}()
}

func main() {

	var pin = "GPIO4"
	var sleepTime = 30

	if len(os.Args) > 1 {
		s := os.Args[1]
		if s != "-" {
			pin = s
		}
	}
	if len(os.Args) > 2 {
		if os.Args[2] != "-" {
			n, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Println("strconv error:", err)
				os.Exit(1)
			}
			if n > 0 {
				sleepTime = n
			}
		}
	}

	err := dht.HostInit()
	if err != nil {
		fmt.Println("HostInit error:", err)
		os.Exit(1)
	}

	//dht2, err := dht.NewDHT("GPIO19", dht.Fahrenheit, "")
	//dht2, err := dht.NewDHT(pin, dht.Fahrenheit, "")
	dht2, err := dht.NewDHT(pin, dht.Celsius, "")
	if err != nil {
		fmt.Println("NewDHT error:", err)
		os.Exit(1)
	}

	stop := make(chan struct{})
	stopped := make(chan struct{})
	var humidity float64
	var temperature float64

	// get sensor reading every 20 seconds in background
	go dht2.ReadBackground(&humidity, &temperature, 20*time.Second, stop, stopped)

	recordMetrics()

	temperatureMetrics(sleepTime, &humidity, &temperature)

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(":2112", nil)
	if err != nil {
		fmt.Println("error for open server:", err)
		os.Exit(1)
	}

	// to stop ReadBackground after done with reading, close the stop channel
	close(stop)

	// can check stopped channel to know when ReadBackground has stopped
	<-stopped
}

/*func main3() {
	err := dht.HostInit()
	if err != nil {
		fmt.Println("HostInit error:", err)
		return
	}

	dht2, err := dht.NewDHT("GPIO4", dht.Fahrenheit, "")
	if err != nil {
		fmt.Println("NewDHT error:", err)
		return
	}

	humidity, temperature, err := dht2.Read()
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Printf("humidity: %v\n", humidity)
	fmt.Printf("temperature: %v\n", temperature)
}*/
