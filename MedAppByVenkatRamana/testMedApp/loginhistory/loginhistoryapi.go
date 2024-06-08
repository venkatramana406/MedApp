package loginhistory

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type LoginHistoryReqStruct struct {
	Login_date  string `json:"login_date"`
	Login_time  string `json:"login_time"`
	Logout_date string `json:"logout_date"`
	Logout_time string `json:"logout_time"`
	User_Id     string `json:"user_id"`
}
type LoginHistoryRespStruct struct {
	LoginHistoryArr []LoginHistoryReqStruct `json:"loginhistoryarr"`
	ErrMsg          string                  `json:"errmsg"`
	Status          string                  `json:"status"`
	Msg             string                  `json:"msg"`
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

func LoginHistoryAPI(w http.ResponseWriter, r *http.Request) {
	
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "GET" {
		log.Println("LoginHistoryAPI(+)")
		var historyRespRec LoginHistoryRespStruct
		historyRespRec.Status = "S"
		data, err := FetchLoginHistory()
		if err != nil {
			historyRespRec.Status = "E"
			historyRespRec.Msg = ""
			historyRespRec.ErrMsg = "SR03" + err.Error()
		} else {
			historyRespRec.LoginHistoryArr = append(historyRespRec.LoginHistoryArr, data...)
			historyRespRec.Status = "S"
			historyRespRec.Msg = "Fetched Successful"
			historyRespRec.ErrMsg = ""
		}
		response, err := json.Marshal(historyRespRec)
		if err != nil {
			fmt.Fprintf(w, "Error taking data "+err.Error())
		} else {
			fmt.Fprintln(w, string(response))
		}
		log.Println("LoginHistoryAPI(-)")
	}
}

func FetchLoginHistory() ([]LoginHistoryReqStruct, error) {
	var HistoryDetailsRec LoginHistoryReqStruct
	var HistoryArrRec []LoginHistoryReqStruct
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("database connection : ", err)
		return HistoryArrRec, err
	} else {
		defer db.Close()
		sqlString := `SELECT
			L.USER_ID,
			IFNULL(LH.LOGIN_DATE,''),
			IFNULL(LH.LOGIN_TIME,''),
			IFNULL(LH.LOGOUT_DATE, 'WORKING'),
			IFNULL(LH.LOGOUT_TIME, 'WORKING')
	    FROM
		MEDAPP_LOGIN_HISTORY LH,MEDAPP_LOGIN L
		WHERE L.LOGIN_ID=LH.LOGIN_ID
		ORDER BY LH.LOGIN_TIME DESC;
		`
		rows, err := db.Query(sqlString)
		if err != nil {
			log.Println("Query execution api : ", err)
			return HistoryArrRec, err
		} else {
			for rows.Next() {
				// log.Println("hi")
				err := rows.Scan(&HistoryDetailsRec.User_Id, &HistoryDetailsRec.Login_date,
					&HistoryDetailsRec.Login_time, &HistoryDetailsRec.Logout_date,
					&HistoryDetailsRec.Logout_time)
				if err != nil {
					log.Println(err)
				} else {
					HistoryArrRec = append(HistoryArrRec, HistoryDetailsRec)
				}
			}
			// var SalesfinalRespRec SalesReportRespStruct
			// SalesfinalRespRec.SalesArr = SalesArrRec
			return HistoryArrRec, err
		}
	}

}
