package adduser

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type AddNewUserReqStruct struct {
	User_id         string `json:"user_id"`
	Password        string `json:"password"`
	Role            string `json:"role"`
	Creater_User_id string `json:"cuser_id"`
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

func AddUserAPI(w http.ResponseWriter, r *http.Request) {
	
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "PUT" {
		log.Println("AddUserAPI(+)")
		var NewUserRec AddNewUserReqStruct
		var finalrespRec CommonRespStruct
		finalrespRec.Status = "S"
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error :", err)
			finalrespRec.Status = "E"
			finalrespRec.ErrMsg = "AD01" + err.Error()
		} else {
			err := json.Unmarshal(body, &NewUserRec)
			if err != nil {
				log.Println("Error :", err)
				finalrespRec.Status = "E"
				finalrespRec.ErrMsg = "AD02" + err.Error()
			} else {
				num, err := InsertNewUser(NewUserRec)
				if err != nil {
					finalrespRec.Status = "E"
					finalrespRec.Msg = ""
					finalrespRec.ErrMsg = "AD03" + err.Error()
				} else {
					if num == 0 {
						finalrespRec.Status = "E"
						finalrespRec.Msg = ""
						finalrespRec.ErrMsg = "AD04" + " No rows inserted"
					} else {
						finalrespRec.Msg = "User Added Successfully"
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
		log.Println("AddUserAPI(-)")
	}
}

func InsertNewUser(newuser AddNewUserReqStruct) (int, error) {
	var InsValue int
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("database connection : ", err)
		return InsValue, err
	} else {
		defer db.Close()
		sqlString := `INSERT INTO MEDAPP_LOGIN
		(USER_ID,PASSWORD,ROLE,CREATED_BY,CREATED_DATE,UPDATED_BY,
		UPDATED_DATE)
		SELECT ?,?,?,?,NOW(),?,NOW()
		WHERE NOT EXISTS (SELECT LOGIN_ID
		FROM MEDAPP_LOGIN
		WHERE USER_ID=?);`

		rows, err := db.Exec(sqlString, newuser.User_id, newuser.Password, newuser.Role, newuser.Creater_User_id, "", newuser.User_id)
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
