package mario_collector

import (
	"bufio"
	mario_collector "code.byted.org/data/datarangers-sdk-go"
	"fmt"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type consumer struct {
	max      int
	tickets  chan *ticket
	reader   []*os.File
	ans      int32
	oknum    int32
	failnum  int32
	register *register
}

type messList struct {
	User   *mario_collector.User    `json:"user,omitempty"`
	Event  []*mario_collector.Event `json:"event,omitempty"`
	Header *mario_collector.Header  `json:"header,omitempty"`
}

func NewConsumer(x int) (*consumer, error) {
	c := &consumer{ans: 0}
	c.max = x
	c.tickets = make(chan *ticket, c.max)
	for i := 0; i < c.max; i++ {
		c.tickets <- &ticket{id: i}
		logFile, err := os.OpenFile("./bak/test"+strconv.Itoa(i), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		c.reader = append(c.reader, logFile)
	}

	c.register = &register{max: x}
	var err error
	if c.register.writer, err = os.OpenFile("./bak/register", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666); err != nil {
		return nil, err
	}

	//先初始化
	if len(c.register.offset) != c.max {
		c.register.offset = []uint32{}
		for i := 0; i < c.max; i++ {
			c.register.offset = append(c.register.offset, 0)
		}
	}

	br := bufio.NewReader(c.register.writer)
	tmp := map[int]uint32{}
	for {
		line, err := br.ReadString('\n')
		line = strings.TrimSpace(line)
		if err == io.EOF {
			break
		}
		err = json.Unmarshal([]byte(line), &tmp)
		if err != nil {
			fmt.Println(err)
		}

		for k, v := range tmp {
			c.register.offset[k] = v
		}
	}

	//初始化
	for i := 0; i < c.max; i++ {
		c.register.preset = append(c.register.preset, 0)
	}
	return c, nil
}

func (this *consumer) execute() {
	//采集线程，
	for i := 0; i < this.max; i++ {
		go this.collcet(i)
	}

	//写入offset, 更新offset ，
	go this.register.writeOffset()
}

type register struct {
	max int
	//下一个记录的起点。
	offset   []uint32
	preset   []uint32
	writer   *os.File
	complete bool
}

func (this *consumer) collcet(j int) {
	br := bufio.NewReader(this.reader[j])
	client, _ := mario_collector.NewAppCollector()
	offset := this.register.offset[j]
	var k uint32
	count := 0
	for k = 0; k < offset; k++ {
		br.ReadLine()
		count += 1
	}
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		var tmpList messList
		if err := json.Unmarshal(a, &tmpList); err != nil {
			fmt.Println(err)
		}
		resp, err := client.Collect(tmpList.User, tmpList.Header, tmpList.Event)
		if err == nil {
			defer resp.Body.Close() // 保证连接复用
			//fmt.Println("response code:", resp.StatusCode) // 查看resp.StatusCode
			//body, _ := ioutil.ReadAll(resp.Body)
			//fmt.Println(string(body)) // 查看resp.Body
			atomic.AddInt32(&this.oknum, 1)
		} else {
			fmt.Println(err)
			atomic.AddInt32(&this.failnum, 1)
		}
		count += 1
		this.register.offset[j] = uint32(count)
	}
	fmt.Println("遍历了行", count)
	atomic.AddInt32(&this.ans, 1)
}

//写个json 进去。
func (this *register) writeOffset() {
	ans := map[int]uint32{}
	for {
		for i := 0; i < len(this.offset); i++ {
			//清零。
			ans[i] = this.offset[i]
			data, _ := json.Marshal(ans)
			this.writer.Truncate(0)
			this.writer.Write(data)
			this.writer.WriteString("\n")
		}
		//复制
		for i, _ := range this.offset {
			this.preset[i] = this.offset[i]
		}
		time.Sleep(1 * time.Second)
		if this.iscomplete() {
			fmt.Println(" - 写入结束 - ")
			this.complete = true
			break
		} else {
			fmt.Println(" - 持续写入中 - ")
		}

	}
}

//两次一样则认为 已经结束。
func (this *register) iscomplete() bool {
	for i, _ := range this.offset {
		if this.preset[i] != this.offset[i] {
			return false
		}
	}
	return true
}
