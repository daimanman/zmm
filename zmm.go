package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"unicode/utf8"
)

var (
	p = flag.Bool("p",false,"百分比显示")
	C = flag.String("C", "N", "统计字符数")
	N = flag.Int("N", 2, "统计第几列中的字符")
	h = flag.Bool("h", false, "show useage ")
	BH = flag.String("BH","","编号名称对应文件")
	T = flag.String("T","a","操作类型")

)
var wg sync.WaitGroup

var useage = `
	-p 非百分比显示
    -C 统计的字符 默认为 'N'
    -N 字符在第几列中 默认统计第二列
	-BF 编号对应文件
	-BH 编号对应文件，第一列为编号 第二列为编号对应的名字
	-T 操作类型 a:计算百分比 b:替换编号

	Example:
	编号转换,bh.txt文件中设置编号对应关系: zmm  -BH="bh.txt" -p c93d8m80p60.phy
	默认输出,未带编号转换: zmm c93d8m80p60.phy

	Note:
	如果在使用中碰到什么问题,请联系作者 QQ:1018793423
`

type DataInfo struct {
	Filename string         //文件名称
	Total    int            //文件N总数
	DataMap  map[string]int //每行N的个数
	SN       string         //个体数
	COL      string         //列数
	Snames   []string
	BhMap *map[string]string //编号对应名称
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

func parseBhMap(bh string) map[string]string{
	bhMap := map[string]string{}
	if bh == ""{
		return bhMap
	}
	f,err := os.Open(bh)
	if err != nil  {
		panic(err)
		return bhMap
	}
	scannerF := bufio.NewScanner(f)
	for scannerF.Scan(){
		line := scannerF.Text()
		fsList := strings.Fields(strings.TrimSpace(line))
		if len(fsList) >= 2{
			bhMap[fsList[0]] = fsList[1]
		}
	}
	if err := scannerF.Err(); err != nil {
		log.Println("读取文件 err ",err)
	}
	return bhMap
}
func (data *DataInfo) show() {
	snames := data.Snames
	dataMap := data.DataMap
	bhMap := data.BhMap
	var n int
	str := strings.Repeat("*", 20)
	fmt.Printf("\n\n%s%s%s\n", str, data.Filename, str)

	fmt.Printf(" Total=%d \n\n", data.Total)

	for _, name := range snames {
		n = dataMap[name]
		bn := dataMap[name+"_BC"]
		nf := float64(n)
		bnf := float64(bn)
		bhName := (*bhMap)[name]
		if bhName == ""{
			bhName = name
		}
		if !*p {
			fmt.Printf("%5s %9d   %9.2f%% \n", bhName, n, (nf/bnf)*100)
		}else {
			fmt.Printf("%5s %9d   %9.4f \n", bhName, n, (nf/bnf))
		}
	}
}

func dealFile(fileFullPath string,bhMap *map[string]string) {
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
		// fmt.Printf("bs line is %d C=%s \n", len(bs), *C)
		if len(bs) > 1 {
			if *N < 1 {
				*N = 1
			}
			n := strings.Count(bs[*N-1], *C)
			// fmt.Printf("count n is %d \n", n)
			if n > 0 {
				total += n
				sname := bs[0]
				data.DataMap[sname] = n
				data.DataMap[sname+"_BC"] = utf8.RuneCountInString(bs[*N-1])
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
	data.BhMap = bhMap
	data.show()

}

func main() {
	flag.Parse()
	if *h {
		fmt.Println(useage)
		return
	}
	bhMap := parseBhMap(*BH)
	files := GetFiles(flag.Args())

	tName := strings.ToLower(*T)
	//计算百分比
	if  tName == "a" {
		lenth := len(files)
		if len(files) == 0 {
			fmt.Printf("未找到文件,请检查参数是否在正确 \n")
			return
		}
		wg.Add(lenth)
		for _, file := range files {
			go dealFile(file,&bhMap)
		}
		wg.Wait()
	}

	//简单的替换编号
	if tName == "b"{
		if len(files) == 0 {
			panic("请指定要转换的文件")
		}

		if *BH == "" {
			panic(fmt.Sprintf("%s %s\n","未指定编号替换文件"," 请通过 -BH 参数指定, Example: -BH=bh.txt"))
		}
		var errFile *os.File
		var errWriter *bufio.Writer

		transFile,err  := os.Open(files[0])
		if err != nil {
			panic(err)
		}
		fileReader := bufio.NewReader(transFile)
		readEnd := false
		for !readEnd{
			line,readErr := fileReader.ReadString('\n')
			if readErr != nil {
				if readErr == io.EOF {
					readEnd = true
				}
			}
			lineStr := strings.TrimSpace(line)
			if lineStr == ""{
				continue
			}
			lineFsList := strings.Fields(lineStr)
			if len(lineFsList) > 0{
				key := lineFsList[0]
				val := bhMap[key]
				if val != "" {
					lineFsList[0] = val
				}else{

					if errFile == nil {
						errFile,_ = os.Create("error.log")
					}

					if errFile != nil && errWriter ==  nil {
						errWriter = bufio.NewWriter(errFile)
						defer errWriter.Flush()
					}

					if errWriter != nil {
						errWriter.WriteString(fmt.Sprintf("未在 %s 文件中找到 %s 对应的编号 \n",*BH,key))
					}


				}
			}
			fmt.Println(strings.Join(lineFsList,"  "))
		}
		defer transFile.Close()



	}





}
