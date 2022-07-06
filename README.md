# arknight-cli

明日方舟查询助手（控制台版）
```
NAME

       ak - arknights tool for cli

SYNOPSIS

       ak <command> [<args>] 

DESCRIPTION

       ...

OPTIONS

       version
           输出 ak 的版本

       help
           输出使用说明

       name <value>
           需要查询的干员名称
       
       alias <value>
           干员的别名，若 --name 的值为空，或找不到该干员（比如名字输错了），会使用别名进行查询

       -A, --auto-correct-name
           自动矫正名字和别名

       -m, --module[=<value>]
           查询干员的模组，默认全部查询，可以指定模组的编号。
           例如，令的第一个模组名为「诗短梦长」，代号「sumy」，因此可以使用如下方式进行查询：
              
              ak -m=all -n 令     // 查询所有的模组（默认）
              ak -m=1 -n 令       // 查询第一个模组
              ak -m=y -n 令       // 查询 y 号模组，令的模组名为
              ak -m=诗短梦长 -n 令 // 查询 y 号模组，

       --level=[value]
           指定干员等级状态，value 的值为 4 位数字，290[1]
                      第 1 位表示干员精英化状态，有效值：0,1,2
                      第 2-3 位表示干员等级，有效值：01-90，
           (optional) 第 4 位表示潜能，有效值：1-6，默认为 1，即一潜

           ak --level=2901 令       // 表示 精二 90级 一潜 令          

       -s, --skill
           表示查询干员技能，可以指定几技能，以及等级

           // 查询某干员 [一技能] [Rank 7] 的描述
           // 查询某干员 [一技能] [Rank 1-7] 的描述
           // 查询某干员 [一技能] [Rank 1-7] 的描述
           // 查询某干员 [一技能] [Rank 1-7] 的描述
           // 查询某干员 [一技能] [Rank 1-7] 的描述
           ak -s
    
    // 查询最符合要求的首位干员所有数据
    ak opr XXX // 默认查询支持 ID name，不支持别名，默认为 --id --name 参数
    ak opr 斯卡蒂 // 斯卡蒂，不会查询出浊心斯卡蒂
    ak opr --name 斯卡蒂 // 精准查找
    ak opr --id 123 // 查出 ID 为 123 的干员
    ak opr --alias 42 // 查出别名为 42 的史尔特尔
    ak opr --alias 小车 // 查出近卫小车、医疗小车等所有小车其中最符号要求的其中之一
    ak opr -ia 42 // 会查出 ID 为 42 的干员，别名为 42 的史尔特尔因为别名优先级较低不会显示，若想同时显示两位干员需要使用查询多个干员的关键字 oprs
    
    // 查询符合要求的所有干员所有数据
    ak oprs XXX // 会查询出与 XXX 相匹配的所有干员数据
    ak oprs 斯卡蒂 // 斯卡蒂、浊心斯卡蒂
    ak oprs --alias 小车 // 查出所有小车
    ak oprs -ia 42 // 查出 ID 为 42 的干员，以及别名为 42 的史尔特尔
    
    // 查询技能
    ak skill XXX // 查询某干员的所有技能，参数为干员 ID 或 name 与 opr 相同，默认 -in
    ak skill --id 42 // 查询 ID 为 42 的干员的所有技能
    ak skill --name 史尔特尔 // 查询 史尔特尔 的所有技能
    ak skill --alias 42 // 查询 史尔特尔 的所有技能
    ak skill --order X XXX // 查询某干员的第 X 个技能，会列出该技能所有等级（R1-7、M1-3）的描述
    ak skill --order 3 史尔特尔 // 查询 史尔特尔 的第 3 个技能
    ak skill --rank X XXX // 查询干员指定技能等级的所有技能描述，未专精用数字 1-7 表示 1-7 级，专精用 m1-m3 表示，或用 a,b,c 表示，无视大小写
    ak skill --rank 7 令 // 查询令所有技能等级 7 级的技能描述
    ak skill --rank m1 令 // 查询令所有技能等级专精 1 级的技能描述
    ak skill --rank a 令 // 查询令所有技能等级专精 1 级的技能描述
    ak skill --order 3 --rank C 令 // 查询令三技能等级专精 3 级的所有技能描述
    
    // 查询技能升级材料
    ak skill-upgrade XXX // 同技能，不过查询的是升级材料，必须要带上 --rank 参数指定等级，不然默认是 --rank 1c 表示从 1 级到专 3 的所有升级材料，--order 默认为 all
    ak skill-upgrade --rank 67 令 // 查询令技能从 6 级升到 7 级所需要的升级材料 
    ak skill-upgrade --rank 7 令 // 上述命令可以简化，省略第一个等级，表示达到该等级所需材料
    ak skill-upgrade --rank 17 令 // 可以跨多级查询，查询令技能从 1 升到 7 级所需要的所有升级材料
    ak skill-upgrade --rank 17 --squeeze 令 // 同上，但相同材料会合并在一起
    ak skill-upgrade --order 1 --rank c 令 // 可以指定几技能，查询令一技能专 2 到专 3 的升级材料
    ak skill-upgrade --rank ac 令 // 涉及到专精技能时，如果不指定是哪一个技能，则会有多个技能的材料混在一起，这样有歧义，所以不支持这样查询。有专精技能时一定要指定是哪一个技能
    // 查询技能专精材料
    ak skill-master X X XXX // 需要三个参数，分别是 几技能、专精几级、干员名称，skill-master 是 skill-upgrade --order X --rank --in XXX 的固定格式，所以只需要按照顺序填写参数即可
    ak skill-master 2 c 令 // 查询令二技能专三的专精材料
    
    // 查询模组
    ak module
    
```

但是在 QQ、浏览器等其他输入框使用中文进行查询的时候，应该用怎样的规则呢？

例如：`ak skill-upgrade --rank 17 --squeeze 令` 查询 令 1-7 级 技能升级所需的材料。
上述可以简化为 `ak skill-upgrade -rs 17 令`
中文或许可以使用
```
#干员 令
#干员 别名 42
#干员 别名 ID 
#技能升级 R 17 合计 令
```

