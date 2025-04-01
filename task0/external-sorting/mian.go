package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
)

var B = 200

func main() {
	//设数据存在./data.bin文件里，二进制方式存储，每个数据4个字节长度，用int32，大端存
	GenerateData()
	//先获取文件长度
	status, err := os.Stat("./data.bin")
	if err != nil {
		log.Fatalf("数据文件打开失败")
		return
	}
	fileLength := status.Size()

	file, err := os.Open("./data.bin")
	if err != nil {
		log.Fatalf("数据文件打开失败")
		return
	}
	defer file.Close()

	var bufLength = int64(B) //每次1024个数据？
	var off int64 = 0        //偏移量

	//生成归并段,要控制同一时段最多10个线程运行
	routineC := make(chan bool, 10) //用管道来控制开了几个线程
	sum := 0                        //记录产生了多少个归并段
	lock := sync.WaitGroup{}
	for i := 0; off < fileLength; i++ {
		sum++
		lock.Add(1)
		go func(i int) {
			defer lock.Done()

			routineC <- true

			offset := int64(i) * bufLength //根据i获取偏移量
			buf := make([]byte, bufLength) //缓冲区

			n, err := file.ReadAt(buf, offset)
			if err != nil && err != io.EOF {
				log.Fatalf("第%d次读取数据失败", i)
				return
			}

			//转数字
			count := n / 4
			nums := make([]int32, count)
			for j := 0; j < count; j++ {
				nums[j] = int32(binary.BigEndian.Uint32(buf[j*4 : j*4+4]))
			}

			//开始排序
			quicksort(nums, 0, len(nums)-1)
			//log.Println("nums:", nums)

			//先把nums写回二进制类型
			buf = make([]byte, n)
			for j, num := range nums {
				binary.BigEndian.PutUint32(buf[j*4:], uint32(num))
			}

			//再写入临时文件
			name := fmt.Sprintf("0-data-%d.tmp", i)
			f, err := os.Create(name)
			if err != nil {
				log.Fatalf("创建临时文件失败: %v", err)
			}
			_, err = f.Write(buf)
			if err != nil {
				log.Fatalf("写入临时文件失败: %v", err)
			}
			f.Close()
			<-routineC //放行
		}(i)
		off += bufLength
	}
	lock.Wait()

	//可以k路归并了，这里k最大取10，多线程
	forNum := 0
	sump := sum
	for {
		forNum++
		lock = sync.WaitGroup{}
		s := 0 //记录开了几个线程，要控制同一时段最多5个线程运行
		routineC = make(chan bool, 5)
		for i := 0; i*10 < sump; i++ {
			s++
			lock.Add(1)
			go func(i int) {
				routineC <- true
				defer lock.Done()

				finalname := fmt.Sprintf("%d-data-%d.tmp", forNum, i)
				//log.Println("将生成文件", finalname)
				k := 10
				if sump-i*10 < 10 {
					k = sump - i*10
				}
				//执行k路归并

				knums := make([][]int32, 0)
				ktimes := make([]int64, 0)
				for j := 0; j < k; j++ {
					ktimes = append(ktimes, 0)

					buf := make([]byte, B) //一次性不一定能读完，但每次最多读1024个数据
					name := fmt.Sprintf("%d-data-%d.tmp", forNum-1, i*10+j)
					file, _ := os.Open(name)
					n, _ := file.ReadAt(buf, ktimes[j]*int64(B))
					ktimes[j]++
					//转数字
					count := n / 4
					knums = append(knums, make([]int32, count))
					for u := 0; u < count; u++ {
						knums[j][u] = int32(binary.BigEndian.Uint32(buf[u*4 : u*4+4]))
					}
					file.Close()
				}
				//log.Println("所有数据:", knums)
				//还要有一个写入缓冲区？
				writeBuf := make([]byte, B)
				cal := 0

				minHeap := make([]int32, k) //最小堆
				origin := make([]int, k)    //每个数据存放在哪个归并段
				numsOff := make([]int, k)   //每个归并段读取了多少个数据(偏移量)

				for u := 0; u < k; u++ {
					minHeap[u] = knums[u][0]
					origin[u] = u
					numsOff[u] = 1
					//与父节点比大小
					f := u
					for f >= 0 {
						if minHeap[f] >= minHeap[(f-1)/2] {
							break
						}
						minHeap[f], minHeap[(f-1)/2] = minHeap[(f-1)/2], minHeap[f]
						origin[f], origin[(f-1)/2] = origin[(f-1)/2], origin[f]
						f = (f - 1) / 2
					}
				}
				//log.Println("最小堆:", minHeap)

				for k > 0 {
					if cal >= B {
						//log.Println("buf:", writeBuf)
						file, err := os.OpenFile(finalname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
						if err != nil {
							log.Println("打开失败,", err)
							file.Close()
							return
						}
						_, err = file.Write(writeBuf)
						if err != nil {
							log.Printf("写入失败,%v", err)
							file.Close()
							return
						}
						//log.Println("一次写入")
						cal = 0
						writeBuf = make([]byte, B)
						file.Close()
					}
					//弹出一个数据
					binary.BigEndian.PutUint32(writeBuf[cal:cal+4], uint32(minHeap[0]))
					cal += 4

					if numsOff[origin[0]] == -1 {
						minHeap[0] = minHeap[k-1]
						origin[0] = origin[k-1]
						minHeap = minHeap[:k-1]
						k--
					} else if numsOff[origin[0]] >= len(knums[origin[0]]) {
						//当前读取到的归并段数据已经排序完成，要检验该归并段是否还有数据没有读取
						buf := make([]byte, B)
						j := origin[0]
						name := fmt.Sprintf("%d-data-%d.tmp", forNum-1, i*10+j)
						file, _ := os.Open(name)
						n, _ := file.ReadAt(buf, ktimes[j]*int64(B))
						if n > 0 {
							ktimes[j]++
							count := n / 4
							knums[j] = make([]int32, count)
							for u := 0; u < count; u++ {
								knums[j][u] = int32(binary.BigEndian.Uint32(buf[u*4 : u*4+4]))
							}
							numsOff[j] = 0
							minHeap[0] = knums[j][numsOff[j]]
							numsOff[j]++
						} else {
							//则该归并段数据已经完全输入完成
							numsOff[origin[0]] = -1 //标记一下
							//把最后一个数据放到开头
							minHeap[0] = minHeap[k-1]
							origin[0] = origin[k-1]
							minHeap = minHeap[:k-1]
							k--
						}
						file.Close()
					} else {
						minHeap[0] = knums[origin[0]][numsOff[origin[0]]]
						numsOff[origin[0]]++
						//log.Println("新数据插入")
					}
					//调整最小堆
					f := 0
					for {
						l := 2*f + 1
						r := 2*f + 2
						smallest := f
						if l < k && minHeap[l] < minHeap[smallest] {
							smallest = l
						}
						if r < k && minHeap[r] < minHeap[smallest] {
							smallest = r
						}
						if smallest == f {
							break
						}
						// 交换值和 origin
						minHeap[f], minHeap[smallest] = minHeap[smallest], minHeap[f]
						origin[f], origin[smallest] = origin[smallest], origin[f]
						f = smallest
					}
					//log.Println("调整最小堆:", minHeap)
				}
				file, err := os.OpenFile(finalname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					log.Println("打开失败,", err)
					return
				}
				_, err = file.Write(writeBuf[:cal])
				if err != nil {
					log.Printf("写入失败,%v", err)
					return
				}
				//log.Println("一次写入")
				file.Close()

				//log.Println("文件写入完成:", finalname)
				//tmp文件可以删了
				for j := 0; ; j++ {
					name := fmt.Sprintf("%d-data-%d.tmp", forNum-1, i*10+j)
					if err := os.Remove(name); err != nil {
						break
					}
				}
				<-routineC
			}(i)
		}
		lock.Wait()
		sump = s
		if s == 1 {
			//说明是最后一次归并，得到最终结果
			oldname := fmt.Sprintf("%d-data-0.tmp", forNum)
			newname := "final"
			if err := os.Rename(oldname, newname); err != nil {
				log.Println("err:", err)
			}
			break
		}
	}

	fmt.Printf("排序结果在final文件中\n")
}

func quicksort(nums []int32, l, r int) {
	if l >= r {
		return
	}
	left := l
	right := r
	star := nums[left]
	for left < right {
		for nums[right] >= star && left < right {
			right--
		}
		if left == right {
			break
		}
		nums[left] = nums[right]
		left++
		for nums[left] <= star && left < right {
			left++
		}
		if left == right {
			break
		}
		nums[right] = nums[left]
		right--
	}
	nums[left] = star
	quicksort(nums, l, right-1)
	quicksort(nums, left+1, r)
}

func GenerateData() {
	file, _ := os.Create("data.bin")
	buf := make([]byte, 4096)
	for i := 0; i < 1024; i++ {
		num := rand.Uint32() % 10000
		binary.BigEndian.PutUint32(buf[i*4:i*4+4], num)
	}
	file.Write(buf)
	file.Close()
}
