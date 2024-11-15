package main

import (
	"bytes"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/cli-runtime/pkg/printers"
	"log"
	"os"
	"strings"
)

func main() {

	// Path to the Kubernetes Services JSON file
	filePath := "services.json"

	// Read JSON file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse JSON into a Kubernetes List object
	serviceList := &corev1.ServiceList{}
	if err := json.Unmarshal(fileData, &serviceList); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Set print options
	opts := printers.PrintOptions{
		Wide:      true,
		NoHeaders: false,
	}
	printTable(opts, populateKubeServiceTable(*serviceList))
}

func printTable(opts printers.PrintOptions, table metav1.Table) {
	// Print the table
	out := bytes.NewBuffer([]byte{})
	printer := printers.NewTablePrinter(opts)
	err := printer.PrintObj(&table, out)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(out.String())
}

func populateKubeServiceTable(serviceList corev1.ServiceList) metav1.Table {
	table := &metav1.Table{}
	table.ColumnDefinitions = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Namespace", Type: "string"},
		{Name: "Type", Type: "string"},
		{Name: "Cluster-IP", Type: "string"},
		{Name: "Ports", Type: "string"},
	}
	for _, svc := range serviceList.Items {
		row := metav1.TableRow{
			Cells: []interface{}{
				svc.Name,
				svc.Namespace,
				svc.Spec.Type,
				svc.Spec.ClusterIP,
				generatePortsList(svc.Spec.Ports),
			},
		}
		table.Rows = append(table.Rows, row)
	}
	return *table
}

func generatePortsList(ports []corev1.ServicePort) string {
	var portEntries []string
	for _, port := range ports {
		portEntries = append(portEntries, fmt.Sprintf("%d/%s", port.Port, port.Protocol))
	}
	return strings.Join(portEntries, ",")
}
