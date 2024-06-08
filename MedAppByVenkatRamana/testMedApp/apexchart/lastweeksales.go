package apexchart

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type DetailStruct struct {
	Bill_date   string `json:"bill_date"`
	Daily_sales string `json:"daily_sales"`
}
type LastWeekSalesStruct struct {
	Lastweek []DetailStruct `json:"last_week"`
	Role     string         `json:"role"`
	ErrMsg   string         `json:"errmsg"`
	Status   string         `json:"status"`
	Msg      string         `json:"msg"`
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

func LastWeekSalesAPI(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "POST" {
		log.Println("LastWeekSaleAPI(+)")
		var resp LastWeekSalesStruct

		WeeklySales, err := getweeksales()
		if err != nil {
			resp.ErrMsg = "DB01" + err.Error()
			resp.Status = "E"
		} else {
			resp.Lastweek = WeeklySales
			resp.Status = "S"
			resp.Msg = "Fetched succesfull"
		}

		Data, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(w, "Error taking data "+err.Error())
		} else {
			fmt.Fprintln(w, string(Data))
		}
	}
	log.Println("LastWeekSaleAPI(-)")
}

func getweeksales() ([]DetailStruct, error) {
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
		BILL_DATE >= CURDATE() - INTERVAL 7 DAY 
		AND BILL_DATE < CURDATE() 
	GROUP BY 
		BILL_DATE;`
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
