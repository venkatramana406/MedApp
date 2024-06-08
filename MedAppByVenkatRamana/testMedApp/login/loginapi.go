package login

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type UserStruct struct {
	Login_id         int    `json:"login_id"`
	Role             string `json:"role"`
	User_Id          string `json:"user_id"`
	Login_History_id int    `json:"login_history_id"`
}
type LoginReqStruct struct {
	User_Id  string `json:"user_id"`
	Password string `json:"password"`
}
type LoginRespStruct struct {
	User_Details UserStruct `json:"userdetails"`
	ErrMsg       string     `json:"errmsg"`
	Status       string     `json:"status"`
	Msg          string     `json:"msg"`
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

func LoginAPI(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Toke,Authorization")
	if r.Method == "POST" {
		log.Println("LoginApi(+)")
		var LoginDetailRec LoginReqStruct
		var lfinalrespRec LoginRespStruct
		lfinalrespRec.Status = "S"
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error :", err)
			lfinalrespRec.Status = "E"
			lfinalrespRec.ErrMsg = "LOG01" + err.Error()
		} else {
			err := json.Unmarshal(body, &LoginDetailRec)
			if err != nil {
				log.Println("Error :", err)
				lfinalrespRec.Status = "E"
				lfinalrespRec.ErrMsg = "LOG02" + err.Error()
			} else {
				data, err := LoginValidation(LoginDetailRec)
				if err != nil {
					lfinalrespRec.Status = "E"
					lfinalrespRec.ErrMsg = "LOG03" + err.Error()
				} else {
					lfinalrespRec.User_Details = data
					lfinalrespRec.Status = "S"
					lfinalrespRec.Msg = "Login Successful"

				}
			}
		}
		data, err := json.Marshal(lfinalrespRec)
		if err != nil {
			fmt.Fprintf(w, "Error taking data "+err.Error())
		} else {
			// log.Println("marshal-", lfinalrespRec.User_Details.Login_History_id)
			fmt.Fprintln(w, string(data))
		}
		log.Println("LoginApi(-)")
	}
}
func LoginValidation(LoginRec LoginReqStruct) (UserStruct, error) {
	var UserDetailRec UserStruct
	db, err := LocalDBConnect()
	if err != nil {
		log.Println("Error :", err)
		return UserDetailRec, err
	} else {
		defer db.Close()
		sqlString := `SELECT LOGIN_ID,USER_ID,ROLE
						FROM MEDAPP_LOGIN
						WHERE USER_ID=? AND PASSWORD=?;`
		rows, err := db.Query(sqlString, LoginRec.User_Id, LoginRec.Password)
		// log.Println("username : ", LoginRec.User_Id)
		// log.Println("password : ", LoginRec.Password)
		if err != nil {
			return UserDetailRec, err
		} else {
			for rows.Next() {
				err := rows.Scan(&UserDetailRec.Login_id, &UserDetailRec.User_Id, &UserDetailRec.Role)
				if err != nil {
					return UserDetailRec, err
				} else {
					value, err := LoginHistoryInsert(UserDetailRec)
					if err != nil {
						return UserDetailRec, err
					} else {
						UserDetailRec.Login_History_id = value
						return UserDetailRec, nil
					}
				}
			}
		}
	}
	return UserDetailRec, http.ErrNotSupported
}
func LoginHistoryInsert(userdetail UserStruct) (int, error) {
	var Login_History_id int
	db, err := LocalDBConnect()
	if err != nil {
		log.Println(" database connect :", err)
		return Login_History_id, err
	} else {
		defer db.Close()
		sqlString := `INSERT  INTO  MEDAPP_LOGIN_HISTORY
		(LOGIN_ID,LOGIN_DATE,LOGIN_TIME,CREATED_BY,
		CREATED_DATE,UPDATED_BY,UPDATED_DATE) 
		VALUES (?,NOW(),NOW(),?,NOW(),?,NOW());`
		_, err := db.Exec(sqlString, userdetail.Login_id, "venkatramana", "")
		// log.Println("username : ", LoginRec.User_Id)
		// log.Println("password : ", LoginRec.Password)
		if err != nil {
			log.Println(" Query Execution Error :", err)
			return Login_History_id, err
		} else {
			// log.Println(" LoginhistoryInsert Query Execution Successful :", err)
			sqlString := `SELECT LOGIN_HISTORY_ID
			FROM MEDAPP_LOGIN_HISTORY 
			WHERE LOGIN_HISTORY_ID=(SELECT COUNT(*)
			FROM MEDAPP_LOGIN_HISTORY );`
			rows, err := db.Query(sqlString)
			// log.Println("username : ", LoginRec.User_Id)
			// log.Println("password : ", LoginRec.Password)
			if err != nil {
				return Login_History_id, err
			} else {
				for rows.Next() {
					err := rows.Scan(&Login_History_id)
					// log.Println("loginhistoryid = ", Login_History_id)
					if err != nil {
						return Login_History_id, err
					} else {

						return Login_History_id, err
					}

				}
			}
			return Login_History_id, err
		}
	}

}
