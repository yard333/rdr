package dump

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"rdr/decoder"

	"github.com/urfave/cli"
)

//解析rdb文件函数
func exportcsv(c *cli.Context, v string) {
	filename := filepath.Base(v)

	if !counters.Check(filename) {
		decodern := decoder.NewDecoder()

		go Decode(c, decodern, v)
		//counter := NewCounter()
		//counter.Count(decodern.Entries)
		//counters.Set(filename, counter)

		//print top
		res := []*decoder.Entry{}

		for e := range decodern.Entries {
			res = append(res, e)
		}

		/*
			for i := 0; i < counter.largestEntries.Len(); i++ {
				entries := *counter.largestEntries
				res = append(res, entries[i])
			}
		*/

		sort.Sort(sort.Reverse(entryHeap(res)))

		num, _ := strconv.Atoi(c.String("top"))
		if num != 0 {
			if num < len(res) {
				res = res[:num]
			}
		}

		//如果未指定导出路径，按日期生成csv文件
		var cpath string
		if c.String("c") == "" {
			cpath = "./" + time.Now().Format("20060102") + ".csv"
		} else {
			cpath = c.String("c")
		}

		f, err := os.Create(cpath)
		if err != nil {
			log.Println(err.Error())
		}

		_, err = f.WriteString("Key,Bytes,Type,DBIndex,NumOfElem,LenOfLargestElem,FieldOfLargestElem")
		if err != nil {
			log.Println(err.Error())
		}

		for j := 0; j < len(res); j++ {
			var v string = res[j].FieldOfLargestElem
			if len(res[j].FieldOfLargestElem) > 512 {
				v = res[j].FieldOfLargestElem[:512] + "..."
			}

			_, err = f.WriteString("\n" + res[j].Key + "," + strconv.FormatUint(res[j].Bytes, 10) + "," + res[j].Type + "," + strconv.FormatUint(res[j].DBIndex, 10) + "," + strconv.FormatUint(res[j].NumOfElem, 10) + "," + strconv.FormatUint(res[j].LenOfLargestElem, 10) + "," + v)
			if err != nil {
				log.Println(err.Error())
			}

		}

		f.Close()
	}
}

func Csv(c *cli.Context) {

	//如果有-f参数直接解析rdb文件，否则根据参数连接到redis执行bgsave后再解析rdb
	if c.String("file") != "" {
		for _, v := range listPathFiles(c.String("file")) {
			exportcsv(c, v)
		}
	} else {
		err := Rdb(c.String("port"), c.String("password"))
		if err != nil {
			log.Println(err.Error())
			return
		}

		fmt.Println()
		fmt.Println("将开始解析rdb文件")
		fmt.Println()

		exportcsv(c, "/usr/local/redis/data/dump.rdb")

		if c.String("rm") == "y" {
			err := os.Remove("/usr/local/redis/data/dump.rdb")
			if err != nil {
				log.Println(err.Error())
			} else {
				fmt.Println("删除/usr/local/redis/data/dump.rdb成功")
			}
		}

	}

}
