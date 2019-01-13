// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	httpstat "github.com/tcnksm/go-httpstat"
)

var Version string
var cfgFile string
var provider string
var regions []string
var timeout int
var limit int
var output string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloudping",
	Short: "Returns the geographically closest region.",
	Long: `cloudping identifies the cloud provider regions geographically closest
and returns them in order of lowest to highest latency.`,
	Run: func(cmd *cobra.Command, args []string) {
		rs := endpoints.AwsPartition().Services()[dynamodb.EndpointsID].Regions()

		if len(regions) > 0 {
			rs = makeEndpoints(regions)
		}

		ch := make(chan pingResult)
		pingResults := make(map[string]int, len(rs))

		for region := range rs {
			go func(ch chan<- pingResult, region string) {
				ch <- pingRegion(region)
			}(ch, region)
		}

		for i := 0; i < len(rs); i++ {
			pr := <-ch
			pingResults[pr.region] = pr.latency
		}

		close(ch)

		orderedResults := sortList(pingResults)

		for index, result := range orderedResults {
			if limit > 0 && index >= limit {
				break
			}

			println(result.Key)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cloudping.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVar(&provider, "provider", "aws", "Cloud provider")
	rootCmd.Flags().StringSliceVar(&regions, "regions", nil, "Limits checks to specific regions")
	rootCmd.Flags().IntVar(&limit, "limit", 0, "Limits the number of regions returned")
	rootCmd.Flags().StringVar(&output, "output", "txt", "Output format. One of: txt, json, yaml")
	rootCmd.Flags().IntVar(&timeout, "timeout", 500, "Timeout for each region in milliseconds")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cloudping" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cloudping")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

type pingResult struct {
	region  string
	latency int
}

// pingRegion returns the total request duration
func pingRegion(region string) pingResult {
	uri := fmt.Sprintf("https://dynamodb.%s.amazonaws.com/ping?x=%s", region, randStringBytes(12))

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Create a httpstat powered context
	var result httpstat.Result

	ctx1, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeout))
	defer cancel()
	req = req.WithContext(ctx1)

	ctx2 := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx2)

	// Send request by default HTTP client
	client := http.DefaultClient
	if _, err := client.Do(req); err != nil {
		return pingResult{region: region, latency: 99999}
	}

	//end := time.Now()
	//
	//// Show the results
	//log.Printf("%s - DNS lookup: %d ms", uri, int(result.DNSLookup/time.Millisecond))
	//log.Printf("%s - TCP connection: %d ms", uri, int(result.TCPConnection/time.Millisecond))
	//log.Printf("%s - TLS handshake: %d ms", uri, int(result.TLSHandshake/time.Millisecond))
	//log.Printf("%s - Server processing: %d ms", uri, int(result.ServerProcessing/time.Millisecond))
	//log.Printf("%s - Content transfer: %d ms", uri, int(result.ContentTransfer(time.Now())/time.Millisecond))
	//log.Printf("%s - Content transfer: %d ms", uri, int(result.Total(end)))

	return pingResult{region: region, latency: int(result.TCPConnection / time.Millisecond)}
}

// A data structure to hold key/value pairs
type Pair struct {
	Key   string
	Value int
}

// A slice of pairs that implements sort.Interface to sort by values
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func sortList(noble map[string]int) PairList {
	p := make(PairList, len(noble))

	i := 0
	for k, v := range noble {
		p[i] = Pair{k, v}
		i++
	}

	sort.Sort(p)

	return p
}

// makeEndpoints creates a map based on a string list
func makeEndpoints(regions []string) map[string]endpoints.Region {
	rs := map[string]endpoints.Region{}

	for _, id := range regions {
		rs[id] = endpoints.Region{}
	}

	return rs
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
