package main

import (
	"flag"
	xylog "guanghuan.com/xiaoyao/common/log"
	crypto "guanghuan.com/xiaoyao/superbman_server/crypto"
	"io/ioutil"
)

type Config struct {
	InFile  string
	OutFile string
	InStr   string
	Encrypt bool
	Debug   bool
}

var DefConfig = Config{
	InFile:  "",
	OutFile: "",
	InStr:   "hello world",
	Encrypt: true,
}

func process_cmd() {
	flag.StringVar(&DefConfig.InFile, "if", DefConfig.InFile, "input file")
	flag.StringVar(&DefConfig.OutFile, "of", DefConfig.OutFile, "output file")
	flag.StringVar(&DefConfig.InStr, "i", DefConfig.InStr, "input string")

	flag.BoolVar(&DefConfig.Encrypt, "enc", DefConfig.Encrypt, "encrypt or decrypt")
	flag.BoolVar(&DefConfig.Debug, "debug", DefConfig.Debug, "enable debug log")
	flag.Parse()
}

func print_config() {
	xylog.Info(`Config:
	input  file   = %s
	output file   = %s
	input  string = %s
	encryption?   = %t
	------------------------
	`, DefConfig.InFile,
		DefConfig.OutFile,
		DefConfig.InStr,
		DefConfig.Encrypt)

}

func apply_config() bool {
	//if DefConfig.InFile != "" {
	//	if DefConfig.OutFile == "" {
	//		DefConfig.OutFile = DefConfig.InFile + ".o"
	//	}
	//}
	//	xylog.EnableDebug(DefConfig.Debug)
	return true
}

func crypt(in []byte, is_encrypt bool) (out []byte, err error) {
	if is_encrypt {
		out, err = crypto.Encrypt(in)
	} else {
		out, err = crypto.Decrypt(in)
	}
	return
}

func main() {
	process_cmd()
	print_config()
	if !apply_config() {
		return
	}

	var (
		in  []byte
		out []byte
		err error
	)

	if DefConfig.InFile != "" {
		in, err = ioutil.ReadFile(DefConfig.InFile)
		if err != nil {
			xylog.ErrorNoId("Error reading file (%s) : %s", DefConfig.InFile, err.Error())
			return
		} else {
			xylog.InfoNoId("Reading file: %s", DefConfig.InFile)
		}
	} else {
		in = []byte(DefConfig.InStr)
		xylog.InfoNoId("Use input string: %s", DefConfig.InStr)
	}

	//	xylog.Debug("in : %x", in)
	out, err = crypt(in, DefConfig.Encrypt)
	if err != nil {
		xylog.ErrorNoId("Error decrypt/encrypt input data: ", err.Error())
		return
	}
	//	xylog.Debug("out: %x", out)

	if DefConfig.OutFile != "" {
		err = ioutil.WriteFile(DefConfig.OutFile, out, 0644)
		if err != nil {
			xylog.ErrorNoId("Error writing output file (%s) : %s", DefConfig.OutFile, err.Error())
			return
		} else {
			xylog.InfoNoId("Writing to file: %s", DefConfig.OutFile)
		}
	}
	xylog.InfoNoId("done")
}
