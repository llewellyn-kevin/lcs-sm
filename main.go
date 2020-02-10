package main

func main() {
	ConnectDB("lcs_sm.db")
	defer db.Close()
	SetRoutes()
}