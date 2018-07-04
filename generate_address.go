package main

import (
	"db"
	"fmt"
	"os"
	"os/exec"
	"shared"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

func exe_cmd(cmd string) string {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	restult := fmt.Sprintf("%s", out)
	restult = strings.Replace(restult, "\n", "", 1)
	return restult

}

type TomlConfig struct {
	Walletclient  string
	Walletdir     string
	Walletconf    string
	Walletcurreny string
	Walletnetwork string
	Args          string
}

var config TomlConfig

func main() {
	file := os.Args[1]
	size := os.Args[2]

	if _, err := toml.DecodeFile(file, &config); err != nil {
		fmt.Println(err)
		return
	}
	i, _ := strconv.Atoi(size)
	for a := 0; a < i; a++ {

		address := generateAddress(config.Walletclient, config.Walletdir, config.Walletconf, "getnewaddress")
		saveAddress(address, config.Walletcurreny)
		fmt.Println(address, "Saved")
	}

}
func generateAddress(currency_client string, curreny_dir string, currency_conf string, args string) string {
	commands := currency_client + "  -datadir=" + curreny_dir + " -conf=" + currency_conf + " " + args
	newtxs := exe_cmd(commands)
	return newtxs
}
func saveAddress(address string, currency string) {
	bucket := db.GetDbConnection(shared.BUCKET)

	bucket.Insert(config.Walletcurreny+address,
		shared.Address{
			Currency: config.Walletcurreny,
			Address:  address,
			Inuse:    false,
			Network:  config.Walletnetwork,
		}, 0)
	bucket.Close()

}
