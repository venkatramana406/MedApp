package logout

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type LogoutReqStruct struct {
	Login_History_Id int    `json:"login_history_id"`
	User_id          string `json:"user_id"`
}
type CommonRespStruct struct {
	ErrMsg string `json:"errmsg"`
	Status string `json:"status"`
	Msg    string `json:"msg"`
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

func LogoutAPI(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "PUT" {
		log.Println("LogoutApi(+)")
		var LogoutReqRec LogoutReqStruct
		var finalrespRec CommonRespStruct
		finalrespRec.Status = "S"
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error :", err)
			finalrespRec.Status = "E"
			finalrespRec.ErrMsg = "LO01" + err.Error()
		} else {
			err := json.Unmarshal(body, &LogoutReqRec)
			if err != nil {
				log.Println("Error :", err)
				finalrespRec.Status = "E"
				finalrespRec.ErrMsg = "LO02" + err.Error()
			} else {
				rows, err := DoLogout(LogoutReqRec)
				if err != nil {
					finalrespRec.Status = "E"
					finalrespRec.ErrMsg = "LO03" + err.Error()
				} else {
					if rows == 0 {
						finalrespRec.Status = "E"
						finalrespRec.Msg = "No Rows Updated"
					} else {
						finalrespRec.Status = "S"
						finalrespRec.Msg = "Logout Successful"
					}

				}
			}
		}
		data, err := json.Marshal(finalrespRec)
		if err != nil {
			fmt.Fprintf(w, "Error taking data "+err.Error())
		} else {
			fmt.Fprintln(w, string(data))
		}
		log.Println("LogoutApi(-)")
	}
}
func DoLogout(logoutDetails LogoutReqStruct) (int, error) {
	var InsValue int
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("Error :", err)
		return InsValue, err
	} else {
		defer db.Close()
		sqlString := `UPDATE MEDAPP_LOGIN_HISTORY
		SET LOGOUT_DATE =NOW(),LOGOUT_TIME=NOW(),UPDATED_BY=?,
		UPDATED_DATE=NOW()  WHERE LOGIN_HISTORY_ID=? ;`
		rows, err := db.Exec(sqlString, "venkatramana", logoutDetails.Login_History_Id)
		if err != nil {
			return InsValue, err
		} else {

			rowsAffected, err := rows.RowsAffected()
			if err != nil {
				log.Println("Error getting rows affected:", err)
				return InsValue, err
			}
			if rowsAffected == 0 {
				return 0, err
			} else {
				return 1, nil
			}
		}
	}
}
