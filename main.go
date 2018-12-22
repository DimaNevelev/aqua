// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/dimanevelev/travers/cmd"
	"log"
)

var (
	// Trace is for full detailed messages.
	Trace *log.Logger

	// Info is for important messages.
	Info *log.Logger

	// Warning is for need to know issue messages.
	Warning *log.Logger

	// Error is for error messages.
	Error *log.Logger
)

func main() {
	cmd.Execute()
}
