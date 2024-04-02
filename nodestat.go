package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetNodeStat(port int, nodeName string) (Stat, error) {
	resp, err := http.DefaultClient.Get(fmt.Sprintf("http://localhost:%d/api/v1/nodes/%s/proxy/stats/summary", port, nodeName))
	if err != nil {
		return Stat{}, err
	}
	defer resp.Body.Close()
	stat := Stat{}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return Stat{}, err
	}
	err = json.Unmarshal(raw, &stat)
	if err != nil {
		return Stat{}, err
	}
	return stat, err
}

type Stat struct {
	Node struct {
		NodeName         string `json:"nodeName"`
		SystemContainers []struct {
			Name      string    `json:"name"`
			StartTime time.Time `json:"startTime"`
			Cpu       struct {
				Time                 time.Time `json:"time"`
				UsageNanoCores       int       `json:"usageNanoCores"`
				UsageCoreNanoSeconds int64     `json:"usageCoreNanoSeconds"`
			} `json:"cpu"`
			Memory struct {
				Time            time.Time `json:"time"`
				AvailableBytes  int64     `json:"availableBytes,omitempty"`
				UsageBytes      int64     `json:"usageBytes"`
				WorkingSetBytes int       `json:"workingSetBytes"`
				RssBytes        int       `json:"rssBytes"`
				PageFaults      int       `json:"pageFaults"`
				MajorPageFaults int       `json:"majorPageFaults"`
			} `json:"memory"`
		} `json:"systemContainers"`
		StartTime time.Time `json:"startTime"`
		Cpu       struct {
			Time                 time.Time `json:"time"`
			UsageNanoCores       int       `json:"usageNanoCores"`
			UsageCoreNanoSeconds int64     `json:"usageCoreNanoSeconds"`
		} `json:"cpu"`
		Memory struct {
			Time            time.Time `json:"time"`
			AvailableBytes  int64     `json:"availableBytes"`
			UsageBytes      int64     `json:"usageBytes"`
			WorkingSetBytes int       `json:"workingSetBytes"`
			RssBytes        int       `json:"rssBytes"`
			PageFaults      int       `json:"pageFaults"`
			MajorPageFaults int       `json:"majorPageFaults"`
		} `json:"memory"`
		Network struct {
			Time       time.Time `json:"time"`
			Name       string    `json:"name"`
			RxBytes    int64     `json:"rxBytes"`
			RxErrors   int       `json:"rxErrors"`
			TxBytes    int64     `json:"txBytes"`
			TxErrors   int       `json:"txErrors"`
			Interfaces []struct {
				Name     string `json:"name"`
				RxBytes  int64  `json:"rxBytes"`
				RxErrors int    `json:"rxErrors"`
				TxBytes  int64  `json:"txBytes"`
				TxErrors int    `json:"txErrors"`
			} `json:"interfaces"`
		} `json:"network"`
		Fs struct {
			Time           time.Time `json:"time"`
			AvailableBytes int64     `json:"availableBytes"`
			CapacityBytes  int64     `json:"capacityBytes"`
			UsedBytes      int64     `json:"usedBytes"`
			InodesFree     int       `json:"inodesFree"`
			Inodes         int       `json:"inodes"`
			InodesUsed     int       `json:"inodesUsed"`
		} `json:"fs"`
		Runtime struct {
			ImageFs struct {
				Time           time.Time `json:"time"`
				AvailableBytes int64     `json:"availableBytes"`
				CapacityBytes  int64     `json:"capacityBytes"`
				UsedBytes      int64     `json:"usedBytes"`
				InodesFree     int       `json:"inodesFree"`
				Inodes         int       `json:"inodes"`
				InodesUsed     int       `json:"inodesUsed"`
			} `json:"imageFs"`
		} `json:"runtime"`
		Rlimit struct {
			Time    time.Time `json:"time"`
			Maxpid  int       `json:"maxpid"`
			Curproc int       `json:"curproc"`
		} `json:"rlimit"`
	} `json:"node"`
	Pods []struct {
		PodRef struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
			Uid       string `json:"uid"`
		} `json:"podRef"`
		StartTime  time.Time `json:"startTime"`
		Containers []struct {
			Name      string    `json:"name"`
			StartTime time.Time `json:"startTime"`
			Cpu       struct {
				Time                 time.Time `json:"time"`
				UsageNanoCores       int       `json:"usageNanoCores"`
				UsageCoreNanoSeconds int64     `json:"usageCoreNanoSeconds"`
			} `json:"cpu"`
			Memory struct {
				Time            time.Time `json:"time"`
				UsageBytes      int       `json:"usageBytes"`
				WorkingSetBytes int       `json:"workingSetBytes"`
				RssBytes        int       `json:"rssBytes"`
				PageFaults      int       `json:"pageFaults"`
				MajorPageFaults int       `json:"majorPageFaults"`
				AvailableBytes  int       `json:"availableBytes,omitempty"`
			} `json:"memory"`
			Rootfs struct {
				Time           time.Time `json:"time"`
				AvailableBytes int64     `json:"availableBytes"`
				CapacityBytes  int64     `json:"capacityBytes"`
				UsedBytes      int       `json:"usedBytes"`
				InodesFree     int       `json:"inodesFree"`
				Inodes         int       `json:"inodes"`
				InodesUsed     int       `json:"inodesUsed"`
			} `json:"rootfs"`
			Logs struct {
				Time           time.Time `json:"time"`
				AvailableBytes int64     `json:"availableBytes"`
				CapacityBytes  int64     `json:"capacityBytes"`
				UsedBytes      int       `json:"usedBytes"`
				InodesFree     int       `json:"inodesFree"`
				Inodes         int       `json:"inodes"`
				InodesUsed     int       `json:"inodesUsed"`
			} `json:"logs"`
		} `json:"containers"`
		Cpu struct {
			Time                 time.Time `json:"time"`
			UsageNanoCores       int       `json:"usageNanoCores"`
			UsageCoreNanoSeconds int64     `json:"usageCoreNanoSeconds"`
		} `json:"cpu"`
		Memory struct {
			Time            time.Time `json:"time"`
			UsageBytes      int       `json:"usageBytes"`
			WorkingSetBytes int       `json:"workingSetBytes"`
			RssBytes        int       `json:"rssBytes"`
			PageFaults      int       `json:"pageFaults"`
			MajorPageFaults int       `json:"majorPageFaults"`
			AvailableBytes  int       `json:"availableBytes,omitempty"`
		} `json:"memory"`
		Network struct {
			Time       time.Time `json:"time"`
			Name       string    `json:"name"`
			RxBytes    int64     `json:"rxBytes"`
			RxErrors   int       `json:"rxErrors"`
			TxBytes    int64     `json:"txBytes"`
			TxErrors   int       `json:"txErrors"`
			Interfaces []struct {
				Name     string `json:"name"`
				RxBytes  int64  `json:"rxBytes"`
				RxErrors int    `json:"rxErrors"`
				TxBytes  int64  `json:"txBytes"`
				TxErrors int    `json:"txErrors"`
			} `json:"interfaces"`
		} `json:"network"`
		Volume []struct {
			Time           time.Time `json:"time"`
			AvailableBytes int64     `json:"availableBytes"`
			CapacityBytes  int64     `json:"capacityBytes"`
			UsedBytes      int       `json:"usedBytes"`
			InodesFree     int       `json:"inodesFree"`
			Inodes         int       `json:"inodes"`
			InodesUsed     int       `json:"inodesUsed"`
			Name           string    `json:"name"`
		} `json:"volume,omitempty"`
		EphemeralStorage struct {
			Time           time.Time `json:"time"`
			AvailableBytes int64     `json:"availableBytes"`
			CapacityBytes  int64     `json:"capacityBytes"`
			UsedBytes      int       `json:"usedBytes"`
			InodesFree     int       `json:"inodesFree"`
			Inodes         int       `json:"inodes"`
			InodesUsed     int       `json:"inodesUsed"`
		} `json:"ephemeral-storage"`
		ProcessStats struct {
			ProcessCount int `json:"process_count"`
		} `json:"process_stats"`
	} `json:"pods"`
}
