package nomp

import (
	"strings"
	"strconv"
	"encoding/json"
)

type GlobalStat struct {
	Workers  uint16 `json:"workers"`
	Hashrate float64 `json:"hashrate"`
}

type Algo struct {
	Workers     uint16 `json:"workers"`
	Hashrate    float64 `json:"hashrate"`
	HashrateStr string `json:"hashrateString"`
}

type Stat struct {
	ValidShares   uint32 `json:"validShares"`
	ValidBlocks   uint32 `json:"validBlocks"`
	InvalidShares uint32 `json:"invalidShares"`
	TotalPaid     float64 `json:"totalPaid"`
}

func (a *Stat) UnmarshalJSON(data []byte) error {
	type Alias Stat
	aux := &struct {
		ValidShares   json.Number `json:"validShares"`
		ValidBlocks   json.Number `json:"validBlocks"`
		InvalidShares json.Number `json:"invalidShares"`
		TotalPaid     json.Number `json:"totalPaid"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if len(aux.ValidShares) == 0 {
		a.ValidShares = 0
	} else {
		val, err := strconv.ParseUint(strings.Trim(string(aux.ValidShares), "\""), 10, 32)
		if err != nil {
			return err
		}
		a.ValidShares = uint32(val)
	}
	if len(aux.ValidBlocks) == 0 {
		a.ValidBlocks = 0
	} else {
		val, err := strconv.ParseUint(strings.Trim(string(aux.ValidBlocks), "\""), 10, 32)
		if err != nil {
			return err
		}
		a.ValidBlocks = uint32(val)
	}
	if len(aux.InvalidShares) == 0 {
		a.InvalidShares = 0
	} else {
		val, err := strconv.ParseUint(strings.Trim(string(aux.InvalidShares), "\""), 10, 32)
		if err != nil {
			return err
		}
		a.InvalidShares = uint32(val)
	}
	if len(aux.TotalPaid) == 0 {
		a.TotalPaid = 0
	} else {
		val, err := strconv.ParseFloat(strings.Trim(string(aux.TotalPaid), "\""), 64)
		if err != nil {
			return err
		}
		a.TotalPaid = val
	}
	return nil
}

type Blocks struct {
	Pending   uint16 `json:"pending"`
	Confirmed uint32 `json:"confirmed"`
	Orphaned  uint32 `json:"orphaned"`
}

type Worker struct {
	Shares        float64 `json:"shares"`
	InvalidShares float64 `json:"invalidshares"`
	Hashrate      float64 `json:"hashrateString"`
}

type worker struct {
	Shares        float64 `json:"shares"`
	InvalidShares float64 `json:"invalidshares"`
	Hashrate      float64 `json:"hashrate,omitempty"`
	HashrateStr   string `json:"hashrateString"`
}

func GetHashrate(hashrateStr string) float64 {
	fields := strings.Split(hashrateStr, " ")
	if len(fields) == 0 {
		return 0.00
	}
	hashrate, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0.00
	}
	switch fields[1] {
	case "KSol/s": hashrate *= 1000
	case "KH": hashrate *= 1000
	case "MH": hashrate *= 1000 * 1000
	case "GH": hashrate *= 1000 * 1000 * 1000
	case "TH": hashrate *= 1000 * 1000 * 1000 * 1000
	case "PH": hashrate *= 1000 * 1000 * 1000 * 1000 * 1000
	}
	return hashrate
}

func (w *Worker) UnmarshalJSON(data []byte) error {
	var worker worker
	err := json.Unmarshal(data, &worker)
	if err != nil {
		return err
	}
	w.Shares = worker.Shares
	w.InvalidShares = worker.InvalidShares
	if worker.Hashrate > 0 {
		w.Hashrate = worker.Hashrate
	} else {
		w.Hashrate = GetHashrate(worker.HashrateStr)
	}
	return nil
}

type Workers map[string]Worker
type Algos map[string]Algo
type Pools map[string]Pool

type Pool struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Algorithm   string `json:"algorithm"`
	Stat        Stat `json:"poolStats"`
	Blocks      Blocks `json:"blocks"`
	Workers     Workers`json:"workers"`
	Hashrate    float64 `json:"hashrate"`
	WorkerCount uint16 `json:"workerCount"`
	HashrateStr string `json:"hashrateString"`
}

func (p *Pool) UnmarshalJSON(data []byte) error {
	type Alias Pool
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if p.Algorithm == "equihash" {
		p.Hashrate /= 500000
	}
	p.fixWorkerHashrate()
	return nil
}

/*
 * NOTE: hashrate is calculated by a following keplet:
 * hashrate = shareMultiplier * shares / website.stats.hashrateWindow
 * We can get the shareMultiplier / website.stats.hashrateWindow value from
 * divide the pool hashrate by sum of the workers shares
 */
func (p *Pool) fixWorkerHashrate() {
	var shares float64 = 0
	for _, w := range p.Workers {
		shares += w.Shares
	}
	shareMultiplier := p.Hashrate / shares
	for idx, w := range p.Workers {
		w.Hashrate = w.Shares * shareMultiplier
		p.Workers[idx] = w
	}
}

type Status struct {
	Time   uint64 `json:"time"`
	Global GlobalStat `json:"global"`
	Algos  Algos `json:"algos"`
	Pools  Pools `json:"pools"`
}

func (client *NompClient) GetPoolStatus() (Status, error) {
	poolstatus := Status{}
	_, err := client.sling.New().Get("stats").ReceiveSuccess(&poolstatus)
	if err != nil {
		return poolstatus, err
	}

	return poolstatus, err
}
