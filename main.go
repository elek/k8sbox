package main

import (
	"context"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"time"
)

func main() {
	c := kong.Parse(&Show{})
	err := c.Run()
	if err != nil {
		log.Fatalf("%++v", err)
	}
}

type Show struct {
	ProxyPort int
}

func (s Show) Run() error {

	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return errors.WithStack(err)

	}
	client, err := clientset.NewForConfig(config)

	podsByNode := make(map[string][]corev1.Pod)
	pods, err := client.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.WithStack(err)
	}
	for _, pod := range pods.Items {
		podsByNode[pod.Spec.NodeName] = append(podsByNode[pod.Spec.NodeName], pod)
	}

	nodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.WithStack(err)
	}

	var summary ResourceStat

	for _, node := range nodes.Items {
		nodeSummary := ResourceStat{
			Allocatable: node.Status.Allocatable,
			Capacity:    node.Status.Capacity,
		}
		fmt.Printf("Node name: %s %s %s %v\n", color.RedString(node.Name), node.Labels["node.kubernetes.io/instance-type"], node.Labels["cloud.google.com/gke-provisioning"], node.Spec.Taints)

		var stat Stat
		if s.ProxyPort != 0 {
			stat, err = GetNodeStat(s.ProxyPort, node.Name)
			if err != nil {
				return err
			}
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		core := "CPU (cores)"
		mem := "MEM (GB)"
		t.AppendHeader(table.Row{"", core, core, core, mem, mem, mem, "NET", "NET"}, table.RowConfig{AutoMerge: true})
		t.AppendHeader(table.Row{"Pod", "used", "request", "limit", "used", "request", "limit", "RX", "TX"})
		for _, pod := range podsByNode[node.Name] {
			var podSummary ResourceStat
			if pod.Namespace == "kube-system" {
				continue
			}
			for _, c := range pod.Spec.Containers {
				podSummary.Limit = add(podSummary.Limit, c.Resources.Limits)
				podSummary.Requested = add(podSummary.Requested, c.Resources.Requests)
			}
			if stat.Node.NodeName != "" {
				for _, sp := range stat.Pods {
					if sp.PodRef.Name == pod.Name && sp.PodRef.Namespace == pod.Namespace {
						durationSec := int64(time.Since(sp.StartTime).Seconds())
						usedCPU := resource.NewScaledQuantity(int64(sp.Cpu.UsageNanoCores), resource.Nano)
						usedMemory := resource.NewQuantity(int64(sp.Memory.UsageBytes), resource.DecimalSI)

						used := corev1.ResourceList{
							corev1.ResourceCPU:    *usedCPU,
							corev1.ResourceMemory: *usedMemory,
						}
						podSummary.Used = add(podSummary.Used, used)
						t.AppendRow([]interface{}{
							text.FgYellow.Sprintf(pod.Name),
							printCPU(podSummary.Used.Cpu(), 0),
							printCPU(podSummary.Requested.Cpu(), 0),
							printCPU(podSummary.Limit.Cpu(), 0),
							printMem(podSummary.Used.Memory(), 0),
							printMem(podSummary.Requested.Memory(), 0),
							printMem(podSummary.Limit.Memory(), 0),
							sp.Network.RxBytes / durationSec,
							sp.Network.TxBytes / durationSec,
						})
						break
					}
				}
			}
			nodeSummary.add(podSummary)
		}

		t.AppendRow([]interface{}{
			text.FgRed.Sprintf("SUMMARY"),
			printCPU(nodeSummary.Used.Cpu(), text.FgRed),
			printCPU(nodeSummary.Requested.Cpu(), text.FgRed),
			printCPU(nodeSummary.Limit.Cpu(), text.FgRed),
			printMem(nodeSummary.Used.Memory(), text.FgRed),
			printMem(nodeSummary.Requested.Memory(), text.FgRed),
			printMem(nodeSummary.Limit.Memory(), text.FgRed),
			"", "",
		})
		c := text.FgGreen
		t.AppendRow([]interface{}{
			"allocatable",
			printCPU(nodeSummary.Allocatable.Cpu(), c),
			printCPU(nodeSummary.Allocatable.Cpu(), c),
			printCPU(nodeSummary.Allocatable.Cpu(), c),
			printMem(nodeSummary.Allocatable.Memory(), c),
			printMem(nodeSummary.Allocatable.Memory(), c),
			printMem(nodeSummary.Allocatable.Memory(), c),
			"", "",
		}, table.RowConfig{AutoMerge: true})

		t.AppendRow([]interface{}{
			"capacity",
			printCPU(nodeSummary.Capacity.Cpu(), text.FgWhite),
			printCPU(nodeSummary.Capacity.Cpu(), text.FgWhite),
			printCPU(nodeSummary.Capacity.Cpu(), text.FgWhite),
			printMem(nodeSummary.Capacity.Memory(), text.FgWhite),
			printMem(nodeSummary.Capacity.Memory(), text.FgWhite),
			printMem(nodeSummary.Capacity.Memory(), text.FgWhite),
			"", "",
		}, table.RowConfig{AutoMerge: true})

		t.Render()

		summary.add(nodeSummary)
		fmt.Println()

	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	core := "CPU (cores)"
	mem := "MEM (GB)"
	t.AppendHeader(table.Row{core, core, core, core, mem, mem, mem, mem}, table.RowConfig{AutoMerge: true})
	t.AppendHeader(table.Row{"allocatable", "used", "request", "limit", "allocatable", "used", "request", "limit"})
	t.AppendRow([]interface{}{
		printCPU(summary.Allocatable.Cpu(), text.FgWhite),
		printCPU(summary.Used.Cpu(), text.FgWhite),
		printCPU(summary.Requested.Cpu(), text.FgWhite),
		printMem(summary.Limit.Memory(), text.FgWhite),
		printMem(summary.Allocatable.Memory(), text.FgWhite),
		printMem(summary.Used.Memory(), text.FgWhite),
		printMem(summary.Requested.Memory(), text.FgWhite),
		printMem(summary.Limit.Memory(), text.FgWhite),
	})
	t.Render()
	return nil
}

func add(target corev1.ResourceList, addendum corev1.ResourceList) corev1.ResourceList {
	result := make(corev1.ResourceList)
	for k, v := range target {
		q := &resource.Quantity{}
		q.Add(v)
		result[k] = *q
	}
	for k, v := range addendum {
		q := &resource.Quantity{}
		existing, ok := result[k]
		if ok {
			q.Add(existing)
		}
		q.Add(v)
		result[k] = *q
	}
	return result
}

type ResourceStat struct {
	Capacity    corev1.ResourceList
	Allocatable corev1.ResourceList
	Used        corev1.ResourceList
	Requested   corev1.ResourceList
	Limit       corev1.ResourceList
}

func (s *ResourceStat) add(summary ResourceStat) {
	s.Capacity = add(s.Capacity, summary.Capacity)
	s.Allocatable = add(s.Allocatable, summary.Allocatable)
	s.Requested = add(s.Requested, summary.Requested)
	s.Limit = add(s.Limit, summary.Limit)
	s.Used = add(s.Used, summary.Used)
}

func printMem(r *resource.Quantity, color text.Color) string {
	return color.Sprintf("%.3f", float64(r.ScaledValue(resource.Mega))/1024)
}

func printCPU(r *resource.Quantity, color text.Color) string {
	return color.Sprintf("%.2f", float64(r.ScaledValue(resource.Milli))/1000)
}
