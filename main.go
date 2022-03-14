// Copyright 2017 XUEQIU.COM
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
	"os"
	"rdr/decoder"
	"rdr/dump"

	"github.com/urfave/cli"

	"fmt"
)

//go:generate go-bindata -prefix "static/" -o=static/static.go -pkg=static -ignore static.go static/...
//go:generate go-bindata -prefix "views/" -o=views/views.go -pkg=views -ignore views.go views/...

// keys is function for command `keys`
// output all keys in rdbfile(s) get from args
func keys(c *cli.Context) {
	if c.NArg() < 1 {
		fmt.Fprintln(c.App.ErrWriter, "keys requires at least 1 argument")
		cli.ShowCommandHelp(c, "keys")
		return
	}
	for _, filepath := range c.Args() {
		decoder := decoder.NewDecoder()
		go dump.Decode(c, decoder, filepath)
		for e := range decoder.Entries {
			fmt.Fprintf(c.App.Writer, "%v\n", e.Key)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "rdr"
	app.Usage = "a tool to parse redis rdbfile"
	app.Version = "v1.0.0"
	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "dump",
			Usage:     "dump statistical information of rdbfile to STDOUT",
			ArgsUsage: "FILE1 [FILE2] [FILE3]...",
			Action:    dump.ToCliWriter,
		},
		cli.Command{
			Name:      "show",
			Usage:     "show statistical information of rdbfile by webpage",
			ArgsUsage: "DIR1 [DIR2] [DIR3] or FILE1 [FILE2] [FILE3]...",
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:  "port, p",
					Value: 8080,
					Usage: "Port for rdr to listen",
				},
			},
			Action: dump.Show,
		},
		cli.Command{
			Name:      "keys",
			Usage:     "get all keys from rdbfile",
			ArgsUsage: "FILE1 [FILE2] [FILE3]...",
			Action:    keys,
		},
		//新增输出csv格式
		cli.Command{
			Name:      "csv",
			Usage:     "print top key with csv formt",
			ArgsUsage: "DIR1 [DIR2] [DIR3] or FILE1 [FILE2] [FILE3]...",
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:  "top, t",
					Value: 0,
					Usage: "the number of top",
				},
				cli.StringFlag{
					Name:  "port,p",
					Value: "6379",
					Usage: "port",
				},
				cli.StringFlag{
					Name:  "password,a",
					Value: "",
					Usage: "password",
				},
				cli.StringFlag{
					Name:  "file,f",
					Value: "",
					Usage: "filepath",
				},
				cli.StringFlag{
					Name:  "rm",
					Value: "n",
					Usage: "is need rm dump.rdb",
				},
				cli.StringFlag{
					Name:  "c",
					Value: "",
					Usage: "export path",
				},
			},
			Action: dump.Csv,
		},
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.ErrWriter, "command %q can not be found.\n", command)
		cli.ShowAppHelp(c)
	}
	app.Run(os.Args)
}
