package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	router := gin.Default()
	router.Static("/public", "./public")
	router.Static("/asin", "./asin")
	router.Static("/trans", "./trans")
	router.LoadHTMLGlob("templates/*")
	router.POST("/upload", func(c *gin.Context) {
		// Source
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		filename := filepath.Base(file.Filename)
		filedir := filepath.Dir(file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		prefix := "http://localhost/asin/"
		cellMap := map[int]string{2: "B", 3: "C", 4: "D", 5: "E", 6: "F", 7: "G", 8: "H", 9: "I", 10: "J"}
		f := excelize.NewFile()
		// Create a new sheet.
		index := f.NewSheet("Sheet1")
		// Set value of a cell.
		f.SetCellValue("Sheet1", "A1", "ASIN")
		f.SetCellValue("Sheet1", "B1", "高清图片1")
		f.SetCellValue("Sheet1", "C1", "高清图片2")
		f.SetCellValue("Sheet1", "D1", "高清图片3")
		f.SetCellValue("Sheet1", "E1", "高清图片4")
		f.SetCellValue("Sheet1", "F1", "高清图片5")
		f.SetCellValue("Sheet1", "G1", "高清图片6")
		f.SetCellValue("Sheet1", "H1", "高清图片7")
		f.SetCellValue("Sheet1", "I1", "高清图片8")
		f.SetCellValue("Sheet1", "J1", "高清图片9")
		f1, err := excelize.OpenFile(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		indexName := f1.GetSheetName(f1.GetActiveSheetIndex())
		cols, err := f1.GetCols(indexName)
		if err != nil {
			return
		}
		value := ""
		allimgs := map[string]string{}
		for i, col := range cols {
			if i == 0 {
				for k, rowCell := range col {
					if rowCell != "ASIN" {
						f.SetCellValue("Sheet1", "A"+strconv.Itoa(k+1), rowCell)
					}
					if k != 0 {
						pic, err := f1.GetCellValue(indexName, "B"+strconv.Itoa(k+1))
						if err != nil {
							fmt.Println(err)
							return
						}
						urlarr := strings.Split(pic, "|")
						length := len(urlarr)
						for i := 1; i <= length; i++ {
							item := i - 1
							imgName := rowCell + "_" + strconv.Itoa(i) + ".jpg"
							allimgs[imgName] = urlarr[item]
							value = prefix + imgName
							f.SetCellValue("Sheet1", cellMap[i+1]+strconv.Itoa(k+1), value)
						}
					}
				}
			}
		}
		start := time.Now()
		DownloadAllFiles(filedir+"/asin/", allimgs)
		end := time.Now()
		fmt.Println(end.Sub(start))
		// Set active sheet of the workbook.
		f.SetActiveSheet(index)
		// Save xlsx file by the given path.
		if err := f.SaveAs("./trans/" + "new-" + filename); err != nil {
			fmt.Println(err)
		}
		c.HTML(
			http.StatusOK,
			"notify.html",
			gin.H{
				"Linkv":    "http://localhost/trans/new-" + filename,
				"FileName": "new-" + filename,
			},
		)
	})
	router.POST("/login", func(c *gin.Context) {
		user := c.PostForm("uname")
		psw := c.PostForm("psw")
		if user == "$username" && psw == "$password" {
			// Call the HTML method of the Context to render a template
			c.HTML(
				http.StatusOK,
				"upload.html",
				gin.H{
					"title": "Home Page",
				},
			)
		} else {
			c.String(http.StatusUnauthorized, fmt.Sprint("用户名或密码错误，请重新输入！"))
		}
	})
	router.Run(":80")
}

func DownloadAllFiles(path string, allimgs map[string]string) {
	workersCount := 20
	wg := &sync.WaitGroup{}
	imgChan := make(chan map[string]string)
	for i := 0; i < workersCount; i++ {
		go func() {
			for imgUrl := range imgChan {
				_ = DownloadFile(path, imgUrl)
				wg.Done()
			}
		}()
	}
	for i, u := range allimgs {
		img := map[string]string{i: u}
		imgChan <- img
		wg.Add(1)
	}
	wg.Wait()
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(path string, imgUrl map[string]string) error {
	for i, u := range imgUrl {
		// Get the data
		resp, err := http.Get(u)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Create the file
		out, err := os.Create(path + i)
		if err != nil {
			return err
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		return err
	}
	return nil
}
