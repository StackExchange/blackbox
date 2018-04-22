package main

import (
	"fmt"

	"github.com/StackExchange/blackbox/pkg/bbutil"
)

func cmdAdminList() error {
	bbu, err := bbutil.New()
	if err != nil {
		return err
	}
	names, err := bbu.Administrators()
	if err != nil {
		return err
	}
	for _, item := range names {
		fmt.Println(item.Name)
	}
	return nil
}
