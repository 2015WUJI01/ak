package cmd

const logo = `
     _         _          _       _     _       
    / \   _ __| | ___ __ (_) __ _| |__ | |_ ___ 
   / _ \ | '__| |/ / '_ \| |/ _` + "`" + ` | '_ \| __/ __|
  / ___ \| |  |   <| | | | | (_| | | | | |_\__ \
 /_/   \_\_|  |_|\_\_| |_|_|\__, |_| |_|\__|___/
                            |___/
`
const version = "0.0.2"
const content = `
  - feat: 新增 update、opr、oprs、alias 四个命令的功能，具体使用可以在 ak [command] help 中查看
  - feat: 当当！弄了一个 logo ！
  - feat: cobra 自带的 help 命令习惯之后感觉挺好用的，新手可能需要点时间适应
`

func versionMsg() string {
	msg := "ak version: " + version
	msg += logo
	msg += "update log: " + content
	return msg
}
