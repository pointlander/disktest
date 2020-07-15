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

	_, err := os.Stat("test.bin")
	if err != nil {
		fmt.Println("generating test.bin")
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
	} else {
		fmt.Println("found test.bin")
	}

	// sync; echo 3 > /proc/sys/vm/drop_caches
	for i := 0; i < 2; i++ {
		fmt.Println("sync")
		err := exec.Command("sync").Run()
		if err != nil {
			panic(err)
		}
		fmt.Println("sleep 1")
		time.Sleep(time.Second)
	}
	fmt.Println("echo 3 > /proc/sys/vm/drop_caches")
	err = exec.Command("sh", "-c", "/usr/bin/echo 3 > /proc/sys/vm/drop_caches").Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("sleep 1")
	time.Sleep(time.Second)

	rand.Seed(1)
	fmt.Println("verifying test.bin")
	bytes := 0
	input, err := os.Open("test.bin")
	if err != nil {
		panic(err)
	}
	var data [8]byte
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
