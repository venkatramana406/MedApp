package apexchart

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type LastMonthSalesStruct struct {
	LastMonth []DetailStruct `json:"last_month"`
	Role      string         `json:"role"`
	ErrMsg    string         `json:"errmsg"`
	Status    string         `json:"status"`
	Msg       string         `json:"msg"`
}

func LastMonthSalesAPI(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "POST" {
		log.Println("LastMonthSalesAPI(+)")
		var resp LastMonthSalesStruct

		MonthlySales, err := getMonthsales()
		if err != nil {
			resp.ErrMsg = "DB01" + err.Error()
			resp.Status = "E"
		} else {
			resp.LastMonth = MonthlySales
			resp.Status = "S"
			resp.Msg = "Fetched succesfull"
		}

		Data, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(w, "Error taking data "+err.Error())
		} else {
			fmt.Fprintln(w, string(Data))
		}
		log.Println("LastMonthSalesAPI(-)")
	}
}

func getMonthsales() ([]DetailStruct, error) {
	var sales []DetailStruct
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("database connection : ", err)
		return sales, err
	} else {
		defer db.Close()
		sqlString := `SELECT 
		DATE_FORMAT(BILL_DATE, '%Y-%m-%d') AS BILL_DATE,
		IFNULL(SUM(NET_PRICE), 0) AS DAILY_SALES 
	FROM 
		MEDAPP_BILL_MASTER 
	WHERE 
		BILL_DATE >= CURDATE() - INTERVAL 1 MONTH 
		AND BILL_DATE < CURDATE() 
	GROUP BY 
		DATE_FORMAT(BILL_DATE, '%Y-%m-%d');`
		records, err := db.Query(sqlString)
		if err != nil {
			log.Println("Query execution api : ", err)
			return sales, err
		} else {

			var sale DetailStruct
			for records.Next() {
				err := records.Scan(&sale.Bill_date, &sale.Daily_sales)
				if err != nil {
					log.Println(err.Error())
				} else {
					sales = append(sales, sale)
				}
			}

		}
		return sales, err
	}
}
