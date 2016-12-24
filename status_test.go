package nomp

import (
	"fmt"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetPoolStatus(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	sampleItem := `{
			   "time": 1474239882,
			   "global": {
			      "workers": 21,
			      "hashrate": 0
			   },
			   "algos": {
			      "test1": {
			         "workers": 1,
			         "hashrate": 2433814.801066667,
			         "hashrateString": "2.43 MH"
			      },
			      "test2": {
			         "workers": 0,
			         "hashrate": 0,
			         "hashrateString": "0.00 KH"
			      }
			   },
			   "pools": {
			      "test1": {
			         "name": "test1",
			         "symbol": "TEST1",
			         "algorithm": "test1",
			         "poolStats": {
			            "validShares": 0,
			            "validBlocks": 0,
			            "invalidShares": "1359059",
			            "totalPaid": "13579727.61959752997063333"
			         },
			         "blocks": {
			            "pending": 0,
			            "confirmed": 6769,
			            "orphaned": 0
			         },
			         "workers": {
			            "worker1": {
			               "shares": 0.17,
			               "invalidshares": 0,
			               "hashrateString": "2.43 MH"
			            }
			         },
			         "hashrate": 2433814.801066667,
			         "workerCount": 1,
			         "hashrateString": "2.43 MH"
			      },
			      "test2": {
			         "name": "test2",
			         "symbol": "TEST2",
			         "algorithm": "test2",
			         "poolStats": {
			            "validShares": "15402335",
			            "validBlocks": "3966",
			            "invalidShares": "388455",
			            "totalPaid": "4591548.98264059998791708"
			         },
			         "blocks": {
			            "pending": 0,
			            "confirmed": 3527,
			            "orphaned": 0
			         },
			         "workers": {
			         },
			         "hashrate": 0,
			         "workerCount": 0,
			         "hashrateString": "0.00 KH"
			      }

			   }
			}`

	expectedItem := Status{
		Time: 1474239882,
		Global: GlobalStat{
			Workers: 21,
			Hashrate: 0,
		},
		Algos: Algos{
			"test1": Algo{
				Workers: 1,
				Hashrate: 2433814.801066667,
				HashrateStr: "2.43 MH",
			},
			"test2": Algo{
				Workers: 0,
				Hashrate: 0,
				HashrateStr: "0.00 KH",
			},
		},
		Pools: Pools{
			"test1": Pool{
				Name: "test1",
				Symbol: "TEST1",
				Algorithm: "test1",
				Stat: Stat{
					ValidShares: 0,
					ValidBlocks: 0,
					InvalidShares: 1359059,
					TotalPaid: 13579727.61959752997063333,
				},
				Blocks: Blocks{
					Pending: 0,
					Confirmed: 6769,
					Orphaned: 0,
				},
				Workers: Workers{
					"worker1": Worker{
						Shares: 0.17,
						InvalidShares: 0,
						Hashrate: 2430000,
					},
				},
				Hashrate: 2433814.801066667,
				WorkerCount: 1,
				HashrateStr: "2.43 MH",
			},
			"test2": Pool{
				Name: "test2",
				Symbol: "TEST2",
				Algorithm: "test2",
				Stat: Stat{
					ValidShares: 15402335,
					ValidBlocks: 3966,
					InvalidShares: 388455,
					TotalPaid: 4591548.98264059998791708,
				},
				Blocks: Blocks{
					Pending: 0,
					Confirmed: 3527,
					Orphaned: 0,
				},
				Workers: Workers{
				},
				Hashrate: 0,
				WorkerCount: 0,
				HashrateStr: "0.00 KH",
			},
		},
	}

	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, sampleItem)
	})

	mposClient := NewNompClient(httpClient, "http://dummy.com/", "")
	poolstatus, err := mposClient.GetPoolStatus()

	assert.Nil(t, err)
	assert.Equal(t, expectedItem, poolstatus)
}
