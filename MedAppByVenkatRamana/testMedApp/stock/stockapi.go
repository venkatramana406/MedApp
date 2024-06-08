package stock

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type StockStruct struct {
	Medicine_Name string `json:"medicine_name"`
	Brand         string `json:"brand"`
}
type StockviewRespStruct struct {
	StockArr []StockStruct `json:"stockarr"`
	ErrMsg   string        `json:"errmsg"`
	Status   string        `json:"status"`
	Msg      string        `json:"msg"`
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

func StockAPI(w http.ResponseWriter, r *http.Request) {
	
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "GET" {
		log.Println("StockAPI(+)")
		var StockiewRespRec StockviewRespStruct
		StockiewRespRec.Status = "S"
		data, err := GetStock()
		if err != nil {
			StockiewRespRec.Status = "E"
			StockiewRespRec.Msg = ""
			StockiewRespRec.ErrMsg = "S01" + err.Error()
		} else {
			StockiewRespRec.StockArr = append(StockiewRespRec.StockArr, data...)
			StockiewRespRec.Status = "S"
			StockiewRespRec.Msg = "Stock detail received"
			StockiewRespRec.ErrMsg = ""
		}
		response, err := json.Marshal(StockiewRespRec)
		if err != nil {
			fmt.Fprintf(w, "Error taking data "+err.Error())
		} else {
			fmt.Fprintln(w, string(response))
		}
		log.Println("StockAPI(-)")
	}
}

func GetStock() ([]StockStruct, error) {
	var GStocksArr []StockStruct
	var GStockDetailsRec StockStruct
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("database connection : ", err)
		return GStocksArr, err
	} else {
		defer db.Close()
		sqlString := `SELECT  M.MEDICINE_NAME, M.BRAND
		FROM   MEDAPP_MEDICINE_MASTER M;`
		rows, err := db.Query(sqlString)
		if err != nil {
			log.Println("Query execution api : ", err)
			return GStocksArr, err
		} else {
			for rows.Next() {
				// log.Println("hi")
				err := rows.Scan(&GStockDetailsRec.Medicine_Name, &GStockDetailsRec.Brand)
				if err != nil {
					log.Println("Error while scanning")
					log.Println(err)
				} else {
					GStocksArr = append(GStocksArr, GStockDetailsRec)
				}
			}

			return GStocksArr, err
		}
	}

}
