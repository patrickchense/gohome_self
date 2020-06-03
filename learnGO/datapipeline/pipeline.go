package datapipeline

import (
	"bufio"
	"github.com/google/uuid"
	"gohome_self/learnGO/log"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

/**
https://towardsdatascience.com/concurrent-data-pipelines-in-golang-85b18c2eecc2
这是一篇关于用Go建立data pipeline的文章
截取其中的精华，大概就是通过 goroutine和chan来传输接受数据
中间使用waitGroup来接受多个对象并等待结束，知道最后保存

我觉得中间使用chan 和 waitGroup的方式值得参看，我自己从来没用在实践中使用过chan和waitGruop
*/

//generate uuid from file
func generateData() <-chan uuid.UUID {
	c := make(chan uuid.UUID)
	const filepath = "guids.txt"
	go func() {
		//generate to chan
		file, _ := os.Open(filepath)
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			line = strings.TrimSuffix(line, "\n")
			guid, err := uuid.Parse(line)
			if err != nil {
				continue
			}
			c <- guid
		}
		close(c) //chan 一定要关闭
	}()

	return c
}

//handle data
type inputData struct {
	id        string
	timestamp int64
}

func prepareData(ic <-chan uuid.UUID) <-chan inputData {
	oc := make(chan inputData)
	go func() {
		for id := range ic {
			input := inputData{id: id.String(), timestamp: time.Now().UnixNano()}
			log.Printf("Data ready for processing: %+v", input)
			oc <- input
		}
		close(oc)
	}()

	return oc
}

//external data
type externalData struct {
	inputData
	relatedIds []string
}

func fetchData(ic <-chan inputData) <-chan externalData {
	oc := make(chan externalData)

	go func() {
		wg := &sync.WaitGroup{}

		for input := range ic {
			wg.Add(1)
			go fetchFromExternalService(input, oc, wg)
		}
		wg.Wait() //避免当前线程提前结束
		close(oc)
	}()

	return oc
}

func fetchFromExternalService(input inputData, oc chan externalData, wg *sync.WaitGroup) {
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

	related := make([]string, 0)
	for i := 0; i < rand.Intn(10); i++ {
		related = append(related, uuid.New().String())
	}
	oc <- externalData{input, related}
	wg.Done()
}

//save
type saveResult struct {
	idsSave   []string
	timestamp int64
}

func saveData(ic <-chan externalData) <-chan saveResult {
	oc := make(chan saveResult)

	go func() {
		const batchSize = 7
		batch := make([]string, 0)
		for input := range ic {
			if len(batch) < batchSize {
				batch = append(batch, input.inputData.id)
			} else {
				oc <- persistBatch(batch)
				batch = make([]string, 0)
			}

			if len(batch) > 0 {
				oc <- persistBatch(batch)
			}

			close(oc)
		}
	}()

	return oc
}

func persistBatch(batch []string) saveResult {
	return saveResult{batch, time.Now().UnixNano()}

}
