package main

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gopkg.in/rroy233/logger.v2"
	"log"
	"os"
	"time"
)

var styleMid int
var cookie []byte

type Config struct {
	Cookie    string `json:"cookie"`
	BeginTime string `json:"beginTime"`
	EndTime   string `json:"endTime"`
}

var beginDate = "2001.01.01" //开始时间
var endDate = "2023.12.31"   //结束时间

func main() {

	defer func() {
		log.Println("程序3s后结束")
		time.Sleep(3 * time.Second)
	}()

	logger.New(&logger.Config{
		StdOutput:      true,
		StoreLocalFile: true,
	})

	var err error

	//读取配置
	cfFile, err := os.ReadFile("./config.json")
	if err != nil {
		logger.Error.Println("读取配置文件失败：", err)
		return
	}
	cf := new(Config)
	if err := json.Unmarshal(cfFile, cf); err != nil {
		logger.Error.Println("解析配置文件失败：", err)
		return
	}

	beginDate = cf.BeginTime
	endDate = cf.EndTime
	cookie = []byte(cf.Cookie)

	//获取第一页，获得总页数
	firstPage, err := MakeRequest(beginDate, endDate, 1)
	if err != nil {
		logger.FATAL.Println(err)
	}

	logger.Info.Printf("导出任务已创建，一共%d页\n", firstPage.PageCount)

	//初始化excel
	f := excelize.NewFile()
	//初始化样式
	styleMid, err = f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	//填入导出信息
	f.SetCellValue("Sheet1", "A1", "导出时间")
	f.SetCellValue("Sheet1", "A2", time.Now().Format(time.DateTime))

	activeSheet := "数据"
	totalPage := firstPage.PageCount
	//totalPage := 2

	initSheet(f, activeSheet)
	excelLine := 2 //表格当前行

	//遍历每一页
	for page := 1; page <= totalPage; page++ {
		if page%10 == 0 {
			logger.Info.Println("每10页休息10s")
			time.Sleep(10 * time.Second)
		} else {
			logger.Info.Println("休息2s")
			time.Sleep(2 * time.Second)
		}

		logger.Info.Printf("第[%d]页 - 开始处理\n", page)
		resp, err := MakeRequest(beginDate, endDate, page)
		if err != nil {
			logger.Info.Println(err)
			return
		}
		logger.Info.Printf("第[%d]页 - 数据获取成功，即将处理 [%s]-[%s]\n", page, resp.BeginDate, resp.EndDate)

		//遍历该页每一项
		for _, group := range resp.Groups {
			//每个group又有一个list，里面即是每一笔帐
			for _, listData := range group.List {
				dateTime := time.Unix(listData.Date.Time/1000, 0)
				data := map[string]interface{}{
					"A": dateTime.Format("2006-01-02 15:04"), //时间
					"B": listData.CategoryName,               //分类
					"C": listData.TranName,                   //类型：收入 支出 报销 转账 还款
					"D": listData.ItemAmount,                 //金额
					"E": listData.BuyerAcount,                //账户1
					"F": listData.SellerAcount,               //账户2
					"G": listData.Memo,                       //备注
					"H": "",                                  //账单标记
					"I": "",                                  //账单图片
				}
				//插入表格
				for k, v := range data {
					err = f.SetCellValue(activeSheet, fmt.Sprintf("%s%d", k, excelLine), v)
					if err != nil {
						logger.FATAL.Println(err)
						return
					}
				}

				excelLine++
			}

		}
		logger.Info.Printf("第[%d]页 - 处理完成 [%s]-[%s]。Excel表格当前行：%d\n", page, resp.BeginDate, resp.EndDate, excelLine)

	}

	outFileName := fmt.Sprintf("out-%s-%s.xlsx", beginDate, endDate)
	if err := f.SaveAs(outFileName); err != nil {
		logger.FATAL.Println(err)
		return
	}
	logger.Info.Println("已保存为" + outFileName)

}

func initSheet(file *excelize.File, sheetName string) {
	header := map[string]string{
		"A1": "时间",
		"B1": "分类",
		"C1": "类型",
		"D1": "金额",
		"E1": "账户1",
		"F1": "账户2",
		"G1": "备注",
		"H1": "账单标记",
		"I1": "账单图片",
	}
	file.NewSheet(sheetName)
	for k, v := range header {
		err := file.SetCellValue(sheetName, k, v)
		if err != nil {
			log.Println(err)
			return
		}
		file.SetColWidth(sheetName, "A", "A", 16)
		file.SetColWidth(sheetName, "B", "B", 12)
		file.SetColWidth(sheetName, "C", "D", 10)

		file.SetColWidth(sheetName, "E", "G", 22)
		file.SetColWidth(sheetName, "H", "I", 15)
		file.SetColStyle(sheetName, "A:I", styleMid)
	}
	return
}
