package main

import (
	"encoding/json"
	"fmt"

	"github.com/coreos/etcd/client"
)

func printResponseKey(resp *client.Response) {
	shellState := getShellState()
	switch shellState.format {
	case "short":
		if resp.Action != "delete" {
			if resp.Node.Value == "" {
				if resp.Node.Dir {
					fmt.Println("Directory:")
					for _, node := range resp.Node.Nodes {
						fmt.Println(node.Key)
					}
				} else {
					fmt.Println("Empty string: Use 'env [j[son]|p[retty]]' to get a better view")
				}
			} else {
				fmt.Println(resp.Node.Value)
			}
		} else {
			fmt.Println("Deleted Value:", resp.PrevNode.Value)
		}
	case "json":
		b, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	case "pretty":
		b, err := json.MarshalIndent(resp, " ", " ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}
}
