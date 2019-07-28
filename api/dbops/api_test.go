package dbops

import(
	"testing"
	)

func clearTables()  {

	dbConn.Exec("TRUNCATE users")
	dbConn.Exec("TRUNCATE video_info")
	dbConn.Exec("TRUNCATE comments")
	dbConn.Exec("TRUNCATE sessions")
}

func TestMain(m *testing.M)  {

	clearTables()
	m.Run()
	clearTables()
}

func TestUserWorkFlow(t *testing.T) {
	t.Run("add",testAddUserCredential)
	t.Run("get",testGetUserCredential)
	t.Run("delete",testDeleteUser)
	t.Run("reGet",testRegetUser)
}

func testAddUserCredential(t *testing.T) {
	err := AddUserCredential("rex","123")
	if err!=nil {
		t.Errorf("Error of Add User : %v", err)
	}
}

func testGetUserCredential(t *testing.T) {
	pwd,err := GetUserCredential("rex")
	if err!=nil {
		t.Errorf("Error of Get User : %v", err)
	}
	t.Logf("User's pwd: %v",pwd)
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("rex", "123")
	if err!=nil {
		t.Errorf("Error of Delete User : %v", err)
	}
}

func testRegetUser(t *testing.T)  {
	pwd,err := GetUserCredential("rex")
	if err!=nil {
		t.Errorf("Error of reGet User : %v", err)
	}
	if pwd!= "" {
		t.Error("Error of reGet User : reGet pwd is not null")
	}
}
