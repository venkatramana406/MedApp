package apexchart

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type RespData struct {
	User_Id     string `json:"user_id"`
	Bill_Amount string `json:"bill_amount"`
}
type SalesManPerformance struct {
	ResultData []RespData `json:"result_data"`
	ErrMsg     string     `json:"errmsg"`
	Status     string     `json:"status"`
	Msg        string     `json:"msg"`
}

func SalesManAPI(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "POST" {
		log.Println("SalesManAPI(+)")
		var resp SalesManPerformance

		TotalSales, err := getsales()
		if err != nil {
			resp.ErrMsg = "DB01" + err.Error()
			resp.Status = "E"
		} else {
			resp.ResultData = TotalSales
			resp.Status = "S"
			resp.Msg = "Fetched succesfull"
		}

		Data, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(w, "Error taking data "+err.Error())
		} else {
			fmt.Fprintln(w, string(Data))
		}

		log.Println("SalesManAPI(-)")
	}
}
func getsales() ([]RespData, error) {
	var sales []RespData
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("database connection : ", err)
		return sales, err
	} else {
		defer db.Close()
		sqlString := `select sum(mbm.bill_amount), ml.user_id
		from medapp_bill_master mbm join medapp_login ml 
		on mbm.login_id =ml.login_id 
		group by ml.login_id ;`
		records, err := db.Query(sqlString)
		if err != nil {
			log.Println("Query execution api : ", err)
			return sales, err
		} else {
			var sale RespData
			for records.Next() {
				err := records.Scan(&sale.Bill_Amount, &sale.User_Id)
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
