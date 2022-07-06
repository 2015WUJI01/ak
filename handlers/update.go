package handlers

import (
	"fmt"
	"log"
	"main/pkg/help"
	"main/pkg/logger"
	"main/ri"
	"main/scheduler"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type updateHelper struct{}

// IncrUpdate 增量更新
func IncrUpdate(c *scheduler.Context) {
	// 获取参数
	interval, oprs := parseArgs(c.PretreatedMessage)

	// -t 参数
	if interval == 0 {
		interval = 15 * time.Minute
	}

	// -o 参数，有则只更新指定的干员，否则更新新增的干员
	if len(oprs) > 0 {
		var oprNames []string
		for _, opr := range oprs {
			oprNames = append(oprNames, opr.Name)
		}
	} else {
		// 增量更新干员名称
		oprs = ri.OC.UpdateOperatorsName()
		logger.Info("update all:")
	}

	var err error
	start := time.Now()
	_, _ = c.Reply("开始批量增量更新...")

	// 更新 items 数据
	err = ri.UpdateItemsData()
	if err != nil {
		_, _ = c.Reply(err.Error())
	}

	// 若当前无指定干员也无新增干员，则全部更新
	if len(oprs) == 0 {
		oprs = ri.OC.Operators
	}

	// 依次更新干员数据
	asyncUpdateOprs(oprs)
	// 刷新 OC
	_ = ri.OC.Update()
	// _, _ = c.Reply("批量增量更新完成")
	_, _ = c.Reply(fmt.Sprintf("批量增量更新完成，更新耗时 %.2fs", help.SpendSeconds(start)))
}

func asyncUpdateOprs(oprs []ri.Operator) {
	length := len(oprs)
	var (
		wg sync.WaitGroup
		mx = &sync.Mutex{}
		ch = make(chan struct{}, 30)
	)

	var (
		current = 0
		all     = len(oprs)
	)
	for i := length - 1; i >= 0; i-- {
		wg.Add(1)
		ch <- struct{}{}
		go func(oprs []ri.Operator, i int, mx *sync.Mutex) {
			_ = oprs[i].Fetch()
			oprs[i].Update()
			mx.Lock()
			current++
			logger.Infof("[%d/%d] No.%d %s has been updated", current, all, oprs[i].UUID, oprs[i].Name)
			mx.Unlock()
			<-ch
			wg.Done()
		}(oprs, i, mx)
	}
	wg.Wait()
	logger.Infof("all operators has been updated")
}

func parseArgs(msg string) (interval time.Duration, oprs []ri.Operator) {

	if r := regexp.MustCompile("-i\\s?([^\\s]+)"); r.MatchString(msg) {
		res := r.FindStringSubmatch(msg)
		i, _ := strconv.Atoi(res[1])
		interval = time.Duration(i) * time.Minute
	}

	if r := regexp.MustCompile("-o\\s?([^\\s]+)"); r.MatchString(msg) {
		res := r.FindStringSubmatch(msg)
		for _, s := range strings.Split(res[1], "，") {
			for _, name := range strings.Split(s, ",") {
				var opr ri.Operator
				var ok bool
				if uuid, err := strconv.Atoi(name); err != nil {
					opr, ok = ri.OC.FindOprByAlias(name)
				} else {
					opr, ok = ri.OC.FindOprByUuid(uuid)
				}
				if ok {
					oprs = append(oprs, opr)
				} else {
					log.Println("参数有误，暂无「" + name + "」相关的 uuid 或干员名称，查询失败")
				}
			}
		}
	}

	return
}
