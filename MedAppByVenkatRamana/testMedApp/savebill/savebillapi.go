package savebill

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type BillDetailReqStruct struct {
	BillDetailsArr []BillDetailStruct `json:"bill_details"`
}
type BillDetailStruct struct {
	BillNo       string  `json:"bill_no"`
	MedicineName string  `json:"medicine_name"`
	Quantity     int     `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	Amount       float64 `json:"amount"`
	UserId       string  `json:"user_id"`
}

type BillMasterReqStruct struct {
	BillNo     string  `json:"bill_no"`
	BillDate   string  `json:"billdate"`
	BillAmount int     `json:"bill_amount"`
	BillGst    float64 `json:"gst"`
	NetPrice   float64 `json:"net_price"`
	UserId     string  `json:"user_id"`
}
type BillSaveRespStruct struct {
	Bill_No int    `json:"bill_no"`
	ErrMsg  string `json:"errmsg"`
	Status  string `json:"status"`
	Msg     string `json:"msg"`
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

func SaveBillAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("SaveBillAPI(+)")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "POST" {
		var BillDetailsRec BillDetailReqStruct
		var SaveBillRespRec BillSaveRespStruct
		SaveBillRespRec.Status = "S"
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error :", err)
			log.Println("unable to read body")
			SaveBillRespRec.Status = "E"
			SaveBillRespRec.ErrMsg = "AD01" + err.Error()
		} else {
			err := json.Unmarshal(body, &BillDetailsRec)
			if err != nil {
				log.Println("Error :", err)
				SaveBillRespRec.Status = "E"
				SaveBillRespRec.ErrMsg = "AD02" + err.Error()
				log.Println("unable to unmarshall")
			} else {
				_, err := InsertBillDetails(BillDetailsRec)
				if err != nil {
					SaveBillRespRec.Status = "E"
					SaveBillRespRec.Msg = ""
					SaveBillRespRec.ErrMsg = "AD03" + err.Error()
				} else {
					SaveBillRespRec.Msg = "bill Added Successfully"
					SaveBillRespRec.Status = "S"
				}
			}
		}
		data, err := json.Marshal(SaveBillRespRec)
		if err != nil {
			fmt.Fprintf(w, "Error taking data "+err.Error())
		} else {
			fmt.Fprintln(w, string(data))
		}
		log.Println("SaveBillAPI(-)")
	}
}

func SaveBillMasterAPI(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "POST" {
		log.Println("SaveBillMasterAPI(+)")
		var billmas BillMasterReqStruct
		var SaveBillRespRec BillSaveRespStruct
		SaveBillRespRec.Status = "S"
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error :", err)
			SaveBillRespRec.Status = "E"
			SaveBillRespRec.ErrMsg = "AD01" + err.Error()
		} else {
			err := json.Unmarshal(body, &billmas)
			if err != nil {
				log.Println("Error :", err)
				SaveBillRespRec.Status = "E"
				SaveBillRespRec.ErrMsg = "AD02" + err.Error()
			} else {

				_, err := InsertMaster(billmas)
				if err != nil {
					SaveBillRespRec.Status = "E"
					SaveBillRespRec.Msg = ""
					SaveBillRespRec.ErrMsg = "AD03" + err.Error()
				} else {
					SaveBillRespRec.Msg = "bill Added Successfully"
					SaveBillRespRec.Status = "S"
				}
			}
		}
		data, err := json.Marshal(SaveBillRespRec)
		if err != nil {
			fmt.Fprintf(w, "Error taking data "+err.Error())
		} else {
			fmt.Fprintln(w, string(data))
		}
		log.Println("SaveBillMasterAPI(-)")
	}
}
func InsertBillDetails(saveDetail BillDetailReqStruct) (int, error) {
	var value int
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("database connection : ", err)
		return value, err
	} else {
		defer db.Close()
		sqlString := `INSERT INTO MEDAPP_BILL_DETAILS (BILL_NO, MEDICINE_MASTER_ID, QUANTITY, UNIT_PRICE, AMOUNT, CREATED_BY, CREATED_DATE, UPDATED_BY, UPDATED_DATE)
		VALUES (?, (SELECT MEDICINE_MASTER_ID FROM MEDAPP_MEDICINE_MASTER M WHERE M.MEDICINE_NAME = ?), ?, ?, ?, ?, NOW(), ?, NOW())`
		for _, bill := range saveDetail.BillDetailsArr {
			result, err := db.Exec(sqlString, bill.BillNo, bill.MedicineName, bill.Quantity, bill.UnitPrice, bill.Amount, bill.UserId, bill.UserId)
			if err != nil {
				log.Println("SB02", err)
				return 0, err
			}
			rowsAffected, err := result.RowsAffected()
			log.Println(rowsAffected)
			if err != nil {
				log.Println("SB03", err)
				return 0, err
			}
			if rowsAffected == 0 {
				return 0, nil
			}

		}
		return 1, nil

	}
}

func InsertMaster(lbill BillMasterReqStruct) (int, error) {
	var value int
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("database connection : ", err)
		return value, err
	} else {
		defer db.Close()
		log.Println(lbill)
		sqlString := `INSERT INTO MEDAPP_BILL_MASTER (BILL_NO, BILL_DATE, BILL_AMOUNT, BILL_GST, NET_PRICE, LOGIN_ID, CREATED_BY, CREATED_DATE, UPDATED_BY, UPDATED_DATE) 
		VALUES (?, NOW(), ?, ?, ?, (SELECT LOGIN_ID FROM MEDAPP_LOGIN WHERE USER_ID=?), ?, NOW(), ?, NOW())`
		result, err := db.Exec(sqlString, lbill.BillNo, lbill.BillAmount, lbill.BillGst, lbill.NetPrice, lbill.UserId, lbill.UserId, lbill.UserId)
		if err != nil {
			log.Println("error in api", err)
			return 0, err
		}
		rowsAffected, err1 := result.RowsAffected()
		if err1 != nil {
			log.Println(err1)
		}
		log.Println(rowsAffected)

		return 1, nil

	}
}
