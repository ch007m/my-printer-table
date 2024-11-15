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

type Config struct {
	FilePath       string
	PrinterOptions printers.PrintOptions
	Succeeded      bool
}

func main() {
	c := &Config{
		FilePath: "services.json",
		PrinterOptions: printers.PrintOptions{
			Wide:      true,
			NoHeaders: false,
		},
	}

	// Read JSON file
	fileData, err := os.ReadFile(c.FilePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse JSON into a Kubernetes List object
	serviceList := &corev1.ServiceList{}
	if err := json.Unmarshal(fileData, &serviceList); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	err = c.PrintTable(populateKubeServiceTable(*serviceList))
	if err != nil {
		fmt.Println("Error:", err)
	}

	if c.Succeeded {
		fmt.Println("Table printed successfully !")
	}

}

func (c *Config) PrintTable(table metav1.Table) error {
	out := bytes.NewBuffer([]byte{})

	printer := printers.NewTablePrinter(c.PrinterOptions)
	err := printer.PrintObj(&table, out)
	if err != nil {
		return err
	}

	fmt.Println(out.String())
	c.Succeeded = true
	return nil
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
