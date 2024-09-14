package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

type Book struct {
	Name   string        `json:"name"`
	Author string        `json:"author"`
	ID     int           `json:"ID"`
	State  []interface{} `json:"state"`
}

const (
	filePath = "./LibraryIndex.json"
	借書天數     = 7
)

func loadBooks() ([]Book, error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("讀取文件錯誤: %v", err)
	}

	var books []Book
	if err := json.Unmarshal(fileData, &books); err != nil {
		return nil, fmt.Errorf("解析JSON錯誤: %v", err)
	}

	return books, nil
}

func saveBooks(books []Book) error {
	newData, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return fmt.Errorf("生成JSON錯誤: %v", err)
	}

	if err := ioutil.WriteFile(filePath, newData, 0644); err != nil {
		return fmt.Errorf("寫入文件錯誤: %v", err)
	}

	return nil
}

func getStateString(state []interface{}, index int) string {
	if index >= len(state) {
		return ""
	}
	switch v := state[index].(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func 新書紀錄() {
	books, err := loadBooks()
	if err != nil {
		fmt.Println(err)
		return
	}

	var newBook Book
	fmt.Print("請輸入新書的書名：")
	fmt.Scanln(&newBook.Name)
	fmt.Print("請輸入作者名：")
	fmt.Scanln(&newBook.Author)
	fmt.Print("請輸入書本的ID或索引碼：")
	fmt.Scanln(&newBook.ID)
	newBook.State = []interface{}{"未借出", "無", "", ""}

	books = append(books, newBook)
	if err := saveBooks(books); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("登記成功！")
}

func 書籍丟失() {
	books, err := loadBooks()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("請輸入丟失的書名：")
	var 丟失書籍 string
	fmt.Scanln(&丟失書籍)

	for i, book := range books {
		if book.Name == 丟失書籍 {
			fmt.Printf("請確認丟失的書籍是否如下：\n書名：%s，作者：%s，書籍ID：%d，借還狀態：%s，借閱人：%s\na.正確  b.重來：",
				book.Name, book.Author, book.ID, getStateString(book.State, 0), getStateString(book.State, 1))
			var 確認 string
			fmt.Scanln(&確認)

			fmt.Print("為避免同名書請確認書本ID：")
			var 書本ID int
			fmt.Scanln(&書本ID)

			if 確認 == "a" && book.ID == 書本ID {
				books = append(books[:i], books[i+1:]...)
				if err := saveBooks(books); err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("成功註銷該書！")
				return
			}
			fmt.Println("已取消操作")
			return
		}
	}
	fmt.Println("書籍已不存在，或檢查是不是打錯字了喔！")
}

func 借書() {
	books, err := loadBooks()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("請問你想找甚麼書呢？")
	var 所找書籍 string
	fmt.Scanln(&所找書籍)

	for i, book := range books {
		if book.Name == 所找書籍 {
			fmt.Printf("想要借這本書嗎？\n書名：%s，作者：%s，書籍ID：%d，借還狀態：%s，借閱人：%s\na.是  b.否：",
				book.Name, book.Author, book.ID, getStateString(book.State, 0), getStateString(book.State, 1))
			var 是否借書 string
			fmt.Scanln(&是否借書)

			if 是否借書 == "a" {
				fmt.Print("為避免同名書請確認書本ID：")
				var 書本ID int
				fmt.Scanln(&書本ID)

				if getStateString(book.State, 0) == "未借出" && book.ID == 書本ID {
					fmt.Print("請輸入你的名字：")
					var 借書人名 string
					fmt.Scanln(&借書人名)

					now := time.Now()
					借書時間 := now.Format("01月02號")
					預計還書時間 := now.AddDate(0, 0, 借書天數).Format("01月02號")

					books[i].State = []interface{}{"已借出", 借書人名, 借書時間, 預計還書時間}
					if err := saveBooks(books); err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println("借書成功！")
					return
				}
				fmt.Printf("好書被搶先啦～%s再來吧！\n", getStateString(book.State, 3))
				return
			}
			fmt.Println("好書不可錯過，歡迎下次來借這本書喔！")
			return
		}
	}
	fmt.Println("請檢查是不是打錯字了喔！")
}

func 還書() {
	books, err := loadBooks()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("請輸入歸還的書名：")
	var 所還書籍 string
	fmt.Scanln(&所還書籍)

	for i, book := range books {
		if book.Name == 所還書籍 {
			fmt.Printf("確定歸還以下書籍：\n書名：%s，作者：%s，書籍ID：%d，借還狀態：%s，借閱人：%s\na.是  b.否：",
				book.Name, book.Author, book.ID, getStateString(book.State, 0), getStateString(book.State, 1))
			var 是否還書 string
			fmt.Scanln(&是否還書)

			fmt.Print("為避免同名書請確認書本ID：")
			var 書本ID int
			fmt.Scanln(&書本ID)

			if 是否還書 == "a" && getStateString(book.State, 0) == "已借出" && book.ID == 書本ID {
				books[i].State = []interface{}{"未借出", "無", "", ""}
				fmt.Println("還書中......")
				if err := saveBooks(books); err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("還書成功！")
				return
			}
			fmt.Println("該書已經歸還，如有疑問請找圖書管理員")
			return
		}
	}
	fmt.Println("請檢查是不是打錯字了喔！")
}

func main() {
	for {
		fmt.Print("需要做甚麼呢？  a.新書紀錄  b.借書  c.還書  d.書籍遺失登記  e.退出  :  ")
		var 執行 string
		fmt.Scanln(&執行)

		switch 執行 {
		case "a":
			新書紀錄()
		case "b":
			借書()
		case "c":
			還書()
		case "d":
			書籍丟失()
		case "e":
			fmt.Println("謝謝使用，再見！")
			return
		default:
			fmt.Println("似乎輸入錯了，再來一次吧！")
		}
	}
}

//各種內心獨白
//以前自學python 學得毫無章法，還沒練習過如何處裡json 檔，要不要練習一下呢？可是花時間去研究這個的話會沒時間學go 怎麼辦……
//不然我先用python 寫完再讓AI 幫我翻譯成go 好了
//還有好多邊界條件沒排掉
//不行了不能再改下去了，再改下去會沒時間睡覺，晚睡就會晚起，晚起9/15就會很晚才去辦台灣人居住證，就會沒時間學go了！！
//AI 翻譯成go 之後有錯誤訊息QAQ
//我筆電(筆記本)目前只能打繁體字，人家真的看得懂嗎？
