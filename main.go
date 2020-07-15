// Copyright 2020 The Disk Test Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

// Gigabyte is one gigabyte
const Gigabyte = 1024 * 1024 * 1024

func main() {
	rand.Seed(1)
	output, err := os.Create("test.bin")
	if err != nil {
		panic(err)
	}
	var data [8]byte
	for i := 0; i < Gigabyte; i++ {
		v, index := rand.Uint64(), 0
		for j := 0; j < 64; j += 8 {
			data[index] = byte(0xff & (v >> j))
			index++
		}
		_, err := output.Write(data[:])
		if err != nil {
			panic(err)
		}
	}
	output.Close()

	// sync; echo 3 > /proc/sys/vm/drop_caches
	for i := 0; i < 2; i++ {
		err := exec.Command("sync").Run()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
	drop, err := os.Open("/proc/sys/vm/drop_caches")
	if err != nil {
		panic(err)
	}
	_, err = drop.Write([]byte{'3'})
	if err != nil {
		panic(err)
	}
	drop.Close()
	time.Sleep(time.Second)

	rand.Seed(1)
	bytes := 0
	input, err := os.Open("test.bin")
	if err != nil {
		panic(err)
	}
	for i := 0; i < Gigabyte; i++ {
		_, err := input.Read(data[:])
		if err != nil {
			panic(err)
		}
		v, index := rand.Uint64(), 0
		for j := 0; j < 64; j += 8 {
			if data[index] != byte(0xff&(v>>j)) {
				bytes++
			}
			index++
		}
	}
	input.Close()
	fmt.Println("corrupt bytes", bytes)
}
