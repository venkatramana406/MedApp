package updatestock

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type UpdateStockReqStruct struct {
	Medicine_Name string  `json:"medicine_name"`
	Brand         string  `json:"brand"`
	Quantity      int     `json:"quantity"`
	Unit_price    float64 `json:"unit_price"`
	User_id       string  `json:"user_id"`
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

func UpdateStockAPI(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "PUT" {
		log.Println("UpdateStockAPI(+)")
		var UpdateMedRec UpdateStockReqStruct
		var finalrespRec CommonRespStruct

		finalrespRec.Status = "S"
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error :", err)
			log.Println("unable to read body")
			finalrespRec.Status = "E"
			finalrespRec.ErrMsg = "US01" + err.Error()
		} else {
			err := json.Unmarshal(body, &UpdateMedRec)
			if err != nil {
				log.Println("Error :", err)
				finalrespRec.Status = "E"
				finalrespRec.ErrMsg = "US02" + err.Error()
				log.Println("unable to unmarshall")
			} else {
				num, err := StockUpdate(UpdateMedRec)
				if err != nil {
					finalrespRec.Status = "E"
					finalrespRec.Msg = ""
					finalrespRec.ErrMsg = "US03" + err.Error()
				} else {
					if num == 0 {
						finalrespRec.Status = "E"
						finalrespRec.Msg = ""
						finalrespRec.ErrMsg = "US04" + " No rows inserted"
					} else {
						finalrespRec.Msg = "Medicine Updated Successfully"
						finalrespRec.Status = "S"
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
		log.Println("UpdateStockAPI(-)")
	}
}

func StockUpdate(UpdateMedRec UpdateStockReqStruct) (int, error) {
	var InsValue int
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("database connection : ", err)
		return InsValue, err
	} else {
		defer db.Close()
		sqlString := `UPDATE  MEDAPP_STOCK SET QUANTITY=IFNULL(QUANTITY,0)+?,UNIT_PRICE =?,
		UPDATED_BY =?,UPDATED_DATE =NOW()
		WHERE MEDICINE_MASTER_ID IN(SELECT MEDICINE_MASTER_ID  FROM MEDAPP_MEDICINE_MASTER M
		WHERE  M.MEDICINE_NAME=? AND M.BRAND=?);`
		rows, err := db.Exec(sqlString, UpdateMedRec.Quantity, UpdateMedRec.Unit_price, UpdateMedRec.User_id, UpdateMedRec.Medicine_Name, UpdateMedRec.Brand)
		if err != nil {
			log.Println("Query execution api : ", err)
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
