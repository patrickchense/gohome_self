package crawler

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"gohome_self/tracingbook/mail"
	"gohome_self/tracingbook/models"
	"io/ioutil"
	"log"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// books sites,
/*

https://www.230book.com/   顶点
https://www.uukanshu.com/b/125477/  UU看书
笔趣阁

这里只处理顶点
*/

type FetchBookInfo struct {
	BookId       uint
	BookName     string
	LatestNumber int64
}

var allBooks map[uint]FetchBookInfo
var dingDianURL = "https://www.230book.com/book/"

func InitDingDian() {
	allBooks = make(map[uint]FetchBookInfo)
	allBooks[6333] = FetchBookInfo{BookId: 6333, BookName: "我真没想重生啊", LatestNumber: 0}
	allBooks[2602] = FetchBookInfo{BookId: 2602, BookName: "是篮球之神啊", LatestNumber: 0}
	allBooks[6454] = FetchBookInfo{BookId: 6454, BookName: "青莲之巅", LatestNumber: 0}
	allBooks[12368] = FetchBookInfo{BookId: 12368, BookName: "万族之劫", LatestNumber: 0}
	allBooks[1738] = FetchBookInfo{BookId: 1738, BookName: "大医凌然", LatestNumber: 0}
	allBooks[2849] = FetchBookInfo{BookId: 2849, BookName: "最初进化", LatestNumber: 0}
}

func FetchBooks() {
	ticker := time.NewTicker(20 * time.Second)
	quit := make(chan struct{})
	rand.Seed(time.Now().UnixNano())
	go func() {
		for {
			for bookId, bookInfo := range allBooks {
				select {
				case <-ticker.C:
					fmt.Printf("start with %s, ln: %d\n", bookInfo.BookName, bookInfo.LatestNumber)
					time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
					updates, err := bookUpdateDingdian(bookInfo.BookId, 0, bookInfo.BookName, bookInfo.LatestNumber)
					if err != nil {
						fmt.Errorf("error fetch book updates, %e", err)
						return
					}
					if len(updates) > 1 {
						sort.Slice(updates, func(i, j int) bool {
							return updates[i].LatestNumber > updates[j].LatestNumber
						})
						fmt.Printf("all len: %d\n", len(updates))
						fmt.Printf("latest cp: %s\n", updates[0].LatestName)
						bookInfo.LatestNumber = updates[0].LatestNumber
						allBooks[bookId] = bookInfo
						mail.NML.NotifyUpdates(updates)
						fmt.Printf("send email done!")
					}
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}
	}()
}

func main() {
	InitDingDian()
	FetchBooks()
	for {
		time.Sleep(100 * time.Minute)
	}
}

/*

books by id,  store the bookId and site

https://www.230book.com/modules/article/search.php  没搞懂

https://www.230book.com/book/6333/  书名和id  GBK编码需要解析

in the head
og:novel:lastest_chapter_name
<meta property="og:novel:lastest_chapter_name" content="688、我做的生意，比抢钱可快多了！"/>
<meta property="og:novel:latest_chapter_url" content="https://www.230book.com/book/6333/3669657.html"/>

normal charpter <li><a href="1101075.html">7、看老子脸色行事</a></li>
latest: <a href="/book/6333/3669657.html">688、我做的生意，比抢钱可快多了！</a>
need to find all the a in body

*/
func bookUpdateDingdian(bookId uint, siteId uint, bookName string, last int64) (updates []models.UpdateItem, err error) {
	c := colly.NewCollector()
	np, _ := regexp.Compile("^[0-9]+[、]{1}")
	latestName := ""
	latestNumber := int64(0)
	accessUrl := dingDianURL + strconv.FormatUint(uint64(bookId), 10) + "/"
	fmt.Printf("access: %s\n", accessUrl)
	c.OnHTML("meta[property]", func(e *colly.HTMLElement) {
		//get latest url and name
		content := e.Attr("property")
		//println("all meta:" + content)
		if content == "og:novel:lastest_chapter_name" {
			latestName, latestNumber = handleDingdianLatest(e.Attr("content"))
			if latestNumber <= last {
				fmt.Printf("no updates " + string(bookId) + " at " + time.Now().String())
				return
			} else {
				fmt.Printf("new updates(%s) %d from %d to %d\n", latestName, bookId, last, latestNumber)
			}
		}
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Response.Headers.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(bytes.NewBuffer(e.Response.Body))
			if err != nil {
				return
			}
			defer reader.Close()
			data, err := ioutil.ReadAll(reader)
			if err != nil {
				return
			}
			println(data)
		} else {
			cpName, err := GbkToUtf8Byte([]byte(e.Text))
			if err != nil {
				return
			}
			if np.Match(cpName) {
				bn := string(cpName)
				ah := e.Attr("href")
				bi, _ := strconv.ParseInt(bn[0:strings.Index(bn, "、")], 10, 64)
				if bi > last {
					item := models.UpdateItem{
						BookId:       bookId,
						BookName:     bookName,
						SiteId:       siteId,
						LatestName:   bn,
						LatestNumber: bi,
						BookUrl:      accessUrl + ah,
					}
					updates = append(updates, item)
					fmt.Printf("parse update book %s, cp: %s", bookName, bn)
				}
			}
		}
	})

	err = c.Visit(accessUrl)
	return
}

func handleDingdianLatestUrl(content string) string {
	c, err := GbkToUtf8([]byte(content))
	if err != nil {
		log.Fatal("convert " + content + " failed" + err.Error())
		return ""
	}
	return c
}

func handleDingdianLatest(content string) (string, int64) {
	c, err := GbkToUtf8([]byte(content))
	if err != nil {
		log.Fatal("convert " + content + " failed" + err.Error())
		return "", 0
	}
	cs := strings.Split(c, "、")
	if latestNumber, err := strconv.ParseInt(cs[0], 10, 64); err == nil {
		return c, latestNumber
	} else {
		log.Fatal("convert " + content + " failed" + err.Error())
	}
	return "", 0
}

func GbkToUtf8(s []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return "", e
	}
	return string(d), nil
}

func GbkToUtf8Byte(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
