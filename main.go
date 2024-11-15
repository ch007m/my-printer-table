package main

import (
	"bytes"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

var (
	podName      = "my-pod-name"
	podNamespace = "my-namespace"
	columns      = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Ready", Type: "string"},
		{Name: "Status", Type: "string"},
		{Name: "Retries", Type: "integer", Priority: 1},
		{Name: "Age", Type: "string", Priority: 1},
	}
	rows = []metav1.TableRow{
		{Cells: []interface{}{"test1", "1/1", "podPhase", int64(5), "20h"}},
		{Cells: []interface{}{"test2", "1/2", "podPhase", int64(30), "21h"}},
		{Cells: []interface{}{"test3", "4/4", "podPhase", int64(1), "22h"}},
	}
	pod = &corev1.Pod{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: podNamespace,
			Labels: map[string]string{
				"first-label":  "12",
				"second-label": "label-value",
			},
		},
	}
)

func main() {

	// Print a row using a pod
	opts := printers.PrintOptions{
		Wide:      true,
		NoHeaders: false,
	}
	printRow(opts, *pod)

	// Print headers and rows using wide
	opts = printers.PrintOptions{
		Wide: true,
	}
	printTable(opts)
}

func printRow(opts printers.PrintOptions, pod corev1.Pod) {
	buf := &bytes.Buffer{}
	printer := printers.NewTablePrinter(opts)
	err := printer.PrintObj(&pod, buf)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(buf.String())
}

func printTable(opts printers.PrintOptions) {
	// Create the table from the columns and rows.
	table := &metav1.Table{
		ColumnDefinitions: columns,
		Rows:              rows,
	}

	// Print the table
	out := bytes.NewBuffer([]byte{})
	printer := printers.NewTablePrinter(opts)
	err := printer.PrintObj(table, out)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(out.String())
}
