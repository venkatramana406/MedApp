package dashboard

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type DashboardReqStruct struct {
	User_Id string `json:"user_id"`
}
type ManagerDashStruct struct {
	Today_sales     float64 `json:"today_sales"`
	Inventory_value float64 `json:"inventory_value"`
}
type BillerDashStruct struct {
	Today_sales     float64 `json:"today_sales"`
	Yesterday_sales float64 `json:"yesterday_sales"`
}
type CommonMessageStruct struct {
	Message string `json:"message"`
}
type DashboardRespStruct struct {
	Manager []ManagerDashStruct `json:"manager_detail"`
	Biller  []BillerDashStruct  `json:"biller_detail"`
	Role    string              `json:"role"`
	ErrMsg  string              `json:"errmsg"`
	Status  string              `json:"status"`
	Msg     string              `json:"msg"`
}

var role string

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
func DashboardRoleAPI(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "POST" {
		log.Println("DashboardRoleAPI(+)")
		var DashRespRec DashboardRespStruct
		db, err := LocalDBConnect()
		if err != nil {
			log.Println("Error in db connection")
		} else {
			defer db.Close()
			sqlString := `SELECT L.ROLE  
		FROM MEDAPP_LOGIN L,MEDAPP_LOGIN_HISTORY LH
		WHERE L.LOGIN_ID=LH.LOGIN_ID ORDER BY LH.LOGIN_HISTORY_ID DESC LIMIT 1;`
			records, err := db.Query(sqlString)
			if err != nil {
				log.Println("ERROR IN FETCHING ROLE : ", err)
			} else {
				for records.Next() {
					err := records.Scan(&role)
					if err != nil {
						DashRespRec.Status = "E"
						DashRespRec.Msg = ""
						DashRespRec.ErrMsg = "DRA01" + err.Error()
					} else {
						DashRespRec.Msg = "Role Fetched successfull"
						DashRespRec.Status = "S"
						DashRespRec.Role = role
						// log.Println("Fetched role")
					}
				}
			}
		}
		respData, err := json.Marshal(DashRespRec)
		if err != nil {
			fmt.Fprintf(w, "Error taking data "+err.Error())

		} else {
			fmt.Fprintln(w, string(respData))

		}
		log.Println("DashboardRoleAPI(-)")
	}
}
func DashboardAPI(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "POST" {
		log.Println("DashboardAPI(+)")
		var UseridRec DashboardReqStruct
		var resp DashboardRespStruct
		db, err := LocalDBConnect()
		if err != nil {
			log.Println("error in db connection")
		} else {
			defer db.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("Error :", err)

			} else {
				err := json.Unmarshal(body, &UseridRec)
				if err != nil {
					log.Println("Error :", err)
				}
				ManagerRec, err := DashManager()
				if err != nil {
					resp.ErrMsg = "DB01" + err.Error()
					resp.Status = "E"
				} else {

					resp.Manager = append(resp.Manager, ManagerRec)
					resp.Status = "S"
					resp.Msg = "Calculated succesfull"
				}
				result, err := DashBiller(UseridRec.User_Id)
				if err != nil {
					resp.ErrMsg = "DB02" + err.Error()
					resp.Status = "E"

				} else {
					resp.Biller = append(resp.Biller, result)
					resp.Status = "S"
					resp.Msg = "Calculated succesfull"
				}
			}
			Data, err := json.Marshal(resp)
			if err != nil {
				fmt.Fprintf(w, "Error taking data "+err.Error())
			} else {
				fmt.Fprintln(w, string(Data))
			}
		}
		log.Println("DashboardAPI(-)")
	}
}

func DashManager() (ManagerDashStruct, error) {
	var ManagerDashRec ManagerDashStruct
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("database connection : ", err)
		return ManagerDashRec, err
	} else {
		defer db.Close()
		sqlString := `SELECT COALESCE(SUM(NET_PRICE), 0) AS TODAY_SALES 
		FROM MEDAPP_BILL_MASTER WHERE BILL_DATE= CURDATE();`
		records, err := db.Query(sqlString)
		if err != nil {
			log.Println("Query execution api : ", err)
			return ManagerDashRec, err
		} else {
			for records.Next() {
				err := records.Scan(&ManagerDashRec.Today_sales)
				if err != nil {
					log.Println(err.Error())
				}
			}
			sqlString := `SELECT COALESCE(SUM(QUANTITY*UNIT_PRICE), 0) AS INVENTORY_VALUE
				FROM MEDAPP_STOCK 
				WHERE MEDICINE_MASTER_ID IN (SELECT MEDICINE_MASTER_ID 
					FROM MEDAPP_MEDICINE_MASTER);`
			rows, err := db.Query(sqlString)
			if err != nil {
				log.Println("Error fetch INVENTORY VALUE ", err.Error())
				return ManagerDashRec, err
			} else {
				for rows.Next() {
					err := rows.Scan(&ManagerDashRec.Inventory_value)
					if err != nil {
						log.Println("Error while scaning today sales")
					}
				}
			}
		}
		return ManagerDashRec, err
	}
}
func DashBiller(userid string) (BillerDashStruct, error) {
	var BillerDashRec BillerDashStruct
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("database connection : ", err)
		return BillerDashRec, err
	} else {
		defer db.Close()
		sqlString := `SELECT NVL(SUM(NET_PRICE), 0) AS YESTERDAY_SALES
                 FROM MEDAPP_BILL_MASTER
                 WHERE BILL_DATE= (SELECT DATE_ADD(CURDATE(),INTERVAL -1 DAY)) AND LOGIN_ID=(SELECT LOGIN_ID FROM MEDAPP_LOGIN WHERE USER_ID=?)`
		records, err := db.Query(sqlString, userid)
		if err != nil {
			log.Println("Query execution api : ", err)
			return BillerDashRec, err
		} else {
			for records.Next() {
				err := records.Scan(&BillerDashRec.Yesterday_sales)
				if err != nil {
					log.Println("Error while scaning today sales")
				}
			}
			sqlString = `SELECT NVL(SUM(NET_PRICE), 0) AS YESTERDAY_SALES
                 FROM MEDAPP_BILL_MASTER
                 WHERE BILL_DATE= CURDATE() AND LOGIN_ID=(SELECT LOGIN_ID FROM MEDAPP_LOGIN WHERE USER_ID=?)`
			rows, err := db.Query(sqlString, userid)
			if err != nil {
				log.Println("Error fetch INVENTORY VALUE ", err.Error())
				return BillerDashRec, err
			} else {
				for rows.Next() {
					err := rows.Scan(&BillerDashRec.Today_sales)
					if err != nil {
						log.Println("Error while scaning today sales")
					}
				}
			}
		}
		return BillerDashRec, err
	}
}
