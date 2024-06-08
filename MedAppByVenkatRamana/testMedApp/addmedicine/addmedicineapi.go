package addmedicine

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type AddNewMedicineReqStruct struct {
	Medicine_Name string `json:"medicine_name"`
	Brand         string `json:"brand"`
	User_id       string `json:"user_id"`
}
type AddNewMedicineRespStruct struct {
	Last_Medicine StockStruct `json:"last_medicine"`
	ErrMsg        string      `json:"errmsg"`
	Status        string      `json:"status"`
	Msg           string      `json:"msg"`
}
type StockStruct struct {
	Medicine_Name string `json:"medicine_name"`
	Brand         string `json:"brand"`
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

func AddMedicineAPI(w http.ResponseWriter, r *http.Request) {
	
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "PUT" {
		log.Println("AddMedicineAPI(+)")
		var MedicineRec AddNewMedicineReqStruct
		var finalrespRec AddNewMedicineRespStruct
		finalrespRec.Status = "S"
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error :", err)
			finalrespRec.Status = "E"
			finalrespRec.ErrMsg = "AD01" + err.Error()
		} else {
			err := json.Unmarshal(body, &MedicineRec)
			if err != nil {
				log.Println("Error :", err)
				finalrespRec.Status = "E"
				finalrespRec.ErrMsg = "AD02" + err.Error()
				
			} else {
				data, err := AddNewMedicine(MedicineRec)
				if err != nil {
					finalrespRec.Status = "E"
					finalrespRec.Msg = ""
					finalrespRec.ErrMsg = "AD03" + err.Error()
				} else {
					finalrespRec.Last_Medicine = data
					finalrespRec.Msg = "Medicine Added Successfully"
					finalrespRec.Status = "S"
				}
			}
		}
		data, err := json.Marshal(finalrespRec)
		if err != nil {
			fmt.Fprintf(w, "Error taking data "+err.Error())
		} else {
			fmt.Fprintln(w, string(data))
		}
		log.Println("AddMedicineAPI(-)")
	}
}

func AddNewMedicine(newmed AddNewMedicineReqStruct) (StockStruct, error) {
	var LastMed StockStruct
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("database connection : ", err)
		return LastMed, err
	} else {
		defer db.Close()
		sqlString := `INSERT  INTO  MEDAPP_MEDICINE_MASTER(MEDICINE_NAME,BRAND,
			CREATED_BY,CREATED_DATE,UPDATED_BY,UPDATED_DATE)
			SELECT  ?,?,?,NOW(),?,NOW()  
			WHERE NOT EXISTS(SELECT MEDICINE_NAME, BRAND
			FROM  MEDAPP_MEDICINE_MASTER  M
			WHERE  M.MEDICINE_NAME = ?  AND  M.BRAND = ? );
			`
		rows, err := db.Exec(sqlString, newmed.Medicine_Name, newmed.Brand, newmed.User_id, "ramana", newmed.Medicine_Name, newmed.Brand)
		if err != nil {
			log.Println("Query execution api : ", err)
			return LastMed, err
		} else {
			rowsAffected, err := rows.RowsAffected()
			if err != nil {
				log.Println("Error getting rows affected:", err)
				return LastMed, err
			}
			if rowsAffected == 0 {
				return LastMed, http.ErrContentLength
			} else {
				sqlString1 := `INSERT  INTO  MEDAPP_STOCK(MEDICINE_MASTER_ID,
					CREATED_BY,CREATED_DATE,UPDATED_BY,UPDATED_DATE)
					vALUES(  (SELECT MEDICINE_MASTER_ID FROM MEDAPP_MEDICINE_MASTER WHERE MEDICINE_NAME=?),?,NOW(),?,NOW())`
				_, err := db.Exec(sqlString1, newmed.Medicine_Name, newmed.User_id, "RAMANA")
				if err != nil {
					log.Println("Error fetch last medicine ", err.Error())
					return LastMed, err
				} else {
					return LastMed, err
				}

			}

		}
	}
}
