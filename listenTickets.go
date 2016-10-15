package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type QueryLeftNewDTOS struct {
	Train_no                 string `json:"train_no"`
	Station_train_code       string `json:"station_train_code"`
	Start_station_telecode   string `json:"start_station_telecode"`
	Start_station_name       string `json:"start_station_name"`
	End_station_telecode     string `json:"end_station_telecode"`
	End_station_name         string `json:"end_station_name"`
	From_station_telecode    string `json:"from_station_telecode"`
	From_station_name        string `json:"from_station_name"`
	To_station_telecode      string `json:"to_station_telecode"`
	To_station_name          string `json:"to_station_name"`
	Start_time               string `json:"start_time"`
	Arrive_time              string `json:"arrive_time"`
	Day_difference           string `json:"day_difference"`
	Train_class_name         string `json:"train_class_name"`
	Lishi                    string `json:"lishi"`
	CanWebBuy                string `json:"canWebBuy"`
	LishiValue               string `json:"lishiValue"`
	Yp_info                  string `json:"yp_info"`
	Control_train_day        string `json:"control_train_day"`
	Start_train_date         string `json:"start_train_date"`
	Seat_feature             string `json:"seat_feature"`
	Yp_ex                    string `json:"yp_ex"`
	Train_seat_feature       string `json:"train_seat_feature"`
	Seat_types               string `json:"seat_types"`
	Location_code            string `json:"location_code"`
	From_station_no          string `json:"from_station_no"`
	To_station_no            string `json:"to_station_no"`
	Control_day              int    `json:"control_day"`
	Sale_time                string `json:"sale_time"`
	Is_support_card          string `json:"is_support_card"`
	Controlled_train_flag    string `json:"controlled_train_flag"`
	Controlled_train_message string `json:"controlled_train_message"`
	Gg_num                   string `json:"gg_num"`
	Gr_num                   string `json:"gr_num"`
	Qt_num                   string `json:"qt_num"`
	Rw_num                   string `json:"rw_num"`
	Rz_num                   string `json:"rz_num"`
	Tz_num                   string `json:"tz_num"`
	Wz_num                   string `json:"wz_num"`
	Yb_num                   string `json:"yb_num"`
	Yw_num                   string `json:"yw_num"`
	Yz_num                   string `json:"yz_num"`
	Ze_num                   string `json:"ze_num"`
	Zy_num                   string `json:"zy_num"`
	Swz_num                  string `json:"swz_num"`
}
type DataS struct {
	QueryLeftNewDTO QueryLeftNewDTOS `json:"queryLeftNewDTO"`
	SecretStr       string           `json:"secretStr"`
	ButtonTextInfo  string           `json:"buttonTextInfo"`
}
type MessagesS interface{}
type ValidateMessagesS interface{}
type Root struct {
	ValidateMessagesShowId string            `json:"validateMessagesShowId"`
	Status                 bool              `json:"status"`
	Httpstatus             int               `json:"httpstatus"`
	Data                   []DataS           `json:"data"`
	Messages               []MessagesS       `json:"messages"`
	ValidateMessages       ValidateMessagesS `json:"validateMessages"`
}

var (
	client *http.Client
	once   sync.Once
)

func checkErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func getTime(h int, m int) int {
	return h*60 + m
}

func main() {
	flag := false
	count := 0
	once.Do(func() {
		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client = &http.Client{Transport: tr}

	})
	for {
		count++
		println("query number:", count)
		requset, errNewRequest := http.NewRequest("GET", "https://kyfw.12306.cn/otn/leftTicket/queryT?leftTicketDTO.train_date=2016-10-02&leftTicketDTO.from_station=NJH&leftTicketDTO.to_station=UAH&purpose_codes=ADULT", nil)
		checkErr(errNewRequest)
		rep, errDo := client.Do(requset)
		checkErr(errDo)
		body, errRead := ioutil.ReadAll(rep.Body)
		rep.Body.Close()
		checkErr(errRead)
		root := &Root{}
		json.Unmarshal(body, root)
		if root.Data == nil {
			continue
		}
		for _, v := range root.Data {
			startTime := v.QueryLeftNewDTO.Start_time
			ts := strings.Split(startTime, ":")
			h, errAtoi := strconv.Atoi(ts[0])
			checkErr(errAtoi)
			m, errAtoi := strconv.Atoi(ts[1])
			checkErr(errAtoi)
			t := h*60 + m
			if t >= getTime(9, 40) && t <= getTime(12, 0) {
				num := v.QueryLeftNewDTO.Ze_num
				if num != "--" && num != "无" {
					println("found ticket ,time:", startTime)
					flag = true
					cmd := exec.Command("cmd.exe", "/c", "start "+"https://kyfw.12306.cn/otn/leftTicket/init")

					err := cmd.Run()
					if err != nil {
						println(err.Error())
					}
					break
				}
			}

		}
		if flag {
			break
		}
		time.Sleep(1 * time.Second)
	}

}
