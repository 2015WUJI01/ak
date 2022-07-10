package cmd

const logo = `
     _         _          _       _     _       
    / \   _ __| | ___ __ (_) __ _| |__ | |_ ___ 
   / _ \ | '__| |/ / '_ \| |/ _` + "`" + ` | '_ \| __/ __|
  / ___ \| |  |   <| | | | | (_| | | | | |_\__ \
 /_/   \_\_|  |_|\_\_| |_|_|\__, |_| |_|\__|___/
                            |___/
`
const version = "0.1.0"
const content = `
  - 删除了大部分用不到的文件，目录变得干净啦~
`

func versionMsg() string {
	msg := "ak version: " + version
	msg += logo
	msg += "update log: " + content
	return msg
}
