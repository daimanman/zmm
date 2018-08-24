package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

var (
	T = flag.Bool("T", false, "标记是否显示title")
	C = flag.String("C", "N", "")
)
var wg sync.WaitGroup

type DataInfo struct {
	Filename string         //文件名称
	Total    int            //文件N总数
	DataMap  map[string]int //每行N的个数
	SN       string         //个体数
	COL      string         //列数
	Snames   []string
}

func TestPath() {
	ps := "dsjkds/dsjdks/*.txt"
	filename := "12djsk.txt"
	p2 := "d.txt"
	m, _ := path.Match("1[2]*.txt", filename)
	fmt.Println(m)
	fmt.Println("dsjkds/dsjdks/*.txt Dir", path.Dir(ps))
	fmt.Println("dsjkds/dsjdks/*.txt Base", path.Base(ps))

	fmt.Println("q.txt Dir", path.Dir(p2))
	fmt.Println("q.txt Base", path.Base(p2))
}

func GetFiles(paths []string) []string {
	fs := make([]string, 0)
	for _, p := range paths {
		dir := path.Dir(p)
		base := path.Base(p)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Println(dir, "文件夹不存在", err.Error())
			continue
		}

		for _, file := range files {
			fname := file.Name()
			isDir := file.IsDir()
			if isDir {
				continue
			}
			match, _ := path.Match(base, fname)
			if match {
				fs = append(fs, path.Join(dir, fname))
			}
		}
	}
	return fs
}

func (data *DataInfo) show() {
	snames := data.Snames
	dataMap := data.DataMap
	total := float64(data.Total)
	var n int
	str := strings.Repeat("*", 20)
	fmt.Printf("\n\n%s%s%s\n", str, data.Filename, str)
	fmt.Printf("SN=%s COL=%s Total=%d \n\n", data.SN, data.COL, data.Total)
	for _, name := range snames {
		n = dataMap[name]
		nf := float64(n)
		fmt.Printf("%s\t %d \t %.3f \n", name, n, nf/total)
	}
}

func dealFile(fileFullPath string) {
	file, err := os.Open(fileFullPath)
	if err != nil {
		fmt.Println(fileFullPath, err.Error())
		return
	}
	defer func() {
		file.Close()
		wg.Done()
	}()
	total := 0
	data := DataInfo{
		Total:    total,
		Filename: file.Name(),
		DataMap:  make(map[string]int),
		Snames:   make([]string, 0),
	}
	rd := bufio.NewReader(file)
	i := 0
	for {
		line, err1 := rd.ReadString('\n')
		if err1 != nil || io.EOF == err1 {
			break
		}
		bs := strings.Fields(line)
		if len(bs) > 1 {
			n := strings.Count(bs[1], *C)
			if n > 0 {
				total += n
				data.DataMap[bs[0]] = n
				data.Snames = append(data.Snames, bs[0])
			}
			if i == 0 {
				data.COL = bs[1]
				data.SN = bs[0]
			}
		}
		i++
	}
	data.Total = total
	data.show()

}

func main() {
	flag.Parse()
	files := GetFiles(flag.Args())
	lenth := len(files)
	if len(files) == 0 {
		fmt.Printf("未找到文件,请检查参数是否在正确 \n")
		return
	}
	wg.Add(lenth)
	for _, file := range files {
		go dealFile(file)
	}
	wg.Wait()

}
