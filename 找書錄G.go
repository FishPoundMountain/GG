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
	借書天數 = 7
	//預定使用者名稱  = "訪客"
	//預定使用者密碼  = "abcd1234"
)

func loadBooks() ([]Book, error) {
	fileData, err := ioutil.ReadFile(filePath)
	var books []Book

	if err != nil {
		return nil, fmt.Errorf("讀取文件錯誤: %v", err)
	}

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
	var newBook Book
	books, err := loadBooks()

	if err != nil {
		fmt.Println(err)
		return
	}

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
	var 丟失書籍 string
	var 確認 string
	var 書本ID int
	books, err := loadBooks()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("請輸入丟失的書名：")
	fmt.Scanln(&丟失書籍)

	for i, book := range books {

		if book.Name == 丟失書籍 {
			fmt.Printf("請確認丟失的書籍：\n書名：%s，作者：%s，書籍ID：%d，借還狀態：%s，借閱人：%s\na.正確  b.重來：",
				book.Name, book.Author, book.ID, getStateString(book.State, 0), getStateString(book.State, 1))
			
			fmt.Scanln(&確認)
			fmt.Print("為避免同名書，請確認書本ID：")
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
	var 所找書籍 string
	var 是否借書 string
	var 書本ID int
	var 借書人名 string
	books, err := loadBooks()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("請問你想借甚麼書呢？")
	fmt.Scanln(&所找書籍)

	for i, book := range books {

		if book.Name == 所找書籍 {
			fmt.Printf("想要借這本書嗎？\n書名：%s，作者：%s，書籍ID：%d，借還狀態：%s，借閱人：%s\na.是  b.否：",
				book.Name, book.Author, book.ID, getStateString(book.State, 0), getStateString(book.State, 1))
			fmt.Scanln(&是否借書)

			if 是否借書 == "a" {
				fmt.Print("為避免同名書，請確認書本ID：")
				fmt.Scanln(&書本ID)

				if getStateString(book.State, 0) == "未借出" && book.ID == 書本ID {
					now := time.Now()

					fmt.Print("請輸入你的名字：")
					fmt.Scanln(&借書人名)
					
					借書時間 := now.Format("01月02號")
					預計還書時間 := now.AddDate(0, 0, 借書天數).Format("01月01號")
					books[i].State = []interface{}{"已借出", 借書人名, 借書時間, 預計還書時間}
					if err := saveBooks(books); err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println("借書成功！")
					return
				}
				fmt.Printf("好書被搶先啦，%s之後再來吧！\n", getStateString(book.State, 3))
				return
			}
			fmt.Println("好書不可錯過，歡迎下次來借這本書喔！")
			return
		}
	}
	fmt.Println("請檢查是不是打錯字了喔！")
}

func 還書() {
	var 所還書籍 string
	var 是否還書 string
	var 書本ID int
	books, err := loadBooks()
	now := time.Now()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("請輸入歸還的書名：")
	fmt.Scanln(&所還書籍)

	for i, book := range books {
		if book.Name == 所還書籍 {
			fmt.Printf("確定歸還以下書籍：\n書名：%s，作者：%s，書籍ID：%d，借還狀態：%s，借閱人：%s\n", book.Name, book.Author, book.ID, getStateString(book.State, 0), getStateString(book.State, 1))
			fmt.Print("為避免同名書，請確認書本ID：")
			fmt.Scanln(&書本ID)
			fmt.Println("是否還書  a.是  b.否：")
			fmt.Scanln(&是否還書)

			if 是否還書 == "a" && getStateString(book.State, 0) == "已借出" && book.ID == 書本ID {
				books[i].State = []interface{}{"未借出", "無", "", ""}
				fmt.Println("還書中......")
				if err := saveBooks(books); err != nil {
					fmt.Println(err)
					return
				}
				if now.Format("01月01號") == getStateString(book.State, 3){
					fmt.Print("養成有借有還得習慣很重要的喔")
				}
				fmt.Println("還書成功！")
				return
			} else {
				fmt.Println("書籍狀態：\n借還狀態：", getStateString(book.State, 0), "  ，借閱人：\n", getStateString(book.State, 1))
				return
			}
			fmt.Println("該書已經歸還，如有疑問請找圖書管理員")
			return
		}
	}
	fmt.Println("請檢查是不是打錯字了喔！")
}

func login() bool {
	var InputName string
	var InputPassWord string
	var UserDict  = map[string]string {"visitor":"abcd1234"}
	YN := true
	fmt.Println("請先登入")
	fmt.Print("使用者名稱 : ")
	fmt.Scanln(&InputPassWord)
	fmt.Print("密碼 : ")
	fmt.Scanln(&InputName)

    for k, v := range(UserDict) {
        fmt.Println("k, v =", k, v)
        if InputName == k && InputPassWord == v{
            fmt.Println("登入成功!")
            YN = true
        } else {
			fmt.Println("登入失敗")
			YN = false
		}
    }
	return YN
}
/*
func 沒還書的人(){
	還書日期 = "%s之後再來吧！\n", getStateString(book.State, 3)
}

func main() {
	if 成功登入 := login(); 成功登入 {
		for {
			fmt.Println("+___________ଲ借書操作板ଲ__________+")
			fmt.Print("|需要做甚麼？  \n|a.新書紀錄\t\t\t|  \n|b.借書\t\t\t|  \n|c.還書\t\t\t|  \n|d.登記遺失書籍\t\t\t|  \n|e.退出 : \t\t\t|")
			fmt.Println("\n|______________________________」")
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
	} else {
		fmt.Println("請再試一次吧!")
	}

}
*/
//ଲ(ⓛ ω ⓛ)ଲ

func main() {
	var 執行 string

	for {
		fmt.Println("+___________ଲ借書操作板ଲ__________+")
		fmt.Print("|需要做甚麼？  \n|\ta.新書紀錄\t\t|  \n|\tb.想要借書\t\t|  \n|\tc.想要還書\t\t|  \n|\td.登記遺失書籍\t\t|  \n|\te.退出 : \t\t|")
		fmt.Println("\n|______________________________」")
		fmt.Scanln(&執行)

		switch 執行 {
			case "a":
				新書紀錄()	
				fmt.Println("\n")		
			case "b":
				借書()
				fmt.Println("\n")		
			case "c":
				還書()
				fmt.Println("\n")		
			case "d":
				書籍丟失()
				fmt.Println("\n")		
			case "e":
				fmt.Println("歡迎再來！")
				fmt.Println("\n")		
				return
			default:
				fmt.Println("似乎輸入錯了，再來一次吧！")
				fmt.Println("\n")		
		}
	}
}
