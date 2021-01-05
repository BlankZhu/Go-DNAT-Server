package main

import (
	"BlankZhu/Go-DNAT-Server/pkg/config"
	memstorage "BlankZhu/Go-DNAT-Server/pkg/db/memory"
	mysqlstorage "BlankZhu/Go-DNAT-Server/pkg/db/mysql"
	"BlankZhu/Go-DNAT-Server/pkg/server"
	"flag"
	"strconv"
)

func main() {
	cfgPathPtr := flag.String("c", "./config.yaml", "path to YAML config file")
	flag.Parse()
	err := config.LoadYAMLConfigFromFile(*cfgPathPtr)
	if err != nil {
		panic(err)
	}
	conf := config.Get()

	// initialize mysql database connection
	mysqlstorage.Init(&conf.MySQL)
	mysql := mysqlstorage.Get()
	err = mysql.Open()
	if err != nil {
		panic(err)
	}
	defer mysql.Close()

	// load record from database to memory
	list, err := mysql.List()
	if err != nil {
		panic(err)
	}
	mem := memstorage.Get()
	for _, v := range list {
		mem.Add(v)
	}

	r := server.SetupRouter()
	r.Run(":" + strconv.Itoa(int(conf.Port)))
}
