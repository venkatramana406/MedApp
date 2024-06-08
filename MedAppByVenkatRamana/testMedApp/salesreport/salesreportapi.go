package salesreport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type SalesReportReqStruct struct {
	From_date string `json:"from_date"`
	To_date   string `json:"to_date"`
}
type SalesReportStruct struct {
	Bill_No       string  `json:"bill_no"`
	Bill_date     string  `json:"bill_date"`
	Medicine_Name string  `json:"medicine_name"`
	Quantity      int     `json:"quantity"`
	Amount        float64 `json:"amount"`
}
type SalesReportRespStruct struct {
	SalesArr []SalesReportStruct `json:"salesarr"`
	ErrMsg   string              `json:"errmsg"`
	Status   string              `json:"status"`
	Msg      string              `json:"msg"`
}

func LocalDBConnect() (*sql.DB, error) {
	log.Println("LocalDBConnect+")
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "root", "192.168.2.5", 3306, "venkatramana")
	db, err := sql.Open("mysql", connString)
	if err != nil {
		log.Println("Open connection failed:", err.Error())
		return db, err
	}
	log.Println("LocalDBConnect-")
	return db, nil
}

func SalesReportAPI(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "PUT" {
		log.Println("SalesReportAPI(+)")
		var DateRangeRec SalesReportReqStruct
		var SalesfinalRespRec SalesReportRespStruct
		SalesfinalRespRec.Status = "S"
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error :", err)
			SalesfinalRespRec.Status = "E"
			SalesfinalRespRec.ErrMsg = "SR01" + err.Error()
		} else {
			err := json.Unmarshal(body, &DateRangeRec)
			if err != nil {
				log.Println("Error :", err)
				SalesfinalRespRec.Status = "E"
				SalesfinalRespRec.ErrMsg = "SR02" + err.Error()
			} else {
				data, err := GetSalesReport(DateRangeRec)
				if err != nil {
					SalesfinalRespRec.Status = "E"
					SalesfinalRespRec.Msg = ""
					SalesfinalRespRec.ErrMsg = "SR03" + err.Error()
				} else {
					SalesfinalRespRec.Status = "S"
					SalesfinalRespRec.Msg = "Fetched Successful"
					SalesfinalRespRec.ErrMsg = ""
					SalesfinalRespRec.SalesArr = append(SalesfinalRespRec.SalesArr, data...)
				}
			}
		}
		data, err := json.Marshal(SalesfinalRespRec)
		if err != nil {
			fmt.Fprintf(w, "Error taking data "+err.Error())
		} else {
			fmt.Fprintln(w, string(data))
		}
		log.Println("SalesReportAPI(-)")
	}
}

func GetSalesReport(DaterangeRec SalesReportReqStruct) ([]SalesReportStruct, error) {
	var SalesDetailsRec SalesReportStruct
	var SalesArrRec []SalesReportStruct
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("database connection : ", err)
		return SalesArrRec, err
	} else {
		defer db.Close()
		sqlString := `SELECT BM.BILL_NO, BM.BILL_DATE, MM.MEDICINE_NAME, BD.QUANTITY, BD.AMOUNT
		FROM MEDAPP_BILL_MASTER BM
		JOIN MEDAPP_BILL_DETAILS BD ON BM.BILL_NO = BD.BILL_NO
		JOIN MEDAPP_MEDICINE_MASTER MM ON BD.MEDICINE_MASTER_ID = MM.MEDICINE_MASTER_ID
		WHERE BM.BILL_DATE BETWEEN ? AND ?;
		`
		rows, err := db.Query(sqlString, DaterangeRec.From_date, DaterangeRec.To_date)
		if err != nil {
			log.Println("Query execution api : ", err)
			return SalesArrRec, err
		} else {
			for rows.Next() {
				log.Println("hi")
				err := rows.Scan(&SalesDetailsRec.Bill_No, &SalesDetailsRec.Bill_date, &SalesDetailsRec.Medicine_Name, &SalesDetailsRec.Quantity, &SalesDetailsRec.Amount)
				if err != nil {
					log.Println("Error while scanning")

				} else {
					SalesArrRec = append(SalesArrRec, SalesDetailsRec)
				}
			}
			// var SalesfinalRespRec SalesReportRespStruct
			// SalesfinalRespRec.SalesArr = SalesArrRec
			return SalesArrRec, err
		}
	}

}
