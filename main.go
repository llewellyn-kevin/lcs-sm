package main

func main() {
	ConnectDB("lcs-sm.db")
	defer db.Close()
	SetRoutes()
}
