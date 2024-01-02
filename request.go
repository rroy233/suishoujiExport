package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Resp struct {
	Income    float64     `json:"income"`
	BeginDate string      `json:"beginDate"`
	Symbol    string      `json:"symbol"`
	PageCount int         `json:"pageCount"`
	EndDate   string      `json:"endDate"`
	PageNo    int         `json:"pageNo"`
	Payout    float64     `json:"payout"`
	Groups    []GroupData `json:"groups"`
}

type GroupData struct {
	Income float64    `json:"income"`
	Payout float64    `json:"payout"`
	List   []ListData `json:"list"`
}

type ListData struct {
	Account         int64    `json:"account"`
	BuyerAcount     string   `json:"buyerAcount"`
	BuyerAcountId   int64    `json:"buyerAcountId"`
	CategoryIcon    string   `json:"categoryIcon"`
	CategoryId      int64    `json:"categoryId"`
	CategoryName    string   `json:"categoryName"`
	Content         string   `json:"content"`
	CurrencyAmount  float64  `json:"currencyAmount"`
	Date            DateData `json:"date"`
	ImgId           int64    `json:"imgId"`
	ItemAmount      float64  `json:"itemAmount"`
	MemberId        int      `json:"memberId"`
	MemberName      string   `json:"memberName"`
	Memo            string   `json:"memo"`
	ProjectId       int      `json:"projectId"`
	ProjectName     string   `json:"projectName"`
	Relation        string   `json:"relation"`
	SId             string   `json:"sId"`
	SellerAcount    string   `json:"sellerAcount"`
	SellerAcountId  int64    `json:"sellerAcountId"`
	TranId          int64    `json:"tranId"`
	TranName        string   `json:"tranName"`
	TranType        int      `json:"tranType"`
	TransferStoreId int      `json:"transferStoreId"`
	Url             string   `json:"url"`
}

type DateData struct {
	Date           int   `json:"date"`
	Day            int   `json:"day"`
	Hours          int   `json:"hours"`
	Minutes        int   `json:"minutes"`
	Month          int   `json:"month"`
	Seconds        int   `json:"seconds"`
	Time           int64 `json:"time"`
	TimezoneOffset int   `json:"timezoneOffset"`
	Year           int   `json:"year"`
}

// MakeRequest 日期不要求都是2位
func MakeRequest(beginDate string, endDate string, page int) (*Resp, error) {
	form := fmt.Sprintf("opt=list2&beginDate=%s&endDate=%s&cids=0&bids=0&sids=0&pids=0&memids=0&order=&isDesc=0&page=%d&note=&mids=0",
		beginDate,
		endDate,
		page,
	)
	body := strings.NewReader(form)
	req, err := http.NewRequest("POST", "https://www.sui.com/tally/new.rmi", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authority", "www.sui.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh-TW;q=0.9,zh;q=0.8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", string(cookie))
	req.Header.Set("Dnt", "1")
	req.Header.Set("Origin", "https://www.sui.com")
	req.Header.Set("Referer", "https://www.sui.com/tally/new.do")
	req.Header.Set("Sec-Ch-Ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respData := new(Resp)
	err = json.Unmarshal(data, respData)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
