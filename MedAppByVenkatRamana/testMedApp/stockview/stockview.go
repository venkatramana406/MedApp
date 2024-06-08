package stockview

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type StockviewReqStruct struct {
	Medicine_Name string  `json:"medicine_name"`
	Brand         string  `json:"brand"`
	Quantity      int     `json:"quantity"`
	Unit_price    float64 `json:"unit_price"`
}
type StockviewRespStruct struct {
	StockviewArr []StockviewReqStruct `json:"stockviewarr"`
	ErrMsg       string               `json:"errmsg"`
	Status       string               `json:"status"`
	Msg          string               `json:"msg"`
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

func StockViewAPI(w http.ResponseWriter, r *http.Request) {
	
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "GET" {
		log.Println("StockViewAPI(+)")
		var StockiewRespRec StockviewRespStruct
		StockiewRespRec.Status = "S"
		data, err := GetStockviewRecords()
		if err != nil {
			StockiewRespRec.Status = "E"
			StockiewRespRec.Msg = ""
			StockiewRespRec.ErrMsg = "SV01" + err.Error()
		} else {
			StockiewRespRec.StockviewArr = append(StockiewRespRec.StockviewArr, data...)
			StockiewRespRec.Status = "S"
			StockiewRespRec.Msg = "Fetched Successful"
			StockiewRespRec.ErrMsg = ""
		}
		response, err := json.Marshal(StockiewRespRec)
		if err != nil {
			fmt.Fprintf(w, "Error taking data "+err.Error())
		} else {
			fmt.Fprintln(w, string(response))
		}
		log.Println("StockViewAPI(-)")
	}
}

func GetStockviewRecords() ([]StockviewReqStruct, error) {
	var StockDetailsRec StockviewReqStruct
	var StocksArr []StockviewReqStruct
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("database connection : ", err)
		return StocksArr, err
	} else {
		defer db.Close()
		sqlString := `SELECT  M.MEDICINE_NAME, M.BRAND, NVL( S.QUANTITY,0),NVL(  S.UNIT_PRICE,0)
		FROM   MEDAPP_STOCK S,   MEDAPP_MEDICINE_MASTER M
		WHERE   S.MEDICINE_MASTER_ID  =   M.MEDICINE_MASTER_ID;`
		rows, err := db.Query(sqlString)
		if err != nil {
			log.Println("Query execution api : ", err)
			return StocksArr, err
		} else {
			for rows.Next() {
				// log.Println("hi")
				err := rows.Scan(&StockDetailsRec.Medicine_Name, &StockDetailsRec.Brand, &StockDetailsRec.Quantity, &StockDetailsRec.Unit_price)
				if err != nil {
					log.Println("Error while scanning")
					log.Println(err)
				} else {
					StocksArr = append(StocksArr, StockDetailsRec)
				}
			}
			// var SalesfinalRespRec SalesReportRespStruct
			// SalesfinalRespRec.SalesArr = SalesArrRec
			return StocksArr, err
		}
	}

}
